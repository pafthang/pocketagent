package probe

import (
	"strings"

	"github.com/pafthang/pocketagent/pkgs/models"
)

// Server validates connectivity and lists tools for a server config.
func Server(cfg ServerConfig) Result {
	transport := strings.ToLower(strings.TrimSpace(cfg.Transport))
	if transport == "" {
		transport = "stdio"
	}

	var client connector
	var err error

	switch transport {
	case "http":
		client, err = newHTTPClient(cfg)
	case "stdio":
		client, err = newStdioClient(cfg)
	default:
		return Result{Error: "unsupported transport: " + transport}
	}
	if err != nil {
		return Result{Error: err.Error()}
	}
	defer client.close()

	tools, err := client.listTools()
	if err != nil {
		return Result{Error: err.Error()}
	}

	out := make([]Tool, 0, len(tools))
	for _, tool := range tools {
		out = append(out, Tool{Name: tool.Name, Description: tool.Description})
	}
	return Result{Connected: true, Tools: out}
}

// FromModel probes a stored MCP server record.
func FromModel(server models.MCPServer) Result {
	return Server(ServerConfig{
		Name:      server.Name,
		Transport: server.Transport,
		Command:   server.Command,
		Args:      server.Args,
		URL:       server.URL,
		Env:       server.Env,
		Enabled:   true,
	})
}
