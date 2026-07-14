package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
)

// AuthMiddleware validates PocketBase JWT and stores the user in context.
func AuthMiddleware(pb *pbclient.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := ExtractToken(c.Request().Header.Get("Authorization"))
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authorization required"})
			}

			session, err := pb.AuthRefresh(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}

			ctx := WithUser(c.Request().Context(), session.User)
			c.SetRequest(c.Request().WithContext(ctx))
			c.Set("auth_token", session.Token)
			return next(c)
		}
	}
}

// ExtractToken parses bearer or raw authorization header value.
func ExtractToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return header
}