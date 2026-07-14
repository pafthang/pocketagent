package spaceapis

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space/auth"
	"github.com/pafthang/pocketagent/internal/space/rbac"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func listTeamsHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		spaceID := c.Param("spaceId")
		role, err := d.Auth.MemberRole(user.ID, spaceID)
		if err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("resolve role", err))
		}
		if err := rbac.RequireRole(role, rbac.ActionTeamRead); err != nil {
			return httpx.RespondError(c, err)
		}

		filter := fmt.Sprintf("space_id = %q", spaceID)
		teams, total, err := d.PB.ListTeams(pbclient.ListOptions{Page: 1, PerPage: 200, Filter: filter})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"teams": teams, "total": total})
	}
}

func createTeamHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		spaceID := c.Param("spaceId")
		role, err := d.Auth.MemberRole(user.ID, spaceID)
		if err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("resolve role", err))
		}
		if err := rbac.RequireRole(role, rbac.ActionTeamWrite); err != nil {
			return httpx.RespondError(c, err)
		}

		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.Name == "" {
			return httpx.RespondError(c, httpx.ErrBadRequest("name is required"))
		}

		team, err := d.PB.CreateTeam(models.Team{
			SpaceID:     spaceID,
			Name:        req.Name,
			Description: req.Description,
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusCreated, team)
	}
}

func getTeamHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		spaceID := c.Param("spaceId")
		role, err := d.Auth.MemberRole(user.ID, spaceID)
		if err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("resolve role", err))
		}
		if err := rbac.RequireRole(role, rbac.ActionTeamRead); err != nil {
			return httpx.RespondError(c, err)
		}

		team, err := d.PB.GetTeam(c.Param("teamId"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if team.SpaceID != spaceID {
			return httpx.RespondError(c, httpx.ErrNotFound("team not found"))
		}
		return c.JSON(http.StatusOK, team)
	}
}

func updateTeamHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		spaceID := c.Param("spaceId")
		role, err := d.Auth.MemberRole(user.ID, spaceID)
		if err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("resolve role", err))
		}
		if err := rbac.RequireRole(role, rbac.ActionTeamWrite); err != nil {
			return httpx.RespondError(c, err)
		}

		team, err := d.PB.GetTeam(c.Param("teamId"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if team.SpaceID != spaceID {
			return httpx.RespondError(c, httpx.ErrNotFound("team not found"))
		}

		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.Name != "" {
			team.Name = req.Name
		}
		if req.Description != "" {
			team.Description = req.Description
		}

		updated, err := d.PB.UpdateTeam(team.ID, team)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, updated)
	}
}

func deleteTeamHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		spaceID := c.Param("spaceId")
		role, err := d.Auth.MemberRole(user.ID, spaceID)
		if err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("resolve role", err))
		}
		if err := rbac.RequireRole(role, rbac.ActionTeamDelete); err != nil {
			return httpx.RespondError(c, err)
		}

		team, err := d.PB.GetTeam(c.Param("teamId"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if team.SpaceID != spaceID {
			return httpx.RespondError(c, httpx.ErrNotFound("team not found"))
		}

		if err := d.PB.DeleteTeam(team.ID); err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func listTeamMembersHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		spaceID := c.Param("spaceId")
		role, err := d.Auth.MemberRole(user.ID, spaceID)
		if err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("resolve role", err))
		}
		if err := rbac.RequireRole(role, rbac.ActionTeamRead); err != nil {
			return httpx.RespondError(c, err)
		}

		team, err := d.PB.GetTeam(c.Param("teamId"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if team.SpaceID != spaceID {
			return httpx.RespondError(c, httpx.ErrNotFound("team not found"))
		}

		filter := fmt.Sprintf("team_id = %q", team.ID)
		members, total, err := d.PB.ListTeamMembers(pbclient.ListOptions{Page: 1, PerPage: 200, Filter: filter})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"members": members, "total": total})
	}
}

func addTeamMemberHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		spaceID := c.Param("spaceId")
		role, err := d.Auth.MemberRole(user.ID, spaceID)
		if err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("resolve role", err))
		}
		if err := rbac.RequireRole(role, rbac.ActionTeamWrite); err != nil {
			return httpx.RespondError(c, err)
		}

		team, err := d.PB.GetTeam(c.Param("teamId"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if team.SpaceID != spaceID {
			return httpx.RespondError(c, httpx.ErrNotFound("team not found"))
		}

		var req struct {
			MemberType string `json:"member_type"`
			MemberID   string `json:"member_id"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.MemberID == "" {
			return httpx.RespondError(c, httpx.ErrBadRequest("member_id is required"))
		}
		if req.MemberType == "" {
			req.MemberType = models.MemberTypeUser
		}
		if req.MemberType != models.MemberTypeUser && req.MemberType != models.MemberTypeAgent {
			return httpx.RespondError(c, httpx.ErrBadRequest("invalid member_type"))
		}

		member, err := d.PB.CreateTeamMember(models.TeamMember{
			TeamID:     team.ID,
			MemberType: req.MemberType,
			MemberID:   req.MemberID,
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusCreated, member)
	}
}

func deleteTeamMemberHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		spaceID := c.Param("spaceId")
		role, err := d.Auth.MemberRole(user.ID, spaceID)
		if err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("resolve role", err))
		}
		if err := rbac.RequireRole(role, rbac.ActionTeamWrite); err != nil {
			return httpx.RespondError(c, err)
		}

		team, err := d.PB.GetTeam(c.Param("teamId"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if team.SpaceID != spaceID {
			return httpx.RespondError(c, httpx.ErrNotFound("team not found"))
		}

		memberID := c.Param("memberId")
		if err := d.PB.DeleteTeamMember(memberID); err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}