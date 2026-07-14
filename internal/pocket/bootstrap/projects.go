package bootstrap

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pocketbase/pocketbase/core"
)

func ensureProjectsCollection(app core.App) error {
	return ensureLockedCollection(app, pbclient.ProjectsCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "title", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "goal"})
		addFieldIfMissing(col, &core.TextField{Name: "description"})
		addFieldIfMissing(col, &core.SelectField{
			Name:     "status",
			Required: true,
			Values: []string{
				models.ProjectDraft,
				models.ProjectPlanning,
				models.ProjectAwaitingApproval,
				models.ProjectApproved,
				models.ProjectExecuting,
				models.ProjectPaused,
				models.ProjectCompleted,
				models.ProjectFailed,
				models.ProjectCancelled,
			},
		})
		addFieldIfMissing(col, &core.JSONField{Name: "plan_json"})
		addFieldIfMissing(col, &core.TextField{Name: "parent_task_id"})
		addFieldIfMissing(col, &core.TextField{Name: "creator_id"})
		addFieldIfMissing(col, &core.TextField{Name: "planner_agent_id"})
		addFieldIfMissing(col, &core.JSONField{Name: "team_agent_ids"})
		addFieldIfMissing(col, &core.JSONField{Name: "tags"})
		addFieldIfMissing(col, &core.TextField{Name: "started_at"})
		addFieldIfMissing(col, &core.TextField{Name: "completed_at"})
		addFieldIfMissing(col, &core.JSONField{Name: "metadata"})
	})
}

func ensureProjectItemsCollection(app core.App) error {
	return ensureLockedCollection(app, pbclient.ProjectItemsCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "project_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "title", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "description"})
		addFieldIfMissing(col, &core.SelectField{
			Name:     "status",
			Required: true,
			Values: []string{
				models.ItemInbox,
				models.ItemAssigned,
				models.ItemInProgress,
				models.ItemReview,
				models.ItemDone,
				models.ItemBlocked,
				models.ItemSkipped,
			},
		})
		addFieldIfMissing(col, &core.SelectField{
			Name:   "priority",
			Values: []string{"low", "medium", "high", "urgent"},
		})
		addFieldIfMissing(col, &core.JSONField{Name: "assignee_ids"})
		addFieldIfMissing(col, &core.TextField{Name: "execution_task_id"})
		addFieldIfMissing(col, &core.NumberField{Name: "sort_order"})
		addFieldIfMissing(col, &core.JSONField{Name: "tags"})
	})
}
