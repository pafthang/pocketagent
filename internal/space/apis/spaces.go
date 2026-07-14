package spaceapis

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space/audit"
	"github.com/pafthang/pocketagent/internal/space/auth"
	"github.com/pafthang/pocketagent/internal/space/rbac"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func listSpacesHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		if super, err := d.Auth.IsSuperAdmin(user.ID); err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("check super admin", err))
		} else if super {
			spaces, total, err := d.PB.ListSpaces(pbclient.ListOptions{Page: 1, PerPage: 200})
			if err != nil {
				return httpx.MapPocketError(c, err)
			}
			return c.JSON(http.StatusOK, map[string]interface{}{"spaces": spaces, "total": total})
		}

		filter := fmt.Sprintf("user_id = %q", user.ID)
		memberships, _, err := d.PB.ListSpaceMembers(pbclient.ListOptions{Page: 1, PerPage: 200, Filter: filter})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		spaces := make([]models.Space, 0, len(memberships))
		for _, member := range memberships {
			space, err := d.PB.GetSpace(member.SpaceID)
			if err != nil {
				continue
			}
			spaces = append(spaces, space)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"spaces": spaces,
			"total":  len(spaces),
		})
	}
}

func createSpaceHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		var req struct {
			Name        string `json:"name"`
			Slug        string `json:"slug"`
			Description string `json:"description"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.Name == "" {
			return httpx.RespondError(c, httpx.ErrBadRequest("name is required"))
		}
		if req.Slug == "" {
			req.Slug = slugify(req.Name)
		}
		if !slugPattern.MatchString(req.Slug) {
			return httpx.RespondError(c, httpx.ErrBadRequest("invalid slug"))
		}
		if req.Slug == models.SystemSpaceSlug {
			return httpx.RespondError(c, httpx.ErrForbidden("cannot create system space"))
		}

		space, err := d.PB.CreateSpace(models.Space{
			Name:        req.Name,
			Slug:        req.Slug,
			Description: req.Description,
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		if _, err := d.PB.CreateSpaceMember(models.SpaceMember{
			SpaceID: space.ID,
			UserID:  user.ID,
			Role:    models.RoleAdmin,
		}); err != nil {
			_ = d.PB.DeleteSpace(space.ID)
			return httpx.MapPocketError(c, err)
		}

		d.Audit.Record(c, models.AuditLog{
			SpaceID:      space.ID,
			ActorID:      user.ID,
			ActorEmail:   user.Email,
			Action:       audit.AuditSpaceCreate,
			ResourceType: "space",
			ResourceID:   space.ID,
		})
		return c.JSON(http.StatusCreated, space)
	}
}

func getSpaceHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionSpaceRead); err != nil {
			return httpx.RespondError(c, err)
		}

		space, err := d.PB.GetSpace(spaceID)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, space)
	}
}

func updateSpaceHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionSpaceWrite); err != nil {
			return httpx.RespondError(c, err)
		}

		existing, err := d.PB.GetSpace(spaceID)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if existing.IsSystem {
			return httpx.RespondError(c, httpx.ErrForbidden("system space cannot be modified"))
		}

		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.Name != "" {
			existing.Name = req.Name
		}
		if req.Description != "" {
			existing.Description = req.Description
		}

		updated, err := d.PB.UpdateSpace(spaceID, existing)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		d.Audit.Record(c, models.AuditLog{
			SpaceID:      spaceID,
			ActorID:      user.ID,
			ActorEmail:   user.Email,
			Action:       audit.AuditSpaceUpdate,
			ResourceType: "space",
			ResourceID:   spaceID,
		})
		return c.JSON(http.StatusOK, updated)
	}
}

func deleteSpaceHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionSpaceDelete); err != nil {
			return httpx.RespondError(c, err)
		}

		space, err := d.PB.GetSpace(spaceID)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if space.IsSystem {
			return httpx.RespondError(c, httpx.ErrForbidden("system space cannot be deleted"))
		}

		if err := d.PB.DeleteSpace(spaceID); err != nil {
			return httpx.MapPocketError(c, err)
		}
		d.Audit.Record(c, models.AuditLog{
			SpaceID:      spaceID,
			ActorID:      user.ID,
			ActorEmail:   user.Email,
			Action:       audit.AuditSpaceDelete,
			ResourceType: "space",
			ResourceID:   spaceID,
		})
		return c.NoContent(http.StatusNoContent)
	}
}