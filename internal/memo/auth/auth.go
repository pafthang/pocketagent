package auth

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// HeaderServiceToken is the internal service-to-service auth header for memo.
const HeaderServiceToken = "X-Memo-Token"

// RequireServiceToken rejects requests without a valid memo service token.
// When expected is empty (dev), auth is skipped.
func RequireServiceToken(expected string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.TrimSpace(expected) == "" {
				return next(c)
			}

			token := strings.TrimSpace(c.Request().Header.Get(HeaderServiceToken))
			if token == "" {
				token = extractBearer(c.Request().Header.Get("Authorization"))
			}
			if subtle.ConstantTimeCompare([]byte(token), []byte(expected)) != 1 {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			return next(c)
		}
	}
}

func extractBearer(header string) string {
	header = strings.TrimSpace(header)
	if len(header) < 7 {
		return ""
	}
	prefix := strings.ToLower(header[:7])
	if prefix != "bearer " {
		return ""
	}
	return strings.TrimSpace(header[7:])
}