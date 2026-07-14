package gate

import (
	gateapis "github.com/pafthang/pocketagent/internal/gate/apis"
	"github.com/pafthang/pocketagent/pkgs/service"
)

func registerRoutes(s *service.Server, deps *gateapis.Deps) {
	gateapis.RegisterRoutes(s.Echo, *deps)
}