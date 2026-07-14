package agent

import (
	"github.com/labstack/echo/v4"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// UserFromContext returns the authenticated user set by middleware.
func UserFromContext(c echo.Context) (models.AuthUser, bool) {
	return apimw.UserFromContext(c)
}

// SpaceIDFromContext returns the validated space ID from middleware.
func SpaceIDFromContext(c echo.Context) (string, bool) {
	return apimw.SpaceIDFromContext(c)
}

// AuthTokenFromContext returns the bearer token.
func AuthTokenFromContext(c echo.Context) string {
	return apimw.AuthTokenFromContext(c)
}
