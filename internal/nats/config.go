package nats

import (
	"path/filepath"

	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/spf13/viper"
)

// Config for embedded NATS (configs/nats.yaml).
type Config struct {
	Service  string `mapstructure:"service"`
	LogLevel string `mapstructure:"log_level"`
	Port     string `mapstructure:"port"`
	HTTPPort string `mapstructure:"http_port"`
	StoreDir string `mapstructure:"store_dir"`
}

// LoadConfig reads configs/nats.yaml.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := common.Load(common.LoaderOptions{
		Service: "nats",
		Defaults: func(v *viper.Viper) {
			v.SetDefault("service", "nats")
			v.SetDefault("port", "4222")
			v.SetDefault("http_port", "8222")
			v.SetDefault("log_level", "info")
			v.SetDefault("store_dir", "data/nats")
		},
	}, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Service == "" {
		cfg.Service = "nats"
	}
	if cfg.Port == "" {
		cfg.Port = "4222"
	}
	if cfg.HTTPPort == "" {
		cfg.HTTPPort = "8222"
	}
	return &cfg, nil
}

func (c *Config) EnvMapWithRoot(root string) map[string]string {
	env := map[string]string{"LOG_LEVEL": c.LogLevel}
	storeDir := c.StoreDir
	if storeDir != "" && !filepath.IsAbs(storeDir) {
		storeDir = filepath.Join(root, storeDir)
	}
	common.SetEnvMap(env, "PORT", c.Port, "HTTP_PORT", c.HTTPPort, "NATS_STORE_DIR", storeDir)
	return env
}
