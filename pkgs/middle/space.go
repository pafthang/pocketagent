package middle

import (
	"github.com/labstack/echo/v4"
	mwspace "github.com/pafthang/pocketagent/pkgs/middle/space"
)

type SpaceOptions = mwspace.Options

func RequireSpace(opts SpaceOptions) echo.MiddlewareFunc { return mwspace.Require(opts) }