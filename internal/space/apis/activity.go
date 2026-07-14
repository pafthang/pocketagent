package spaceapis

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space/activity"
	"github.com/pafthang/pocketagent/internal/space/auth"
	"github.com/pafthang/pocketagent/internal/space/rbac"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func listActivityHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionTaskRead); err != nil {
			return httpx.RespondError(c, err)
		}

		limit := 50
		if p := c.QueryParam("limit"); p != "" {
			if n, err := strconv.Atoi(p); err == nil && n > 0 && n <= 200 {
				limit = n
			}
		}
		perPage := limit
		if perPage < 100 {
			perPage = 100
		}

		taskEvents, _, err := d.PB.ListTaskEvents(spaceID, pbclient.ListOptions{Page: 1, PerPage: perPage})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		includeAudit := rbac.RoleAllows(role, rbac.ActionAuditRead)
		var audits []models.AuditLog
		if includeAudit {
			audits, _, err = d.PB.ListAuditLogs(spaceID, pbclient.ListOptions{Page: 1, PerPage: perPage})
			if err != nil {
				return httpx.MapPocketError(c, err)
			}
		}

		entries := activity.BuildFeed(taskEvents, audits, includeAudit, limit)
		return c.JSON(http.StatusOK, models.ActivityListResponse{
			Entries: entries,
			Total:   len(entries),
		})
	}
}
