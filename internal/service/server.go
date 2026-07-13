package service

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/common"
)

// Server is a base microservice template
type Server struct {
	Echo   *echo.Echo
	Config *common.Config
	Logger *common.Logger
	Name   string
}

func New(name string) *Server {
	cfg := common.LoadConfig()
	logger := common.NewLogger(name)

	e := echo.New()
	e.HideBanner = true

	common.SetupMiddleware(e, name)

	return &Server{
		Echo:   e,
		Config: cfg,
		Logger: logger,
		Name:   name,
	}
}

// Start starts the server with graceful shutdown
func (s *Server) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go common.GracefulShutdown(cancel, 10*time.Second)

	s.Logger.Printf("Starting %s on port %s", s.Name, s.Config.Port)
	return s.Echo.Start(":" + s.Config.Port)
}

// AddHealth adds /health endpoint
func (s *Server) AddHealth() {
	s.Echo.GET("/health", common.HealthHandler)
}
