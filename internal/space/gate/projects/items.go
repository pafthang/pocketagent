package projectapis

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/httpx"
)

func listProjectItemsHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		if _, err := loadProjectInSpace(pb, spaceID, c.Param("id")); err != nil {
			return httpx.MapPocketError(c, err)
		}

		page, perPage := parsePageParams(c)
		filter := pbclient.ProjectItemsFilter(spaceID, c.Param("id"))
		if status := strings.TrimSpace(c.QueryParam("status")); status != "" {
			filter += fmt.Sprintf(" && status = %q", status)
		}

		items, total, err := pb.ListProjectItems(pbclient.ListOptions{
			Page: page, PerPage: perPage, Filter: filter,
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"items": items,
			"total": total,
		})
	}
}

func createProjectItemHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		project, err := loadProjectInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		var req CreateProjectItemRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		if strings.TrimSpace(req.Title) == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is required"})
		}

		item := req.ToModel(spaceID, project.ID)
		if item.Priority == "" {
			item.Priority = "medium"
		}

		stored, err := pb.CreateProjectItem(item)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusCreated, stored)
	}
}

func patchProjectItemHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		if _, err := loadProjectInSpace(pb, spaceID, c.Param("id")); err != nil {
			return httpx.MapPocketError(c, err)
		}
		existing, err := loadProjectItemInSpace(pb, spaceID, c.Param("id"), c.Param("itemId"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		var req PatchProjectItemRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		updated := existing
		req.ApplyPatch(&updated)
		if updated.Title == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is required"})
		}

		stored, err := pb.UpdateProjectItem(existing.ID, updated)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, stored)
	}
}

func deleteProjectItemHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		if _, err := loadProjectInSpace(pb, spaceID, c.Param("id")); err != nil {
			return httpx.MapPocketError(c, err)
		}
		item, err := loadProjectItemInSpace(pb, spaceID, c.Param("id"), c.Param("itemId"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if err := pb.DeleteProjectItem(item.ID); err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}