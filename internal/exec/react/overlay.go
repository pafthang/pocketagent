package react

import (
	"strings"

	"github.com/pafthang/pocketagent/internal/exec/tools"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

// ToolExecutor runs a named tool with arguments.
type ToolExecutor func(toolName, args string) string

// ToolOverlay adds per-task tools on top of the base executor registry.
type ToolOverlay struct {
	Catalog []ollama.Tool
	Run     ToolExecutor
}

// ToolRunner adapts a tools registry to a ToolExecutor.
func ToolRunner(registry tools.Registry) ToolExecutor {
	return registry.Execute
}

func mergeToolExecutor(base, overlay ToolExecutor) ToolExecutor {
	if overlay == nil {
		return base
	}
	if base == nil {
		return overlay
	}
	return func(toolName, args string) string {
		if overlay != nil {
			if out := overlay(toolName, args); !strings.HasPrefix(out, "Unknown tool:") {
				return out
			}
		}
		return base(toolName, args)
	}
}

func filteredToolExecutor(base ToolExecutor, allowed []string) ToolExecutor {
	if base == nil || len(allowed) == 0 {
		return base
	}

	allowedSet := make(map[string]struct{}, len(allowed))
	for _, name := range allowed {
		allowedSet[name] = struct{}{}
	}

	return func(toolName, args string) string {
		if _, ok := allowedSet[toolName]; !ok {
			return "Tool not allowed for this agent: " + toolName
		}
		return base(toolName, args)
	}
}