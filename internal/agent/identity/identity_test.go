package identity

import (
	"strings"
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestFromAgent(t *testing.T) {
	agent := models.Agent{
		SystemPrompt: "I am Agent X",
		Config: map[string]interface{}{
			ConfigKey: map[string]interface{}{
				"soul":         "be kind",
				"style":        "concise",
				"instructions": "use tools",
			},
		},
	}
	files := FromAgent(agent)
	if files.IdentityFile != "I am Agent X" {
		t.Fatalf("identity_file: %q", files.IdentityFile)
	}
	if files.SoulFile != "be kind" || files.StyleFile != "concise" || files.InstructionsFile != "use tools" {
		t.Fatalf("unexpected blocks: %+v", files)
	}
}

func TestApplyPatchPartial(t *testing.T) {
	agent := models.Agent{
		SystemPrompt: "old",
		Config: map[string]interface{}{
			ConfigKey: map[string]interface{}{"soul": "keep"},
		},
	}
	soul := "new soul"
	patched, updated := ApplyPatch(agent, Patch{SoulFile: &soul})
	if patched.SystemPrompt != "old" {
		t.Fatal("identity_file should be unchanged")
	}
	blocks := patched.Config[ConfigKey].(map[string]interface{})
	if blocks["soul"] != "new soul" {
		t.Fatalf("soul not updated: %v", blocks)
	}
	if len(updated) != 1 || updated[0] != "soul_file" {
		t.Fatalf("updated: %v", updated)
	}
}

func TestCompileAgentPrompt(t *testing.T) {
	prompt := CompileAgentPrompt(models.IdentityFiles{
		IdentityFile:     "who",
		SoulFile:         "values",
		StyleFile:        "tone",
		InstructionsFile: "rules",
	})
	if prompt == "" {
		t.Fatal("expected compiled prompt")
	}
	for _, want := range []string{"# Identity", "# Soul", "# Style", "# Instructions"} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("missing %q in %q", want, prompt)
		}
	}
}
