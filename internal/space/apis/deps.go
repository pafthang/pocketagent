package spaceapis

import (
	"log/slog"

	"github.com/pafthang/pocketagent/internal/space/audit"
	"github.com/pafthang/pocketagent/internal/space/rbac"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
)

// HandlerConfig holds handler-facing settings (subset of service config).
type HandlerConfig struct {
	PublicBaseURL            string
	RequireEmailVerification bool
	InviteTTLHours           int
	VerificationTTLHours     int
}

// Deps wires space HTTP handlers.
type Deps struct {
	PB    *pbclient.Client
	Auth  *rbac.Authorizer
	Audit *audit.Auditor
	Cfg   HandlerConfig
	Log   *slog.Logger
}