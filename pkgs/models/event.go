package models

// Task event types streamed to gate WebSocket clients.
const (
	EventConnected          = "connected"
	EventQueued             = "queued"
	EventOrchestrating      = "orchestrating"
	EventSubtaskDispatched  = "subtask_dispatched"
	EventSubtaskStarted     = "subtask_started"
	EventSubtaskCompleted   = "subtask_completed"
	EventSubtaskResult      = "subtask_result"
	EventCompleted          = "completed"
	EventFailed             = "failed"
	EventTimeout            = "timeout"
	EventLLMToken           = "llm_token"
	EventCancelled          = "cancelled"
	EventSupervisorDelegated = "supervisor_delegated"
)

// TaskEvent is a real-time task progress update.
type TaskEvent struct {
	TaskID  string `json:"task_id"`
	SpaceID string `json:"space_id,omitempty"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Step    int    `json:"step,omitempty"`
	Message string `json:"message,omitempty"`
	Result  string `json:"result,omitempty"`
}

// NewTaskEvent builds an event for a root task correlation ID.
func NewTaskEvent(taskID, eventType, status, message string) TaskEvent {
	return TaskEvent{
		TaskID:  taskID,
		Type:    eventType,
		Status:  status,
		Message: message,
	}
}