package memo

import (
	"github.com/pafthang/pocketagent/internal/memo/svc"
	"github.com/pafthang/pocketagent/pkgs/service"
)

func registerRoutes(s *service.Server, deps *ServiceDeps) {
	svc.RegisterRoutes(s, deps.Manager, deps.Config.ServiceToken)
}