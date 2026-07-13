package common

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

// NewSlogLogger creates structured logger with configurable level
func NewSlogLogger(service string) *slog.Logger {
	level := getLogLevel()

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})).With("service", service)
}

func getLogLevel() slog.Level {
	lvl := viper.GetString("LOG_LEVEL")
	switch lvl {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// LogWithCorrelation adds correlation ID to log
func LogWithCorrelation(logger *slog.Logger, ctx context.Context) *slog.Logger {
	if corrID := GetCorrelationID(ctx); corrID != "" {
		return logger.With("correlation_id", corrID)
	}
	return logger
}
