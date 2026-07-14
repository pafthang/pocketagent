package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// CtrlServiceDef describes how ctrl should run a child service.
type CtrlServiceDef struct {
	Package   string   `mapstructure:"package"`
	WaitPort  int      `mapstructure:"wait_port"`
	DependsOn []string `mapstructure:"depends_on"`
}

// CtrlConfig is loaded from configs/ctrl.yaml.
type CtrlConfig struct {
	StartTimeoutSec int                       `mapstructure:"start_timeout_sec"`
	StopTimeoutSec  int                       `mapstructure:"stop_timeout_sec"`
	Services        map[string]CtrlServiceDef `mapstructure:"services"`
}

// LoadCtrlConfig reads orchestration settings for cmd/ctrl.
func LoadCtrlConfig() (*CtrlConfig, error) {
	v := viper.New()
	v.SetDefault("start_timeout_sec", 30)
	v.SetDefault("stop_timeout_sec", 10)

	dir, err := FindConfigsDir()
	if err != nil {
		return nil, err
	}

	v.SetConfigName("ctrl")
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read ctrl config: %w", err)
	}

	var cfg CtrlConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parse ctrl config: %w", err)
	}

	if len(cfg.Services) == 0 {
		return nil, fmt.Errorf("ctrl config: no services defined")
	}

	return &cfg, nil
}