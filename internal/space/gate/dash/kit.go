package dashboardapis

import "time"

// BuiltinKit returns the default command center kit definition.
func BuiltinKit() map[string]interface{} {
	now := time.Now().UTC().Format(time.RFC3339)
	return map[string]interface{}{
		"id": BuiltinKitID,
		"config": map[string]interface{}{
			"meta": map[string]interface{}{
				"name":        "Space Dashboard",
				"author":      "PocketAgent",
				"version":     "1.0.0",
				"description": "Built-in overview of agents, tasks, and activity",
				"category":    "operations",
				"tags":        []string{"builtin", "tasks", "agents"},
				"icon":        "layout-dashboard",
				"built_in":    true,
			},
			"layout": map[string]interface{}{
				"columns": 2,
				"sections": []map[string]interface{}{
					{
						"title": "Overview",
						"span":  "full",
						"panels": []map[string]interface{}{
							{
								"id":   "metrics",
								"type": "metrics-row",
								"items": []map[string]interface{}{
									{"label": "Agents", "source": "api:stats", "field": "agents_total", "format": "number"},
									{"label": "Running", "source": "api:stats", "field": "tasks_running", "format": "number"},
									{"label": "Queued", "source": "api:stats", "field": "tasks_queued", "format": "number"},
									{"label": "Completed", "source": "api:stats", "field": "tasks_completed", "format": "number"},
								},
							},
						},
					},
					{
						"title": "Agents",
						"span":  "left",
						"panels": []map[string]interface{}{
							{"id": "agents", "type": "agent-roster", "source": "agents"},
						},
					},
					{
						"title": "Tasks",
						"span":  "right",
						"panels": []map[string]interface{}{
							{
								"id":     "kanban",
								"type":   "kanban",
								"source": "kanban",
								"columns": []map[string]interface{}{
									{"key": "queued", "label": "Queued", "color": "gray"},
									{"key": "running", "label": "Running", "color": "blue"},
									{"key": "completed", "label": "Completed", "color": "green"},
									{"key": "failed", "label": "Failed", "color": "orange"},
								},
								"card_fields": []string{"priority"},
							},
							{
								"id":     "recent",
								"type":   "table",
								"source": "recent_tasks",
								"columns": []map[string]interface{}{
									{"key": "title", "label": "Task"},
									{"key": "status", "label": "Status"},
									{"key": "agent_id", "label": "Agent"},
									{"key": "updated_at", "label": "Updated"},
								},
							},
						},
					},
					{
						"title": "Activity",
						"span":  "full",
						"panels": []map[string]interface{}{
							{"id": "feed", "type": "feed", "source": "activity", "max_items": 20},
						},
					},
				},
			},
			"workflows": map[string]interface{}{},
		},
		"user_values":  map[string]interface{}{},
		"installed_at": now,
		"active":       true,
	}
}

// BuiltinCatalogEntry describes the built-in kit in the store.
func BuiltinCatalogEntry() map[string]interface{} {
	return map[string]interface{}{
		"id":          BuiltinKitID,
		"name":        "Space Dashboard",
		"description": "Agents, running tasks, kanban, and recent activity",
		"icon":        "layout-dashboard",
		"category":    "operations",
		"author":      "PocketAgent",
		"tags":        []string{"builtin"},
		"preview":     "Metrics, agent roster, task kanban, activity feed",
		"installed":   true,
	}
}