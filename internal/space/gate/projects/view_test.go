package projectapis

import (
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestProgress(t *testing.T) {
	p := Progress([]models.ProjectItem{
		{Status: models.ItemDone},
		{Status: models.ItemInProgress},
		{Status: models.ItemInbox},
	})
	if p.Total != 3 || p.Completed != 1 || p.InProgress != 1 || p.HumanPending != 1 || p.Percent != 33 {
		t.Fatalf("unexpected progress: %+v", p)
	}
}

func TestNormalizeTitleFromGoal(t *testing.T) {
	if NormalizeTitle("", "Build a thing") != "Build a thing" {
		t.Fatal("expected goal as title")
	}
}
