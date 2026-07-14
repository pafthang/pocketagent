package pocket

import (
	"path/filepath"

	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/spf13/viper"
)

// Config for embedded PocketBase (configs/pocket.yaml).
//
// Credentials — two env families (often the same values in dev):
//
//   - POCKETBASE_SUPERUSER_EMAIL/PASSWORD — pocket service only. On first start,
//     bootstrap creates the auth user, PocketBase superuser record, and admin
//     membership in the system "admin" space (see bootstrap_superuser.go).
//
//   - POCKETBASE_ADMIN_EMAIL/PASSWORD — gate, agent, space, exec, task. Used by
//     pocket/client to authenticate as a PocketBase admin for CRUD, RBAC, and
//     workers. In compose, ADMIN_* defaults to SUPERUSER_* when unset.
//
// Production: set both explicitly; do not rely on dev defaults.
type Config struct {
	Service           string `mapstructure:"service"`
	LogLevel          string `mapstructure:"log_level"`
	Port              string `mapstructure:"port"`
	DataDir           string `mapstructure:"data_dir"`
	SuperuserEmail    string `mapstructure:"superuser_email"`
	SuperuserPassword string `mapstructure:"superuser_password"`
}

// LoadConfig reads configs/pocket.yaml.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := common.Load(common.LoaderOptions{
		Service: "pocket",
		Defaults: func(v *viper.Viper) {
			v.SetDefault("service", "pocket")
			v.SetDefault("port", "8090")
			v.SetDefault("log_level", "info")
			v.SetDefault("data_dir", "data")
		},
	}, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Service == "" {
		cfg.Service = "pocket"
	}
	cfg.SuperuserEmail, cfg.SuperuserPassword = common.ResolvePocketBaseSuperuser(cfg.SuperuserEmail, cfg.SuperuserPassword)
	if err := validatePocketSecrets(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validatePocketSecrets(cfg *Config) error {
	if err := common.ValidateRequiredSecret("POCKETBASE_SUPERUSER_EMAIL", cfg.SuperuserEmail); err != nil {
		return err
	}
	return common.ValidateRequiredSecret("POCKETBASE_SUPERUSER_PASSWORD", cfg.SuperuserPassword)
}

func (c *Config) ListenAddr() string {
	return ":" + c.Port
}

func (c *Config) EnvMapWithRoot(root string) map[string]string {
	env := map[string]string{"LOG_LEVEL": c.LogLevel}
	dataDir := c.DataDir
	if dataDir != "" && !filepath.IsAbs(dataDir) {
		dataDir = filepath.Join(root, dataDir)
	}
	common.SetEnvMap(env,
		"PORT", c.Port,
		"POCKETBASE_DATA_DIR", dataDir,
		"POCKETBASE_SUPERUSER_EMAIL", c.SuperuserEmail,
		"POCKETBASE_SUPERUSER_PASSWORD", c.SuperuserPassword,
	)
	return env
}
