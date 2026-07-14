package memo

import (
	"log"

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
	s.AddHealth(common.Deps{MemoStore: deps.Manager.Ping})
	s.AddMetrics()

	log.Printf("memo ready (dir: %s, default: %s, per-space collections, min_similarity: %.2f)",
		cfg.DataDir, cfg.Collection, cfg.RAGMinSimilarity)

	registerRoutes(s, deps)

	return s.Start()
}