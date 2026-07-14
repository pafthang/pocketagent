package common

import (
	"github.com/labstack/echo/v4"
	hp "github.com/pafthang/pocketagent/pkgs/common/health"
)

type HealthStatus = hp.Status
type HealthResponse = hp.Response
type Deps = hp.Deps

func HealthHandler(deps Deps) echo.HandlerFunc { return hp.Handler(deps) }

func PocketHealthURL(base string) string   { return hp.PocketURL(base) }
func OllamaHealthURL(base string) string   { return hp.OllamaURL(base) }
func MemoHealthURL(base string) string     { return hp.MemoURL(base) }
func ServiceHealthURL(base string) string  { return hp.ServiceURL(base) }
func NATSMonitoringURL(base string) string { return hp.NATSMonitoringURL(base) }