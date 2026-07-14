package task

import (
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/pafthang/pocketagent/pkgs/service"
)

func Run() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	w, err := service.NewWorker(cfg.Service, cfg.NatsURL, cfg.LogLevel)
	if err != nil {
		return err
	}

	w.ServeHealth(cfg.HealthPort, common.Deps{
		PocketBaseURL: cfg.PocketBaseURL,
		OllamaURL:     cfg.OllamaURL,
		MemoURL:       cfg.MemoURL,
	})

	if _, err := wireWorker(w, cfg); err != nil {
		return err
	}

	return w.Run()
}
