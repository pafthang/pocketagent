package common

import (
	"context"
	"time"
)

// RetryConfig holds retry settings
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
	Backoff     float64
}

// DefaultRetryConfig returns sensible defaults
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Backoff:     2.0,
	}
}

// Retry executes fn with exponential backoff
func Retry(ctx context.Context, fn func() error, cfg RetryConfig) error {
	var lastErr error

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := fn(); err != nil {
			lastErr = err
			time.Sleep(cfg.Delay)
			cfg.Delay = time.Duration(float64(cfg.Delay) * cfg.Backoff)
			continue
		}
		return nil
	}

	return lastErr
}
