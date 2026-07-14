package files

import (
	"time"

	fileapis "github.com/pafthang/pocketagent/internal/files/apis"
	"github.com/pafthang/pocketagent/internal/space"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
	"github.com/pafthang/pocketagent/pkgs/service"
)

func registerRoutes(s *service.Server, deps *fileapis.Deps, cfg *Config) {
	rbac := apimw.NewPocketRBAC(deps.PB, time.Duration(cfg.AuthorizeCacheSecs)*time.Second)
	tenant := s.Echo.Group("", rbac.AuthMiddleware(), apimw.RequireSpace(apimw.SpaceOptions{AllowQueryFallback: true}))
	fileRead := rbac.RequireAction(space.ActionFileRead)
	fileWrite := rbac.RequireAction(space.ActionFileWrite)
	fileapis.RegisterRoutes(tenant, deps, fileRead, fileWrite)
}
