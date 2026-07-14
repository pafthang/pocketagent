package gate

import (
	"time"

	agentclient "github.com/pafthang/pocketagent/internal/agent/client"
	filesclient "github.com/pafthang/pocketagent/internal/files/client"
	gateapis "github.com/pafthang/pocketagent/internal/gate/apis"
	memoapis "github.com/pafthang/pocketagent/internal/memo/apis"
	memoclient "github.com/pafthang/pocketagent/internal/memo/client"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	spaceclient "github.com/pafthang/pocketagent/internal/space/client"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

func buildDeps(cfg *Config) (*gateapis.Deps, error) {
	natsClient, err := natsclient.NewClient(cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	pbClient, err := pbclient.NewServiceClient(cfg.PocketBaseURL, cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	if err != nil {
		natsClient.Close()
		return nil, err
	}

	return &gateapis.Deps{
		NATS:       natsClient,
		PB:         pbClient,
		Space:      spaceclient.New(cfg.SpaceURL),
		Agent:      agentclient.New(cfg.AgentURL),
		Files:      filesclient.New(cfg.FilesURL),
		RBAC:       apimw.NewPocketRBAC(pbClient, time.Duration(cfg.AuthorizeCacheSecs)*time.Second),
		Memo:       memoapis.Deps{Memo: memoclient.New(cfg.MemoURL, cfg.MemoServiceToken), Ollama: ollama.NewConfigured(cfg.OllamaURL, cfg.EmbedModel)},
		EmbedModel: cfg.EmbedModel,
		OllamaURL:  cfg.OllamaURL,
		LLMModel:   cfg.LLMModel,
		RateLimit:  cfg.RateLimit,
	}, nil
}