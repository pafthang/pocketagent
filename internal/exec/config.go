package exec

import (
	"strconv"

	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/spf13/viper"
)

// Config for the exec worker (configs/exec.yaml).
type Config struct {
	Service              string `mapstructure:"service"`
	LogLevel             string `mapstructure:"log_level"`
	NatsURL              string `mapstructure:"nats_url"`
	OllamaURL            string `mapstructure:"ollama_url"`
	MemoURL              string `mapstructure:"memo_url"`
	MemoServiceToken     string `mapstructure:"memo_service_token"`
	PocketBaseURL        string `mapstructure:"pocketbase_url"`
	PocketBaseAdminEmail string `mapstructure:"pocketbase_admin_email"`
	PocketBaseAdminPass  string `mapstructure:"pocketbase_admin_password"`
	EmbedModel           string `mapstructure:"embed_model"`
	LLMModel             string `mapstructure:"llm_model"`
	HealthPort           string `mapstructure:"health_port"`
	StreamLLMTokens      bool   `mapstructure:"stream_llm_tokens"`
	SearchProvider       string `mapstructure:"search_provider"`
	SerperAPIKey         string `mapstructure:"serper_api_key"`
	TavilyAPIKey         string `mapstructure:"tavily_api_key"`
	CodeExecEnabled      bool   `mapstructure:"code_exec_enabled"`
	CodeExecTimeoutSec   int    `mapstructure:"code_exec_timeout_sec"`
	MCPServersJSON       string `mapstructure:"mcp_servers"`
}

// LoadConfig reads configs/exec.yaml.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := common.Load(common.LoaderOptions{
		Service: "exec",
		Defaults: func(v *viper.Viper) {
			v.SetDefault("service", "exec")
			v.SetDefault("log_level", "info")
			v.SetDefault("nats_url", "nats://127.0.0.1:4222")
			v.SetDefault("ollama_url", "http://127.0.0.1:11434")
			v.SetDefault("memo_url", "http://127.0.0.1:8082")
			v.SetDefault("pocketbase_url", "http://127.0.0.1:8090")
			v.SetDefault("embed_model", "nomic-embed-text")
			v.SetDefault("llm_model", "llama3.1")
			v.SetDefault("health_port", "9084")
			v.SetDefault("stream_llm_tokens", true)
			v.SetDefault("search_provider", "auto")
			v.SetDefault("code_exec_enabled", false)
			v.SetDefault("code_exec_timeout_sec", 10)
		},
	}, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Service == "" {
		cfg.Service = "exec"
	}
	cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass = common.ResolvePocketBaseAdmin(cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	cfg.MemoServiceToken = common.ResolveMemoServiceToken(cfg.MemoServiceToken)
	if err := validateExecSecrets(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validateExecSecrets(cfg *Config) error {
	if err := common.ValidateRequiredSecret("POCKETBASE_ADMIN_EMAIL", cfg.PocketBaseAdminEmail); err != nil {
		return err
	}
	return common.ValidateRequiredSecret("POCKETBASE_ADMIN_PASSWORD", cfg.PocketBaseAdminPass)
}

func (c *Config) EnvMapWithRoot(root string) map[string]string {
	env := map[string]string{"LOG_LEVEL": c.LogLevel}
	common.SetEnvMap(env,
		"NATS_URL", c.NatsURL,
		"OLLAMA_URL", c.OllamaURL,
		"MEMO_URL", c.MemoURL,
		"MEMO_SERVICE_TOKEN", c.MemoServiceToken,
		"POCKETBASE_URL", c.PocketBaseURL,
		"POCKETBASE_ADMIN_EMAIL", c.PocketBaseAdminEmail,
		"POCKETBASE_ADMIN_PASSWORD", c.PocketBaseAdminPass,
		"EMBED_MODEL", c.EmbedModel,
		"LLM_MODEL", c.LLMModel,
		"HEALTH_PORT", c.HealthPort,
	)
	if c.StreamLLMTokens {
		env["STREAM_LLM_TOKENS"] = "true"
	}
	common.SetEnvMap(env,
		"SEARCH_PROVIDER", c.SearchProvider,
		"SERPER_API_KEY", c.SerperAPIKey,
		"TAVILY_API_KEY", c.TavilyAPIKey,
		"MCP_SERVERS", c.MCPServersJSON,
	)
	if c.CodeExecEnabled {
		env["CODE_EXEC_ENABLED"] = "true"
	}
	if c.CodeExecTimeoutSec > 0 {
		env["CODE_EXEC_TIMEOUT_SEC"] = strconv.Itoa(c.CodeExecTimeoutSec)
	}
	_ = root
	return env
}
