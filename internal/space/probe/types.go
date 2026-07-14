package probe

// ServerConfig defines an MCP server connection for probing.
type ServerConfig struct {
	Name      string            `json:"name"`
	Transport string            `json:"transport"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	URL       string            `json:"url"`
	Env       map[string]string `json:"env"`
	Enabled   bool              `json:"enabled"`
}

// Tool describes an MCP tool returned by a probe.
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Result is the outcome of probing an MCP server.
type Result struct {
	Connected bool   `json:"connected"`
	Error     string `json:"error,omitempty"`
	Tools     []Tool `json:"tools,omitempty"`
}