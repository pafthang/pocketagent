package agentapis

import (
	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

// RegisterRoutes wires agent HTTP endpoints.
func RegisterRoutes(tenant *echo.Group, pb *pbclient.Client, rbac *apimw.PocketRBAC) {
	readAction := rbac.RequireAction(space.ActionAgentRead)
	writeAction := rbac.RequireAction(space.ActionAgentWrite)

	tenant.GET("/agents", ListHandler(pb), readAction)
	tenant.POST("/agents", CreateHandler(pb), writeAction)
	tenant.GET("/agents/:id/runtime-config", getRuntimeConfigHandler(pb), readAction)
	tenant.GET("/agents/:id/identity", getAgentIdentityHandler(pb), readAction)
	tenant.PUT("/agents/:id/identity", putAgentIdentityHandler(pb), writeAction)
	tenant.GET("/agents/:id", GetHandler(pb), readAction)
	tenant.PUT("/agents/:id", UpdateHandler(pb), writeAction)
	tenant.DELETE("/agents/:id", DeleteHandler(pb), writeAction)
}
