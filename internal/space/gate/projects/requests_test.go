package projectapis

import (
	"testing"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestCreateProjectRequestToModel(t *testing.T) {
	project := CreateProjectRequest{
		Goal:        " Ship it ",
		Description: "details",
		Tags:        []string{"a"},
	}.ToModel("space-1", "user-1", "My Project")

	if project.SpaceID != "space-1" || project.CreatorID != "user-1" {
		t.Fatalf("ids: %+v", project)
	}
	if project.Title != "My Project" || project.Goal != "Ship it" {
		t.Fatalf("fields: %+v", project)
	}
	if project.Status != models.ProjectDraft {
		t.Fatalf("status = %q", project.Status)
	}
}

func TestPatchProjectRequestApplyPatch(t *testing.T) {
	title := "New title"
	goal := " updated "
	project := models.Project{Title: "Old", Goal: "g"}
	PatchProjectRequest{Title: &title, Goal: &goal}.ApplyPatch(&project, func(t, g string) string {
		return t
	})
	if project.Title != "New title" || project.Goal != "updated" {
		t.Fatalf("project = %+v", project)
	}
}
