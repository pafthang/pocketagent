package bootstrap

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pocketbase/pocketbase/core"
)

func ensureInviteAuditCollections(app core.App) error {
	if err := ensureCollectionSpaceInvites(app); err != nil {
		return err
	}
	if err := ensureCollectionAuditLogs(app); err != nil {
		return err
	}
	return ensureCollectionEmailVerifications(app)
}

func ensureCollectionSpaceInvites(app core.App) error {
	return ensureLockedCollection(app, pbclient.SpaceInvitesCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "email", Required: true})
		addFieldIfMissing(col, &core.SelectField{
			Name:     "role",
			Required: true,
			Values:   []string{models.RoleAdmin, models.RoleEditor, models.RoleViewer},
		})
		addFieldIfMissing(col, &core.TextField{Name: "token_hash", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "invited_by"})
		addFieldIfMissing(col, &core.SelectField{
			Name:     "status",
			Required: true,
			Values: []string{
				models.InvitePending,
				models.InviteAccepted,
				models.InviteRevoked,
				models.InviteExpired,
			},
		})
		addFieldIfMissing(col, &core.DateField{Name: "expires_at", Required: true})
	})
}

func ensureCollectionAuditLogs(app core.App) error {
	return ensureLockedCollection(app, pbclient.AuditLogsCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "actor_id"})
		addFieldIfMissing(col, &core.TextField{Name: "actor_email"})
		addFieldIfMissing(col, &core.TextField{Name: "action", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "resource_type"})
		addFieldIfMissing(col, &core.TextField{Name: "resource_id"})
		addFieldIfMissing(col, &core.JSONField{Name: "metadata"})
		addFieldIfMissing(col, &core.TextField{Name: "ip_address"})
	})
}

func ensureCollectionEmailVerifications(app core.App) error {
	return ensureLockedCollection(app, pbclient.EmailVerificationsCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "user_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "email", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "token_hash", Required: true})
		addFieldIfMissing(col, &core.SelectField{
			Name:     "status",
			Required: true,
			Values:   []string{models.VerificationPending, models.VerificationDone},
		})
		addFieldIfMissing(col, &core.DateField{Name: "expires_at", Required: true})
	})
}
