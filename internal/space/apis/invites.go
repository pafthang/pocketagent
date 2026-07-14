package spaceapis

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space/audit"
	"github.com/pafthang/pocketagent/internal/space/auth"
	"github.com/pafthang/pocketagent/internal/space/rbac"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func createInviteHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionInviteWrite); err != nil {
			return httpx.RespondError(c, err)
		}

		var req struct {
			Email string `json:"email"`
			Role  string `json:"role"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		if req.Email == "" {
			return httpx.RespondError(c, httpx.ErrBadRequest("email is required"))
		}
		if req.Role == "" {
			req.Role = models.RoleViewer
		}
		if !isValidRole(req.Role) {
			return httpx.RespondError(c, httpx.ErrBadRequest("invalid role"))
		}

		if existingUser, found, err := d.PB.FindUserByEmail(req.Email); err != nil {
			return httpx.MapPocketError(c, err)
		} else if found {
			filter := fmt.Sprintf("space_id = %q && user_id = %q", spaceID, existingUser.ID)
			members, _, err := d.PB.ListSpaceMembers(pbclient.ListOptions{Page: 1, PerPage: 1, Filter: filter})
			if err != nil {
				return httpx.MapPocketError(c, err)
			}
			if len(members) > 0 {
				return httpx.RespondError(c, httpx.ErrConflict("user is already a member"))
			}
		}

		pendingFilter := fmt.Sprintf(`email = %q && status = %q`, req.Email, models.InvitePending)
		pending, _, err := d.PB.ListSpaceInvites(spaceID, pbclient.ListOptions{Page: 1, PerPage: 1, Filter: pendingFilter})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if len(pending) > 0 {
			return httpx.RespondError(c, httpx.ErrConflict("pending invite already exists"))
		}

		token, err := auth.GenerateToken()
		if err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("generate token", err))
		}

		expiresAt := time.Now().UTC().Add(time.Duration(d.Cfg.InviteTTLHours) * time.Hour).Format(time.RFC3339)
		invite, err := d.PB.CreateSpaceInvite(models.SpaceInvite{
			SpaceID:   spaceID,
			Email:     req.Email,
			Role:      req.Role,
			InvitedBy: user.ID,
			ExpiresAt: expiresAt,
		}, auth.HashToken(token))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		d.Audit.Record(c, models.AuditLog{
			SpaceID:      spaceID,
			ActorID:      user.ID,
			ActorEmail:   user.Email,
			Action:       audit.AuditInviteCreate,
			ResourceType: "invite",
			ResourceID:   invite.ID,
			Metadata:     map[string]interface{}{"email": req.Email, "role": req.Role},
		})
		logInviteLink(d.Log, d.Cfg.PublicBaseURL, token)

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"invite":     invite,
			"token":      token,
			"accept_url": d.Cfg.PublicBaseURL + "/invites/accept",
			"expires_at": expiresAt,
		})
	}
}

func listInvitesHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionInviteRead); err != nil {
			return httpx.RespondError(c, err)
		}

		invites, total, err := d.PB.ListSpaceInvites(spaceID, pbclient.ListOptions{Page: 1, PerPage: 200})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"invites": invites, "total": total})
	}
}

func revokeInviteHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionInviteWrite); err != nil {
			return httpx.RespondError(c, err)
		}

		inviteID := c.Param("inviteId")
		invite, err := d.PB.GetSpaceInvite(inviteID)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if invite.SpaceID != spaceID {
			return httpx.RespondError(c, httpx.ErrNotFound("invite not found"))
		}
		if invite.Status != models.InvitePending {
			return httpx.RespondError(c, httpx.ErrConflict("invite is not pending"))
		}

		updated, err := d.PB.UpdateSpaceInvite(inviteID, map[string]interface{}{"status": models.InviteRevoked})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		d.Audit.Record(c, models.AuditLog{
			SpaceID:      spaceID,
			ActorID:      user.ID,
			ActorEmail:   user.Email,
			Action:       audit.AuditInviteRevoke,
			ResourceType: "invite",
			ResourceID:   inviteID,
		})
		return c.JSON(http.StatusOK, updated)
	}
}

func previewInviteHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Param("token")
		if token == "" {
			return httpx.RespondError(c, httpx.ErrBadRequest("token is required"))
		}

		invite, err := d.PB.FindSpaceInviteByTokenHash(auth.HashToken(token))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if inviteExpired(invite.ExpiresAt) {
			_, _ = d.PB.UpdateSpaceInvite(invite.ID, map[string]interface{}{"status": models.InviteExpired})
			return httpx.RespondError(c, httpx.ErrBadRequest("invite expired"))
		}

		space, err := d.PB.GetSpace(invite.SpaceID)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		return c.JSON(http.StatusOK, models.InvitePreview{
			SpaceID:   invite.SpaceID,
			SpaceName: space.Name,
			Email:     invite.Email,
			Role:      invite.Role,
			Status:    invite.Status,
			ExpiresAt: invite.ExpiresAt,
		})
	}
}

func acceptInviteHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Token    string `json:"token"`
			Password string `json:"password"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.Token == "" || req.Password == "" {
			return httpx.RespondError(c, httpx.ErrBadRequest("token and password are required"))
		}
		if len(req.Password) < 8 {
			return httpx.RespondError(c, httpx.ErrBadRequest("password must be at least 8 characters"))
		}

		invite, err := d.PB.FindSpaceInviteByTokenHash(auth.HashToken(req.Token))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if inviteExpired(invite.ExpiresAt) {
			_, _ = d.PB.UpdateSpaceInvite(invite.ID, map[string]interface{}{"status": models.InviteExpired})
			return httpx.RespondError(c, httpx.ErrBadRequest("invite expired"))
		}

		session, userID, err := resolveInviteUser(d.PB, invite.Email, req.Password)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		filter := fmt.Sprintf("space_id = %q && user_id = %q", invite.SpaceID, userID)
		existing, _, err := d.PB.ListSpaceMembers(pbclient.ListOptions{Page: 1, PerPage: 1, Filter: filter})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if len(existing) == 0 {
			if _, err := d.PB.CreateSpaceMember(models.SpaceMember{
				SpaceID: invite.SpaceID,
				UserID:  userID,
				Role:    invite.Role,
			}); err != nil {
				return httpx.MapPocketError(c, err)
			}
			d.Audit.Record(c, models.AuditLog{
				SpaceID:      invite.SpaceID,
				ActorID:      userID,
				ActorEmail:   invite.Email,
				Action:       audit.AuditMemberAdd,
				ResourceType: "member",
				Metadata:     map[string]interface{}{"role": invite.Role, "via": "invite"},
			})
		}

		if _, err := d.PB.UpdateSpaceInvite(invite.ID, map[string]interface{}{"status": models.InviteAccepted}); err != nil {
			return httpx.MapPocketError(c, err)
		}

		d.Audit.Record(c, models.AuditLog{
			SpaceID:      invite.SpaceID,
			ActorID:      userID,
			ActorEmail:   invite.Email,
			Action:       audit.AuditInviteAccept,
			ResourceType: "invite",
			ResourceID:   invite.ID,
		})

		return c.JSON(http.StatusOK, session)
	}
}

func resolveInviteUser(pb *pbclient.Client, email, password string) (models.AuthSession, string, error) {
	if _, found, err := pb.FindUserByEmail(email); err != nil {
		return models.AuthSession{}, "", err
	} else if found {
		session, err := pb.AuthWithPassword(email, password)
		if err != nil {
			return models.AuthSession{}, "", err
		}
		return session, session.User.ID, nil
	}

	if _, err := pb.RegisterUser(email, password); err != nil {
		return models.AuthSession{}, "", err
	}
	session, err := pb.AuthWithPassword(email, password)
	if err != nil {
		return models.AuthSession{}, "", err
	}
	return session, session.User.ID, nil
}

func inviteExpired(expiresAt string) bool {
	if expiresAt == "" {
		return true
	}
	t, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return true
	}
	return time.Now().UTC().After(t)
}
