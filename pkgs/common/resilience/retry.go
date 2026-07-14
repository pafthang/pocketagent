package resilience

import (
	"context"
	"errors"
	"net"
	"time"
)

// RetryConfig holds retry settings.
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
	Backoff     float64
}

// DefaultRetryConfig returns sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Backoff:     2.0,
	}
}

// HTTPStatusError carries a non-2xx HTTP response.
type HTTPStatusError struct {
	StatusCode int
	Body       string
}

func (e *HTTPStatusError) Error() string {
	if e.Body != "" {
		return e.Body
	}
	return "http error"
}

// NewHTTPStatusError builds a retry-classified HTTP error.
func NewHTTPStatusError(statusCode int, body string) error {
	return &HTTPStatusError{StatusCode: statusCode, Body: body}
}

// IsRetryable reports whether an error is worth retrying.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if errors.Is(err, ErrCircuitOpen) {
		return false
	}

	var httpErr *HTTPStatusError
	if errors.As(err, &httpErr) {
		switch httpErr.StatusCode {
		case 429, 500, 502, 503, 504:
			return true
		default:
			return httpErr.StatusCode >= 500
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout() || netErr.Temporary()
	}

	var opErr *net.OpError
	return errors.As(err, &opErr)
}

// Retry executes fn with exponential backoff for retryable errors.
func Retry(ctx context.Context, fn func() error, cfg RetryConfig) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}

	var lastErr error
	delay := cfg.Delay

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := fn(); err != nil {
			lastErr = err
			if attempt == cfg.MaxAttempts || !IsRetryable(err) {
				return err
			}
			if err := sleepWithContext(ctx, delay); err != nil {
				return err
			}
			delay = time.Duration(float64(delay) * cfg.Backoff)
			continue
		}
		return nil
	}

	return lastErr
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}