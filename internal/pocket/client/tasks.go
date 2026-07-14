package client

import (
	"fmt"

	"github.com/pafthang/pocketagent/pkgs/models"
)



// CreateTask stores a new task record.
func (c *Client) CreateTask(task models.Task) (models.Task, error) {
	record, err := c.CreateRecord(TasksCollection, taskRecordData(task))
	if err != nil {
		return models.Task{}, err
	}
	return taskFromRecord(record), nil
}

// GetTask returns a task by PocketBase record ID.
func (c *Client) GetTask(id string) (models.Task, error) {
	record, err := c.GetRecord(TasksCollection, id)
	if err != nil {
		return models.Task{}, err
	}
	return taskFromRecord(record), nil
}

// GetTaskByCorrelationID returns a task by correlation ID.
func (c *Client) GetTaskByCorrelationID(correlationID string) (models.Task, error) {
	filter := fmt.Sprintf("correlation_id = %q", correlationID)
	records, _, err := c.ListRecordsOpts(TasksCollection, ListOptions{Page: 1, PerPage: 1, Filter: filter})
	if err != nil {
		return models.Task{}, err
	}
	if len(records) == 0 {
		return models.Task{}, &APIError{StatusCode: 404, Message: "task not found"}
	}
	return taskFromRecord(records[0]), nil
}

// ListTasks returns tasks with optional filter.
func (c *Client) ListTasks(opts ListOptions) ([]models.Task, int, error) {
	records, total, err := c.ListRecordsOpts(TasksCollection, opts)
	if err != nil {
		return nil, 0, err
	}

	tasks := make([]models.Task, 0, len(records))
	for _, record := range records {
		tasks = append(tasks, taskFromRecord(record))
	}
	return tasks, total, nil
}

// UpdateTask patches a task by record ID.
func (c *Client) UpdateTask(id string, task models.Task) (models.Task, error) {
	record, err := c.UpdateRecord(TasksCollection, id, taskRecordData(task))
	if err != nil {
		return models.Task{}, err
	}
	return taskFromRecord(record), nil
}

// UpdateTaskByCorrelationID patches a task matched by correlation ID.
func (c *Client) UpdateTaskByCorrelationID(correlationID string, task models.Task) (models.Task, error) {
	existing, err := c.GetTaskByCorrelationID(correlationID)
	if err != nil {
		return models.Task{}, err
	}
	return c.UpdateTask(existing.ID, mergeTaskUpdate(existing, task))
}

// IsTaskCancelled reports whether a task with the given correlation ID is cancelled.
func (c *Client) IsTaskCancelled(correlationID string) (bool, error) {
	task, err := c.GetTaskByCorrelationID(correlationID)
	if err != nil {
		return false, err
	}
	return task.Status == models.TaskCancelled, nil
}

// CancelTask marks a root task and its pending subtasks as cancelled.
func (c *Client) CancelTask(correlationID string) (models.Task, error) {
	existing, err := c.GetTaskByCorrelationID(correlationID)
	if err != nil {
		return models.Task{}, err
	}
	if existing.Status.IsTerminal() {
		return models.Task{}, &APIError{StatusCode: 409, Message: "task cannot be cancelled"}
	}

	updated, err := c.UpdateTaskByCorrelationID(correlationID, models.Task{
		Status: models.TaskCancelled,
		Error:  "cancelled by user",
	})
	if err != nil {
		return models.Task{}, err
	}
	if err := c.cancelPendingSubtasks(correlationID); err != nil {
		return updated, err
	}
	return updated, nil
}

func (c *Client) cancelPendingSubtasks(parentCorrID string) error {
	filter := fmt.Sprintf(`parent_id = %q && (status = "queued" || status = "running")`, parentCorrID)
	tasks, _, err := c.ListTasks(ListOptions{Page: 1, PerPage: 100, Filter: filter})
	if err != nil {
		return err
	}
	for _, task := range tasks {
		if _, err := c.UpdateTaskByCorrelationID(task.CorrelationID, models.Task{
			Status: models.TaskCancelled,
			Error:  "cancelled by parent",
		}); err != nil {
			return err
		}
	}
	return nil
}

// ListSubtasks returns child tasks for a parent correlation ID.
func (c *Client) ListSubtasks(parentCorrID string) ([]models.Task, error) {
	filter := fmt.Sprintf("parent_id = %q", parentCorrID)
	tasks, _, err := c.ListTasks(ListOptions{Page: 1, PerPage: 100, Filter: filter})
	return tasks, err
}

func mergeTaskUpdate(existing, patch models.Task) models.Task {
	merged := existing
	if patch.Status != "" {
		merged.Status = patch.Status
	}
	if patch.Result != "" {
		merged.Result = patch.Result
	}
	if patch.Error != "" {
		merged.Error = patch.Error
	}
	if patch.AgentID != "" {
		merged.AgentID = patch.AgentID
	}
	if patch.Prompt != "" {
		merged.Prompt = patch.Prompt
	}
	if patch.Tools != nil {
		merged.Tools = patch.Tools
	}
	if patch.SkillID != "" {
		merged.SkillID = patch.SkillID
	}
	return merged
}

func taskRecordData(task models.Task) map[string]interface{} {
	data := map[string]interface{}{
		"correlation_id": task.CorrelationID,
		"space_id":       task.SpaceID,
		"prompt":         task.Prompt,
	}
	if task.UserID != "" {
		data["user_id"] = task.UserID
	}
	if task.AgentID != "" {
		data["agent_id"] = task.AgentID
	}
	if task.Status != "" {
		data["status"] = string(task.Status)
	}
	if task.Result != "" {
		data["result"] = task.Result
	}
	if task.Error != "" {
		data["error"] = task.Error
	}
	if task.ParentID != nil {
		data["parent_id"] = *task.ParentID
	}
	if task.Workflow != "" {
		data["workflow"] = task.Workflow
	}
	if len(task.WorkerAgentIDs) > 0 {
		data["worker_agent_ids"] = task.WorkerAgentIDs
	}
	if len(task.Tools) > 0 {
		data["tools"] = task.Tools
	}
	if task.SkillID != "" {
		data["skill_id"] = task.SkillID
	}
	return data
}

func taskFromRecord(record map[string]interface{}) models.Task {
	task := models.Task{
		ID:            stringField(record, "id"),
		CorrelationID: stringField(record, "correlation_id"),
		SpaceID:       stringField(record, "space_id"),
		UserID:        stringField(record, "user_id"),
		AgentID:       stringField(record, "agent_id"),
		Prompt:        stringField(record, "prompt"),
		Status:        models.TaskStatus(stringField(record, "status")),
		Result:        stringField(record, "result"),
		Error:         stringField(record, "error"),
		CreatedAt:     stringField(record, "created"),
		UpdatedAt:     stringField(record, "updated"),
	}
	if parentID := stringField(record, "parent_id"); parentID != "" {
		task.ParentID = &parentID
	}
	task.Workflow = stringField(record, "workflow")
	task.WorkerAgentIDs = stringSliceField(record, "worker_agent_ids")
	task.Tools = stringSliceField(record, "tools")
	task.SkillID = stringField(record, "skill_id")
	return task
}
