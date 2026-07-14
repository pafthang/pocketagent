package models

// DashboardMetrics holds space-level counters for the command center.
type DashboardMetrics struct {
	AgentsTotal    int `json:"agents_total"`
	TasksTotal     int `json:"tasks_total"`
	TasksQueued    int `json:"tasks_queued"`
	TasksRunning   int `json:"tasks_running"`
	TasksCompleted int `json:"tasks_completed"`
	TasksFailed    int `json:"tasks_failed"`
	TasksCancelled int `json:"tasks_cancelled"`
}

// DashboardAgent is an agent row for the command center roster.
type DashboardAgent struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Description string   `json:"description"`
	Model       string   `json:"model,omitempty"`
	Status      string   `json:"status"`
	Level       string   `json:"level"`
	Specialties []string `json:"specialties"`
	CurrentTask string   `json:"current_task_id,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

// DashboardKanbanCard is a task card for kanban panels.
type DashboardKanbanCard struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	Priority  string `json:"priority,omitempty"`
	AgentID   string `json:"agent_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// DashboardTaskRow is a compact task row for tables.
type DashboardTaskRow struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	AgentID   string `json:"agent_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// DashboardFeedItem is an activity feed row.
type DashboardFeedItem struct {
	Message   string `json:"message"`
	Type      string `json:"type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	AgentID   string `json:"agent_id,omitempty"`
	TaskID    string `json:"task_id,omitempty"`
}

// DashboardSummary is the aggregated command center payload.
type DashboardSummary struct {
	Metrics      DashboardMetrics                 `json:"metrics"`
	Agents       []DashboardAgent                 `json:"agents"`
	RunningTasks []Task                           `json:"running_tasks"`
	RecentTasks  []DashboardTaskRow               `json:"recent_tasks"`
	Kanban       map[string][]DashboardKanbanCard `json:"kanban"`
	Activity     []DashboardFeedItem              `json:"activity"`
	GeneratedAt  string                           `json:"generated_at"`
}