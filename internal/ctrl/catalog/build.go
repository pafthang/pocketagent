package catalog

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pafthang/pocketagent/pkgs/common"
)

// Build resolves ctrl.yaml services into runnable process definitions.
func Build(root string, ctrlCfg *common.CtrlConfig) ([]Service, error) {
	defs := make([]Service, 0, len(ctrlCfg.Services))

	for name, svc := range ctrlCfg.Services {
		cfgPath, err := common.ConfigFilePath(name)
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(cfgPath); err != nil {
			return nil, fmt.Errorf("missing config for %q: %s", name, cfgPath)
		}

		env, err := loadServiceEnv(name, root)
		if err != nil {
			return nil, err
		}

		defs = append(defs, Service{
			Name:       name,
			Package:    svc.Package,
			WaitPort:   svc.WaitPort,
			HealthPort: parseHealthPort(env),
			DependsOn:  svc.DependsOn,
			Env:        env,
		})
	}

	return defs, nil
}

func parseHealthPort(env map[string]string) int {
	if env == nil {
		return 0
	}
	raw, ok := env["HEALTH_PORT"]
	if !ok || raw == "" {
		return 0
	}
	port, err := strconv.Atoi(raw)
	if err != nil || port <= 0 {
		return 0
	}
	return port
}