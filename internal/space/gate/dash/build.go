package dashboardapis

import (
	"fmt"
	"time"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space/activity"
	"github.com/pafthang/pocketagent/pkgs/models"
)

const BuiltinKitID = "builtin-dashboard"

// BuildOptions tunes list sizes for dashboard aggregation.
type BuildOptions struct {
	RecentLimit   int
	ActivityLimit int
}

// BuildSummary aggregates agents, tasks, and recent activity for a space.
func BuildSummary(pb *pbclient.Client, spaceID string, opts BuildOptions) (models.DashboardSummary, error) {
	if opts.RecentLimit <= 0 {
		opts.RecentLimit = 25
	}
	if opts.ActivityLimit <= 0 {
		opts.ActivityLimit = 20
	}

	agents, _, err := pb.ListAgentsFilter(pbclient.ListOptions{
		Page: 1, PerPage: 200, Filter: fmt.Sprintf("space_id = %q", spaceID),
	})
	if err != nil {
		return models.DashboardSummary{}, err
	}

	metrics, err := countMetrics(pb, spaceID)
	if err != nil {
		return models.DashboardSummary{}, err
	}
	metrics.AgentsTotal = len(agents)

	runningTasks, _, err := pb.ListTasks(pbclient.ListOptions{
		Page:    1,
		PerPage: opts.RecentLimit,
		Filter:  fmt.Sprintf(`space_id = %q && status = %q`, spaceID, models.TaskRunning),
	})
	if err != nil {
		return models.DashboardSummary{}, err
	}

	recentRecords, _, err := pb.ListTasks(pbclient.ListOptions{
		Page:    1,
		PerPage: opts.RecentLimit,
		Filter:  fmt.Sprintf("space_id = %q", spaceID),
	})
	if err != nil {
		return models.DashboardSummary{}, err
	}

	runningAgents := runningAgentSet(runningTasks)
	roster := make([]models.DashboardAgent, 0, len(agents))
	for _, agent := range agents {
		roster = append(roster, mapAgent(agent, runningAgents[agent.ID]))
	}

	kanbanTasks, _, err := pb.ListTasks(pbclient.ListOptions{
		Page:    1,
		PerPage: 100,
		Filter: fmt.Sprintf(
			`space_id = %q && (status = %q || status = %q || status = %q || status = %q || status = %q)`,
			spaceID,
			models.TaskQueued,
			models.TaskRunning,
			models.TaskCompleted,
			models.TaskFailed,
			models.TaskCancelled,
		),
	})
	if err != nil {
		return models.DashboardSummary{}, err
	}

	taskEvents, _, err := pb.ListTaskEvents(spaceID, pbclient.ListOptions{Page: 1, PerPage: opts.ActivityLimit})
	if err != nil {
		return models.DashboardSummary{}, err
	}

	return models.DashboardSummary{
		Metrics:      metrics,
		Agents:       roster,
		RunningTasks: runningTasks,
		RecentTasks:  mapTaskRows(recentRecords),
		Kanban:       mapKanban(kanbanTasks),
		Activity:     mapActivity(taskEvents),
		GeneratedAt:  time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// KitData maps a summary into PawKit-compatible panel data keys.
func KitData(summary models.DashboardSummary) map[string]interface{} {
	stats := map[string]interface{}{
		"agents_total":    summary.Metrics.AgentsTotal,
		"tasks_total":     summary.Metrics.TasksTotal,
		"tasks_queued":    summary.Metrics.TasksQueued,
		"tasks_running":   summary.Metrics.TasksRunning,
		"tasks_completed": summary.Metrics.TasksCompleted,
		"tasks_failed":    summary.Metrics.TasksFailed,
		"tasks_cancelled": summary.Metrics.TasksCancelled,
	}

	agents := make([]map[string]interface{}, 0, len(summary.Agents))
	for _, agent := range summary.Agents {
		agents = append(agents, map[string]interface{}{
			"id":              agent.ID,
			"name":            agent.Name,
			"role":            agent.Role,
			"description":     agent.Description,
			"session_key":     "",
			"backend":         "pocketagent",
			"status":          agent.Status,
			"level":           agent.Level,
			"current_task_id": nilIfEmpty(agent.CurrentTask),
			"specialties":     agent.Specialties,
			"last_heartbeat":  nil,
			"created_at":      agent.CreatedAt,
			"updated_at":      agent.UpdatedAt,
			"metadata":        map[string]interface{}{"model": agent.Model},
		})
	}

	activityItems := make([]map[string]interface{}, 0, len(summary.Activity))
	for _, item := range summary.Activity {
		activityItems = append(activityItems, map[string]interface{}{
			"message":    item.Message,
			"type":       item.Type,
			"created_at": item.CreatedAt,
			"agent_id":   item.AgentID,
		})
	}

	recent := make([]map[string]interface{}, 0, len(summary.RecentTasks))
	for _, row := range summary.RecentTasks {
		recent = append(recent, map[string]interface{}{
			"id":         row.ID,
			"title":      row.Title,
			"status":     row.Status,
			"agent_id":   row.AgentID,
			"created_at": row.CreatedAt,
			"updated_at": row.UpdatedAt,
		})
	}

	kanban := make(map[string]interface{}, len(summary.Kanban))
	for key, cards := range summary.Kanban {
		items := make([]map[string]interface{}, 0, len(cards))
		for _, card := range cards {
			items = append(items, map[string]interface{}{
				"id":         card.ID,
				"title":      card.Title,
				"status":     card.Status,
				"priority":   card.Priority,
				"agent_id":   card.AgentID,
				"created_at": card.CreatedAt,
			})
		}
		kanban[key] = items
	}

	return map[string]interface{}{
		"api:stats":    stats,
		"agents":       agents,
		"kanban":       kanban,
		"recent_tasks": recent,
		"activity":     activityItems,
	}
}

func countMetrics(pb *pbclient.Client, spaceID string) (models.DashboardMetrics, error) {
	var metrics models.DashboardMetrics
	statuses := []models.TaskStatus{
		models.TaskQueued,
		models.TaskRunning,
		models.TaskCompleted,
		models.TaskFailed,
		models.TaskCancelled,
	}
	for _, status := range statuses {
		_, total, err := pb.ListTasks(pbclient.ListOptions{
			Page:    1,
			PerPage: 1,
			Filter:  fmt.Sprintf(`space_id = %q && status = %q`, spaceID, status),
		})
		if err != nil {
			return metrics, err
		}
		metrics.TasksTotal += total
		switch status {
		case models.TaskQueued:
			metrics.TasksQueued = total
		case models.TaskRunning:
			metrics.TasksRunning = total
		case models.TaskCompleted:
			metrics.TasksCompleted = total
		case models.TaskFailed:
			metrics.TasksFailed = total
		case models.TaskCancelled:
			metrics.TasksCancelled = total
		}
	}
	return metrics, nil
}

func runningAgentSet(tasks []models.Task) map[string]string {
	out := make(map[string]string)
	for _, task := range tasks {
		if task.AgentID != "" {
			out[task.AgentID] = task.CorrelationID
		}
	}
	return out
}

func mapAgent(agent models.Agent, currentTask string) models.DashboardAgent {
	status := "idle"
	if currentTask != "" {
		status = "active"
	}
	role := agent.Description
	if role == "" {
		role = "Agent"
	}
	return models.DashboardAgent{
		ID:          agent.ID,
		Name:        agent.Name,
		Role:        role,
		Description: agent.Description,
		Model:       agent.Model,
		Status:      status,
		Level:       "specialist",
		Specialties: append([]string(nil), agent.Tools...),
		CurrentTask: currentTask,
		CreatedAt:   agent.CreatedAt,
		UpdatedAt:   agent.UpdatedAt,
	}
}

func mapTaskRows(tasks []models.Task) []models.DashboardTaskRow {
	out := make([]models.DashboardTaskRow, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, models.DashboardTaskRow{
			ID:        task.CorrelationID,
			Title:     truncate(task.Prompt, 120),
			Status:    string(task.Status),
			AgentID:   task.AgentID,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		})
	}
	return out
}

func mapKanban(tasks []models.Task) map[string][]models.DashboardKanbanCard {
	out := map[string][]models.DashboardKanbanCard{
		"queued":    {},
		"running":   {},
		"completed": {},
		"failed":    {},
	}
	for _, task := range tasks {
		key := kanbanColumn(task.Status)
		out[key] = append(out[key], models.DashboardKanbanCard{
			ID:        task.CorrelationID,
			Title:     truncate(task.Prompt, 100),
			Status:    string(task.Status),
			Priority:  "medium",
			AgentID:   task.AgentID,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		})
	}
	return out
}

func kanbanColumn(status models.TaskStatus) string {
	switch status {
	case models.TaskQueued:
		return "queued"
	case models.TaskRunning:
		return "running"
	case models.TaskCompleted, models.TaskDegraded:
		return "completed"
	default:
		return "failed"
	}
}

func mapActivity(events []models.StoredTaskEvent) []models.DashboardFeedItem {
	out := make([]models.DashboardFeedItem, 0, len(events))
	for _, event := range events {
		entry := activity.FromTaskEvent(event)
		out = append(out, models.DashboardFeedItem{
			Message:   entry.Content,
			Type:      entry.Type,
			CreatedAt: entry.Timestamp,
			TaskID:    event.TaskID,
		})
	}
	return out
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
