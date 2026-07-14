package space

import (
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/pafthang/pocketagent/pkgs/service"
)

func Run() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	deps, err := buildDeps(cfg)
	if err != nil {
		return err
	}

	s := service.NewServer(cfg.Service, cfg.ListenAddr(), cfg.LogLevel)
	s.AddHealth(common.Deps{PocketBaseURL: cfg.PocketBaseURL})
	s.AddMetrics()

	registerRoutes(s, deps, cfg)

	return s.Start()
}