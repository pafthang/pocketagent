package common

import (
	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/common/httpmw"
	"github.com/pafthang/pocketagent/pkgs/common/metrics"
)

type RateLimitConfig = httpmw.RateLimitConfig

func SetupMiddleware(e *echo.Echo, serviceName string) { httpmw.Setup(e, serviceName) }
func MetricsMiddleware(serviceName string) echo.MiddlewareFunc {
	return metrics.EchoMiddleware(serviceName)
}
func RateLimiter(perMinute int, burst int) echo.MiddlewareFunc {
	return httpmw.RateLimiter(perMinute, burst)
}
func APIRateLimiter(cfg RateLimitConfig) echo.MiddlewareFunc  { return httpmw.APIRateLimiter(cfg) }
func AuthRateLimiter(cfg RateLimitConfig) echo.MiddlewareFunc { return httpmw.AuthRateLimiter(cfg) }