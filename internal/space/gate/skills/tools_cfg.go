package skillapis

import "github.com/pafthang/pocketagent/internal/exec/tools"

func defaultToolConfig() tools.Config {
	return tools.LoadFromEnv()
}