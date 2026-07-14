package space

import (
	"github.com/pafthang/pocketagent/internal/space/auth"
	"github.com/pafthang/pocketagent/internal/space/rbac"
)

// RBAC action names (re-exported for gate, agent, middleware).
const (
	ActionSpaceRead    = rbac.ActionSpaceRead
	ActionSpaceWrite   = rbac.ActionSpaceWrite
	ActionSpaceDelete  = rbac.ActionSpaceDelete
	ActionMemberRead   = rbac.ActionMemberRead
	ActionMemberWrite  = rbac.ActionMemberWrite
	ActionTeamRead     = rbac.ActionTeamRead
	ActionTeamWrite    = rbac.ActionTeamWrite
	ActionTeamDelete   = rbac.ActionTeamDelete
	ActionAgentRead    = rbac.ActionAgentRead
	ActionAgentWrite   = rbac.ActionAgentWrite
	ActionTaskRead     = rbac.ActionTaskRead
	ActionTaskWrite    = rbac.ActionTaskWrite
	ActionMemoryRead   = rbac.ActionMemoryRead
	ActionMemoryWrite  = rbac.ActionMemoryWrite
	ActionMCPRead      = rbac.ActionMCPRead
	ActionMCPWrite     = rbac.ActionMCPWrite
	ActionSkillRead    = rbac.ActionSkillRead
	ActionSkillWrite   = rbac.ActionSkillWrite
	ActionProjectRead  = rbac.ActionProjectRead
	ActionProjectWrite = rbac.ActionProjectWrite
	ActionFileRead     = rbac.ActionFileRead
	ActionFileWrite    = rbac.ActionFileWrite
	ActionInviteRead   = rbac.ActionInviteRead
	ActionInviteWrite  = rbac.ActionInviteWrite
	ActionAuditRead    = rbac.ActionAuditRead
)

type (
	// Authorizer evaluates space-scoped permissions.
	Authorizer = rbac.Authorizer
	// CachedAuthorizer evaluates RBAC in-process with optional TTL caching.
	CachedAuthorizer = rbac.CachedAuthorizer
)

var (
	NewAuthorizer       = rbac.NewAuthorizer
	NewCachedAuthorizer = rbac.NewCachedAuthorizer
	AuthMiddleware      = auth.AuthMiddleware
)