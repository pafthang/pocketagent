package tools

// ToolInfo describes a tool exposed to agents.
type ToolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Server      string `json:"server,omitempty"`
}

// BuiltinToolInfos returns static builtin tool metadata for a config.
func BuiltinToolInfos(cfg Config) []ToolInfo {
	catalog := builtinCatalog(cfg)
	out := make([]ToolInfo, 0, len(catalog))
	for _, tool := range catalog {
		out = append(out, ToolInfo{
			Name:        tool.Function.Name,
			Description: tool.Function.Description,
			Source:      "builtin",
		})
	}
	return out
}