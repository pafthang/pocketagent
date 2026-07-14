package models

// ActivityEntry is a normalized space activity feed item.
type ActivityEntry struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Content   string                 `json:"content"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp string                 `json:"timestamp"`
	Source    string                 `json:"source,omitempty"`
	TaskID    string                 `json:"task_id,omitempty"`
}

// ActivityListResponse is returned by GET /spaces/:id/activity.
type ActivityListResponse struct {
	Entries []ActivityEntry `json:"entries"`
	Total   int             `json:"total"`
}

// StoredTaskEvent is a persisted task progress event.
type StoredTaskEvent struct {
	ID        string `json:"id"`
	SpaceID   string `json:"space_id"`
	TaskID    string `json:"task_id"`
	EventType string `json:"event_type"`
	Status    string `json:"status,omitempty"`
	Step      int    `json:"step,omitempty"`
	Message   string `json:"message,omitempty"`
	Result    string `json:"result,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}