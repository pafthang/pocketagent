package common

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/pkgs/common/runtime"
)

func RunMain(service string, fn func() error) { runtime.RunMain(service, fn) }
func WaitForSignal()                        { runtime.WaitForSignal() }
func GracefulShutdown(cancel context.CancelFunc, timeout time.Duration) {
	runtime.GracefulShutdown(cancel, timeout)
}
func GracefulNATSShutdown(nc *nats.Conn, timeout time.Duration) {
	runtime.GracefulNATSShutdown(nc, timeout)
}
func WithTimeout(fn func() error, timeout time.Duration) error {
	return runtime.WithTimeout(fn, timeout)
}