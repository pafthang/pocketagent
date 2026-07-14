package taskapis

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pafthang/pocketagent/pkgs/httpx"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	taskschedule "github.com/pafthang/pocketagent/internal/task/schedule"
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func createScheduleHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var schedule models.Schedule
		if err := c.Bind(&schedule); err != nil {
			return err
		}

		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}
		if schedule.Name == "" || schedule.Prompt == "" || schedule.CronExpr == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "name, prompt and cron_expr are required"})
		}
		if err := common.GuardPrompt(nil, common.LoadPromptGuardConfig(), schedule.Prompt); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		if err := taskschedule.ValidateCronExpr(schedule.CronExpr); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		schedule.SpaceID = spaceID
		schedule.Enabled = true
		nextRun, err := taskschedule.NextCronRun(schedule.CronExpr, time.Now().UTC())
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		schedule.NextRunAt = nextRun.UTC().Format(time.RFC3339)

		if schedule.AgentID != "" {
			agent, err := pb.GetAgent(schedule.AgentID)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid agent_id"})
			}
			if agent.SpaceID != "" && agent.SpaceID != spaceID {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "agent does not belong to this space"})
			}
		}

		stored, err := pb.CreateSchedule(schedule)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusCreated, stored)
	}
}

func listSchedulesHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
		filter := fmt.Sprintf("space_id = %q", spaceID)

		schedules, total, err := pb.ListSchedules(pbclient.ListOptions{
			Page:    page,
			PerPage: perPage,
			Filter:  filter,
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"schedules": schedules,
			"total":     total,
		})
	}
}

func getScheduleHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		schedule, err := pb.GetSchedule(c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if schedule.SpaceID != spaceID {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "schedule not found"})
		}
		return c.JSON(http.StatusOK, schedule)
	}
}

func updateScheduleHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		existing, err := pb.GetSchedule(c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if existing.SpaceID != spaceID {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "schedule not found"})
		}

		var patch models.Schedule
		if err := c.Bind(&patch); err != nil {
			return err
		}

		if patch.CronExpr != "" {
			if err := taskschedule.ValidateCronExpr(patch.CronExpr); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}
			nextRun, err := taskschedule.NextCronRun(patch.CronExpr, time.Now().UTC())
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}
			patch.NextRunAt = nextRun.UTC().Format(time.RFC3339)
		}

		updated, err := pb.UpdateSchedule(existing.ID, patch)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, updated)
	}
}

func deleteScheduleHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := apimw.SpaceIDFromContext(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": apimw.HeaderSpaceID + " header is required"})
		}

		schedule, err := pb.GetSchedule(c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if schedule.SpaceID != spaceID {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "schedule not found"})
		}

		if err := pb.DeleteSchedule(schedule.ID); err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
