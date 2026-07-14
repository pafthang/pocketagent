package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/pafthang/pocketagent/internal/exec/mcp"
	"github.com/pafthang/pocketagent/internal/exec/react"
	"github.com/pafthang/pocketagent/internal/exec/tools"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pafthang/pocketagent/pkgs/ollama"
	"github.com/pafthang/pocketagent/pkgs/service"
)

// Deps holds NATS task handler dependencies.
type Deps struct {
	Executor *react.Executor
	Pocket   *pbclient.Client
	ToolCfg  tools.Config
}

// Task returns the exec NATS consumer for ReAct task execution.
func Task(w *service.Worker, deps Deps) func(ctx context.Context, task models.Task) error {
	return func(ctx context.Context, task models.Task) error {
		logger := common.LogWithCorrelation(w.Log, ctx)
		start := time.Now()
		status := "ok"
		defer func() {
			common.TaskDuration.WithLabelValues("exec", status).Observe(time.Since(start).Seconds())
		}()

		corrID := common.GetCorrelationID(ctx)
		rootID := common.RootCorrelationID(corrID)
		step := common.SubtaskIndex(rootID, corrID)

		if cancelled, err := isTaskCancelled(deps.Pocket, corrID, rootID); err != nil {
			status = "error"
			return fmt.Errorf("check task status: %w", err)
		} else if cancelled {
			logger.Info("skipping cancelled subtask", "correlation_id", corrID, "root_id", rootID)
			_ = w.PublishEvent(ctx, rootID, spaceEvent(task, models.TaskEvent{
				TaskID: rootID, Type: models.EventCancelled, Status: "cancelled",
				Step: step, Message: "subtask skipped (cancelled)",
			}))
			if corrID != rootID {
				_, _ = deps.Pocket.UpdateTaskByCorrelationID(corrID, models.Task{
					Status: models.TaskCancelled,
					Error:  "cancelled by parent",
				})
			}
			return nil
		}

		if corrID != rootID {
			_, _ = deps.Pocket.UpdateTaskByCorrelationID(corrID, models.Task{Status: models.TaskRunning})
		}

		if task.Prompt == "" {
			status = "error"
			_ = w.PublishEvent(ctx, rootID, spaceEvent(task, models.TaskEvent{
				TaskID: rootID, Type: models.EventFailed, Status: "error",
				Step: step, Message: "empty task prompt",
			}))
			return fmt.Errorf("empty task prompt")
		}
		if err := common.GuardPrompt(logger, common.LoadPromptGuardConfig(), task.Prompt); err != nil {
			status = "error"
			_ = w.PublishEvent(ctx, rootID, spaceEvent(task, models.TaskEvent{
				TaskID: rootID, Type: models.EventFailed, Status: "error",
				Step: step, Message: err.Error(),
			}))
			if corrID != rootID {
				_, _ = deps.Pocket.UpdateTaskByCorrelationID(corrID, models.Task{
					Status: models.TaskFailed,
					Error:  err.Error(),
				})
			}
			return err
		}

		_ = w.PublishEvent(ctx, rootID, spaceEvent(task, models.TaskEvent{
			TaskID: rootID, Type: models.EventSubtaskStarted, Status: "running",
			Step: step, Message: fmt.Sprintf("ReAct started: %s", truncate(task.Prompt, 80)),
		}))

		agent, err := resolveAgent(deps.Pocket, task)
		if err != nil {
			status = "error"
			_ = w.PublishEvent(ctx, rootID, spaceEvent(task, models.TaskEvent{
				TaskID: rootID, Type: models.EventFailed, Status: "error",
				Step: step, Message: err.Error(),
			}))
			return err
		}

		logger.Info("starting ReAct task",
			"prompt", task.Prompt,
			"agent_id", task.AgentID,
			"space_id", task.SpaceID,
			"model", agentModelName(agent, deps.Executor.LLMModel),
		)

		if deps.Executor.StreamLLM {
			deps.Executor.OnStream = func(chunk ollama.ChatStreamEvent) {
				delta := chunk.Message.Content
				if delta == "" {
					delta = chunk.Message.Thinking
				}
				if delta == "" {
					return
				}
				_ = w.PublishEvent(ctx, rootID, spaceEvent(task, models.TaskEvent{
					TaskID:  rootID,
					Type:    models.EventLLMToken,
					Status:  "streaming",
					Step:    step,
					Message: delta,
				}))
			}
		}

		var overlay *react.ToolOverlay
		spaceServers, err := mcp.LoadSpaceServers(deps.Pocket, task.SpaceID)
		if err != nil {
			status = "error"
			_ = w.PublishEvent(ctx, rootID, spaceEvent(task, models.TaskEvent{
				TaskID: rootID, Type: models.EventFailed, Status: "error",
				Step: step, Message: err.Error(),
			}))
			return fmt.Errorf("load space mcp servers: %w", err)
		}
		if len(spaceServers) > 0 {
			spaceSet := tools.BuildMCPOnly(spaceServers)
			defer spaceSet.Close()
			overlay = &react.ToolOverlay{
				Catalog: spaceSet.Catalog,
				Run:     react.ToolRunner(spaceSet.Registry),
			}
		}

		userProfile := ""
		if task.SpaceID != "" && task.UserID != "" {
			if profile, err := deps.Pocket.GetSpaceProfile(task.SpaceID, task.UserID); err != nil {
				logger.Warn("load space profile", "error", err, "user_id", task.UserID)
			} else {
				userProfile = profile.Content
			}
		}

		result, err := deps.Executor.ExecuteWithOverlay(ctx, task, agent, overlay, userProfile)
		if err != nil {
			status = "error"
			_ = w.PublishEvent(ctx, rootID, spaceEvent(task, models.TaskEvent{
				TaskID: rootID, Type: models.EventFailed, Status: "error",
				Step: step, Message: err.Error(),
			}))
			if corrID != rootID {
				_, _ = deps.Pocket.UpdateTaskByCorrelationID(corrID, models.Task{
					Status: models.TaskFailed,
					Error:  err.Error(),
				})
			}
			return fmt.Errorf("react execution: %w", err)
		}

		if cancelled, err := isTaskCancelled(deps.Pocket, corrID, rootID); err != nil {
			status = "error"
			return fmt.Errorf("check task status: %w", err)
		} else if cancelled {
			logger.Info("aborting cancelled subtask before publish", "correlation_id", corrID)
			_ = w.PublishEvent(ctx, rootID, spaceEvent(task, models.TaskEvent{
				TaskID: rootID, Type: models.EventCancelled, Status: "cancelled",
				Step: step, Message: "subtask aborted (cancelled)",
			}))
			if corrID != rootID {
				_, _ = deps.Pocket.UpdateTaskByCorrelationID(corrID, models.Task{
					Status: models.TaskCancelled,
					Error:  "cancelled by parent",
				})
			}
			return nil
		}

		resultKey := corrID
		if resultKey == "" {
			resultKey = task.AgentID
		}
		if resultKey == "" {
			resultKey = "unknown"
		}

		if err := w.Publish(natsclient.SubjectResultsPrefix+resultKey, []byte(result.FinalAnswer)); err != nil {
			status = "error"
			_ = w.PublishEvent(ctx, rootID, spaceEvent(task, models.TaskEvent{
				TaskID: rootID, Type: models.EventFailed, Status: "error",
				Step: step, Message: err.Error(),
			}))
			return fmt.Errorf("publish result: %w", err)
		}

		_ = w.PublishEvent(ctx, rootID, spaceEvent(task, models.TaskEvent{
			TaskID: rootID, Type: models.EventSubtaskCompleted, Status: "ok",
			Step: step, Result: result.FinalAnswer,
			Message: fmt.Sprintf("ReAct done (%d steps, %d tools)", result.Steps, len(result.ToolCalls)),
		}))

		logger.Info("ReAct completed",
			"final_answer", result.FinalAnswer,
			"steps", result.Steps,
			"tool_calls", len(result.ToolCalls),
		)

		return nil
	}
}

func agentModelName(agent models.Agent, fallback string) string {
	if agent.Model != "" {
		return agent.Model
	}
	return fallback
}

func isTaskCancelled(pb *pbclient.Client, corrID, rootID string) (bool, error) {
	if cancelled, err := pb.IsTaskCancelled(rootID); err != nil {
		return false, err
	} else if cancelled {
		return true, nil
	}
	if corrID != rootID {
		return pb.IsTaskCancelled(corrID)
	}
	return false, nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func spaceEvent(task models.Task, event models.TaskEvent) models.TaskEvent {
	event.SpaceID = task.SpaceID
	return event
}