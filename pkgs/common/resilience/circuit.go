package resilience

import (
	"errors"
	"sync"
	"time"
)

// ErrCircuitOpen is returned when the breaker rejects calls.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitState is the breaker state.
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// CircuitBreaker trips after consecutive failures and recovers after a timeout.
type CircuitBreaker struct {
	mu           sync.Mutex
	state        CircuitState
	failureCount int
	threshold    int
	timeout      time.Duration
	lastFailure  time.Time
}

// NewCircuitBreaker creates a breaker with failure threshold and open timeout.
func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	if threshold <= 0 {
		threshold = 5
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &CircuitBreaker{
		threshold: threshold,
		timeout:   timeout,
		state:     CircuitClosed,
	}
}

// Call runs fn when the breaker allows traffic.
func (cb *CircuitBreaker) Call(fn func() error) error {
	if err := cb.beforeCall(); err != nil {
		return err
	}

	err := fn()
	cb.afterCall(err)
	return err
}

func (cb *CircuitBreaker) beforeCall() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitOpen:
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = CircuitHalfOpen
			return nil
		}
		return ErrCircuitOpen
	default:
		return nil
	}
}

func (cb *CircuitBreaker) afterCall(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailure = time.Now()
		if cb.failureCount >= cb.threshold {
			cb.state = CircuitOpen
		}
		return
	}

	cb.failureCount = 0
	cb.state = CircuitClosed
}

// State returns the current breaker state.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}