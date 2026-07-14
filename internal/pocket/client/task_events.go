package client

import (
	"fmt"

	"github.com/pafthang/pocketagent/pkgs/models"
)



// CreateTaskEventRecord persists a task progress event.
func (c *Client) CreateTaskEventRecord(event models.TaskEvent) (models.StoredTaskEvent, error) {
	data := map[string]interface{}{
		"space_id":   event.SpaceID,
		"task_id":    event.TaskID,
		"event_type": event.Type,
	}
	if event.Status != "" {
		data["status"] = event.Status
	}
	if event.Step != 0 {
		data["step"] = event.Step
	}
	if event.Message != "" {
		data["message"] = event.Message
	}
	if event.Result != "" {
		data["result"] = event.Result
	}
	record, err := c.CreateRecord(TaskEventsCollection, data)
	if err != nil {
		return models.StoredTaskEvent{}, err
	}
	return taskEventFromRecord(record), nil
}

// ListTaskEvents returns task events for a space.
func (c *Client) ListTaskEvents(spaceID string, opts ListOptions) ([]models.StoredTaskEvent, int, error) {
	filter := fmt.Sprintf("space_id = %q", spaceID)
	if opts.Filter != "" {
		filter = filter + " && (" + opts.Filter + ")"
	}
	opts.Filter = filter
	records, total, err := c.ListRecordsOpts(TaskEventsCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.StoredTaskEvent, 0, len(records))
	for _, record := range records {
		out = append(out, taskEventFromRecord(record))
	}
	return out, total, nil
}

func taskEventFromRecord(record map[string]interface{}) models.StoredTaskEvent {
	step := 0
	if v, ok := record["step"]; ok {
		switch n := v.(type) {
		case float64:
			step = int(n)
		case int:
			step = n
		}
	}
	return models.StoredTaskEvent{
		ID:        stringField(record, "id"),
		SpaceID:   stringField(record, "space_id"),
		TaskID:    stringField(record, "task_id"),
		EventType: stringField(record, "event_type"),
		Status:    stringField(record, "status"),
		Step:      step,
		Message:   stringField(record, "message"),
		Result:    stringField(record, "result"),
		CreatedAt: stringField(record, "created"),
	}
}
