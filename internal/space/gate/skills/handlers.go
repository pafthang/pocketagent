package skillapis

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/httpx"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/exec/tools"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	taskapis "github.com/pafthang/pocketagent/internal/task/apis"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func RegisterRoutes(tenant *echo.Group, pb *pbclient.Client, nc *natsclient.Client, readAction, writeAction echo.MiddlewareFunc) {
	toolCfg := defaultToolConfig()
	tenant.GET("/skills/catalog", listSkillsCatalogHandler(pb), readAction)

	tenant.GET("/skills", listSkillsHandler(pb), readAction)
	tenant.POST("/skills", createSkillHandler(pb), writeAction)
	tenant.GET("/skills/:id", getSkillHandler(pb), readAction)
	tenant.PATCH("/skills/:id", patchSkillHandler(pb), writeAction)
	tenant.DELETE("/skills/:id", deleteSkillHandler(pb), writeAction)
	tenant.POST("/skills/:id/run", runSkillHandler(pb, nc, toolCfg), writeAction)
	tenant.GET("/tools", ListToolsHandler(pb, toolCfg), readAction)
}

func listSkillsHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 200
		}

		skills, total, err := pb.ListSkills(pbclient.ListOptions{
			Page:    page,
			PerPage: perPage,
			Filter:  fmt.Sprintf("space_id = %q", spaceID),
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"skills":   skills,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		})
	}
}

func createSkillHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		var req CreateSkillRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		skill, err := req.ToModel(spaceID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		if err := ensureUniqueSkillName(pb, spaceID, skill.Name, ""); err != nil {
			return httpx.MapPocketError(c, err)
		}

		stored, err := pb.CreateSkill(skill)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusCreated, stored)
	}
}

func getSkillHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		skill, err := loadSkillInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, skill)
	}
}

func patchSkillHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		existing, err := loadSkillInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		var req PatchSkillRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		updated := existing
		req.ApplyPatch(&updated)
		if updated.Name == "" || updated.Prompt == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "name and prompt are required"})
		}
		if updated.Name != existing.Name {
			if err := ensureUniqueSkillName(pb, spaceID, updated.Name, existing.ID); err != nil {
				return httpx.MapPocketError(c, err)
			}
		}

		stored, err := pb.UpdateSkillRecord(updated)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, stored)
	}
}

func deleteSkillHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		skill, err := loadSkillInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if err := pb.DeleteSkill(skill.ID); err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func runSkillHandler(pb *pbclient.Client, nc *natsclient.Client, toolCfg tools.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		skill, err := loadSkillInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		var req RunSkillRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		if len(skill.Tools) > 0 {
			if err := tools.ValidateToolAllowList(pb, spaceID, toolCfg, skill.Tools); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}
		}

		prompt := composeSkillPrompt(skill, req.Input)
		task := models.Task{
			AgentID: req.AgentID,
			Prompt:  prompt,
			SkillID: skill.ID,
		}
		if len(skill.Tools) > 0 {
			task.Tools = append([]string{}, skill.Tools...)
		}

		return taskapis.PublishTaskWithTools(c, nc, pb, task, toolCfg)
	}
}

func listSkillsCatalogHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		names, err := listInstalledSkillNames(pb, spaceID)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		entries, err := loadSkillsCatalog(nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, catalogToResponses(entries, installedSkillNames(names)))
	}
}

func composeSkillPrompt(skill models.Skill, input string) string {
	prompt := skill.Prompt
	input = strings.TrimSpace(input)
	if input == "" {
		return prompt
	}
	if strings.Contains(prompt, "{{input}}") {
		return strings.ReplaceAll(prompt, "{{input}}", input)
	}
	return prompt + "\n\n" + input
}

func listInstalledSkillNames(pb *pbclient.Client, spaceID string) ([]string, error) {
	skills, _, err := pb.ListSkills(pbclient.ListOptions{
		Page:    1,
		PerPage: 200,
		Filter:  fmt.Sprintf("space_id = %q", spaceID),
	})
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(skills))
	for _, skill := range skills {
		names = append(names, skill.Name)
	}
	return names, nil
}

func ensureUniqueSkillName(pb *pbclient.Client, spaceID, name, excludeID string) error {
	existing, err := pbclient.FindSkillByName(pb, spaceID, name)
	if err != nil {
		var apiErr *pbclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}
	if excludeID != "" && existing.ID == excludeID {
		return nil
	}
	return fmt.Errorf("skill %q already exists in this space", name)
}

func loadSkillInSpace(pb *pbclient.Client, spaceID, id string) (models.Skill, error) {
	skill, err := pb.GetSkill(id)
	if err != nil {
		return models.Skill{}, err
	}
	if skill.SpaceID != spaceID {
		return models.Skill{}, &pbclient.APIError{StatusCode: http.StatusNotFound, Message: "skill not found"}
	}
	return skill, nil
}
