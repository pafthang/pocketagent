package models

const (
	WorkflowDefault    = "default"
	WorkflowSupervisor = "supervisor"
)

// Schedule is a cron-driven recurring task definition.
type Schedule struct {
	ID              string   `json:"id"`
	SpaceID         string   `json:"space_id"`
	Name            string   `json:"name"`
	AgentID         string   `json:"agent_id,omitempty"`
	Prompt          string   `json:"prompt"`
	CronExpr        string   `json:"cron_expr"`
	Workflow        string   `json:"workflow,omitempty"`
	WorkerAgentIDs  []string `json:"worker_agent_ids,omitempty"`
	Enabled         bool     `json:"enabled"`
	LastRunAt       string   `json:"last_run_at,omitempty"`
	NextRunAt       string   `json:"next_run_at,omitempty"`
	LastTaskID      string   `json:"last_task_id,omitempty"`
	CreatedAt       string   `json:"created_at,omitempty"`
	UpdatedAt       string   `json:"updated_at,omitempty"`
}