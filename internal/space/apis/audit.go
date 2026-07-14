package spaceapis

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/space/auth"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/internal/space/rbac"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
)

func listAuditLogsHandler(d Deps) echo.HandlerFunc {
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
		if err := rbac.RequireRole(role, rbac.ActionAuditRead); err != nil {
			return httpx.RespondError(c, err)
		}

		page := 1
		perPage := 50
		if p := c.QueryParam("page"); p != "" {
			if n, err := strconv.Atoi(p); err == nil && n > 0 {
				page = n
			}
		}
		if p := c.QueryParam("per_page"); p != "" {
			if n, err := strconv.Atoi(p); err == nil && n > 0 && n <= 200 {
				perPage = n
			}
		}

		logs, total, err := d.PB.ListAuditLogs(spaceID, pbclient.ListOptions{Page: page, PerPage: perPage})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"logs": logs, "total": total})
	}
}