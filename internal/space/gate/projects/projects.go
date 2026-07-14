package projectapis

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

func listProjectsHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		page, perPage := parsePageParams(c)
		filter := pbclient.ProjectsFilter(spaceID)
		if status := strings.TrimSpace(c.QueryParam("status")); status != "" {
			filter += fmt.Sprintf(" && status = %q", status)
		}

		projects, total, err := pb.ListProjects(pbclient.ListOptions{
			Page: page, PerPage: perPage, Filter: filter,
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		if c.QueryParam("format") == "mc" {
			mcProjects := make([]map[string]interface{}, 0, len(projects))
			for _, project := range projects {
				itemIDs, _ := projectItemIDs(pb, spaceID, project.ID)
				mcProjects = append(mcProjects, ToMCProject(project, itemIDs))
			}
			return c.JSON(http.StatusOK, map[string]interface{}{
				"projects": mcProjects,
				"count":    total,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"projects": projects,
			"total":    total,
		})
	}
}

func createProjectHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		var req CreateProjectRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		title := NormalizeTitle(req.Title, req.Goal)
		goal := strings.TrimSpace(req.Goal)
		if title == "Untitled project" && goal == "" && strings.TrimSpace(req.Description) == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "title or goal is required"})
		}

		creatorID := ""
		if user, ok := apimw.UserFromContext(c); ok {
			creatorID = user.ID
		}
		project := req.ToModel(spaceID, creatorID, title)
		if err := validatePlannerAgent(pb, spaceID, project.PlannerAgentID); err != nil {
			return httpx.MapPocketError(c, err)
		}

		stored, err := pb.CreateProject(project)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		if c.QueryParam("format") == "mc" {
			return c.JSON(http.StatusCreated, map[string]interface{}{
				"project": ToMCProject(stored, nil),
			})
		}
		return c.JSON(http.StatusCreated, stored)
	}
}

func getProjectHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		project, err := loadProjectInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		items, _, err := pb.ListProjectItems(pbclient.ListOptions{
			Page: 1, PerPage: 500, Filter: pbclient.ProjectItemsFilter(spaceID, project.ID),
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		if c.QueryParam("format") == "mc" || c.QueryParam("include") == "items" {
			itemIDs := make([]string, 0, len(items))
			mcTasks := make([]map[string]interface{}, 0, len(items))
			for _, item := range items {
				itemIDs = append(itemIDs, item.ID)
				mcTasks = append(mcTasks, ToMCTask(item, project.ID, project.CreatorID))
			}
			return c.JSON(http.StatusOK, map[string]interface{}{
				"project":  ToMCProject(project, itemIDs),
				"tasks":    mcTasks,
				"items":    items,
				"progress": Progress(items),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"project": project,
			"items":   items,
		})
	}
}

func patchProjectHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		existing, err := loadProjectInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		var req PatchProjectRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		updated := existing
		req.ApplyPatch(&updated, NormalizeTitle)
		if updated.Title == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is required"})
		}
		if err := validatePlannerAgent(pb, spaceID, updated.PlannerAgentID); err != nil {
			return httpx.MapPocketError(c, err)
		}

		stored, err := pb.UpdateProject(existing.ID, updated)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, stored)
	}
}

func deleteProjectHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		project, err := loadProjectInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if err := pb.DeleteProjectItemsByProject(project.ID); err != nil {
			return httpx.MapPocketError(c, err)
		}
		if err := pb.DeleteProject(project.ID); err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}