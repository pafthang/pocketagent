package mcpapis

import (
	"strings"

	"github.com/pafthang/pocketagent/pkgs/models"
)

// MCPServerInput is the body for creating MCP servers.
type MCPServerInput struct {
	Name      string            `json:"name"`
	Transport string            `json:"transport"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	URL       string            `json:"url"`
	Env       map[string]string `json:"env"`
	Enabled   *bool             `json:"enabled"`
}

// PatchMCPServerRequest is the gate API body for PATCH /mcp/servers/:id.
type PatchMCPServerRequest struct {
	Name      *string            `json:"name"`
	Transport *string            `json:"transport"`
	Command   *string            `json:"command"`
	Args      *[]string          `json:"args"`
	URL       *string            `json:"url"`
	Env       *map[string]string `json:"env"`
	Enabled   *bool              `json:"enabled"`
}

// ApplyPatch mutates an MCP server record in place.
func (r PatchMCPServerRequest) ApplyPatch(server *models.MCPServer) {
	if r.Name != nil {
		server.Name = strings.TrimSpace(*r.Name)
	}
	if r.Transport != nil {
		server.Transport = strings.TrimSpace(*r.Transport)
	}
	if r.Command != nil {
		server.Command = strings.TrimSpace(*r.Command)
	}
	if r.Args != nil {
		server.Args = *r.Args
	}
	if r.URL != nil {
		server.URL = strings.TrimSpace(*r.URL)
	}
	if r.Env != nil {
		server.Env = *r.Env
	}
	if r.Enabled != nil {
		server.Enabled = *r.Enabled
	}
}

// InstallMCPPresetRequest is the body for POST /mcp/presets/install.
type InstallMCPPresetRequest struct {
	PresetID  string            `json:"preset_id"`
	Env       map[string]string `json:"env"`
	ExtraArgs []string          `json:"extra_args"`
}
