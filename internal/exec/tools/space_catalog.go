package tools

import (
	"fmt"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
)

// CollectSpaceTools returns builtin and MCP tools available in a space.
func CollectSpaceTools(pb *pbclient.Client, spaceID string, toolCfg Config) ([]ToolInfo, error) {
	result := BuiltinToolInfos(toolCfg)

	servers, _, err := pb.ListMCPServers(pbclient.ListOptions{
		Page:    1,
		PerPage: 100,
		Filter:  fmt.Sprintf("space_id = %q && enabled = true", spaceID),
	})
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		probe := ProbeMCPServer(server)
		if !probe.Connected {
			continue
		}
		for _, tool := range probe.Tools {
			result = append(result, ToolInfo{
				Name:        PublicToolName(server.Name, tool.Name),
				Description: fmt.Sprintf("MCP/%s: %s", server.Name, tool.Description),
				Source:      "mcp",
				Server:      server.Name,
			})
		}
	}
	return result, nil
}

// ValidateToolAllowList ensures requested tools exist in the space catalog.
func ValidateToolAllowList(pb *pbclient.Client, spaceID string, toolCfg Config, allowList []string) error {
	if len(allowList) == 0 {
		return nil
	}
	available, err := availableToolNames(pb, spaceID, toolCfg)
	if err != nil {
		return err
	}
	for _, name := range allowList {
		if _, ok := available[name]; !ok {
			return fmt.Errorf("tool %q is not available in this space", name)
		}
	}
	return nil
}

func availableToolNames(pb *pbclient.Client, spaceID string, toolCfg Config) (map[string]struct{}, error) {
	infos, err := CollectSpaceTools(pb, spaceID, toolCfg)
	if err != nil {
		return nil, err
	}
	names := make(map[string]struct{}, len(infos))
	for _, info := range infos {
		names[info.Name] = struct{}{}
	}
	return names, nil
}
