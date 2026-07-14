package agent

import (
	agentapis "github.com/pafthang/pocketagent/internal/agent/apis"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
	"github.com/pafthang/pocketagent/pkgs/service"
)

func registerRoutes(s *service.Server, deps *RouteDeps) {
	tenant := s.Echo.Group("", deps.RBAC.AuthMiddleware(), apimw.RequireSpace(apimw.SpaceOptions{}))
	agentapis.RegisterRoutes(tenant, deps.PB, deps.RBAC)
}