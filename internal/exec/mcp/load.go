package mcp

import (
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/internal/exec/tools"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// LoadSpaceServers returns enabled MCP servers for a space.
func LoadSpaceServers(pb *pbclient.Client, spaceID string) ([]tools.MCPServerConfig, error) {
	if pb == nil || strings.TrimSpace(spaceID) == "" {
		return nil, nil
	}

	servers, _, err := pb.ListMCPServers(pbclient.ListOptions{
		Page:    1,
		PerPage: 100,
		Filter:  fmt.Sprintf("space_id = %q && enabled = true", spaceID),
	})
	if err != nil {
		return nil, err
	}

	out := make([]tools.MCPServerConfig, 0, len(servers))
	for _, server := range servers {
		out = append(out, serverToConfig(server))
	}
	return out, nil
}

func serverToConfig(server models.MCPServer) tools.MCPServerConfig {
	transport := strings.ToLower(strings.TrimSpace(server.Transport))
	if transport == "" {
		transport = "stdio"
	}
	return tools.MCPServerConfig{
		Name:      server.Name,
		Transport: transport,
		Command:   server.Command,
		Args:      server.Args,
		URL:       server.URL,
		Env:       server.Env,
		Enabled:   server.Enabled,
	}
}