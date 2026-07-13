package common

import "net/http"

// HealthHandler returns standard health check
func HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": c.Get("service_name").(string),
	})
}
