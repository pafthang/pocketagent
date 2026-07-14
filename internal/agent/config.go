package agent

import (
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/spf13/viper"
)

// Config for the agent service (configs/agent.yaml).
type Config struct {
	Service              string `mapstructure:"service"`
	LogLevel             string `mapstructure:"log_level"`
	Port                 string `mapstructure:"port"`
	PocketBaseURL        string `mapstructure:"pocketbase_url"`
	PocketBaseAdminEmail string `mapstructure:"pocketbase_admin_email"`
	PocketBaseAdminPass  string `mapstructure:"pocketbase_admin_password"`
	SpaceURL             string `mapstructure:"space_url"`
	AuthorizeCacheSecs   int    `mapstructure:"authorize_cache_seconds"`
}

// LoadConfig reads configs/agent.yaml.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := common.Load(common.LoaderOptions{
		Service: "agent",
		Defaults: func(v *viper.Viper) {
			v.SetDefault("service", "agent")
			v.SetDefault("port", "8081")
			v.SetDefault("log_level", "info")
			v.SetDefault("pocketbase_url", "http://127.0.0.1:8090")
			v.SetDefault("space_url", "http://127.0.0.1:8083")
			v.SetDefault("authorize_cache_seconds", 30)
		},
	}, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Service == "" {
		cfg.Service = "agent"
	}
	cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass = common.ResolvePocketBaseAdmin(cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	if err := validateAgentSecrets(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validateAgentSecrets(cfg *Config) error {
	if err := common.ValidateRequiredSecret("POCKETBASE_ADMIN_EMAIL", cfg.PocketBaseAdminEmail); err != nil {
		return err
	}
	return common.ValidateRequiredSecret("POCKETBASE_ADMIN_PASSWORD", cfg.PocketBaseAdminPass)
}

func (c *Config) ListenAddr() string {
	return ":" + c.Port
}

func (c *Config) EnvMapWithRoot(root string) map[string]string {
	env := map[string]string{"LOG_LEVEL": c.LogLevel}
	common.SetEnvMap(env,
		"PORT", c.Port,
		"POCKETBASE_URL", c.PocketBaseURL,
		"POCKETBASE_ADMIN_EMAIL", c.PocketBaseAdminEmail,
		"POCKETBASE_ADMIN_PASSWORD", c.PocketBaseAdminPass,
		"SPACE_URL", c.SpaceURL,
	)
	_ = root
	return env
}
