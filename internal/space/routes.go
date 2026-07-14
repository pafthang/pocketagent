package space

import (
	spaceapis "github.com/pafthang/pocketagent/internal/space/apis"
	"github.com/pafthang/pocketagent/pkgs/service"
)

func registerRoutes(s *service.Server, deps *RouteDeps, cfg *Config) {
	spaceapis.RegisterRoutes(s.Echo, buildAPIDeps(s, deps, cfg), cfg.RateLimit, AuthMiddleware(deps.PB))
}