package bootstrap

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pocketbase/pocketbase/core"
)

func ensureAgentsCollection(app core.App) error {
	return ensureLockedCollection(app, pbclient.AgentsCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id"})
		addFieldIfMissing(col, &core.TextField{Name: "name", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "description"})
		addFieldIfMissing(col, &core.TextField{Name: "model"})
		addFieldIfMissing(col, &core.TextField{Name: "system_prompt"})
		addFieldIfMissing(col, &core.JSONField{Name: "tools"})
		addFieldIfMissing(col, &core.JSONField{Name: "config"})
	})
}