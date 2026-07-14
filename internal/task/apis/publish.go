package taskapis

import (
	"context"
	"net/http"

	apimw "github.com/pafthang/pocketagent/pkgs/middle"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/exec/tools"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space/activity"
	taskcore "github.com/pafthang/pocketagent/internal/task"
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// PublishTaskWithTools validates and enqueues a task for orchestration.
func PublishTaskWithTools(c echo.Context, nc *natsclient.Client, pb *pbclient.Client, task models.Task, toolCfg tools.Config) error {
	if task.Prompt == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "prompt is required"})
	}
	if err := common.GuardPrompt(nil, common.LoadPromptGuardConfig(), task.Prompt); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	spaceID, ok := apimw.SpaceIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
	}
	task.SpaceID = spaceID
	if user, ok := apimw.UserFromContext(c); ok && user.ID != "" {
		task.UserID = user.ID
	}

	if len(task.Tools) > 0 {
		if err := tools.ValidateToolAllowList(pb, spaceID, toolCfg, task.Tools); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}

	if task.AgentID != "" {
		agentRecord, err := pb.GetAgent(task.AgentID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid agent_id"})
		}
		if agentRecord.SpaceID != "" && agentRecord.SpaceID != spaceID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "agent does not belong to this space"})
		}
	}
	if task.Workflow == models.WorkflowSupervisor && len(task.WorkerAgentIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "worker_agent_ids required for supervisor workflow"})
	}
	for _, workerID := range task.WorkerAgentIDs {
		worker, err := pb.GetAgent(workerID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid worker_agent_id"})
		}
		if worker.SpaceID != "" && worker.SpaceID != spaceID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "worker agent does not belong to this space"})
		}
	}

	ctx := context.Background()
	stored, _, err := taskcore.PersistAndEnqueue(ctx, pb, nc, task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	queued := models.NewTaskEvent(stored.CorrelationID, models.EventQueued, "queued", "task accepted by orchestrator")
	queued.SpaceID = spaceID
	activity.Record(pb, queued)

	return c.JSON(http.StatusCreated, stored)
}
