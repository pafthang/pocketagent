package httpmw

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pafthang/pocketagent/pkgs/common/metrics"
	"github.com/pafthang/pocketagent/pkgs/common/telemetry"
)

// Setup adds common middleware to Echo.
func Setup(e *echo.Echo, serviceName string) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(telemetry.EchoMiddleware(serviceName))
	e.Use(metrics.EchoMiddleware(serviceName))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("service_name", serviceName)
			return next(c)
		}
	})
}