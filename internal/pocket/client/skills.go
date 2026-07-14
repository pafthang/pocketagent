package client

import (
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/models"
)



// CreateSkill stores a new skill record.
func (c *Client) CreateSkill(skill models.Skill) (models.Skill, error) {
	record, err := c.CreateRecord(SkillsCollection, skillRecordData(skill))
	if err != nil {
		return models.Skill{}, err
	}
	return skillFromRecord(record), nil
}

// GetSkill returns a skill by ID.
func (c *Client) GetSkill(id string) (models.Skill, error) {
	record, err := c.GetRecord(SkillsCollection, id)
	if err != nil {
		return models.Skill{}, err
	}
	return skillFromRecord(record), nil
}

// ListSkills returns skills with optional filter.
func (c *Client) ListSkills(opts ListOptions) ([]models.Skill, int, error) {
	records, total, err := c.ListRecordsOpts(SkillsCollection, opts)
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.Skill, 0, len(records))
	for _, record := range records {
		out = append(out, skillFromRecord(record))
	}
	return out, total, nil
}

// UpdateSkillRecord replaces a skill record.
func (c *Client) UpdateSkillRecord(skill models.Skill) (models.Skill, error) {
	record, err := c.UpdateRecord(SkillsCollection, skill.ID, skillRecordData(skill))
	if err != nil {
		return models.Skill{}, err
	}
	return skillFromRecord(record), nil
}

// DeleteSkill removes a skill by ID.
func (c *Client) DeleteSkill(id string) error {
	return c.DeleteRecord(SkillsCollection, id)
}

func skillRecordData(skill models.Skill) map[string]interface{} {
	data := map[string]interface{}{
		"space_id":    skill.SpaceID,
		"name":        skill.Name,
		"description": skill.Description,
		"prompt":      skill.Prompt,
	}
	if skill.Category != "" {
		data["category"] = skill.Category
	}
	if skill.Tools != nil {
		data["tools"] = skill.Tools
	}
	if skill.ArgumentHint != "" {
		data["argument_hint"] = skill.ArgumentHint
	}
	if skill.CatalogID != "" {
		data["catalog_id"] = skill.CatalogID
	}
	return data
}

func skillFromRecord(record map[string]interface{}) models.Skill {
	skill := models.Skill{
		ID:           stringField(record, "id"),
		SpaceID:      stringField(record, "space_id"),
		Name:         stringField(record, "name"),
		Description:  stringField(record, "description"),
		Prompt:       stringField(record, "prompt"),
		Category:     stringField(record, "category"),
		ArgumentHint: stringField(record, "argument_hint"),
		CatalogID:    stringField(record, "catalog_id"),
		CreatedAt:    stringField(record, "created"),
		UpdatedAt:    stringField(record, "updated"),
	}
	if tools, ok := record["tools"].([]interface{}); ok {
		skill.Tools = make([]string, 0, len(tools))
		for _, t := range tools {
			if s, ok := t.(string); ok {
				skill.Tools = append(skill.Tools, s)
			}
		}
	}
	return skill
}

func FindSkillByName(c *Client, spaceID, name string) (models.Skill, error) {
	name = strings.TrimSpace(name)
	skills, _, err := c.ListSkills(ListOptions{
		Page:    1,
		PerPage: 1,
		Filter:  fmt.Sprintf("space_id = %q && name = %q", spaceID, name),
	})
	if err != nil {
		return models.Skill{}, err
	}
	if len(skills) == 0 {
		return models.Skill{}, &APIError{StatusCode: 404, Message: "skill not found"}
	}
	return skills[0], nil
}
