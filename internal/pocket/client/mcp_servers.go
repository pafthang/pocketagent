package client

import (
	"fmt"

	"github.com/pafthang/pocketagent/pkgs/models"
)



// CreateMCPServer stores a new MCP server record.
func (c *Client) CreateMCPServer(server models.MCPServer) (models.MCPServer, error) {
	record, err := c.CreateRecord(MCPServersCollection, mcpServerRecordData(server))
	if err != nil {
		return models.MCPServer{}, err
	}
	return mcpServerFromRecord(record), nil
}

// GetMCPServer returns an MCP server by ID.
func (c *Client) GetMCPServer(id string) (models.MCPServer, error) {
	record, err := c.GetRecord(MCPServersCollection, id)
	if err != nil {
		return models.MCPServer{}, err
	}
	return mcpServerFromRecord(record), nil
}

// ListMCPServers returns MCP servers with optional filter.
func (c *Client) ListMCPServers(opts ListOptions) ([]models.MCPServer, int, error) {
	records, total, err := c.ListRecordsOpts(MCPServersCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.MCPServer, 0, len(records))
	for _, record := range records {
		out = append(out, mcpServerFromRecord(record))
	}
	return out, total, nil
}

// UpdateMCPServer patches an MCP server by record ID.
func (c *Client) UpdateMCPServer(id string, patch models.MCPServer) (models.MCPServer, error) {
	existing, err := c.GetMCPServer(id)
	if err != nil {
		return models.MCPServer{}, err
	}
	record, err := c.UpdateRecord(MCPServersCollection, id, mcpServerRecordData(mergeMCPServerUpdate(existing, patch)))
	if err != nil {
		return models.MCPServer{}, err
	}
	return mcpServerFromRecord(record), nil
}

// UpdateMCPServerRecord replaces an MCP server record.
func (c *Client) UpdateMCPServerRecord(server models.MCPServer) (models.MCPServer, error) {
	record, err := c.UpdateRecord(MCPServersCollection, server.ID, mcpServerRecordData(server))
	if err != nil {
		return models.MCPServer{}, err
	}
	return mcpServerFromRecord(record), nil
}

// DeleteMCPServer removes an MCP server by ID.
func (c *Client) DeleteMCPServer(id string) error {
	return c.DeleteRecord(MCPServersCollection, id)
}

func mergeMCPServerUpdate(existing, patch models.MCPServer) models.MCPServer {
	merged := existing
	if patch.Name != "" {
		merged.Name = patch.Name
	}
	if patch.Transport != "" {
		merged.Transport = patch.Transport
	}
	if patch.Command != "" {
		merged.Command = patch.Command
	}
	if patch.Args != nil {
		merged.Args = patch.Args
	}
	if patch.URL != "" {
		merged.URL = patch.URL
	}
	if patch.Env != nil {
		merged.Env = patch.Env
	}
	return merged
}

func mcpServerRecordData(server models.MCPServer) map[string]interface{} {
	data := map[string]interface{}{
		"space_id":  server.SpaceID,
		"name":      server.Name,
		"transport": server.Transport,
		"enabled":   server.Enabled,
	}
	if server.Command != "" {
		data["command"] = server.Command
	}
	if server.Args != nil {
		data["args"] = server.Args
	}
	if server.URL != "" {
		data["url"] = server.URL
	}
	if server.Env != nil {
		data["env"] = server.Env
	}
	return data
}

func mcpServerFromRecord(record map[string]interface{}) models.MCPServer {
	server := models.MCPServer{
		ID:        stringField(record, "id"),
		SpaceID:   stringField(record, "space_id"),
		Name:      stringField(record, "name"),
		Transport: stringField(record, "transport"),
		Command:   stringField(record, "command"),
		URL:       stringField(record, "url"),
		Enabled:   boolField(record, "enabled"),
		CreatedAt: stringField(record, "created"),
		UpdatedAt: stringField(record, "updated"),
	}
	server.Args = stringSliceField(record, "args")
	server.Env = stringMapField(record, "env")
	return server
}

func stringMapField(record map[string]interface{}, key string) map[string]string {
	raw, ok := record[key].(map[string]interface{})
	if !ok {
		return nil
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		out[k] = fmt.Sprint(v)
	}
	return out
}
