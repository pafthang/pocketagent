package context

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/models"
)

const (
	HeaderSpaceID = "X-Space-Id"

	UserKey    = "auth_user"
	SpaceIDKey = "space_id"
	TokenKey   = "auth_token"
)

// SetUser stores the authenticated user in echo context.
func SetUser(c echo.Context, user models.AuthUser) { c.Set(UserKey, user) }

// SetSpaceID stores the validated tenant space ID in echo context.
func SetSpaceID(c echo.Context, spaceID string) { c.Set(SpaceIDKey, spaceID) }

// SetToken stores the bearer token in echo context.
func SetToken(c echo.Context, token string) { c.Set(TokenKey, token) }

// UserFromContext returns the authenticated user set by AuthMiddleware.
func UserFromContext(c echo.Context) (models.AuthUser, bool) {
	user, ok := c.Get(UserKey).(models.AuthUser)
	return user, ok
}

// SpaceIDFromContext returns the validated space ID from RequireSpace.
func SpaceIDFromContext(c echo.Context) (string, bool) {
	spaceID, ok := c.Get(SpaceIDKey).(string)
	return spaceID, ok && spaceID != ""
}

// AuthTokenFromContext returns the bearer token.
func AuthTokenFromContext(c echo.Context) string {
	token, _ := c.Get(TokenKey).(string)
	return token
}

// ExtractBearer reads Authorization header or token query param.
func ExtractBearer(c echo.Context) string {
	auth := strings.TrimSpace(c.Request().Header.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:])
	}
	if auth != "" {
		return auth
	}
	return strings.TrimSpace(c.QueryParam("token"))
}