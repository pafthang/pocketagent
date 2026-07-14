package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	memoclient "github.com/pafthang/pocketagent/internal/memo/client"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	"github.com/pafthang/pocketagent/internal/task/decompose"
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pafthang/pocketagent/pkgs/service"
)

const cancelPollInterval = 2 * time.Second

// Deps holds orchestrator handler dependencies.
type Deps struct {
	Worker     *service.Worker
	Store      *Store
	Decomposer *decompose.Decomposer
	Memory     *memoclient.Client
	Embedder   memoclient.Embedder
	TimeoutSec int
}

// Handler returns a NATS consumer for high-level task orchestration.
func Handler(d Deps) func(ctx context.Context, task models.Task) error {
	return func(ctx context.Context, task models.Task) error {
		logger := common.LogWithCorrelation(d.Worker.Log, ctx)
		corrID := taskCorrelationID(ctx, task)

		if task.Prompt == "" {
			d.Store.markFailed(corrID, "empty task prompt")
			return fmt.Errorf("empty task prompt")
		}
		if err := common.GuardPrompt(logger, common.LoadPromptGuardConfig(), task.Prompt); err != nil {
			d.Store.markFailed(corrID, err.Error())
			return err
		}
		if corrID == "" {
			return fmt.Errorf("missing correlation id")
		}

		if d.Store.isCancelled(corrID) {
			logger.Info("skipping cancelled task", "correlation_id", corrID)
			return nil
		}

		d.Store.markRunning(corrID)

		logger.Info("orchestrating task", "prompt", task.Prompt, "agent_id", task.AgentID, "space_id", task.SpaceID)

		_ = d.Worker.PublishEvent(ctx, corrID, spaceEvent(task, models.NewTaskEvent(
			corrID, models.EventOrchestrating, "running",
			"decomposing into subtasks",
		)))

		plans := d.Decomposer.Plan(ctx, task)

		_ = d.Worker.PublishEvent(ctx, corrID, spaceEvent(task, models.TaskEvent{
			TaskID:  corrID,
			Type:    models.EventOrchestrating,
			Status:  "running",
			Message: fmt.Sprintf("split into %d subtask(s)", len(plans)),
		}))

		results := make(map[int]string)
		var mu sync.Mutex
		var wg sync.WaitGroup

		sub, err := d.Worker.Subscribe(natsclient.SubjectResults, func(rctx context.Context, msg *nats.Msg) {
			resultCorrID := common.GetCorrelationID(rctx)
			idx := common.SubtaskIndex(corrID, resultCorrID)
			if idx < 0 {
				_ = msg.Ack()
				return
			}

			result := string(msg.Data)

			mu.Lock()
			results[idx] = result
			mu.Unlock()
			wg.Done()

			d.Store.markSubtaskCompleted(resultCorrID, result)

			_ = d.Worker.PublishEvent(ctx, corrID, spaceEvent(task, models.TaskEvent{
				TaskID:  corrID,
				Type:    models.EventSubtaskResult,
				Status:  "ok",
				Step:    idx,
				Result:  result,
				Message: fmt.Sprintf("subtask %d result received", idx),
			}))
			_ = msg.Ack()
		})
		if err != nil {
			return fmt.Errorf("subscribe results: %w", err)
		}
		defer sub.Unsubscribe()

		for i, plan := range plans {
			wg.Add(1)

			subCorrID := subtaskCorrelationID(corrID, i)
			subCtx := common.WithCorrelationID(ctx, subCorrID)

			subtask := models.Task{
				CorrelationID: subCorrID,
				SpaceID:       task.SpaceID,
				AgentID:       plan.AgentID,
				Prompt:        plan.Prompt,
			}
			d.Store.createSubtask(corrID, subtask)

			subject := natsclient.SubjectTasksPrefix + subCorrID

			if err := d.Worker.PublishJSON(subCtx, subject, subtask); err != nil {
				wg.Done()
				d.Store.markSubtaskFailed(subCorrID, err.Error())
				logger.Warn("failed to publish subtask", "index", i, "error", err)
				continue
			}

			eventType := models.EventSubtaskDispatched
			message := fmt.Sprintf("dispatched subtask %d", i)
			if plan.AgentID != "" && plan.AgentID != task.AgentID {
				eventType = models.EventSupervisorDelegated
				message = fmt.Sprintf("delegated subtask %d to agent %s", i, plan.AgentID)
			}

			_ = d.Worker.PublishEvent(ctx, corrID, spaceEvent(task, models.TaskEvent{
				TaskID:  corrID,
				Type:    eventType,
				Status:  "running",
				Step:    i,
				Message: message,
			}))
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		timeout := time.Duration(d.TimeoutSec) * time.Second
		if timeout <= 0 {
			timeout = 30 * time.Second
		}

		deadline := time.After(timeout)
		cancelPoll := time.NewTicker(cancelPollInterval)
		defer cancelPoll.Stop()

		timedOut := false
		cancelled := false

	waitLoop:
		for {
			select {
			case <-done:
				logger.Info("all subtasks completed", "count", len(plans))
				break waitLoop
			case <-cancelPoll.C:
				if d.Store.isCancelled(corrID) {
					cancelled = true
					logger.Info("task cancelled during orchestration", "correlation_id", corrID)
					break waitLoop
				}
			case <-deadline:
				timedOut = true
				logger.Warn("subtask timeout", "received", len(results), "expected", len(plans))
				_ = d.Worker.PublishEvent(ctx, corrID, spaceEvent(task, models.NewTaskEvent(
					corrID, models.EventTimeout, "degraded",
					fmt.Sprintf("timeout: %d/%d results", len(results), len(plans)),
				)))
				break waitLoop
			}
		}

		if cancelled {
			cancelledCount := d.Store.cancelPendingSubtasks(corrID)
			d.Store.markCancelled(corrID, "cancelled")
			_ = d.Worker.PublishEvent(ctx, corrID, spaceEvent(task, models.NewTaskEvent(
				corrID, models.EventCancelled, "cancelled",
				fmt.Sprintf("task cancelled (%d pending subtask(s) stopped)", cancelledCount),
			)))
			return nil
		}

		final := buildFinalAnswer(results)

		if err := d.Memory.StoreScopedWithMeta(ctx, d.Embedder, task.SpaceID, corrID, final, map[string]string{
			"source":   "task_result",
			"task_id":  corrID,
			"space_id": task.SpaceID,
			"agent_id": task.AgentID,
			"prompt":   truncateMeta(task.Prompt, 200),
		}); err != nil {
			logger.Warn("failed to save to memory", "error", err)
		} else {
			logger.Info("final answer saved to memory", "answer", final)
		}

		eventType := models.EventCompleted
		status := "completed"
		taskStatus := models.TaskCompleted
		if timedOut {
			status = "degraded"
			taskStatus = models.TaskDegraded
		}

		d.Store.markCompleted(corrID, final, taskStatus)

		_ = d.Worker.PublishEvent(ctx, corrID, spaceEvent(task, models.TaskEvent{
			TaskID:  corrID,
			Type:    eventType,
			Status:  status,
			Result:  final,
			Message: "orchestration finished",
		}))

		return nil
	}
}

func spaceEvent(task models.Task, event models.TaskEvent) models.TaskEvent {
	event.SpaceID = task.SpaceID
	return event
}

func taskCorrelationID(ctx context.Context, task models.Task) string {
	if task.CorrelationID != "" {
		return task.CorrelationID
	}
	return common.GetCorrelationID(ctx)
}
