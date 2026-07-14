package catalog

import (
	"fmt"

	"github.com/pafthang/pocketagent/internal/agent"
	"github.com/pafthang/pocketagent/internal/exec"
	"github.com/pafthang/pocketagent/internal/files"
	"github.com/pafthang/pocketagent/internal/gate"
	"github.com/pafthang/pocketagent/internal/memo"
	"github.com/pafthang/pocketagent/internal/nats"
	"github.com/pafthang/pocketagent/internal/pocket"
	"github.com/pafthang/pocketagent/internal/space"
	"github.com/pafthang/pocketagent/internal/task"
)

type envProvider interface {
	EnvMapWithRoot(root string) map[string]string
}

func loadServiceEnv(name, root string) (map[string]string, error) {
	provider, err := configProvider(name)
	if err != nil {
		return nil, err
	}
	return provider.EnvMapWithRoot(root), nil
}

func configProvider(name string) (envProvider, error) {
	switch name {
	case "gate":
		return gate.LoadConfig()
	case "agent":
		return agent.LoadConfig()
	case "files":
		return files.LoadConfig()
	case "memo":
		return memo.LoadConfig()
	case "exec":
		return exec.LoadConfig()
	case "task":
		return task.LoadConfig()
	case "nats":
		return nats.LoadConfig()
	case "pocket":
		return pocket.LoadConfig()
	case "space":
		return space.LoadConfig()
	default:
		return nil, fmt.Errorf("unknown service %q", name)
	}
}