package orchestrator

import (
	"fmt"
	"log/slog"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// Store persists task lifecycle updates in PocketBase.
type Store struct {
	pb  *pbclient.Client
	log *slog.Logger
}

func NewStore(pb *pbclient.Client, log *slog.Logger) *Store {
	return &Store{pb: pb, log: log}
}

func (s *Store) markRunning(correlationID string) {
	s.update(correlationID, models.Task{Status: models.TaskRunning})
}

func (s *Store) markCompleted(correlationID, result string, status models.TaskStatus) {
	s.update(correlationID, models.Task{Status: status, Result: result})
}

func (s *Store) markFailed(correlationID, message string) {
	s.update(correlationID, models.Task{Status: models.TaskFailed, Error: message})
}

func (s *Store) markCancelled(correlationID, message string) {
	s.update(correlationID, models.Task{Status: models.TaskCancelled, Error: message})
}

func (s *Store) createSubtask(parentCorrID string, sub models.Task) {
	if s == nil || s.pb == nil || parentCorrID == "" || sub.CorrelationID == "" {
		return
	}
	parentID := parentCorrID
	sub.ParentID = &parentID
	sub.Status = models.TaskQueued
	if _, err := s.pb.CreateTask(sub); err != nil && s.log != nil {
		s.log.Warn("failed to create subtask record",
			"parent_id", parentCorrID,
			"correlation_id", sub.CorrelationID,
			"error", err,
		)
	}
}

func (s *Store) markSubtaskCompleted(correlationID, result string) {
	s.update(correlationID, models.Task{Status: models.TaskCompleted, Result: result})
}

func (s *Store) markSubtaskFailed(correlationID, message string) {
	s.update(correlationID, models.Task{Status: models.TaskFailed, Error: message})
}

func (s *Store) markSubtaskCancelled(correlationID, message string) {
	s.update(correlationID, models.Task{Status: models.TaskCancelled, Error: message})
}

func (s *Store) isCancelled(correlationID string) bool {
	if correlationID == "" || s == nil || s.pb == nil {
		return false
	}
	task, err := s.pb.GetTaskByCorrelationID(correlationID)
	if err != nil {
		return false
	}
	return task.Status == models.TaskCancelled
}

func (s *Store) cancelPendingSubtasks(parentCorrID string) int {
	subtasks, err := s.listSubtasks(parentCorrID)
	if err != nil {
		if s.log != nil {
			s.log.Warn("failed to list subtasks for cancel", "parent_id", parentCorrID, "error", err)
		}
		return 0
	}
	cancelled := 0
	for _, sub := range subtasks {
		if sub.Status == models.TaskQueued || sub.Status == models.TaskRunning {
			s.markSubtaskCancelled(sub.CorrelationID, "cancelled by parent")
			cancelled++
		}
	}
	return cancelled
}

func (s *Store) listSubtasks(parentCorrID string) ([]models.Task, error) {
	if s == nil || s.pb == nil || parentCorrID == "" {
		return nil, fmt.Errorf("store unavailable")
	}
	filter := fmt.Sprintf("parent_id = %q", parentCorrID)
	tasks, _, err := s.pb.ListTasks(pbclient.ListOptions{Page: 1, PerPage: 100, Filter: filter})
	return tasks, err
}

func (s *Store) update(correlationID string, patch models.Task) {
	if correlationID == "" || s == nil || s.pb == nil {
		return
	}
	if _, err := s.pb.UpdateTaskByCorrelationID(correlationID, patch); err != nil && s.log != nil {
		s.log.Warn("failed to persist task update",
			"correlation_id", correlationID,
			"status", patch.Status,
			"error", err,
		)
	}
}
