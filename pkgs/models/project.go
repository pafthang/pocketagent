package models

// Project lifecycle statuses (planning workflow).
const (
	ProjectDraft             = "draft"
	ProjectPlanning          = "planning"
	ProjectAwaitingApproval  = "awaiting_approval"
	ProjectApproved        = "approved"
	ProjectExecuting       = "executing"
	ProjectPaused          = "paused"
	ProjectCompleted       = "completed"
	ProjectFailed          = "failed"
	ProjectCancelled       = "cancelled"
)

// ProjectItem kanban statuses (distinct from execution Task.Status).
const (
	ItemInbox       = "inbox"
	ItemAssigned    = "assigned"
	ItemInProgress  = "in_progress"
	ItemReview      = "review"
	ItemDone        = "done"
	ItemBlocked     = "blocked"
	ItemSkipped     = "skipped"
)

// Project is a space-scoped goal with optional planning and execution linkage.
type Project struct {
	ID             string                 `json:"id"`
	SpaceID        string                 `json:"space_id"`
	Title          string                 `json:"title"`
	Goal           string                 `json:"goal,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Status         string                 `json:"status"`
	PlanJSON       map[string]interface{} `json:"plan_json,omitempty"`
	ParentTaskID   string                 `json:"parent_task_id,omitempty"`
	CreatorID      string                 `json:"creator_id,omitempty"`
	PlannerAgentID string                 `json:"planner_agent_id,omitempty"`
	TeamAgentIDs   []string               `json:"team_agent_ids,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	StartedAt      string                 `json:"started_at,omitempty"`
	CompletedAt    string                 `json:"completed_at,omitempty"`
	CreatedAt      string                 `json:"created_at,omitempty"`
	UpdatedAt      string                 `json:"updated_at,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// ProjectItem is a kanban card belonging to a project.
type ProjectItem struct {
	ID              string   `json:"id"`
	SpaceID         string   `json:"space_id"`
	ProjectID       string   `json:"project_id"`
	Title           string   `json:"title"`
	Description     string   `json:"description,omitempty"`
	Status          string   `json:"status"`
	Priority        string   `json:"priority,omitempty"`
	AssigneeIDs     []string `json:"assignee_ids,omitempty"`
	ExecutionTaskID string   `json:"execution_task_id,omitempty"`
	SortOrder       int      `json:"sort_order,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	CreatedAt       string   `json:"created_at,omitempty"`
	UpdatedAt       string   `json:"updated_at,omitempty"`
}

// ProjectProgress summarizes kanban item counts for a project.
type ProjectProgress struct {
	Total         int `json:"total"`
	Completed     int `json:"completed"`
	Skipped       int `json:"skipped"`
	InProgress    int `json:"in_progress"`
	Blocked       int `json:"blocked"`
	HumanPending  int `json:"human_pending"`
	Percent       int `json:"percent"`
}