package ollama

// Tool definition for LLM tool calling
type Tool struct {
	Type      string                 `json:"type"`
	Function  ToolFunction           `json:"function"`
}

type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Example tools
func GetExampleTools() []Tool {
	return []Tool{
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "search_web",
				Description: "Search the web for information",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"query"},
				},
			},
		},
	}
}
