package decompose

import (
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestFallbackSupervisorPlan(t *testing.T) {
	workers := []string{"agent-a", "agent-b"}
	plans := fallbackSupervisorPlan("do A and B", workers, 4)
	if len(plans) < 2 {
		t.Fatalf("expected multiple plans, got %d", len(plans))
	}
	if plans[0].AgentID == "" {
		t.Fatal("expected agent assignment")
	}
}

func TestIsSupervisorWorkflow(t *testing.T) {
	if !isSupervisorWorkflow(models.Task{Workflow: models.WorkflowSupervisor}) {
		t.Fatal("expected supervisor workflow")
	}
	if !isSupervisorWorkflow(models.Task{WorkerAgentIDs: []string{"w1"}}) {
		t.Fatal("worker list should enable supervisor mode")
	}
}
