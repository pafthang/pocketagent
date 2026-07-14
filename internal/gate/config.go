package gate

import (
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/spf13/viper"
)

// Config for the gate HTTP facade (configs/gate.yaml).
type Config struct {
	Service              string                 `mapstructure:"service"`
	LogLevel             string                 `mapstructure:"log_level"`
	Port                 string                 `mapstructure:"port"`
	NatsURL              string                 `mapstructure:"nats_url"`
	PocketBaseURL        string                 `mapstructure:"pocketbase_url"`
	PocketBaseAdminEmail string                 `mapstructure:"pocketbase_admin_email"`
	PocketBaseAdminPass  string                 `mapstructure:"pocketbase_admin_password"`
	SpaceURL             string                 `mapstructure:"space_url"`
	AgentURL             string                 `mapstructure:"agent_url"`
	FilesURL             string                 `mapstructure:"files_url"`
	MemoURL              string                 `mapstructure:"memo_url"`
	MemoServiceToken     string                 `mapstructure:"memo_service_token"`
	OllamaURL            string                 `mapstructure:"ollama_url"`
	EmbedModel           string                 `mapstructure:"embed_model"`
	LLMModel             string                 `mapstructure:"llm_model"`
	AuthorizeCacheSecs   int                    `mapstructure:"authorize_cache_seconds"`
	RateLimit            common.RateLimitConfig `mapstructure:",squash"`
}

// LoadConfig reads configs/gate.yaml.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := common.Load(common.LoaderOptions{
		Service: "gate",
		Defaults: func(v *viper.Viper) {
			v.SetDefault("service", "gate")
			v.SetDefault("port", "8080")
			v.SetDefault("log_level", "info")
			v.SetDefault("nats_url", "nats://127.0.0.1:4222")
			v.SetDefault("pocketbase_url", "http://127.0.0.1:8090")
			v.SetDefault("space_url", "http://127.0.0.1:8083")
			v.SetDefault("agent_url", "http://127.0.0.1:8081")
			v.SetDefault("files_url", "http://127.0.0.1:8086")
			v.SetDefault("memo_url", "http://127.0.0.1:8082")
			v.SetDefault("ollama_url", "http://127.0.0.1:11434")
			v.SetDefault("embed_model", "nomic-embed-text")
			v.SetDefault("llm_model", "llama3.1")
			v.SetDefault("authorize_cache_seconds", 30)
		},
	}, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Service == "" {
		cfg.Service = "gate"
	}
	cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass = common.ResolvePocketBaseAdmin(cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	cfg.MemoServiceToken = common.ResolveMemoServiceToken(cfg.MemoServiceToken)
	if err := validateGateSecrets(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validateGateSecrets(cfg *Config) error {
	if err := common.ValidateRequiredSecret("POCKETBASE_ADMIN_EMAIL", cfg.PocketBaseAdminEmail); err != nil {
		return err
	}
	return common.ValidateRequiredSecret("POCKETBASE_ADMIN_PASSWORD", cfg.PocketBaseAdminPass)
}

// ListenAddr returns the HTTP listen address.
func (c *Config) ListenAddr() string {
	return ":" + c.Port
}

// EnvMapWithRoot converts config to env vars for child processes.
func (c *Config) EnvMapWithRoot(root string) map[string]string {
	env := map[string]string{"LOG_LEVEL": c.LogLevel}
	common.SetEnvMap(env,
		"PORT", c.Port,
		"NATS_URL", c.NatsURL,
		"POCKETBASE_URL", c.PocketBaseURL,
		"POCKETBASE_ADMIN_EMAIL", c.PocketBaseAdminEmail,
		"POCKETBASE_ADMIN_PASSWORD", c.PocketBaseAdminPass,
		"SPACE_URL", c.SpaceURL,
		"AGENT_URL", c.AgentURL,
		"FILES_URL", c.FilesURL,
		"MEMO_URL", c.MemoURL,
		"MEMO_SERVICE_TOKEN", c.MemoServiceToken,
		"OLLAMA_URL", c.OllamaURL,
		"EMBED_MODEL", c.EmbedModel,
		"LLM_MODEL", c.LLMModel,
	)
	_ = root
	return env
}
