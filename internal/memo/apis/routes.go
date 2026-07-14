package memoapis

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes wires tenant-scoped memory endpoints exposed via gate.
func RegisterRoutes(tenant *echo.Group, deps Deps, embedModel, ollamaURL string, readAction, writeAction echo.MiddlewareFunc) {
	tenant.GET("/memory/stats", memoryStatsHandler(deps), readAction)
	tenant.GET("/memory/settings", memorySettingsHandler(embedModel, ollamaURL), readAction)
	tenant.POST("/memory/settings", saveMemorySettingsHandler(), writeAction)
	tenant.POST("/memory/search", searchMemoryHandler(deps), readAction)
	tenant.GET("/memory", listMemoryHandler(deps), readAction)
	tenant.POST("/memory", ingestMemoryHandler(deps), writeAction)
	tenant.GET("/memory/:id", getMemoryHandler(deps), readAction)
	tenant.DELETE("/memory/:id", deleteMemoryHandler(deps), writeAction)
}