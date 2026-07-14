package common

import (
	"context"
	"log/slog"

	"github.com/pafthang/pocketagent/pkgs/common/log"
)

type Logger = log.Logger

func NewLogger(service string) *Logger     { return log.New(service) }
func NewSlogLogger(service string) *slog.Logger { return log.NewSlog(service) }
func LogWithCorrelation(logger *slog.Logger, ctx context.Context) *slog.Logger {
	return log.WithCorrelation(logger, ctx)
}