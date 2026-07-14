package rbac

import (
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestMemoryPermissionsByRole(t *testing.T) {
	cases := []struct {
		role    string
		action  string
		allowed bool
	}{
		{models.RoleAdmin, ActionMemoryRead, true},
		{models.RoleAdmin, ActionMemoryWrite, true},
		{models.RoleEditor, ActionMemoryRead, true},
		{models.RoleEditor, ActionMemoryWrite, true},
		{models.RoleViewer, ActionMemoryRead, true},
		{models.RoleViewer, ActionMemoryWrite, false},
		{models.RoleEditor, ActionMCPWrite, true},
		{models.RoleViewer, ActionMCPRead, true},
		{models.RoleViewer, ActionMCPWrite, false},
		{models.RoleEditor, ActionSkillWrite, true},
		{models.RoleViewer, ActionSkillRead, true},
		{models.RoleViewer, ActionSkillWrite, false},
	}

	for _, tc := range cases {
		got := roleAllows(tc.role, tc.action)
		if got != tc.allowed {
			t.Fatalf("role=%s action=%s got=%v want=%v", tc.role, tc.action, got, tc.allowed)
		}
	}
}
