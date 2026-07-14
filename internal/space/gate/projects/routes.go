package projectapis

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes wires project CRUD, items, planning, and streaming endpoints.
func RegisterRoutes(tenant *echo.Group, deps Deps, readAction, writeAction echo.MiddlewareFunc) {
	pb := deps.PB
	registerPlanningRoutes(tenant, deps, readAction, writeAction)

	tenant.GET("/projects", listProjectsHandler(pb), readAction)
	tenant.POST("/projects", createProjectHandler(pb), writeAction)
	tenant.GET("/projects/:id", getProjectHandler(pb), readAction)
	tenant.PATCH("/projects/:id", patchProjectHandler(pb), writeAction)
	tenant.DELETE("/projects/:id", deleteProjectHandler(pb), writeAction)

	tenant.GET("/projects/:id/items", listProjectItemsHandler(pb), readAction)
	tenant.POST("/projects/:id/items", createProjectItemHandler(pb), writeAction)
	tenant.PATCH("/projects/:id/items/:itemId", patchProjectItemHandler(pb), writeAction)
	tenant.DELETE("/projects/:id/items/:itemId", deleteProjectItemHandler(pb), writeAction)
}