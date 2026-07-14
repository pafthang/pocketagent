package worker

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/common"
)

// ServeHealth starts a background HTTP server for /health when healthPort is set.
func (w *Worker) ServeHealth(healthPort string, deps common.Deps) {
	if healthPort == "" {
		return
	}

	deps.Service = w.Name
	deps.NATS = w.NC
	if deps.JetStream == nil {
		deps.JetStream = w.JS
	}
	if deps.DLQWarnCount == 0 {
		deps.DLQWarnCount = 1
	}

	e := echo.New()
	e.HideBanner = true
	e.GET("/health", common.HealthHandler(deps))
	e.GET("/metrics", echo.WrapHandler(common.MetricsHandler()))

	addr := ":" + healthPort
	w.healthServer = e
	w.healthAddr = addr

	go func() {
		w.Log.Info("health endpoint listening", "addr", addr)
		if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			w.Log.Warn("health server stopped", "error", err)
		}
	}()
}

func (w *Worker) shutdownHealth(timeout time.Duration) {
	if w.healthServer == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := w.healthServer.Shutdown(ctx); err != nil {
		w.Log.Warn("health shutdown failed", "error", err)
	}
}