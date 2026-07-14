package gateapis

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/gate/proxy"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func registerAuthRoutes(e *echo.Echo, deps Deps) {
	auth := e.Group("")
	if deps.RateLimit.EffectiveEnabled() {
		auth.Use(common.AuthRateLimiter(deps.RateLimit))
	}
	auth.POST("/auth/register", proxy.FixedPath(deps.Space, http.MethodPost, "/auth/register", false))
	auth.POST("/auth/login", proxy.FixedPath(deps.Space, http.MethodPost, "/auth/login", false))
	auth.POST("/auth/refresh", proxy.FixedPath(deps.Space, http.MethodPost, "/auth/refresh", false))
	auth.POST("/auth/verify-email", proxy.FixedPath(deps.Space, http.MethodPost, "/auth/verify-email", false))
	auth.GET("/invites/:token", func(c echo.Context) error {
		token := c.Param("token")
		return proxy.FixedPath(deps.Space, http.MethodGet, "/invites/"+token, false)(c)
	})
	auth.POST("/invites/accept", proxy.FixedPath(deps.Space, http.MethodPost, "/invites/accept", false))
}