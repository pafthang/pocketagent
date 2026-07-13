package models

type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
)

type Task struct {
	ID          string            `json:"id"`
	AgentID     string            `json:"agent_id"`
	Prompt      string            `json:"prompt"`
	Status      TaskStatus        `json:"status"`
	Result      string            `json:"result,omitempty"`
	Error       string            `json:"error,omitempty"`
	ParentID    *string           `json:"parent_id,omitempty"` // for hierarchical tasks
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}
