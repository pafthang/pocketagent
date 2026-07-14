package httpsrv

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/common"
)

// Server is a base HTTP microservice.
type Server struct {
	Echo       *echo.Echo
	Log        *slog.Logger
	Name       string
	listenAddr string
}

// New creates an HTTP server with common middleware.
func New(name, listenAddr, logLevel string) *Server {
	e := echo.New()
	e.HideBanner = true

	common.SetupMiddleware(e, name)

	if logLevel != "" {
		_ = os.Setenv("LOG_LEVEL", logLevel)
	}

	return &Server{
		Echo:       e,
		Log:        common.NewSlogLogger(name),
		Name:       name,
		listenAddr: listenAddr,
	}
}

// ListenAddr returns the HTTP listen address.
func (s *Server) ListenAddr() string {
	return s.listenAddr
}

// Start starts the server with graceful shutdown.
func (s *Server) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = common.ShutdownTelemetry(shutdownCtx)
		if err := s.Echo.Shutdown(shutdownCtx); err != nil {
			s.Log.Warn("shutdown failed", "error", err)
		}
	}()

	s.Log.Info("starting", "addr", s.ListenAddr())
	err := s.Echo.Start(s.ListenAddr())
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// AddHealth adds /health with optional dependency checks.
func (s *Server) AddHealth(deps common.Deps) {
	deps.Service = s.Name
	s.Echo.GET("/health", common.HealthHandler(deps))
}

// AddMetrics exposes Prometheus metrics at /metrics.
func (s *Server) AddMetrics() {
	s.Echo.GET("/metrics", echo.WrapHandler(common.MetricsHandler()))
}