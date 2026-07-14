package bootstrap

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pocketbase/pocketbase/core"
)

func ensureSchedulesCollection(app core.App) error {
	return ensureLockedCollection(app, pbclient.SchedulesCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "name", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "agent_id"})
		addFieldIfMissing(col, &core.TextField{Name: "prompt", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "cron_expr", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "workflow"})
		addFieldIfMissing(col, &core.JSONField{Name: "worker_agent_ids"})
		addFieldIfMissing(col, &core.BoolField{Name: "enabled"})
		addFieldIfMissing(col, &core.DateField{Name: "last_run_at"})
		addFieldIfMissing(col, &core.DateField{Name: "next_run_at"})
		addFieldIfMissing(col, &core.TextField{Name: "last_task_id"})
	})
}