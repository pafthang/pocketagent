package models

// Space roles within a tenant.
const (
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

// System space slug for super-admins.
const SystemSpaceSlug = "admin"

// Team member kinds.
const (
	MemberTypeUser  = "user"
	MemberTypeAgent = "agent"
)

// Space is a tenant boundary for agents, teams, and tasks.
type Space struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	IsSystem    bool   `json:"is_system"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// SpaceMember links a user to a space with a role.
type SpaceMember struct {
	ID        string `json:"id"`
	SpaceID   string `json:"space_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Team groups users and agents inside a space.
type Team struct {
	ID          string `json:"id"`
	SpaceID     string `json:"space_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// TeamMember links a user or agent to a team.
type TeamMember struct {
	ID         string `json:"id"`
	TeamID     string `json:"team_id"`
	MemberType string `json:"member_type"`
	MemberID   string `json:"member_id"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

// Invite statuses.
const (
	InvitePending  = "pending"
	InviteAccepted = "accepted"
	InviteRevoked  = "revoked"
	InviteExpired  = "expired"
)

// Email verification statuses.
const (
	VerificationPending = "pending"
	VerificationDone    = "verified"
)

// AuthUser is a minimal PocketBase user identity.
type AuthUser struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

// AuthSession is returned after login or refresh.
type AuthSession struct {
	Token string   `json:"token"`
	User  AuthUser `json:"user"`
}

// AuthorizeRequest checks whether an action is allowed in a space.
type AuthorizeRequest struct {
	SpaceID      string `json:"space_id"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type,omitempty"`
	ResourceID   string `json:"resource_id,omitempty"`
}

// AuthorizeResponse is the result of a policy check.
type AuthorizeResponse struct {
	Allowed bool   `json:"allowed"`
	Role    string `json:"role,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

// SpaceInvite invites a user to join a space by email.
type SpaceInvite struct {
	ID        string `json:"id"`
	SpaceID   string `json:"space_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	InvitedBy string `json:"invited_by,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// InvitePreview is returned for a valid invite token (public).
type InvitePreview struct {
	SpaceID   string `json:"space_id"`
	SpaceName string `json:"space_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

// AuditLog records a space-scoped action.
type AuditLog struct {
	ID           string                 `json:"id"`
	SpaceID      string                 `json:"space_id"`
	ActorID      string                 `json:"actor_id,omitempty"`
	ActorEmail   string                 `json:"actor_email,omitempty"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type,omitempty"`
	ResourceID   string                 `json:"resource_id,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	IPAddress    string                 `json:"ip_address,omitempty"`
	CreatedAt    string                 `json:"created_at,omitempty"`
}