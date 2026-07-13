package common

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
)

// HealthStatus represents dependency status
type HealthStatus struct {
	Status  string `json:"status"`
	Latency string `json:"latency,omitempty"`
}

// HealthResponse full health check response
type HealthResponse struct {
	Status       string                 `json:"status"`
	Dependencies map[string]HealthStatus `json:"dependencies"`
}

// CheckDependencies performs health checks on critical services
func CheckDependencies(natsConn *nats.Conn, pocketbaseURL, ollamaURL string) HealthResponse {
	deps := make(map[string]HealthStatus)
	overall := "healthy"

	// NATS
	start := time.Now()
	if natsConn != nil && natsConn.IsConnected() {
		deps["nats"] = HealthStatus{Status: "up", Latency: time.Since(start).String()}
	} else {
		deps["nats"] = HealthStatus{Status: "down"}
		overall = "degraded"
	}

	// TODO: add PocketBase and Ollama checks

	return HealthResponse{
		Status:       overall,
		Dependencies: deps,
	}
}

// HealthHandler returns enhanced health check
func HealthHandler(c echo.Context) error {
	// In real usage pass actual connections
	resp := CheckDependencies(nil, "", "")
	return c.JSON(http.StatusOK, resp)
}
