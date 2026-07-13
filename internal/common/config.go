package common

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds common configuration
type Config struct {
	ServiceName string
	Port        string
	NatsURL     string
	LogLevel    string
	OllamaURL   string
}

// LoadConfig loads configuration from .env + env vars
func LoadConfig() *Config {
	_ = godotenv.Load() // load .env if exists

	viper.SetDefault("port", "8080")
	viper.SetDefault("nats_url", "nats://nats:4222")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("ollama_url", "http://ollama:11434")

	viper.AutomaticEnv()

	return &Config{
		ServiceName: viper.GetString("SERVICE_NAME"),
		Port:        viper.GetString("PORT"),
		NatsURL:     viper.GetString("NATS_URL"),
		LogLevel:    viper.GetString("LOG_LEVEL"),
		OllamaURL:   viper.GetString("OLLAMA_URL"),
	}
}
