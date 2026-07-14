package models

// Planning phase keys (Deep Work compatible).
const (
	PlanPhaseGoalAnalysis = "goal_analysis"
	PlanPhaseResearch     = "research"
	PlanPhasePRD          = "prd"
	PlanPhaseTasks        = "tasks"
	PlanPhaseTeam         = "team"
)

// ProjectPlanCommand enqueues async project planning.
type ProjectPlanCommand struct {
	ProjectID string `json:"project_id"`
	SpaceID   string `json:"space_id"`
}