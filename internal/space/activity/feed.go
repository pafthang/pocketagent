package activity

import (
	"sort"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/models"
)

const (
	SourceTask  = "task"
	SourceAudit = "audit"
)

// BuildFeed merges task events and audit logs into a sorted activity list.
func BuildFeed(taskEvents []models.StoredTaskEvent, audits []models.AuditLog, includeAudit bool, limit int) []models.ActivityEntry {
	capacity := len(taskEvents)
	if includeAudit {
		capacity += len(audits)
	}
	entries := make([]models.ActivityEntry, 0, capacity)

	for _, event := range taskEvents {
		entries = append(entries, FromTaskEvent(event))
	}
	if includeAudit {
		for _, log := range audits {
			entries = append(entries, FromAuditLog(log))
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp > entries[j].Timestamp
	})

	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}
	return entries
}

// FromTaskEvent maps a stored task event to a feed entry.
func FromTaskEvent(event models.StoredTaskEvent) models.ActivityEntry {
	entryType, content := mapTaskEvent(event)
	data := map[string]interface{}{
		"event_type": event.EventType,
		"status":     event.Status,
	}
	if event.Step != 0 {
		data["step"] = event.Step
	}
	if event.Result != "" {
		data["result"] = event.Result
	}
	return models.ActivityEntry{
		ID:        "task-" + event.ID,
		Type:      entryType,
		Content:   content,
		Data:      data,
		Timestamp: event.CreatedAt,
		Source:    SourceTask,
		TaskID:    event.TaskID,
	}
}

// FromAuditLog maps an audit log to a feed entry.
func FromAuditLog(log models.AuditLog) models.ActivityEntry {
	content := log.Action
	if log.ActorEmail != "" {
		content = log.ActorEmail + ": " + log.Action
	}
	data := map[string]interface{}{
		"action": log.Action,
	}
	if log.ActorID != "" {
		data["actor_id"] = log.ActorID
	}
	if log.ActorEmail != "" {
		data["actor_email"] = log.ActorEmail
	}
	if log.ResourceType != "" {
		data["resource_type"] = log.ResourceType
	}
	if log.ResourceID != "" {
		data["resource_id"] = log.ResourceID
	}
	if log.Metadata != nil {
		data["metadata"] = log.Metadata
	}
	return models.ActivityEntry{
		ID:        "audit-" + log.ID,
		Type:      "status",
		Content:   content,
		Data:      data,
		Timestamp: log.CreatedAt,
		Source:    SourceAudit,
	}
}

func mapTaskEvent(event models.StoredTaskEvent) (entryType, content string) {
	switch event.EventType {
	case models.EventLLMToken:
		entryType = "thinking"
		content = event.Message
	case models.EventFailed:
		entryType = "error"
		content = firstNonEmpty(event.Message, event.Result, event.Status)
	case models.EventSubtaskResult:
		entryType = "tool_result"
		content = truncate(firstNonEmpty(event.Result, event.Message), 120)
	default:
		entryType = "status"
		content = firstNonEmpty(event.Message, event.Result, event.Status, event.EventType)
	}
	content = strings.TrimSpace(content)
	if content == "" {
		content = event.EventType
	}
	return entryType, content
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
