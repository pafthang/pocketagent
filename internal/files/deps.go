package files

import (
	fileapis "github.com/pafthang/pocketagent/internal/files/apis"
	"github.com/pafthang/pocketagent/internal/files/blob"
	memoclient "github.com/pafthang/pocketagent/internal/memo/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

func buildDeps(cfg *Config) (*fileapis.Deps, error) {
	store, err := blob.NewBackend(cfg.StoreConfig())
	if err != nil {
		return nil, err
	}

	pb, err := pbclient.NewServiceClient(cfg.PocketBaseURL, cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	if err != nil {
		return nil, err
	}

	return &fileapis.Deps{
		PB:     pb,
		Store:  store,
		Memo:   memoclient.New(cfg.MemoURL, cfg.MemoServiceToken),
		Ollama: ollama.NewConfigured(cfg.OllamaURL, cfg.EmbedModel),
	}, nil
}
