package audit

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// Audit action names.
const (
	AuditSpaceCreate   = "space.create"
	AuditSpaceUpdate   = "space.update"
	AuditSpaceDelete   = "space.delete"
	AuditMemberAdd     = "member.add"
	AuditMemberUpdate  = "member.update"
	AuditMemberRemove  = "member.remove"
	AuditInviteCreate  = "invite.create"
	AuditInviteAccept  = "invite.accept"
	AuditInviteRevoke  = "invite.revoke"
	AuditAuthRegister  = "auth.register"
	AuditAuthLogin     = "auth.login"
	AuditEmailVerified = "auth.email_verified"
)

// Auditor persists space-scoped audit events.
type Auditor struct {
	pb  *pbclient.Client
	log *slog.Logger
}

func NewAuditor(pb *pbclient.Client, log *slog.Logger) *Auditor {
	return &Auditor{pb: pb, log: log}
}

func (a *Auditor) Record(c echo.Context, entry models.AuditLog) {
	if a == nil || a.pb == nil {
		return
	}
	if entry.IPAddress == "" && c != nil {
		entry.IPAddress = c.RealIP()
	}
	if _, err := a.pb.CreateAuditLog(entry); err != nil && a.log != nil {
		a.log.Warn("audit log failed", "action", entry.Action, "error", err)
	}
}

func (a *Auditor) RecordCtx(ctx context.Context, entry models.AuditLog) {
	_ = ctx
	if a == nil || a.pb == nil {
		return
	}
	if _, err := a.pb.CreateAuditLog(entry); err != nil && a.log != nil {
		a.log.Warn("audit log failed", "action", entry.Action, "error", err)
	}
}
