package mcpapis

import (
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestPatchMCPServerRequestApplyPatch(t *testing.T) {
	name := "new-name"
	enabled := false
	server := models.MCPServer{Name: "old", Enabled: true}
	PatchMCPServerRequest{Name: &name, Enabled: &enabled}.ApplyPatch(&server)
	if server.Name != "new-name" || server.Enabled {
		t.Fatalf("server = %+v", server)
	}
}
