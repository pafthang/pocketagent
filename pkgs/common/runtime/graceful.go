package runtime

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

// WaitForSignal blocks until SIGINT or SIGTERM.
func WaitForSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// GracefulShutdown waits for interrupt and shuts down.
func GracefulShutdown(cancel context.CancelFunc, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()
	time.Sleep(timeout)
}

// GracefulNATSShutdown properly closes NATS connection.
func GracefulNATSShutdown(nc *nats.Conn, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if nc != nil {
		nc.Drain()
		time.Sleep(timeout)
		nc.Close()
	}
}