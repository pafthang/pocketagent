package skillapis

import (
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/models"
)

// CreateSkillRequest is the gate API body for POST /skills.
type CreateSkillRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Prompt       string   `json:"prompt"`
	Category     string   `json:"category"`
	Tools        []string `json:"tools"`
	ArgumentHint string   `json:"argument_hint"`
	CatalogID    string   `json:"catalog_id"`
}

// ToModel builds a skill record with required-field validation.
func (r CreateSkillRequest) ToModel(spaceID string) (models.Skill, error) {
	name := strings.TrimSpace(r.Name)
	prompt := strings.TrimSpace(r.Prompt)
	if name == "" || prompt == "" {
		return models.Skill{}, fmt.Errorf("name and prompt are required")
	}
	return models.Skill{
		SpaceID:      spaceID,
		Name:         name,
		Description:  strings.TrimSpace(r.Description),
		Prompt:       prompt,
		Category:     strings.TrimSpace(r.Category),
		Tools:        r.Tools,
		ArgumentHint: strings.TrimSpace(r.ArgumentHint),
		CatalogID:    strings.TrimSpace(r.CatalogID),
	}, nil
}

// PatchSkillRequest is the gate API body for PATCH /skills/:id.
type PatchSkillRequest struct {
	Name         *string   `json:"name"`
	Description  *string   `json:"description"`
	Prompt       *string   `json:"prompt"`
	Category     *string   `json:"category"`
	Tools        *[]string `json:"tools"`
	ArgumentHint *string   `json:"argument_hint"`
}

// ApplyPatch mutates a skill record in place.
func (r PatchSkillRequest) ApplyPatch(skill *models.Skill) {
	if r.Name != nil {
		skill.Name = strings.TrimSpace(*r.Name)
	}
	if r.Description != nil {
		skill.Description = strings.TrimSpace(*r.Description)
	}
	if r.Prompt != nil {
		skill.Prompt = strings.TrimSpace(*r.Prompt)
	}
	if r.Category != nil {
		skill.Category = strings.TrimSpace(*r.Category)
	}
	if r.Tools != nil {
		skill.Tools = *r.Tools
	}
	if r.ArgumentHint != nil {
		skill.ArgumentHint = strings.TrimSpace(*r.ArgumentHint)
	}
}

// RunSkillRequest is the body for POST /skills/:id/run.
type RunSkillRequest struct {
	AgentID string `json:"agent_id"`
	Input   string `json:"input"`
}


