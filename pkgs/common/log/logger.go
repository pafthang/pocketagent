package log

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/pafthang/pocketagent/pkgs/common/trace"
	"github.com/spf13/viper"
)

// Logger wraps the standard library logger for service use.
type Logger struct {
	*log.Logger
}

// New creates a prefixed standard logger.
func New(service string) *Logger {
	return &Logger{log.New(os.Stdout, "["+service+"] ", log.LstdFlags)}
}

// NewSlog creates a structured logger with configurable level.
func NewSlog(service string) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel(),
	})).With("service", service)
}

func logLevel() slog.Level {
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

// WithCorrelation adds correlation ID to log records.
func WithCorrelation(logger *slog.Logger, ctx context.Context) *slog.Logger {
	if corrID := trace.GetCorrelationID(ctx); corrID != "" {
		return logger.With("correlation_id", corrID)
	}
	return logger
}