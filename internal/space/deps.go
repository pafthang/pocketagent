package space

import (
	spaceapis "github.com/pafthang/pocketagent/internal/space/apis"
	"github.com/pafthang/pocketagent/internal/space/audit"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space/rbac"
	"github.com/pafthang/pocketagent/pkgs/service"
)

// RouteDeps holds HTTP routing dependencies for the space service.
type RouteDeps struct {
	PB   *pbclient.Client
	Auth *rbac.Authorizer
}

func buildDeps(cfg *Config) (*RouteDeps, error) {
	pb, err := pbclient.NewServiceClient(cfg.PocketBaseURL, cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	if err != nil {
		return nil, err
	}
	return &RouteDeps{
		PB:   pb,
		Auth: rbac.NewAuthorizer(pb),
	}, nil
}

func buildAPIDeps(s *service.Server, deps *RouteDeps, cfg *Config) spaceapis.Deps {
	return spaceapis.Deps{
		PB:    deps.PB,
		Auth:  deps.Auth,
		Audit: audit.NewAuditor(deps.PB, s.Log),
		Cfg: spaceapis.HandlerConfig{
			PublicBaseURL:            cfg.PublicBaseURL,
			RequireEmailVerification: cfg.RequireEmailVerification,
			InviteTTLHours:           cfg.InviteTTLHours,
			VerificationTTLHours:     cfg.VerificationTTLHours,
		},
		Log: s.Log,
	}
}