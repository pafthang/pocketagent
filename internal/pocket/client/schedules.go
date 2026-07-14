package client

import (
	"fmt"
	"time"

	"github.com/pafthang/pocketagent/pkgs/models"
)



// CreateSchedule stores a recurring task schedule.
func (c *Client) CreateSchedule(schedule models.Schedule) (models.Schedule, error) {
	record, err := c.CreateRecord(SchedulesCollection, scheduleRecordData(schedule))
	if err != nil {
		return models.Schedule{}, err
	}
	return scheduleFromRecord(record), nil
}

// GetSchedule returns a schedule by record ID.
func (c *Client) GetSchedule(id string) (models.Schedule, error) {
	record, err := c.GetRecord(SchedulesCollection, id)
	if err != nil {
		return models.Schedule{}, err
	}
	return scheduleFromRecord(record), nil
}

// ListSchedules returns schedules with optional filter.
func (c *Client) ListSchedules(opts ListOptions) ([]models.Schedule, int, error) {
	records, total, err := c.ListRecordsOpts(SchedulesCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.Schedule, 0, len(records))
	for _, record := range records {
		out = append(out, scheduleFromRecord(record))
	}
	return out, total, nil
}

// UpdateSchedule patches a schedule by record ID.
func (c *Client) UpdateSchedule(id string, patch models.Schedule) (models.Schedule, error) {
	existing, err := c.GetSchedule(id)
	if err != nil {
		return models.Schedule{}, err
	}
	record, err := c.UpdateRecord(SchedulesCollection, id, scheduleRecordData(mergeScheduleUpdate(existing, patch)))
	if err != nil {
		return models.Schedule{}, err
	}
	return scheduleFromRecord(record), nil
}

func mergeScheduleUpdate(existing, patch models.Schedule) models.Schedule {
	merged := existing
	if patch.Name != "" {
		merged.Name = patch.Name
	}
	if patch.AgentID != "" {
		merged.AgentID = patch.AgentID
	}
	if patch.Prompt != "" {
		merged.Prompt = patch.Prompt
	}
	if patch.CronExpr != "" {
		merged.CronExpr = patch.CronExpr
	}
	if patch.Workflow != "" {
		merged.Workflow = patch.Workflow
	}
	if patch.WorkerAgentIDs != nil {
		merged.WorkerAgentIDs = patch.WorkerAgentIDs
	}
	if patch.LastRunAt != "" {
		merged.LastRunAt = patch.LastRunAt
	}
	if patch.NextRunAt != "" {
		merged.NextRunAt = patch.NextRunAt
	}
	if patch.LastTaskID != "" {
		merged.LastTaskID = patch.LastTaskID
	}
	return merged
}

// DeleteSchedule removes a schedule.
func (c *Client) DeleteSchedule(id string) error {
	return c.DeleteRecord(SchedulesCollection, id)
}

// ListDueSchedules returns enabled schedules due at or before the given time.
func (c *Client) ListDueSchedules(until time.Time) ([]models.Schedule, error) {
	filter := fmt.Sprintf(`enabled = true && next_run_at <= %q`, until.UTC().Format(time.RFC3339))
	records, _, err := c.ListRecordsOpts(SchedulesCollection, ListOptions{
		Page:    1,
		PerPage: 100,
		Filter:  filter,
	})
	if err != nil {
		return nil, err
	}
	out := make([]models.Schedule, 0, len(records))
	for _, record := range records {
		out = append(out, scheduleFromRecord(record))
	}
	return out, nil
}

func scheduleRecordData(schedule models.Schedule) map[string]interface{} {
	data := map[string]interface{}{
		"space_id":  schedule.SpaceID,
		"name":      schedule.Name,
		"prompt":    schedule.Prompt,
		"cron_expr": schedule.CronExpr,
		"enabled":   schedule.Enabled,
	}
	if schedule.AgentID != "" {
		data["agent_id"] = schedule.AgentID
	}
	if schedule.Workflow != "" {
		data["workflow"] = schedule.Workflow
	}
	if len(schedule.WorkerAgentIDs) > 0 {
		data["worker_agent_ids"] = schedule.WorkerAgentIDs
	}
	if schedule.LastRunAt != "" {
		data["last_run_at"] = schedule.LastRunAt
	}
	if schedule.NextRunAt != "" {
		data["next_run_at"] = schedule.NextRunAt
	}
	if schedule.LastTaskID != "" {
		data["last_task_id"] = schedule.LastTaskID
	}
	return data
}

func scheduleFromRecord(record map[string]interface{}) models.Schedule {
	schedule := models.Schedule{
		ID:         stringField(record, "id"),
		SpaceID:    stringField(record, "space_id"),
		Name:       stringField(record, "name"),
		AgentID:    stringField(record, "agent_id"),
		Prompt:     stringField(record, "prompt"),
		CronExpr:   stringField(record, "cron_expr"),
		Workflow:   stringField(record, "workflow"),
		LastRunAt:  stringField(record, "last_run_at"),
		NextRunAt:  stringField(record, "next_run_at"),
		LastTaskID: stringField(record, "last_task_id"),
		CreatedAt:  stringField(record, "created"),
		UpdatedAt:  stringField(record, "updated"),
	}
	if enabled, ok := record["enabled"].(bool); ok {
		schedule.Enabled = enabled
	}
	schedule.WorkerAgentIDs = stringSliceField(record, "worker_agent_ids")
	return schedule
}
