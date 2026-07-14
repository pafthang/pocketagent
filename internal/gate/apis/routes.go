package gateapis

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/gate/proxy"
	memoapis "github.com/pafthang/pocketagent/internal/memo/apis"
	"github.com/pafthang/pocketagent/internal/space"
	spacegateapis "github.com/pafthang/pocketagent/internal/space/gate"
	taskapis "github.com/pafthang/pocketagent/internal/task/apis"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

// RegisterRoutes wires the gate HTTP facade.
func RegisterRoutes(e *echo.Echo, deps Deps) {
	registerAuthRoutes(e, deps)

	spaceAPI := e.Group("", deps.RBAC.AuthMiddleware())
	spaceAPI.GET("/spaces", proxy.Space(deps.Space, http.MethodGet, "/spaces", true))
	spaceAPI.POST("/spaces", proxy.Space(deps.Space, http.MethodPost, "/spaces", true))
	spaceAPI.GET("/spaces/*", proxy.Space(deps.Space, http.MethodGet, "/spaces", true))
	spaceAPI.POST("/spaces/*", proxy.Space(deps.Space, http.MethodPost, "/spaces", true))
	spaceAPI.PATCH("/spaces/*", proxy.Space(deps.Space, http.MethodPatch, "/spaces", true))
	spaceAPI.DELETE("/spaces/*", proxy.Space(deps.Space, http.MethodDelete, "/spaces", true))
	spaceAPI.POST("/authorize", proxy.Space(deps.Space, http.MethodPost, "/authorize", true))
	spaceAPI.POST("/auth/request-verification", proxy.FixedPath(deps.Space, http.MethodPost, "/auth/request-verification", true))

	tenant := e.Group("", deps.RBAC.AuthMiddleware(), apimw.RequireSpace(apimw.SpaceOptions{AllowQueryFallback: true}))

	agentRead := deps.RBAC.RequireAction(space.ActionAgentRead)
	agentWrite := deps.RBAC.RequireAction(space.ActionAgentWrite)
	tenant.GET("/agents", proxy.Agent(deps.Agent, http.MethodGet, "/agents", true), agentRead)
	tenant.POST("/agents", proxy.Agent(deps.Agent, http.MethodPost, "/agents", true), agentWrite)
	tenant.GET("/agents/*", proxy.Agent(deps.Agent, http.MethodGet, "/agents", true), agentRead)
	tenant.POST("/agents/*", proxy.Agent(deps.Agent, http.MethodPost, "/agents", true), agentWrite)
	tenant.PUT("/agents/*", proxy.Agent(deps.Agent, http.MethodPut, "/agents", true), agentWrite)
	tenant.PATCH("/agents/*", proxy.Agent(deps.Agent, http.MethodPatch, "/agents", true), agentWrite)
	tenant.DELETE("/agents/*", proxy.Agent(deps.Agent, http.MethodDelete, "/agents", true), agentWrite)

	taskRead := deps.RBAC.RequireAction(space.ActionTaskRead)
	taskWrite := deps.RBAC.RequireAction(space.ActionTaskWrite)
	taskapis.RegisterRoutes(tenant, deps.NATS, deps.PB, taskRead, taskWrite)

	memoryRead := deps.RBAC.RequireAction(space.ActionMemoryRead)
	memoryWrite := deps.RBAC.RequireAction(space.ActionMemoryWrite)
	memoapis.RegisterRoutes(tenant, deps.Memo, deps.EmbedModel, deps.OllamaURL, memoryRead, memoryWrite)

	spacegateapis.RegisterRoutes(tenant, spacegateapis.Deps{
		PB:         deps.PB,
		NATS:       deps.NATS,
		OllamaURL:  deps.OllamaURL,
		EmbedModel: deps.EmbedModel,
		LLMModel:   deps.LLMModel,
	}, spacegateapis.Actions{
		MCPRead:      deps.RBAC.RequireAction(space.ActionMCPRead),
		MCPWrite:     deps.RBAC.RequireAction(space.ActionMCPWrite),
		SkillRead:    deps.RBAC.RequireAction(space.ActionSkillRead),
		SkillWrite:   deps.RBAC.RequireAction(space.ActionSkillWrite),
		TaskRead:     taskRead,
		ProjectRead:  deps.RBAC.RequireAction(space.ActionProjectRead),
		ProjectWrite: deps.RBAC.RequireAction(space.ActionProjectWrite),
	})

	fileRead := deps.RBAC.RequireAction(space.ActionFileRead)
	fileWrite := deps.RBAC.RequireAction(space.ActionFileWrite)
	tenant.GET("/files/*", proxy.Files(deps.Files, http.MethodGet, "/files", true), fileRead)
	tenant.POST("/files/*", proxy.Files(deps.Files, http.MethodPost, "/files", true), fileWrite)
	tenant.DELETE("/files/*", proxy.Files(deps.Files, http.MethodDelete, "/files", true), fileWrite)
	tenant.GET("/projects/:id/files/*", proxy.ProjectFiles(deps.Files, http.MethodGet, true), fileRead)
	tenant.POST("/projects/:id/files/*", proxy.ProjectFiles(deps.Files, http.MethodPost, true), fileWrite)
}