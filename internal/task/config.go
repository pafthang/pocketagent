package task

import (
	"strconv"
	"time"

	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/spf13/viper"
)

// Config for the task orchestrator (configs/task.yaml).
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
	MaxSubtasks          int    `mapstructure:"max_subtasks"`
	TimeoutSec           int    `mapstructure:"timeout_sec"`
	HealthPort           string `mapstructure:"health_port"`
	SchedulerIntervalSec int    `mapstructure:"scheduler_interval_sec"`
}

// LoadConfig reads configs/task.yaml.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := common.Load(common.LoaderOptions{
		Service: "task",
		Defaults: func(v *viper.Viper) {
			v.SetDefault("service", "task")
			v.SetDefault("log_level", "info")
			v.SetDefault("nats_url", "nats://127.0.0.1:4222")
			v.SetDefault("ollama_url", "http://127.0.0.1:11434")
			v.SetDefault("memo_url", "http://127.0.0.1:8082")
			v.SetDefault("pocketbase_url", "http://127.0.0.1:8090")
			v.SetDefault("embed_model", "nomic-embed-text")
			v.SetDefault("llm_model", "llama3.1")
			v.SetDefault("max_subtasks", 4)
			v.SetDefault("timeout_sec", 30)
			v.SetDefault("health_port", "9085")
			v.SetDefault("scheduler_interval_sec", 60)
		},
	}, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Service == "" {
		cfg.Service = "task"
	}
	cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass = common.ResolvePocketBaseAdmin(cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	cfg.MemoServiceToken = common.ResolveMemoServiceToken(cfg.MemoServiceToken)
	if err := validateTaskSecrets(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validateTaskSecrets(cfg *Config) error {
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
	if c.MaxSubtasks > 0 {
		env["TASK_MAX_SUBTASKS"] = strconv.Itoa(c.MaxSubtasks)
	}
	if c.TimeoutSec > 0 {
		env["TASK_TIMEOUT_SEC"] = strconv.Itoa(c.TimeoutSec)
	}
	if c.SchedulerIntervalSec > 0 {
		env["SCHEDULER_INTERVAL_SEC"] = strconv.Itoa(c.SchedulerIntervalSec)
	}
	_ = root
	return env
}

func (c *Config) SchedulerInterval() time.Duration {
	if c == nil || c.SchedulerIntervalSec <= 0 {
		return time.Minute
	}
	return time.Duration(c.SchedulerIntervalSec) * time.Second
}
