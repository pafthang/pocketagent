package common

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/common/telemetry"
)

func InitTelemetry(serviceName string)             { telemetry.Init(serviceName) }
func TelemetryEnabled() bool                       { return telemetry.Enabled() }
func TelemetryMiddleware(serviceName string) echo.MiddlewareFunc {
	return telemetry.EchoMiddleware(serviceName)
}
func ShutdownTelemetry(ctx context.Context) error { return telemetry.Shutdown(ctx) }