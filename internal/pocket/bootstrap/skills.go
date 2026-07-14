package bootstrap

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pocketbase/pocketbase/core"
)

func ensureSkillsCollection(app core.App) error {
	return ensureLockedCollection(app, pbclient.SkillsCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "name", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "description"})
		addFieldIfMissing(col, &core.TextField{Name: "prompt", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "category"})
		addFieldIfMissing(col, &core.JSONField{Name: "tools"})
		addFieldIfMissing(col, &core.TextField{Name: "argument_hint"})
		addFieldIfMissing(col, &core.TextField{Name: "catalog_id"})
	})
}