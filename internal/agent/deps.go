package agent

import (
	"time"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

// RouteDeps holds HTTP routing dependencies for the agent service.
type RouteDeps struct {
	PB   *pbclient.Client
	RBAC *apimw.PocketRBAC
}

func buildDeps(cfg *Config) (*RouteDeps, error) {
	pb, err := pbclient.NewServiceClient(cfg.PocketBaseURL, cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	if err != nil {
		return nil, err
	}
	return &RouteDeps{
		PB:   pb,
		RBAC: apimw.NewPocketRBAC(pb, time.Duration(cfg.AuthorizeCacheSecs)*time.Second),
	}, nil
}