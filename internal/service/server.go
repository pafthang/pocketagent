package service

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/common"
)

// BaseServer provides common server setup
type BaseServer struct {
	Echo   *echo.Echo
	Config *common.Config
	Logger *common.Logger
}

func NewBaseServer(name string) *BaseServer {
	cfg := common.LoadConfig()
	logger := common.NewLogger(name)

	e := echo.New()
	e.HideBanner = true

	return &BaseServer{
		Echo:   e,
		Config: cfg,
		Logger: logger,
	}
}

func (s *BaseServer) Start(port string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go common.GracefulShutdown(cancel, 5*time.Second)

	return s.Echo.Start(":" + port)
}
