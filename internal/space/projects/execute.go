package projects

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/common"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// ExecutionDeps wires project item execution.
type ExecutionDeps struct {
	PB *pbclient.Client
	NC *natsclient.Client
}

// SpawnItemExecutions enqueues a task per actionable kanban item.
func SpawnItemExecutions(ctx context.Context, c echo.Context, deps ExecutionDeps, project models.Project) ([]string, error) {
	items, _, err := deps.PB.ListProjectItems(pbclient.ListOptions{
		Page:    1,
		PerPage: 500,
		Filter:  pbclient.ProjectItemsFilter(project.SpaceID, project.ID),
	})
	if err != nil {
		return nil, err
	}

	corrIDs := make([]string, 0)
	for _, item := range items {
		if !itemActionable(item.Status) {
			continue
		}
		corrID, err := spawnItemExecution(ctx, c, deps, project, item)
		if err != nil {
			return corrIDs, err
		}
		corrIDs = append(corrIDs, corrID)
	}
	return corrIDs, nil
}

// SpawnSingleItemExecution re-queues one project item (retry).
func SpawnSingleItemExecution(ctx context.Context, c echo.Context, deps ExecutionDeps, project models.Project, item models.ProjectItem) (string, error) {
	return spawnItemExecution(ctx, c, deps, project, item)
}

func spawnItemExecution(ctx context.Context, c echo.Context, deps ExecutionDeps, project models.Project, item models.ProjectItem) (string, error) {
	if deps.NC == nil {
		return "", fmt.Errorf("nats client not configured")
	}

	agentID := project.PlannerAgentID
	if len(item.AssigneeIDs) > 0 {
		agentID = strings.TrimSpace(item.AssigneeIDs[0])
	}

	task := models.Task{
		SpaceID: project.SpaceID,
		AgentID: agentID,
		Prompt:  itemExecutionPrompt(project, item),
		Status:  models.TaskQueued,
	}
	if user, ok := apimw.UserFromContext(c); ok {
		task.UserID = user.ID
	}

	corrID := fmt.Sprintf("pitem-%s-%d", item.ID, time.Now().UnixNano())
	task.CorrelationID = corrID

	if _, err := deps.PB.CreateTask(task); err != nil {
		return "", err
	}
	runCtx := common.WithCorrelationID(ctx, corrID)
	if err := deps.NC.PublishOrchestrator(runCtx, task); err != nil {
		return "", err
	}

	item.ExecutionTaskID = corrID
	item.Status = models.ItemInProgress
	if _, err := deps.PB.UpdateProjectItem(item.ID, item); err != nil {
		return corrID, err
	}
	return corrID, nil
}

func itemExecutionPrompt(project models.Project, item models.ProjectItem) string {
	var b strings.Builder
	b.WriteString("Complete this project kanban item.\n\n")
	if project.Goal != "" {
		b.WriteString("Project goal: ")
		b.WriteString(project.Goal)
		b.WriteString("\n\n")
	}
	b.WriteString("Item: ")
	b.WriteString(item.Title)
	b.WriteString("\n")
	if strings.TrimSpace(item.Description) != "" {
		b.WriteString("\nDetails:\n")
		b.WriteString(item.Description)
		b.WriteString("\n")
	}
	return b.String()
}

func itemActionable(status string) bool {
	switch status {
	case models.ItemInbox, models.ItemAssigned, models.ItemBlocked, models.ItemReview:
		return true
	default:
		return false
	}
}
