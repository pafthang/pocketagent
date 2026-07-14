package agent

import "github.com/pafthang/pocketagent/pkgs/ollama/api"

// ToolsForAgent returns tool definitions allowed for an agent.
// Empty allowed list means all catalog tools.
func ToolsForAgent(allowed []string, catalog []api.Tool) []api.Tool {
	if len(allowed) == 0 {
		return catalog
	}

	allowedSet := make(map[string]struct{}, len(allowed))
	for _, name := range allowed {
		allowedSet[name] = struct{}{}
	}

	filtered := make([]api.Tool, 0, len(allowed))
	for _, tool := range catalog {
		if _, ok := allowedSet[tool.Function.Name]; ok {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}