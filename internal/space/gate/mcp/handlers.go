package mcpapis

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space/probe"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/pkgs/models"
)

type mcpServerStatus struct {
	Connected  bool   `json:"connected"`
	Connecting bool   `json:"connecting,omitempty"`
	ToolCount  int    `json:"tool_count"`
	Error      string `json:"error"`
	Transport  string `json:"transport"`
	Enabled    bool   `json:"enabled"`
}

func RegisterRoutes(tenant *echo.Group, pb *pbclient.Client, readAction, writeAction echo.MiddlewareFunc) {
	tenant.GET("/mcp/servers", listMCPServersHandler(pb), readAction)
	tenant.POST("/mcp/servers", createMCPServerHandler(pb), writeAction)
	tenant.GET("/mcp/servers/:id", getMCPServerHandler(pb), readAction)
	tenant.PATCH("/mcp/servers/:id", patchMCPServerHandler(pb), writeAction)
	tenant.DELETE("/mcp/servers/:id", deleteMCPServerHandler(pb), writeAction)
	tenant.POST("/mcp/servers/:id/test", testMCPServerByIDHandler(pb), readAction)
	tenant.GET("/mcp/status", mcpStatusHandler(pb), readAction)
	tenant.GET("/mcp/presets", listMCPPresentsHandler(pb), readAction)
	tenant.POST("/mcp/presets/install", installMCPPresetHandler(pb), writeAction)
}

func listMCPServersHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		page, _ := strconv.Atoi(c.QueryParam("page"))
		perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
		servers, total, err := pb.ListMCPServers(pbclient.ListOptions{
			Page:    page,
			PerPage: perPage,
			Filter:  fmt.Sprintf("space_id = %q", spaceID),
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"servers": servers,
			"total":   total,
		})
	}
}

func createMCPServerHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		var req MCPServerInput
		if err := c.Bind(&req); err != nil {
			return err
		}

		server, err := buildMCPServerFromRequest(spaceID, req.Name, req.Transport, req.Command, req.Args, req.URL, req.Env, req.Enabled, true)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		if err := ensureUniqueMCPServerName(pb, spaceID, server.Name, ""); err != nil {
			return httpx.MapPocketError(c, err)
		}

		stored, err := pb.CreateMCPServer(server)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusCreated, stored)
	}
}

func getMCPServerHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		server, err := loadMCPServerInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, server)
	}
}

func patchMCPServerHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		existing, err := loadMCPServerInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		var req PatchMCPServerRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		updated := existing
		req.ApplyPatch(&updated)

		if updated.Name != existing.Name {
			if err := ensureUniqueMCPServerName(pb, spaceID, updated.Name, existing.ID); err != nil {
				return httpx.MapPocketError(c, err)
			}
		}

		stored, err := pb.UpdateMCPServerRecord(updated)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, stored)
	}
}

func deleteMCPServerHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		server, err := loadMCPServerInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if err := pb.DeleteMCPServer(server.ID); err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func testMCPServerByIDHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		server, err := loadMCPServerInSpace(pb, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, ProbeMCPServer(server))
	}
}

func listMCPPresentsHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		servers, _, err := pb.ListMCPServers(pbclient.ListOptions{
			Page:    1,
			PerPage: 200,
			Filter:  fmt.Sprintf("space_id = %q", spaceID),
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		names := make([]string, 0, len(servers))
		for _, server := range servers {
			names = append(names, server.Name)
		}

		presets, err := loadMCPPresents(installedServerNames(names))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, presets)
	}
}

func mcpStatusHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		servers, _, err := pb.ListMCPServers(pbclient.ListOptions{
			Page:    1,
			PerPage: 200,
			Filter:  fmt.Sprintf("space_id = %q", spaceID),
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		status := make(map[string]mcpServerStatus, len(servers))
		for _, server := range servers {
			entry := mcpServerStatus{
				Transport: normalizeTransport(server.Transport),
				Enabled:   server.Enabled,
			}
			if server.Enabled {
				probe := ProbeMCPServer(server)
				entry.Connected = probe.Connected
				entry.ToolCount = len(probe.Tools)
				entry.Error = probe.Error
			}
			status[server.Name] = entry
		}
		return c.JSON(http.StatusOK, status)
	}
}

func installMCPPresetHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		var req InstallMCPPresetRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		if strings.TrimSpace(req.PresetID) == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "preset_id is required"})
		}

		presets, err := loadMCPPresents(nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		var preset *mcpPresetResponse
		for i := range presets {
			if presets[i].ID == req.PresetID {
				preset = &presets[i]
				break
			}
		}
		if preset == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "preset not found"})
		}

		args := []string{"-y", preset.Package}
		if len(req.ExtraArgs) > 0 {
			args = append(args, req.ExtraArgs...)
		} else if preset.NeedsArgs {
			args = append(args, ".")
		}

		server, err := buildMCPServerFromRequest(
			spaceID,
			preset.Name,
			preset.Transport,
			"npx",
			args,
			preset.URL,
			req.Env,
			boolPtr(true),
			true,
		)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		if err := ensureUniqueMCPServerName(pb, spaceID, server.Name, ""); err != nil {
			return httpx.MapPocketError(c, err)
		}

		stored, err := pb.CreateMCPServer(server)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		probe := ProbeMCPServer(stored)
		resp := map[string]interface{}{"status": "installed"}
		if probe.Connected {
			resp["connected"] = true
		}
		if probe.Error != "" {
			resp["error"] = probe.Error
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func ProbeMCPServer(server models.MCPServer) probe.Result {
	return probe.FromModel(server)
}

func buildMCPServerFromRequest(spaceID, name, transport, command string, args []string, url string, env map[string]string, enabled *bool, defaultEnabled bool) (models.MCPServer, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.MCPServer{}, fmt.Errorf("name is required")
	}

	transport = normalizeTransport(transport)
	switch transport {
	case "stdio":
		if strings.TrimSpace(command) == "" {
			return models.MCPServer{}, fmt.Errorf("command is required for stdio transport")
		}
	case "http":
		if strings.TrimSpace(url) == "" {
			return models.MCPServer{}, fmt.Errorf("url is required for http transport")
		}
	default:
		return models.MCPServer{}, fmt.Errorf("transport must be stdio or http")
	}

	isEnabled := defaultEnabled
	if enabled != nil {
		isEnabled = *enabled
	}

	return models.MCPServer{
		SpaceID:   spaceID,
		Name:      name,
		Transport: transport,
		Command:   strings.TrimSpace(command),
		Args:      args,
		URL:       strings.TrimSpace(url),
		Env:       env,
		Enabled:   isEnabled,
	}, nil
}

func normalizeTransport(transport string) string {
	transport = strings.ToLower(strings.TrimSpace(transport))
	if transport == "" {
		return "stdio"
	}
	return transport
}

func ensureUniqueMCPServerName(pb *pbclient.Client, spaceID, name, excludeID string) error {
	existing, err := findMCPServerByName(pb, spaceID, name)
	if err != nil {
		var apiErr *pbclient.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}
	if excludeID != "" && existing.ID == excludeID {
		return nil
	}
	return fmt.Errorf("mcp server %q already exists in this space", name)
}

func findMCPServerByName(pb *pbclient.Client, spaceID, name string) (models.MCPServer, error) {
	name = strings.TrimSpace(name)
	servers, _, err := pb.ListMCPServers(pbclient.ListOptions{
		Page:    1,
		PerPage: 1,
		Filter:  fmt.Sprintf("space_id = %q && name = %q", spaceID, name),
	})
	if err != nil {
		return models.MCPServer{}, err
	}
	if len(servers) == 0 {
		return models.MCPServer{}, &pbclient.APIError{StatusCode: http.StatusNotFound, Message: "mcp server not found"}
	}
	return servers[0], nil
}

func loadMCPServerInSpace(pb *pbclient.Client, spaceID, id string) (models.MCPServer, error) {
	server, err := pb.GetMCPServer(id)
	if err != nil {
		return models.MCPServer{}, err
	}
	if server.SpaceID != spaceID {
		return models.MCPServer{}, &pbclient.APIError{StatusCode: http.StatusNotFound, Message: "mcp server not found"}
	}
	return server, nil
}

func boolPtr(v bool) *bool {
	return &v
}
