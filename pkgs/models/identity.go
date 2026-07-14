package models

// IdentityFiles maps PocketPaw-style persona markdown to storage fields.
type IdentityFiles struct {
	IdentityFile      string `json:"identity_file"`
	SoulFile          string `json:"soul_file"`
	StyleFile         string `json:"style_file"`
	InstructionsFile  string `json:"instructions_file"`
	UserFile          string `json:"user_file"`
}

// IdentitySaveResponse reports which identity sections were updated.
type IdentitySaveResponse struct {
	OK      bool     `json:"ok"`
	Updated []string `json:"updated"`
	AgentID string   `json:"agent_id,omitempty"`
}