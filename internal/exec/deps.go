package exec

import (
	"encoding/json"
	"os"
	"time"

	"github.com/pafthang/pocketagent/internal/exec/react"
	"github.com/pafthang/pocketagent/internal/exec/tools"
	memoclient "github.com/pafthang/pocketagent/internal/memo/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

// WorkerDeps holds runtime dependencies for the exec worker.
type WorkerDeps struct {
	Config   *Config
	Pocket   *pbclient.Client
	Ollama   *ollama.Client
	Memory   *memoclient.Client
	Executor *react.Executor
	ToolCfg  tools.Config
	ToolSet  *tools.Set
}

func buildDeps(cfg *Config) (*WorkerDeps, error) {
	ollamaClient := ollama.NewConfigured(cfg.OllamaURL, cfg.EmbedModel)
	memoryClient := memoclient.New(cfg.MemoURL, cfg.MemoServiceToken)
	baseToolCfg := toolConfig(cfg)
	baseToolCfg.MCPServers = nil
	toolSet := tools.Build(baseToolCfg)

	executor := react.New(ollamaClient, toolSet.Catalog, cfg.LLMModel).
		WithMemory(memoryClient).
		WithStreaming(cfg.StreamLLMTokens, nil)
	executor.ExecuteTool = react.ToolRunner(toolSet.Registry)

	pb, err := pbclient.NewServiceClient(cfg.PocketBaseURL, cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	if err != nil {
		toolSet.Close()
		return nil, err
	}

	return &WorkerDeps{
		Config:   cfg,
		Pocket:   pb,
		Ollama:   ollamaClient,
		Memory:   memoryClient,
		Executor: executor,
		ToolCfg:  baseToolCfg,
		ToolSet:  toolSet,
	}, nil
}

func toolConfig(cfg *Config) tools.Config {
	toolsCfg := tools.LoadFromEnv()
	if cfg == nil {
		return toolsCfg
	}

	if cfg.SearchProvider != "" {
		toolsCfg.SearchProvider = cfg.SearchProvider
	}
	if cfg.SerperAPIKey != "" {
		toolsCfg.SerperAPIKey = cfg.SerperAPIKey
	}
	if cfg.TavilyAPIKey != "" {
		toolsCfg.TavilyAPIKey = cfg.TavilyAPIKey
	}
	if cfg.CodeExecEnabled {
		toolsCfg.CodeExecEnabled = true
	}
	if cfg.CodeExecTimeoutSec > 0 {
		toolsCfg.CodeExecTimeout = time.Duration(cfg.CodeExecTimeoutSec) * time.Second
	}
	if cfg.MCPServersJSON != "" {
		_ = json.Unmarshal([]byte(cfg.MCPServersJSON), &toolsCfg.MCPServers)
	}

	if os.Getenv("CODE_EXEC_ENABLED") == "" && !cfg.CodeExecEnabled && !toolsCfg.CodeExecEnabled {
		toolsCfg.CodeExecEnabled = tools.DefaultConfig().CodeExecEnabled
	}

	return toolsCfg
}