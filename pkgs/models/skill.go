package models

// Skill is a space-scoped prompt shortcut with optional tool subset.
type Skill struct {
	ID           string   `json:"id"`
	SpaceID      string   `json:"space_id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Prompt       string   `json:"prompt"`
	Category     string   `json:"category,omitempty"`
	Tools        []string `json:"tools,omitempty"`
	ArgumentHint string   `json:"argument_hint,omitempty"`
	CatalogID    string   `json:"catalog_id,omitempty"`
	CreatedAt    string   `json:"created_at,omitempty"`
	UpdatedAt    string   `json:"updated_at,omitempty"`
}