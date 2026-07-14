package ctrl

import (
	"github.com/pafthang/pocketagent/internal/ctrl/catalog"
	"github.com/pafthang/pocketagent/internal/ctrl/supervisor"
	"github.com/pafthang/pocketagent/pkgs/common"
)

// RuntimeDeps holds everything needed to run the local orchestrator.
type RuntimeDeps struct {
	Root       string
	ConfigDir  string
	Config     *Config
	Services   []catalog.Service
	Supervisor *supervisor.Supervisor
}

func buildDeps(configDirFlag string) (*RuntimeDeps, error) {
	root, err := common.FindProjectRoot()
	if err != nil {
		return nil, err
	}

	configDir, err := common.InitRuntimeDirs(root, configDirFlag)
	if err != nil {
		return nil, err
	}

	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	services, err := catalog.Build(root, cfg)
	if err != nil {
		return nil, err
	}

	return &RuntimeDeps{
		Root:       root,
		ConfigDir:  configDir,
		Config:     cfg,
		Services:   services,
		Supervisor: supervisor.New(root, configDir, cfg),
	}, nil
}