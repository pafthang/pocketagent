package dashboardapis

import (
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestKanbanColumn(t *testing.T) {
	if kanbanColumn(models.TaskRunning) != "running" {
		t.Fatal("expected running column")
	}
	if kanbanColumn(models.TaskDegraded) != "completed" {
		t.Fatal("expected completed column for degraded")
	}
}

func TestMapKanbanGroupsTasks(t *testing.T) {
	board := mapKanban([]models.Task{
		{CorrelationID: "t1", Prompt: "one", Status: models.TaskQueued},
		{CorrelationID: "t2", Prompt: "two", Status: models.TaskRunning},
	})
	if len(board["queued"]) != 1 || len(board["running"]) != 1 {
		t.Fatalf("unexpected board: %+v", board)
	}
}

func TestKitDataIncludesStats(t *testing.T) {
	data := KitData(models.DashboardSummary{
		Metrics: models.DashboardMetrics{AgentsTotal: 2, TasksRunning: 1},
	})
	stats, ok := data["api:stats"].(map[string]interface{})
	if !ok || stats["agents_total"] != 2 {
		t.Fatalf("unexpected stats: %+v", data["api:stats"])
	}
}
