package resilience

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreakerOpensAfterThreshold(t *testing.T) {
	cb := NewCircuitBreaker(2, time.Minute)
	errFail := errors.New("boom")

	_ = cb.Call(func() error { return errFail })
	_ = cb.Call(func() error { return errFail })

	if cb.State() != CircuitOpen {
		t.Fatalf("expected open state, got %v", cb.State())
	}

	err := cb.Call(func() error { return nil })
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreakerRecoversAfterTimeout(t *testing.T) {
	cb := NewCircuitBreaker(1, 20*time.Millisecond)
	_ = cb.Call(func() error { return errors.New("fail") })

	time.Sleep(30 * time.Millisecond)

	if err := cb.Call(func() error { return nil }); err != nil {
		t.Fatalf("expected recovery, got %v", err)
	}
	if cb.State() != CircuitClosed {
		t.Fatalf("expected closed state, got %v", cb.State())
	}
}