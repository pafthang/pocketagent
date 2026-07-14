package space

import (
	"strconv"

	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/spf13/viper"
)

// Config for the space service (configs/space.yaml).
type Config struct {
	Service                  string                 `mapstructure:"service"`
	LogLevel                 string                 `mapstructure:"log_level"`
	Port                     string                 `mapstructure:"port"`
	PocketBaseURL            string                 `mapstructure:"pocketbase_url"`
	PocketBaseAdminEmail     string                 `mapstructure:"pocketbase_admin_email"`
	PocketBaseAdminPass      string                 `mapstructure:"pocketbase_admin_password"`
	PublicBaseURL            string                 `mapstructure:"public_base_url"`
	RequireEmailVerification bool                   `mapstructure:"require_email_verification"`
	InviteTTLHours           int                    `mapstructure:"invite_ttl_hours"`
	VerificationTTLHours     int                    `mapstructure:"verification_ttl_hours"`
	RateLimit                common.RateLimitConfig `mapstructure:",squash"`
}

// LoadConfig reads configs/space.yaml.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := common.Load(common.LoaderOptions{
		Service: "space",
		Defaults: func(v *viper.Viper) {
			v.SetDefault("service", "space")
			v.SetDefault("port", "8083")
			v.SetDefault("log_level", "info")
			v.SetDefault("pocketbase_url", "http://127.0.0.1:8090")
			v.SetDefault("public_base_url", "http://127.0.0.1:8080")
			v.SetDefault("invite_ttl_hours", 168)
			v.SetDefault("verification_ttl_hours", 24)
		},
	}, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Service == "" {
		cfg.Service = "space"
	}
	if cfg.InviteTTLHours <= 0 {
		cfg.InviteTTLHours = 168
	}
	if cfg.VerificationTTLHours <= 0 {
		cfg.VerificationTTLHours = 24
	}
	cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass = common.ResolvePocketBaseAdmin(cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	if err := validateSpaceSecrets(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validateSpaceSecrets(cfg *Config) error {
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
		"PUBLIC_BASE_URL", c.PublicBaseURL,
	)
	if c.RequireEmailVerification {
		env["REQUIRE_EMAIL_VERIFICATION"] = "true"
	}
	if c.InviteTTLHours > 0 {
		env["INVITE_TTL_HOURS"] = strconv.Itoa(c.InviteTTLHours)
	}
	if c.VerificationTTLHours > 0 {
		env["VERIFICATION_TTL_HOURS"] = strconv.Itoa(c.VerificationTTLHours)
	}
	_ = root
	return env
}
