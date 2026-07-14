package skillapis

import (
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestComposeSkillPromptWithPlaceholder(t *testing.T) {
	prompt := composeSkillPrompt(models.Skill{Prompt: "Hello {{input}}"}, "world")
	if prompt != "Hello world" {
		t.Fatalf("unexpected prompt: %q", prompt)
	}
}

func TestComposeSkillPromptAppend(t *testing.T) {
	prompt := composeSkillPrompt(models.Skill{Prompt: "Do work"}, "details")
	if prompt != "Do work\n\ndetails" {
		t.Fatalf("unexpected prompt: %q", prompt)
	}
}

func TestSearchCatalog(t *testing.T) {
	results, err := searchCatalog("summarize")
	if err != nil {
		t.Fatalf("search catalog: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected catalog matches")
	}
}
