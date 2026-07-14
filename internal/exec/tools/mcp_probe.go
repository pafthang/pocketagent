package tools

import (
	"github.com/pafthang/pocketagent/internal/space/probe"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// ProbeTool describes an MCP tool returned by a probe.
type ProbeTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProbeResult is the outcome of probing an MCP server.
type ProbeResult struct {
	Connected bool        `json:"connected"`
	Error     string      `json:"error,omitempty"`
	Tools     []ProbeTool `json:"tools,omitempty"`
}

// ProbeServer validates connectivity and lists tools for a server config.
func ProbeServer(server MCPServerConfig) ProbeResult {
	return probeResultFrom(probe.Server(probe.ServerConfig{
		Name:      server.Name,
		Transport: server.Transport,
		Command:   server.Command,
		Args:      server.Args,
		URL:       server.URL,
		Env:       server.Env,
		Enabled:   server.Enabled,
	}))
}

// ProbeMCPServer validates connectivity and lists tools for a stored MCP server.
func ProbeMCPServer(server models.MCPServer) ProbeResult {
	return probeResultFrom(probe.FromModel(server))
}

func probeResultFrom(r probe.Result) ProbeResult {
	out := ProbeResult{Connected: r.Connected, Error: r.Error}
	if len(r.Tools) > 0 {
		out.Tools = make([]ProbeTool, 0, len(r.Tools))
		for _, tool := range r.Tools {
			out.Tools = append(out.Tools, ProbeTool{Name: tool.Name, Description: tool.Description})
		}
	}
	return out
}
