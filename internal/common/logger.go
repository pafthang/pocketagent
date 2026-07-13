package common

import (
	"context"
	"log/slog"
	"os"
)

// NewSlogLogger creates structured logger with slog
func NewSlogLogger(service string) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})).With("service", service)
}

// LogWithCorrelation adds correlation ID to log
func LogWithCorrelation(logger *slog.Logger, ctx context.Context) *slog.Logger {
	if corrID := GetCorrelationID(ctx); corrID != "" {
		return logger.With("correlation_id", corrID)
	}
	return logger
}
