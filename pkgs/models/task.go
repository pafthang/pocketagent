package models

type TaskStatus string

const (
	TaskQueued    TaskStatus = "queued"
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
	TaskDegraded  TaskStatus = "degraded"
	TaskFailed    TaskStatus = "failed"
	TaskCancelled TaskStatus = "cancelled"
)

// IsTerminal reports whether a task status is finished.
func (s TaskStatus) IsTerminal() bool {
	switch s {
	case TaskCompleted, TaskDegraded, TaskFailed, TaskCancelled:
		return true
	default:
		return false
	}
}

type Task struct {
	ID            string     `json:"id"`
	CorrelationID string     `json:"correlation_id,omitempty"`
	SpaceID       string     `json:"space_id,omitempty"`
	UserID        string     `json:"user_id,omitempty"`
	AgentID       string     `json:"agent_id,omitempty"`
	Prompt        string     `json:"prompt"`
	Status        TaskStatus `json:"status,omitempty"`
	Result        string     `json:"result,omitempty"`
	Error         string     `json:"error,omitempty"`
	ParentID        *string  `json:"parent_id,omitempty"`
	Workflow        string   `json:"workflow,omitempty"`
	WorkerAgentIDs  []string `json:"worker_agent_ids,omitempty"`
	Tools           []string `json:"tools,omitempty"`
	SkillID         string   `json:"skill_id,omitempty"`
	CreatedAt       string   `json:"created_at,omitempty"`
	UpdatedAt     string     `json:"updated_at,omitempty"`
}
