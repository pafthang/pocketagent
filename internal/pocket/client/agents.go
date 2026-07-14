package client

import (
	"fmt"

	"github.com/pafthang/pocketagent/pkgs/models"
)



// CreateAgent stores a new agent record.
func (c *Client) CreateAgent(agent models.Agent) (models.Agent, error) {
	record, err := c.CreateRecord(AgentsCollection, agentRecordData(agent))
	if err != nil {
		return models.Agent{}, err
	}
	return agentFromRecord(record), nil
}

// GetAgent returns an agent by ID.
func (c *Client) GetAgent(id string) (models.Agent, error) {
	record, err := c.GetRecord(AgentsCollection, id)
	if err != nil {
		return models.Agent{}, err
	}
	return agentFromRecord(record), nil
}

// ListAgents returns agents with pagination metadata.
func (c *Client) ListAgents(page, perPage int) ([]models.Agent, int, error) {
	return c.ListAgentsFilter(ListOptions{Page: page, PerPage: perPage})
}

// ListAgentsFilter returns agents with optional PocketBase filter.
func (c *Client) ListAgentsFilter(opts ListOptions) ([]models.Agent, int, error) {
	records, total, err := c.ListRecordsOpts(AgentsCollection, opts)
	if err != nil {
		return nil, 0, err
	}

	agents := make([]models.Agent, 0, len(records))
	for _, record := range records {
		agents = append(agents, agentFromRecord(record))
	}
	return agents, total, nil
}

// UpdateAgent patches an existing agent.
func (c *Client) UpdateAgent(id string, agent models.Agent) (models.Agent, error) {
	record, err := c.UpdateRecord(AgentsCollection, id, agentRecordData(agent))
	if err != nil {
		return models.Agent{}, err
	}
	return agentFromRecord(record), nil
}

// DeleteAgent removes an agent by ID.
func (c *Client) DeleteAgent(id string) error {
	return c.DeleteRecord(AgentsCollection, id)
}

func agentRecordData(agent models.Agent) map[string]interface{} {
	data := map[string]interface{}{
		"space_id":      agent.SpaceID,
		"name":          agent.Name,
		"description":   agent.Description,
		"model":         agent.Model,
		"system_prompt": agent.SystemPrompt,
	}
	if agent.Tools != nil {
		data["tools"] = agent.Tools
	}
	if agent.Config != nil {
		data["config"] = agent.Config
	}
	return data
}

func agentFromRecord(record map[string]interface{}) models.Agent {
	agent := models.Agent{
		ID:           stringField(record, "id"),
		SpaceID:      stringField(record, "space_id"),
		Name:         stringField(record, "name"),
		Description:  stringField(record, "description"),
		Model:        stringField(record, "model"),
		SystemPrompt: stringField(record, "system_prompt"),
		CreatedAt:    stringField(record, "created"),
		UpdatedAt:    stringField(record, "updated"),
	}

	if tools, ok := record["tools"].([]interface{}); ok {
		agent.Tools = make([]string, 0, len(tools))
		for _, t := range tools {
			if s, ok := t.(string); ok {
				agent.Tools = append(agent.Tools, s)
			}
		}
	}

	if cfg, ok := record["config"].(map[string]interface{}); ok {
		agent.Config = cfg
	}

	return agent
}

func stringField(record map[string]interface{}, key string) string {
	if v, ok := record[key]; ok {
		return fmt.Sprint(v)
	}
	return ""
}
