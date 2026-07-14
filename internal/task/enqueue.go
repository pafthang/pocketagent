package task

import (
	"context"
	"fmt"
	"time"

	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pafthang/pocketagent/pkgs/service"
)

// OrchestratorPublisher enqueues tasks for orchestration.
type OrchestratorPublisher interface {
	PublishOrchestrator(ctx context.Context, task models.Task) error
	PublishEvent(ctx context.Context, corrID string, event models.TaskEvent) error
}

// WorkerPublisher adapts service.Worker for PersistAndEnqueue.
type WorkerPublisher struct {
	Worker *service.Worker
}

func (p WorkerPublisher) PublishOrchestrator(ctx context.Context, task models.Task) error {
	return p.Worker.PublishJSON(ctx, natsclient.SubjectOrchestrator, task)
}

func (p WorkerPublisher) PublishEvent(ctx context.Context, corrID string, event models.TaskEvent) error {
	return p.Worker.PublishEvent(ctx, corrID, event)
}

// PersistAndEnqueue stores a task and publishes it to the orchestrator.
func PersistAndEnqueue(ctx context.Context, pb *pbclient.Client, pub OrchestratorPublisher, task models.Task) (models.Task, string, error) {
	corrID := task.CorrelationID
	if corrID == "" {
		corrID = fmt.Sprintf("task-%d", time.Now().UnixNano())
	}
	task.CorrelationID = corrID
	if task.Status == "" {
		task.Status = models.TaskQueued
	}

	stored, err := pb.CreateTask(task)
	if err != nil {
		return models.Task{}, "", err
	}

	ctx = common.WithCorrelationID(ctx, corrID)
	if err := pub.PublishOrchestrator(ctx, task); err != nil {
		_, _ = pb.UpdateTaskByCorrelationID(corrID, models.Task{
			Status: models.TaskFailed,
			Error:  err.Error(),
		})
		return models.Task{}, "", err
	}

	queued := models.NewTaskEvent(corrID, models.EventQueued, "queued", "task accepted by orchestrator")
	queued.SpaceID = task.SpaceID
	_ = pub.PublishEvent(ctx, corrID, queued)

	return stored, corrID, nil
}

// EnqueueTask persists and publishes a task to the orchestrator via the worker.
func EnqueueTask(ctx context.Context, w *service.Worker, pb *pbclient.Client, task models.Task) (string, error) {
	_, corrID, err := PersistAndEnqueue(ctx, pb, WorkerPublisher{Worker: w}, task)
	return corrID, err
}
