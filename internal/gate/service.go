package gate

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
	defer deps.NATS.Close()

	s := service.NewServer(cfg.Service, cfg.ListenAddr(), cfg.LogLevel)
	if cfg.RateLimit.EffectiveEnabled() {
		s.Echo.Use(common.APIRateLimiter(cfg.RateLimit))
	}

	s.AddHealth(common.Deps{
		NATS:          deps.NATS.Conn(),
		PocketBaseURL: cfg.PocketBaseURL,
		SpaceURL:      cfg.SpaceURL,
		AgentURL:      cfg.AgentURL,
		FilesURL:      cfg.FilesURL,
		MemoURL:       cfg.MemoURL,
		OllamaURL:     cfg.OllamaURL,
	})
	s.AddMetrics()

	registerRoutes(s, deps)

	return s.Start()
}