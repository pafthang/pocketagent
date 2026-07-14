package projectapis

import (
	"strings"

	"github.com/pafthang/pocketagent/pkgs/models"
)

// ToMCProject maps a project to the front MCProject shape.
func ToMCProject(p models.Project, itemIDs []string) map[string]interface{} {
	out := map[string]interface{}{
		"id":               p.ID,
		"title":            p.Title,
		"description":      p.Description,
		"status":           p.Status,
		"planner_agent_id": nilIfEmptyString(p.PlannerAgentID),
		"team_agent_ids":   p.TeamAgentIDs,
		"task_ids":         itemIDs,
		"prd_document_id":  nil,
		"creator_id":       nilIfEmptyString(p.CreatorID),
		"tags":             p.Tags,
		"started_at":       nilIfEmptyString(p.StartedAt),
		"completed_at":     nilIfEmptyString(p.CompletedAt),
		"created_at":       p.CreatedAt,
		"updated_at":       p.UpdatedAt,
		"metadata":         p.Metadata,
	}
	if p.Goal != "" {
		if out["metadata"] == nil {
			out["metadata"] = map[string]interface{}{}
		}
		if meta, ok := out["metadata"].(map[string]interface{}); ok {
			meta["goal"] = p.Goal
		}
	}
	if p.ParentTaskID != "" {
		if out["metadata"] == nil {
			out["metadata"] = map[string]interface{}{}
		}
		if meta, ok := out["metadata"].(map[string]interface{}); ok {
			meta["parent_task_id"] = p.ParentTaskID
		}
	}
	return out
}

// ToMCTask maps a project kanban item to the front MCTask shape.
func ToMCTask(item models.ProjectItem, projectID, creatorID string) map[string]interface{} {
	return map[string]interface{}{
		"id":                 item.ID,
		"title":              item.Title,
		"description":        item.Description,
		"status":             item.Status,
		"priority":           defaultPriority(item.Priority),
		"assignee_ids":       item.AssigneeIDs,
		"creator_id":         nilIfEmptyString(creatorID),
		"parent_task_id":     nil,
		"blocked_by":         []string{},
		"tags":               item.Tags,
		"due_date":           nil,
		"started_at":         nil,
		"completed_at":       nil,
		"created_at":         item.CreatedAt,
		"updated_at":         item.UpdatedAt,
		"project_id":         projectID,
		"task_type":          "agent",
		"blocks":             []string{},
		"active_description": "",
		"estimated_minutes":  nil,
		"output":             nil,
		"retry_count":        0,
		"max_retries":        0,
		"timeout_minutes":    nil,
		"error_message":      nil,
		"metadata": map[string]interface{}{
			"execution_task_id": nilIfEmptyString(item.ExecutionTaskID),
		},
	}
}

// Progress computes kanban progress for project items.
func Progress(items []models.ProjectItem) models.ProjectProgress {
	var progress models.ProjectProgress
	progress.Total = len(items)
	for _, item := range items {
		switch item.Status {
		case models.ItemDone:
			progress.Completed++
		case models.ItemSkipped:
			progress.Skipped++
		case models.ItemInProgress, models.ItemReview:
			progress.InProgress++
		case models.ItemBlocked:
			progress.Blocked++
		case models.ItemAssigned, models.ItemInbox:
			progress.HumanPending++
		}
	}
	if progress.Total > 0 {
		progress.Percent = (progress.Completed * 100) / progress.Total
	}
	return progress
}

// NormalizeTitle picks a display title from request fields.
func NormalizeTitle(title, goal string) string {
	title = strings.TrimSpace(title)
	if title != "" {
		return title
	}
	goal = strings.TrimSpace(goal)
	if goal == "" {
		return "Untitled project"
	}
	if len(goal) > 120 {
		return goal[:120] + "..."
	}
	return goal
}

func defaultPriority(priority string) string {
	priority = strings.TrimSpace(priority)
	if priority == "" {
		return "medium"
	}
	return priority
}

func nilIfEmptyString(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
