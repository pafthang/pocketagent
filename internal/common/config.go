package common

import (
	"github.com/spf13/viper"
)

// Config holds common configuration

type Config struct {
	ServiceName string
	Port        string
	NatsURL     string
	LogLevel    string
}

func LoadConfig() *Config {
	viper.SetDefault("port", "8080")
	viper.SetDefault("nats_url", "nats://nats:4222")
	viper.SetDefault("log_level", "info")

	viper.AutomaticEnv()

	return &Config{
		ServiceName: viper.GetString("SERVICE_NAME"),
		Port:        viper.GetString("PORT"),
		NatsURL:     viper.GetString("NATS_URL"),
		LogLevel:    viper.GetString("LOG_LEVEL"),
	}
}
