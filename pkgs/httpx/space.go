package httpx

import (
	"github.com/labstack/echo/v4"
	mwctx "github.com/pafthang/pocketagent/pkgs/middle/context"
)

// RequireSpaceID returns the tenant space ID set by RequireSpace middleware.
func RequireSpaceID(c echo.Context) (string, bool) {
	return mwctx.SpaceIDFromContext(c)
}