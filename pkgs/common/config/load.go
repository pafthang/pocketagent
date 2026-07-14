package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// LoaderOptions configures service YAML loading.
type LoaderOptions struct {
	Service  string
	Defaults func(*viper.Viper)
}

// Load reads configs/<service>.yaml into dest with env overrides.
func Load(opts LoaderOptions, dest any) error {
	_ = godotenv.Load()

	v := viper.New()
	bindCommonEnv(v)

	if opts.Defaults != nil {
		opts.Defaults(v)
	}

	dir, err := FindConfigsDir()
	if err != nil {
		return err
	}

	v.SetConfigName(opts.Service)
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return v.Unmarshal(dest)
}