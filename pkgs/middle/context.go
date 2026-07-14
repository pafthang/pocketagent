package middle

import (
	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/middle/context"
	"github.com/pafthang/pocketagent/pkgs/models"
)

const HeaderSpaceID = context.HeaderSpaceID

func UserFromContext(c echo.Context) (models.AuthUser, bool) { return context.UserFromContext(c) }
func SpaceIDFromContext(c echo.Context) (string, bool)       { return context.SpaceIDFromContext(c) }
func AuthTokenFromContext(c echo.Context) string             { return context.AuthTokenFromContext(c) }
func ExtractBearer(c echo.Context) string                    { return context.ExtractBearer(c) }
