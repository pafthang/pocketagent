package spaceapis

import (
	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/common"
)

// RegisterRoutes wires space HTTP endpoints.
func RegisterRoutes(e *echo.Echo, deps Deps, rateLimit common.RateLimitConfig, authMiddleware echo.MiddlewareFunc) {
	authGroup := e.Group("")
	if rateLimit.EffectiveEnabled() {
		authGroup.Use(common.AuthRateLimiter(rateLimit))
	}
	authGroup.POST("/auth/register", registerHandler(deps))
	authGroup.POST("/auth/login", loginHandler(deps))
	authGroup.POST("/auth/refresh", refreshHandler(deps))
	authGroup.POST("/auth/verify-email", verifyEmailHandler(deps))

	e.GET("/invites/:token", previewInviteHandler(deps))
	e.POST("/invites/accept", acceptInviteHandler(deps))

	api := e.Group("", authMiddleware)
	api.POST("/auth/request-verification", requestVerificationHandler(deps))
	api.GET("/spaces", listSpacesHandler(deps))
	api.POST("/spaces", createSpaceHandler(deps))
	api.GET("/spaces/:spaceId", getSpaceHandler(deps))
	api.PATCH("/spaces/:spaceId", updateSpaceHandler(deps))
	api.DELETE("/spaces/:spaceId", deleteSpaceHandler(deps))

	api.GET("/spaces/:spaceId/members", listMembersHandler(deps))
	api.POST("/spaces/:spaceId/members", addMemberHandler(deps))
	api.PATCH("/spaces/:spaceId/members/:memberId", updateMemberHandler(deps))
	api.DELETE("/spaces/:spaceId/members/:memberId", deleteMemberHandler(deps))

	api.GET("/spaces/:spaceId/invites", listInvitesHandler(deps))
	api.POST("/spaces/:spaceId/invites", createInviteHandler(deps))
	api.DELETE("/spaces/:spaceId/invites/:inviteId", revokeInviteHandler(deps))

	api.GET("/spaces/:spaceId/audit-logs", listAuditLogsHandler(deps))
	api.GET("/spaces/:spaceId/activity", listActivityHandler(deps))

	api.GET("/spaces/:spaceId/teams", listTeamsHandler(deps))
	api.POST("/spaces/:spaceId/teams", createTeamHandler(deps))
	api.GET("/spaces/:spaceId/teams/:teamId", getTeamHandler(deps))
	api.PATCH("/spaces/:spaceId/teams/:teamId", updateTeamHandler(deps))
	api.DELETE("/spaces/:spaceId/teams/:teamId", deleteTeamHandler(deps))

	api.GET("/spaces/:spaceId/teams/:teamId/members", listTeamMembersHandler(deps))
	api.POST("/spaces/:spaceId/teams/:teamId/members", addTeamMemberHandler(deps))
	api.DELETE("/spaces/:spaceId/teams/:teamId/members/:memberId", deleteTeamMemberHandler(deps))

	api.POST("/authorize", authorizeHandler(deps))
}
