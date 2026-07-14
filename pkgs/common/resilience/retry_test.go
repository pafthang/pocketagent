package resilience

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"
)

func TestIsRetryableHTTPStatus(t *testing.T) {
	if !IsRetryable(NewHTTPStatusError(503, "unavailable")) {
		t.Fatal("503 should be retryable")
	}
	if IsRetryable(NewHTTPStatusError(400, "bad request")) {
		t.Fatal("400 should not be retryable")
	}
	if !IsRetryable(NewHTTPStatusError(429, "rate limit")) {
		t.Fatal("429 should be retryable")
	}
}

func TestRetryStopsOnPermanentError(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	err := Retry(ctx, func() error {
		attempts++
		return NewHTTPStatusError(400, "bad request")
	}, RetryConfig{MaxAttempts: 3, Delay: time.Millisecond})

	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestRetryHonorsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Retry(ctx, func() error {
		return &net.OpError{Err: errors.New("timeout"), Op: "dial"}
	}, DefaultRetryConfig())

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}