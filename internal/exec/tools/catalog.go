package tools

import "github.com/pafthang/pocketagent/pkgs/ollama"

func builtinCatalog(cfg Config) []ollama.Tool {
	tools := []ollama.Tool{
		{
			Type: "function",
			Function: ollama.ToolFunction{
				Name:        "search_web",
				Description: "Search the web (Serper, Tavily, or DuckDuckGo depending on configuration)",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{"type": "string"},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: ollama.ToolFunction{
				Name:        "scrape_page",
				Description: "Fetch and extract text content from a web page URL",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"url": map[string]interface{}{"type": "string"},
					},
					"required": []string{"url"},
				},
			},
		},
	}

	if cfg.CodeExecEnabled {
		tools = append(tools, ollama.Tool{
			Type: "function",
			Function: ollama.ToolFunction{
				Name:        "code_exec",
				Description: "Execute Python code in a sandboxed subprocess (timeout, temp dir, no network isolation)",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"code": map[string]interface{}{
							"type":        "string",
							"description": "Python source code to execute",
						},
						"language": map[string]interface{}{
							"type":        "string",
							"description": "Language (python)",
						},
					},
					"required": []string{"code"},
				},
			},
		})
	}

	return tools
}
