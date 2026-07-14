package gateapis

import (
	"github.com/labstack/echo/v4"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	dashboardapis "github.com/pafthang/pocketagent/internal/space/gate/dash"
	mcpapis "github.com/pafthang/pocketagent/internal/space/gate/mcp"
	projectapis "github.com/pafthang/pocketagent/internal/space/gate/projects"
	skillapis "github.com/pafthang/pocketagent/internal/space/gate/skills"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

// Deps holds shared dependencies for gate-mounted space domain routes.
type Deps struct {
	PB         *pbclient.Client
	NATS       *natsclient.Client
	OllamaURL  string
	EmbedModel string
	LLMModel   string
}

// Actions carries RBAC middleware for tenant routes.
type Actions struct {
	MCPRead, MCPWrite         echo.MiddlewareFunc
	SkillRead, SkillWrite     echo.MiddlewareFunc
	TaskRead                  echo.MiddlewareFunc
	ProjectRead, ProjectWrite echo.MiddlewareFunc
}

// RegisterRoutes wires projects, skills, MCP, and dashboard endpoints on gate.
func RegisterRoutes(tenant *echo.Group, deps Deps, actions Actions) {
	mcpapis.RegisterRoutes(tenant, deps.PB, actions.MCPRead, actions.MCPWrite)
	skillapis.RegisterRoutes(tenant, deps.PB, deps.NATS, actions.SkillRead, actions.SkillWrite)
	dashboardapis.RegisterRoutes(tenant, deps.PB, actions.TaskRead)
	projectapis.RegisterRoutes(tenant, projectapis.Deps{
		PB:       deps.PB,
		NC:       deps.NATS,
		Ollama:   ollama.NewConfigured(deps.OllamaURL, deps.EmbedModel),
		LLMModel: deps.LLMModel,
	}, actions.ProjectRead, actions.ProjectWrite)
}
