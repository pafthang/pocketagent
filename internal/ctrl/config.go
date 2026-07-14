package ctrl

import "github.com/pafthang/pocketagent/pkgs/common"

// Config is orchestration settings from configs/ctrl.yaml.
type Config = common.CtrlConfig

// LoadConfig reads configs/ctrl.yaml.
func LoadConfig() (*Config, error) {
	return common.LoadCtrlConfig()
}