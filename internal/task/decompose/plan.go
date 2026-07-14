package decompose

// SubtaskPlan is a unit of parallel work with an assigned agent.
type SubtaskPlan struct {
	Prompt  string
	AgentID string
}