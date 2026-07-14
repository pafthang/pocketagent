package taskapis

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pafthang/pocketagent/pkgs/httpx"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
)

func getTaskHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		task, err := pb.GetTaskByCorrelationID(c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if task.SpaceID != spaceID {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
		}

		if c.QueryParam("include") == "subtasks" {
			subtasks, err := pb.ListSubtasks(task.CorrelationID)
			if err != nil {
				return httpx.MapPocketError(c, err)
			}
			return c.JSON(http.StatusOK, map[string]interface{}{
				"task":     task,
				"subtasks": subtasks,
			})
		}

		return c.JSON(http.StatusOK, task)
	}
}

func deleteTaskHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		task, err := pb.GetTaskByCorrelationID(c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if task.SpaceID != spaceID {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
		}

		updated, err := pb.CancelTask(task.CorrelationID)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		return c.JSON(http.StatusOK, updated)
	}
}

func listTasksHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("per_page"))

		filter := fmt.Sprintf("space_id = %q", spaceID)
		tasks, total, err := pb.ListTasks(pbclient.ListOptions{
			Page:    page,
			PerPage: perPage,
			Filter:  filter,
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"tasks": tasks,
			"total": total,
		})
	}
}
