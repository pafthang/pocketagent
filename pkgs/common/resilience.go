package common

import (
	"context"
	"time"

	"github.com/pafthang/pocketagent/pkgs/common/resilience"
)

var ErrCircuitOpen = resilience.ErrCircuitOpen

type CircuitState = resilience.CircuitState
type CircuitBreaker = resilience.CircuitBreaker
type RetryConfig = resilience.RetryConfig
type HTTPStatusError = resilience.HTTPStatusError

const (
	CircuitClosed   = resilience.CircuitClosed
	CircuitOpen     = resilience.CircuitOpen
	CircuitHalfOpen = resilience.CircuitHalfOpen
)

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return resilience.NewCircuitBreaker(threshold, timeout)
}
func DefaultRetryConfig() RetryConfig { return resilience.DefaultRetryConfig() }
func NewHTTPStatusError(statusCode int, body string) error {
	return resilience.NewHTTPStatusError(statusCode, body)
}
func IsRetryable(err error) bool { return resilience.IsRetryable(err) }
func Retry(ctx context.Context, fn func() error, cfg RetryConfig) error {
	return resilience.Retry(ctx, fn, cfg)
}