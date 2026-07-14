package react

import "github.com/pafthang/pocketagent/pkgs/models"

// EffectiveAllowedTools returns the tool allow-list for a task run.
// Task-level tools override the agent list when present.
func EffectiveAllowedTools(task models.Task, agent models.Agent) []string {
	if len(task.Tools) > 0 {
		return append([]string{}, task.Tools...)
	}
	if len(agent.Tools) > 0 {
		return append([]string{}, agent.Tools...)
	}
	return nil
}