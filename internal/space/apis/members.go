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

func listMembersHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionMemberRead); err != nil {
			return httpx.RespondError(c, err)
		}

		filter := fmt.Sprintf("space_id = %q", spaceID)
		members, total, err := d.PB.ListSpaceMembers(pbclient.ListOptions{Page: 1, PerPage: 200, Filter: filter})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"members": members, "total": total})
	}
}

func addMemberHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionMemberWrite); err != nil {
			return httpx.RespondError(c, err)
		}

		var req struct {
			UserID string `json:"user_id"`
			Role   string `json:"role"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.UserID == "" {
			return httpx.RespondError(c, httpx.ErrBadRequest("user_id is required"))
		}
		if req.Role == "" {
			req.Role = models.RoleViewer
		}
		if !isValidRole(req.Role) {
			return httpx.RespondError(c, httpx.ErrBadRequest("invalid role"))
		}

		filter := fmt.Sprintf("space_id = %q && user_id = %q", spaceID, req.UserID)
		existing, _, err := d.PB.ListSpaceMembers(pbclient.ListOptions{Page: 1, PerPage: 1, Filter: filter})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if len(existing) > 0 {
			return httpx.RespondError(c, httpx.ErrConflict("user is already a member"))
		}

		member, err := d.PB.CreateSpaceMember(models.SpaceMember{
			SpaceID: spaceID,
			UserID:  req.UserID,
			Role:    req.Role,
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		d.Audit.Record(c, models.AuditLog{
			SpaceID:      spaceID,
			ActorID:      user.ID,
			ActorEmail:   user.Email,
			Action:       audit.AuditMemberAdd,
			ResourceType: "member",
			ResourceID:   member.ID,
			Metadata:     map[string]interface{}{"user_id": req.UserID, "role": req.Role},
		})
		return c.JSON(http.StatusCreated, member)
	}
}

func updateMemberHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionMemberWrite); err != nil {
			return httpx.RespondError(c, err)
		}

		memberID := c.Param("memberId")
		member, err := d.PB.GetSpaceMember(memberID)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if member.SpaceID != spaceID {
			return httpx.RespondError(c, httpx.ErrNotFound("member not found"))
		}

		var req struct {
			Role string `json:"role"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		if !isValidRole(req.Role) {
			return httpx.RespondError(c, httpx.ErrBadRequest("invalid role"))
		}

		member.Role = req.Role
		updated, err := d.PB.UpdateSpaceMember(memberID, member)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		d.Audit.Record(c, models.AuditLog{
			SpaceID:      spaceID,
			ActorID:      user.ID,
			ActorEmail:   user.Email,
			Action:       audit.AuditMemberUpdate,
			ResourceType: "member",
			ResourceID:   memberID,
			Metadata:     map[string]interface{}{"role": req.Role},
		})
		return c.JSON(http.StatusOK, updated)
	}
}

func deleteMemberHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionMemberWrite); err != nil {
			return httpx.RespondError(c, err)
		}

		memberID := c.Param("memberId")
		member, err := d.PB.GetSpaceMember(memberID)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if member.SpaceID != spaceID {
			return httpx.RespondError(c, httpx.ErrNotFound("member not found"))
		}

		if err := d.PB.DeleteSpaceMember(memberID); err != nil {
			return httpx.MapPocketError(c, err)
		}
		d.Audit.Record(c, models.AuditLog{
			SpaceID:      spaceID,
			ActorID:      user.ID,
			ActorEmail:   user.Email,
			Action:       audit.AuditMemberRemove,
			ResourceType: "member",
			ResourceID:   memberID,
		})
		return c.NoContent(http.StatusNoContent)
	}
}