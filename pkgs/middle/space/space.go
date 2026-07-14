package space

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	mwctx "github.com/pafthang/pocketagent/pkgs/middle/context"
)

// Options configures RequireSpace behavior.
type Options struct {
	AllowQueryFallback bool
}

// Require ensures X-Space-Id is present and stores it in context.
func Require(opts Options) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			spaceID := strings.TrimSpace(c.Request().Header.Get(mwctx.HeaderSpaceID))
			if spaceID == "" && opts.AllowQueryFallback {
				spaceID = strings.TrimSpace(c.QueryParam("space_id"))
			}
			if spaceID == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": mwctx.HeaderSpaceID + " header is required",
				})
			}
			mwctx.SetSpaceID(c, spaceID)
			return next(c)
		}
	}
}