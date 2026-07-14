package react

import (
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestEffectiveAllowedToolsTaskOverridesAgent(t *testing.T) {
	got := EffectiveAllowedTools(
		models.Task{Tools: []string{"code_exec"}},
		models.Agent{Tools: []string{"search_web"}},
	)
	if len(got) != 1 || got[0] != "code_exec" {
		t.Fatalf("unexpected tools: %#v", got)
	}
}

func TestEffectiveAllowedToolsFallsBackToAgent(t *testing.T) {
	got := EffectiveAllowedTools(
		models.Task{},
		models.Agent{Tools: []string{"search_web", "scrape_page"}},
	)
	if len(got) != 2 {
		t.Fatalf("unexpected tools: %#v", got)
	}
}

func TestEffectiveAllowedToolsEmptyMeansAll(t *testing.T) {
	got := EffectiveAllowedTools(models.Task{}, models.Agent{})
	if got != nil {
		t.Fatalf("expected nil allow-list, got %#v", got)
	}
}