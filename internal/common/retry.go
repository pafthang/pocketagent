package common

import (
	"context"
	"time"
)

// RetryConfig holds retry settings
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
	Backoff     time.Duration // multiplier
}

// DefaultRetryConfig returns sensible defaults
defaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Backoff:     2,
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
