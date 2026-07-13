package models

type Agent struct {
	ID          string                 `json:"id" db:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Model       string                 `json:"model"` // ollama model
	SystemPrompt string                `json:"system_prompt"`
	Tools       []string               `json:"tools"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}
