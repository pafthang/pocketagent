package bootstrap

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pocketbase/pocketbase/core"
)

func ensureTaskEventsCollection(app core.App) error {
	return ensureLockedCollection(app, pbclient.TaskEventsCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "task_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "event_type", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "status"})
		addFieldIfMissing(col, &core.NumberField{Name: "step"})
		addFieldIfMissing(col, &core.TextField{Name: "message"})
		addFieldIfMissing(col, &core.TextField{Name: "result"})
	})
}