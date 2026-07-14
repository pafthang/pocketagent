package audit

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// SystemAuditLog builds an audit entry for system-wide auth events.
func SystemAuditLog(pb *pbclient.Client, action, actorID, actorEmail string, metadata map[string]interface{}) models.AuditLog {
	spaceID := ""
	if pb != nil {
		if admin, err := pb.GetSpaceBySlug(models.SystemSpaceSlug); err == nil {
			spaceID = admin.ID
		}
	}
	return models.AuditLog{
		SpaceID:    spaceID,
		ActorID:    actorID,
		ActorEmail: actorEmail,
		Action:     action,
		Metadata:   metadata,
	}
}
