package rbac

import (
	"fmt"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// Action names for authorization checks.
const (
	ActionSpaceRead    = "space:read"
	ActionSpaceWrite   = "space:write"
	ActionSpaceDelete  = "space:delete"
	ActionMemberRead   = "member:read"
	ActionMemberWrite  = "member:write"
	ActionTeamRead     = "team:read"
	ActionTeamWrite    = "team:write"
	ActionTeamDelete   = "team:delete"
	ActionAgentRead    = "agent:read"
	ActionAgentWrite   = "agent:write"
	ActionTaskRead     = "task:read"
	ActionTaskWrite    = "task:write"
	ActionMemoryRead   = "memory:read"
	ActionMemoryWrite  = "memory:write"
	ActionMCPRead      = "mcp:read"
	ActionMCPWrite     = "mcp:write"
	ActionSkillRead    = "skill:read"
	ActionSkillWrite   = "skill:write"
	ActionProjectRead  = "project:read"
	ActionProjectWrite = "project:write"
	ActionFileRead     = "file:read"
	ActionFileWrite    = "file:write"
	ActionInviteRead   = "invite:read"
	ActionInviteWrite  = "invite:write"
	ActionAuditRead    = "audit:read"
)

// Authorizer evaluates space-scoped permissions.
type Authorizer struct {
	pb *pbclient.Client
}

func NewAuthorizer(pb *pbclient.Client) *Authorizer {
	return &Authorizer{pb: pb}
}

// IsSuperAdmin checks membership in the system admin space.
func (a *Authorizer) IsSuperAdmin(userID string) (bool, error) {
	adminSpace, err := a.pb.GetSpaceBySlug(models.SystemSpaceSlug)
	if err != nil {
		return false, err
	}

	filter := fmt.Sprintf("space_id = %q && user_id = %q && role = %q", adminSpace.ID, userID, models.RoleAdmin)
	members, _, err := a.pb.ListSpaceMembers(pbclient.ListOptions{Page: 1, PerPage: 1, Filter: filter})
	if err != nil {
		return false, err
	}
	return len(members) > 0, nil
}

// MemberRole returns the user's role in a space, or empty if not a member.
func (a *Authorizer) MemberRole(userID, spaceID string) (string, error) {
	if super, err := a.IsSuperAdmin(userID); err != nil {
		return "", err
	} else if super {
		return models.RoleAdmin, nil
	}

	filter := fmt.Sprintf("space_id = %q && user_id = %q", spaceID, userID)
	members, _, err := a.pb.ListSpaceMembers(pbclient.ListOptions{Page: 1, PerPage: 1, Filter: filter})
	if err != nil {
		return "", err
	}
	if len(members) == 0 {
		return "", nil
	}
	return members[0].Role, nil
}

// Authorize checks whether a user may perform an action in a space.
func (a *Authorizer) Authorize(userID, spaceID, action string) (models.AuthorizeResponse, error) {
	role, err := a.MemberRole(userID, spaceID)
	if err != nil {
		return models.AuthorizeResponse{}, err
	}
	if role == "" {
		return models.AuthorizeResponse{Allowed: false, Reason: "not a space member"}, nil
	}

	allowed := roleAllows(role, action)
	resp := models.AuthorizeResponse{
		Allowed: allowed,
		Role:    role,
	}
	if !allowed {
		resp.Reason = fmt.Sprintf("role %q cannot %s", role, action)
	}
	return resp, nil
}

// RoleAllows reports whether role may perform action.
func RoleAllows(role, action string) bool {
	return roleAllows(role, action)
}

func roleAllows(role, action string) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == action {
			return true
		}
	}
	return false
}

var rolePermissions = map[string][]string{
	models.RoleAdmin: {
		ActionSpaceRead, ActionSpaceWrite, ActionSpaceDelete,
		ActionMemberRead, ActionMemberWrite,
		ActionTeamRead, ActionTeamWrite, ActionTeamDelete,
		ActionAgentRead, ActionAgentWrite,
		ActionTaskRead, ActionTaskWrite,
		ActionMemoryRead, ActionMemoryWrite,
		ActionMCPRead, ActionMCPWrite,
		ActionSkillRead, ActionSkillWrite,
		ActionProjectRead, ActionProjectWrite,
		ActionFileRead, ActionFileWrite,
		ActionInviteRead, ActionInviteWrite,
		ActionAuditRead,
	},
	models.RoleEditor: {
		ActionSpaceRead,
		ActionMemberRead,
		ActionTeamRead, ActionTeamWrite,
		ActionAgentRead, ActionAgentWrite,
		ActionTaskRead, ActionTaskWrite,
		ActionMemoryRead, ActionMemoryWrite,
		ActionMCPRead, ActionMCPWrite,
		ActionSkillRead, ActionSkillWrite,
		ActionProjectRead, ActionProjectWrite,
		ActionFileRead, ActionFileWrite,
	},
	models.RoleViewer: {
		ActionSpaceRead,
		ActionMemberRead,
		ActionTeamRead,
		ActionAgentRead,
		ActionTaskRead,
		ActionMemoryRead,
		ActionMCPRead,
		ActionSkillRead,
		ActionProjectRead,
		ActionFileRead,
	},
}

// RequireRole returns an error if the role cannot perform the action.
func RequireRole(role, action string) error {
	if role == "" {
		return httpx.ErrForbidden("not a space member")
	}
	if !roleAllows(role, action) {
		return httpx.ErrForbidden(fmt.Sprintf("role %q cannot %s", role, action))
	}
	return nil
}
