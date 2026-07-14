package handler

import (
	"fmt"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func resolveAgent(pb *pbclient.Client, task models.Task) (models.Agent, error) {
	if task.AgentID == "" {
		return models.Agent{}, nil
	}

	agent, err := pb.GetAgent(task.AgentID)
	if err != nil {
		return models.Agent{}, fmt.Errorf("load agent %s: %w", task.AgentID, err)
	}

	if task.SpaceID != "" && agent.SpaceID != "" && agent.SpaceID != task.SpaceID {
		return models.Agent{}, fmt.Errorf("agent %s does not belong to space %s", task.AgentID, task.SpaceID)
	}

	return agent, nil
}