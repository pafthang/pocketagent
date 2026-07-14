package models

// SpaceProfile stores per-user context injected into agent prompts.
type SpaceProfile struct {
	ID        string `json:"id"`
	SpaceID   string `json:"space_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}