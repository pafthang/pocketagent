package bootstrap

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pocketbase/pocketbase/core"
)

func ensureTasksCollection(app core.App) error {
	return ensureLockedCollection(app, pbclient.TasksCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "correlation_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "user_id"})
		addFieldIfMissing(col, &core.TextField{Name: "agent_id"})
		addFieldIfMissing(col, &core.TextField{Name: "prompt", Required: true})
		addFieldIfMissing(col, &core.SelectField{
			Name:     "status",
			Required: true,
			Values: []string{
				string(models.TaskQueued),
				string(models.TaskRunning),
				string(models.TaskCompleted),
				string(models.TaskDegraded),
				string(models.TaskFailed),
				string(models.TaskCancelled),
			},
		})
		addFieldIfMissing(col, &core.TextField{Name: "result"})
		addFieldIfMissing(col, &core.TextField{Name: "error"})
		addFieldIfMissing(col, &core.TextField{Name: "parent_id"})
		addFieldIfMissing(col, &core.TextField{Name: "workflow"})
		addFieldIfMissing(col, &core.JSONField{Name: "worker_agent_ids"})
		addFieldIfMissing(col, &core.JSONField{Name: "tools"})
		addFieldIfMissing(col, &core.TextField{Name: "skill_id"})
	})
}
