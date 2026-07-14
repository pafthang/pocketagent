package rbac

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space"
	mwctx "github.com/pafthang/pocketagent/pkgs/middle/context"
)

// PocketRBAC validates sessions and permissions in-process via PocketBase.
// Use on gate/agent tenant routes to avoid HTTP hops to the space service.
type PocketRBAC struct {
	pb   *pbclient.Client
	auth *space.CachedAuthorizer
}

// New wires direct PocketBase auth refresh and space RBAC checks.
func New(pb *pbclient.Client, authorizeCacheTTL time.Duration) *PocketRBAC {
	return &PocketRBAC{
		pb:   pb,
		auth: space.NewCachedAuthorizer(pb, authorizeCacheTTL),
	}
}

// AuthMiddleware validates JWTs via PocketBase auth-refresh (no space HTTP hop).
func (r *PocketRBAC) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := mwctx.ExtractBearer(c)
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authorization required"})
			}

			session, err := r.pb.AuthRefresh(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}

			mwctx.SetUser(c, session.User)
			mwctx.SetToken(c, token)
			return next(c)
		}
	}
}

// RequireAction checks RBAC in-process using the shared space Authorizer.
func (r *PocketRBAC) RequireAction(action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			spaceID, ok := mwctx.SpaceIDFromContext(c)
			if !ok {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": mwctx.HeaderSpaceID + " header is required",
				})
			}

			user, ok := mwctx.UserFromContext(c)
			if !ok || user.ID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			resp, err := r.auth.Authorize(user.ID, spaceID, action)
			if err != nil {
				return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
			}
			if !resp.Allowed {
				msg := resp.Reason
				if msg == "" {
					msg = "forbidden"
				}
				return c.JSON(http.StatusForbidden, map[string]string{"error": msg})
			}

			return next(c)
		}
	}
}