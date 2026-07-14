package projectapis

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space/projects"
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func registerPlanningRoutes(tenant *echo.Group, deps Deps, readAction, writeAction echo.MiddlewareFunc) {
	tenant.GET("/ws/project/:projectId", func(c echo.Context) error {
		return wsProjectStream(c, deps.NC, deps.PB)
	}, readAction)

	tenant.POST("/projects/parse-goal", parseGoalHandler(deps), writeAction)
	tenant.POST("/projects/start", startProjectHandler(deps), writeAction)
	tenant.POST("/projects/:id/plan", planProjectHandler(deps), writeAction)
	tenant.GET("/projects/:id/plan", getPlanHandler(deps.PB), readAction)
	tenant.POST("/projects/:id/approve", approveProjectHandler(deps), writeAction)
	tenant.POST("/projects/:id/pause", transitionProjectHandler(deps.PB, models.ProjectPaused), writeAction)
	tenant.POST("/projects/:id/resume", transitionProjectHandler(deps.PB, models.ProjectExecuting), writeAction)
	tenant.POST("/projects/:id/cancel", transitionProjectHandler(deps.PB, models.ProjectCancelled), writeAction)
	tenant.POST("/projects/:id/tasks/:taskId/skip", skipProjectTaskHandler(deps.PB), writeAction)
	tenant.POST("/projects/:id/tasks/:taskId/retry", retryProjectTaskHandler(deps), writeAction)
}

func parseGoalHandler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req ParseGoalRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		analysis, err := projects.ParseGoal(c.Request().Context(), deps.Ollama, deps.LLMModel, req.Description)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success":       true,
			"goal_analysis": analysis,
		})
	}
}

func startProjectHandler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		var req StartProjectRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		goal := strings.TrimSpace(req.Description)
		if goal == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "description is required"})
		}

		project := models.Project{
			SpaceID:        spaceID,
			Title:          projectplannerTitle(req.Title, goal),
			Goal:           goal,
			Status:         models.ProjectPlanning,
			PlannerAgentID: strings.TrimSpace(req.PlannerAgentID),
			TeamAgentIDs:   req.TeamAgentIDs,
			PlanJSON: map[string]interface{}{
				"current_phase": models.PlanPhaseGoalAnalysis,
				"phase_message": "Starting planning",
				"phases":        map[string]interface{}{},
			},
		}
		if user, ok := apimw.UserFromContext(c); ok {
			project.CreatorID = user.ID
		}
		if err := validatePlannerAgent(deps.PB, spaceID, project.PlannerAgentID); err != nil {
			return httpx.MapPocketError(c, err)
		}

		stored, err := deps.PB.CreateProject(project)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if err := enqueuePlanning(c.Request().Context(), deps, stored.ID, spaceID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"success":    true,
			"project_id": stored.ID,
			"project":    stored,
			"status":     stored.Status,
			"message":    "planning started",
		})
	}
}

func planProjectHandler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		project, err := loadProjectInSpace(deps.PB, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if project.Status == models.ProjectPlanning {
			return c.JSON(http.StatusConflict, map[string]string{"error": "planning already in progress"})
		}
		project.Status = models.ProjectPlanning
		if _, err := deps.PB.UpdateProject(project.ID, project); err != nil {
			return httpx.MapPocketError(c, err)
		}
		if err := enqueuePlanning(c.Request().Context(), deps, project.ID, spaceID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusAccepted, map[string]string{"status": "planning"})
	}
}

func getPlanHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		project, err := loadProjectInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		plan := project.PlanJSON
		if plan == nil {
			plan = map[string]interface{}{}
		}
		currentPhase, _ := plan["current_phase"].(string)
		phaseMessage, _ := plan["phase_message"].(string)
		planningDone, _ := plan["planning_done"].(bool)
		if project.Status == models.ProjectAwaitingApproval {
			planningDone = true
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"project_id":    project.ID,
			"status":        project.Status,
			"current_phase": currentPhase,
			"phase_message": phaseMessage,
			"planning_done": planningDone,
			"plan":          plan,
			"plan_json":     plan,
		})
	}
}

func approveProjectHandler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		project, err := loadProjectInSpace(deps.PB, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		switch project.Status {
		case models.ProjectAwaitingApproval, models.ProjectApproved, models.ProjectDraft:
		default:
			return c.JSON(http.StatusConflict, map[string]string{
				"error": fmt.Sprintf("cannot approve project in status %q", project.Status),
			})
		}

		corrID, err := spawnProjectExecution(c, deps, project)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		itemCorrIDs, err := projects.SpawnItemExecutions(
			c.Request().Context(),
			c,
			projects.ExecutionDeps{PB: deps.PB, NC: deps.NC},
			project,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		project.Status = models.ProjectExecuting
		project.ParentTaskID = corrID
		now := time.Now().UTC().Format(time.RFC3339)
		project.StartedAt = now
		updated, err := deps.PB.UpdateProject(project.ID, project)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success":        true,
			"project":        updated,
			"parent_task_id": corrID,
			"item_task_ids":  itemCorrIDs,
		})
	}
}

