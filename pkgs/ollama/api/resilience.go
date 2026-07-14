package api

import (
	"context"
	"time"

	"github.com/pafthang/pocketagent/pkgs/common"
)

// Breaker is a circuit breaker for Ollama calls.
type Breaker = common.CircuitBreaker

func (c *Client) breaker() *Breaker {
	if c == nil {
		return nil
	}
	if c.Breaker == nil {
		c.Breaker = common.NewCircuitBreaker(5, 30*time.Second)
	}
	return c.Breaker
}

func (c *Client) callWithResilience(ctx context.Context, fn func() error) error {
	return c.breaker().Call(func() error {
		return common.Retry(ctx, fn, common.DefaultRetryConfig())
	})
}