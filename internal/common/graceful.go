package common

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GracefulShutdown waits for interrupt and shuts down
func GracefulShutdown(cancel context.CancelFunc, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()
	time.Sleep(timeout)
}
