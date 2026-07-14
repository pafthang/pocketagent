package ollama

import (
	"github.com/pafthang/pocketagent/pkgs/ollama/agent"
	"github.com/pafthang/pocketagent/pkgs/ollama/api"
)

type Tool = api.Tool
type ToolFunction = api.ToolFunction

// GetExampleTools returns the default builtin tool catalog for tests and fallbacks.
func GetExampleTools() []Tool { return api.ExampleTools() }

// ToolsForAgent returns tool definitions allowed for an agent.
func ToolsForAgent(allowed []string, catalog []Tool) []Tool {
	return agent.ToolsForAgent(allowed, catalog)
}