func transitionProjectHandler(pb *pbclient.Client, status string) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		project, err := loadProjectInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		project.Status = status
		if status == models.ProjectCompleted || status == models.ProjectCancelled {
			project.CompletedAt = time.Now().UTC().Format(time.RFC3339)
		}
		updated, err := pb.UpdateProject(project.ID, project)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if status == models.ProjectCancelled {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"project": updated,
				"type":    "dw_project_cancelled",
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "project": updated})
	}
}

func skipProjectTaskHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		if _, err := loadProjectInSpace(pb, spaceID, c.Param("id")); err != nil {
			return httpx.MapPocketError(c, err)
		}
		item, err := loadProjectItemInSpace(pb, spaceID, c.Param("id"), c.Param("taskId"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		item.Status = models.ItemSkipped
		updated, err := pb.UpdateProjectItem(item.ID, item)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "item": updated})
	}
}

func retryProjectTaskHandler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		if _, err := loadProjectInSpace(deps.PB, spaceID, c.Param("id")); err != nil {
			return httpx.MapPocketError(c, err)
		}
		project, err := loadProjectInSpace(deps.PB, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		item, err := loadProjectItemInSpace(deps.PB, spaceID, c.Param("id"), c.Param("taskId"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		corrID, err := projects.SpawnSingleItemExecution(
			c.Request().Context(),
			c,
			projects.ExecutionDeps{PB: deps.PB, NC: deps.NC},
			project,
			item,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		updated, err := deps.PB.GetProjectItem(item.ID)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success":           true,
			"item":              updated,
			"execution_task_id": corrID,
		})
	}
}

func enqueuePlanning(ctx context.Context, deps Deps, projectID, spaceID string) error {
	if deps.NC == nil {
		return fmt.Errorf("nats client not configured")
	}
	return deps.NC.PublishProjectPlan(ctx, models.ProjectPlanCommand{
		ProjectID: projectID,
		SpaceID:   spaceID,
	})
}

func spawnProjectExecution(c echo.Context, deps Deps, project models.Project) (string, error) {
	if deps.NC == nil {
		return "", fmt.Errorf("nats client not configured")
	}
	prompt := executionPrompt(project)
	workers := project.TeamAgentIDs
	workflow := ""
	if len(workers) > 0 {
		workflow = models.WorkflowSupervisor
	}
	task := models.Task{
		SpaceID:        project.SpaceID,
		AgentID:        project.PlannerAgentID,
		Prompt:         prompt,
		Workflow:       workflow,
		WorkerAgentIDs: workers,
		Status:         models.TaskQueued,
	}
	if user, ok := apimw.UserFromContext(c); ok {
		task.UserID = user.ID
	}
	corrID := fmt.Sprintf("project-%s-%d", project.ID, time.Now().UnixNano())
	task.CorrelationID = corrID

	if _, err := deps.PB.CreateTask(task); err != nil {
		return "", err
	}
	ctx := common.WithCorrelationID(context.Background(), corrID)
	if err := deps.NC.PublishOrchestrator(ctx, task); err != nil {
		return "", err
	}
	return corrID, nil
}

func executionPrompt(project models.Project) string {
	var b strings.Builder
	b.WriteString("Execute this approved project plan.\n\n")
	if project.Goal != "" {
		b.WriteString("Goal: ")
		b.WriteString(project.Goal)
		b.WriteString("\n\n")
	}
	if project.PlanJSON != nil {
		if phases, ok := project.PlanJSON["phases"].(map[string]interface{}); ok {
			if prd, ok := phases[models.PlanPhasePRD].(map[string]interface{}); ok {
				if content, ok := prd["content"].(string); ok && content != "" {
					b.WriteString("PRD:\n")
					b.WriteString(content)
					b.WriteString("\n\n")
				}
			}
		}
	}
	b.WriteString("Work through the project kanban items and produce deliverables.")
	return b.String()
}

func projectplannerTitle(title, goal string) string {
	title = strings.TrimSpace(title)
	if title != "" {
		return title
	}
	goal = strings.TrimSpace(goal)
	if len(goal) > 120 {
		return goal[:120] + "..."
	}
	if goal != "" {
		return goal
	}
	return "Untitled project"
}
