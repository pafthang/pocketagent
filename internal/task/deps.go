package task

import (
	memoclient "github.com/pafthang/pocketagent/internal/memo/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/task/decompose"
	"github.com/pafthang/pocketagent/internal/task/orchestrator"
	"github.com/pafthang/pocketagent/pkgs/ollama"
	"github.com/pafthang/pocketagent/pkgs/service"
)

// WorkerDeps holds runtime dependencies for the task orchestrator worker.
type WorkerDeps struct {
	Config     *Config
	Pocket     *pbclient.Client
	Store      *orchestrator.Store
	Decomposer *decompose.Decomposer
	Memory     *memoclient.Client
	Ollama     *ollama.Client
}

func buildDeps(w *service.Worker, cfg *Config) (*WorkerDeps, error) {
	pb, err := pbclient.NewServiceClient(cfg.PocketBaseURL, cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	if err != nil {
		return nil, err
	}

	ollamaClient := ollama.NewConfigured(cfg.OllamaURL, cfg.EmbedModel)

	return &WorkerDeps{
		Config:     cfg,
		Pocket:     pb,
		Store:      orchestrator.NewStore(pb, w.Log),
		Decomposer: decompose.New(ollamaClient, pb, cfg.LLMModel, cfg.MaxSubtasks, w.Log),
		Memory:     memoclient.New(cfg.MemoURL, cfg.MemoServiceToken),
		Ollama:     ollamaClient,
	}, nil
}
