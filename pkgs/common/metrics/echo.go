package metrics

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// EchoMiddleware records Prometheus request counters.
func EchoMiddleware(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			status := c.Response().Status
			if status == 0 {
				status = 200
			}

			RequestsTotal.WithLabelValues(
				serviceName,
				c.Request().Method,
				strconv.Itoa(status),
			).Inc()

			return err
		}
	}
}