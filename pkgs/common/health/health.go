package health

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
)

const checkTimeout = 2 * time.Second

// Status represents dependency status.
type Status struct {
	Status  string `json:"status"`
	Latency string `json:"latency,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Response is the full health check response.
type Response struct {
	Service      string            `json:"service,omitempty"`
	Status       string            `json:"status"`
	Dependencies map[string]Status `json:"dependencies"`
}

// Deps lists optional dependencies to probe for /health.
type Deps struct {
	Service       string
	NATS          *nats.Conn
	JetStream     nats.JetStreamContext
	PocketBaseURL string
	SpaceURL      string
	AgentURL      string
	FilesURL      string
	OllamaURL     string
	MemoURL       string
	MemoStore     func() error
	DLQWarnCount  uint64
}

// Handler returns a handler that checks configured dependencies.
func Handler(deps Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := deps.Check(c.Request().Context())

		if name, ok := c.Get("service_name").(string); ok && name != "" {
			resp.Service = name
		} else if deps.Service != "" {
			resp.Service = deps.Service
		}

		code := http.StatusOK
		if resp.Status != "healthy" {
			code = http.StatusServiceUnavailable
		}

		return c.JSON(code, resp)
	}
}

// Check runs all configured dependency probes.
func (d Deps) Check(ctx context.Context) Response {
	deps := make(map[string]Status)
	overall := "healthy"

	run := func(name string, status Status) {
		deps[name] = status
		if status.Status != "up" {
			overall = "degraded"
		}
	}

	if d.NATS != nil {
		run("nats", checkNATSConn(d.NATS))
	}
	if d.PocketBaseURL != "" {
		run("pocketbase", checkHTTP(ctx, PocketURL(d.PocketBaseURL)))
	}
	if d.OllamaURL != "" {
		run("ollama", checkHTTP(ctx, OllamaURL(d.OllamaURL)))
	}
	if d.SpaceURL != "" {
		run("space", checkHTTP(ctx, ServiceURL(d.SpaceURL)))
	}
	if d.AgentURL != "" {
		run("agent", checkHTTP(ctx, ServiceURL(d.AgentURL)))
	}
	if d.FilesURL != "" {
		run("files", checkHTTP(ctx, ServiceURL(d.FilesURL)))
	}
	if d.MemoURL != "" {
		run("memo", checkHTTP(ctx, MemoURL(d.MemoURL)))
	}
	if d.MemoStore != nil {
		run("chromem", checkMemoStore(d.MemoStore))
	}
	if d.JetStream != nil {
		run("dlq", checkDLQ(d.JetStream, d.DLQWarnCount))
	}
	if len(deps) == 0 {
		overall = "healthy"
	}

	return Response{
		Service:      d.Service,
		Status:       overall,
		Dependencies: deps,
	}
}

func checkNATSConn(nc *nats.Conn) Status {
	start := time.Now()
	if nc != nil && nc.IsConnected() {
		return Status{Status: "up", Latency: time.Since(start).String()}
	}
	return Status{Status: "down", Error: "not connected"}
}

func checkMemoStore(ping func() error) Status {
	start := time.Now()
	if err := ping(); err != nil {
		return Status{Status: "down", Error: err.Error()}
	}
	return Status{Status: "up", Latency: time.Since(start).String()}
}

func checkDLQ(js nats.JetStreamContext, warnCount uint64) Status {
	start := time.Now()
	stats, err := dlqStreamDepth(js)
	if err != nil {
		return Status{Status: "down", Error: err.Error()}
	}
	status := Status{
		Status:  "up",
		Latency: time.Since(start).String(),
	}
	if warnCount > 0 && stats.Messages >= warnCount {
		status.Status = "down"
		status.Error = fmt.Sprintf("%d dead-letter messages (threshold %d)", stats.Messages, warnCount)
	}
	return status
}

func checkHTTP(ctx context.Context, url string) Status {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Status{Status: "down", Error: err.Error()}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Status{Status: "down", Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return Status{Status: "up", Latency: time.Since(start).String()}
	}

	return Status{
		Status: "down",
		Error:  fmt.Sprintf("HTTP %d", resp.StatusCode),
	}
}

// PocketURL returns the PocketBase health probe URL.
func PocketURL(base string) string {
	return strings.TrimRight(base, "/") + "/api/health"
}

// OllamaURL returns the Ollama health probe URL.
func OllamaURL(base string) string {
	return strings.TrimRight(base, "/") + "/api/tags"
}

// MemoURL returns the memo service health probe URL.
func MemoURL(base string) string {
	return ServiceURL(base)
}

// ServiceURL returns the standard /health probe URL for HTTP services.
func ServiceURL(base string) string {
	return strings.TrimRight(base, "/") + "/health"
}

// NATSMonitoringURL returns the NATS server monitoring health probe URL.
func NATSMonitoringURL(base string) string {
	return strings.TrimRight(base, "/") + "/healthz"
}