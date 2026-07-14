package bootstrap

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pocketbase/pocketbase/core"
)

func ensureMCPServersCollection(app core.App) error {
	return ensureLockedCollection(app, pbclient.MCPServersCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "name", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "transport", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "command"})
		addFieldIfMissing(col, &core.JSONField{Name: "args"})
		addFieldIfMissing(col, &core.TextField{Name: "url"})
		addFieldIfMissing(col, &core.JSONField{Name: "env"})
		addFieldIfMissing(col, &core.BoolField{Name: "enabled"})
	})
}