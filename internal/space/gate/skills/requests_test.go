package skillapis

import "testing"

func TestCreateSkillRequestToModel(t *testing.T) {
	skill, err := CreateSkillRequest{
		Name:   "lint",
		Prompt: "Run linter",
		Tools:  []string{"shell"},
	}.ToModel("space-1")
	if err != nil {
		t.Fatal(err)
	}
	if skill.SpaceID != "space-1" || skill.Name != "lint" {
		t.Fatalf("skill = %+v", skill)
	}
}

func TestCreateSkillRequestRequiresFields(t *testing.T) {
	if _, err := (CreateSkillRequest{}).ToModel("space-1"); err == nil {
		t.Fatal("expected validation error")
	}
}