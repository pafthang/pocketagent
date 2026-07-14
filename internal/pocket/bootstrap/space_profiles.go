package bootstrap

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pocketbase/pocketbase/core"
)

func ensureSpaceProfilesCollection(app core.App) error {
	return ensureLockedCollection(app, pbclient.SpaceProfilesCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "user_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "content"})
	})
}