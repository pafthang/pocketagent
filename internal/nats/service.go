package nats

import (
	"os"
	"strconv"

	natssrv "github.com/nats-io/nats-server/v2/server"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func Run() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port == 0 {
		port = 4222
	}

	httpPort, err := strconv.Atoi(cfg.HTTPPort)
	if err != nil || httpPort == 0 {
		httpPort = 8222
	}

	logger := common.NewSlogLogger(cfg.Service)

	storeDir := cfg.StoreDir
	if storeDir == "" {
		storeDir = "data/nats"
	}
	_ = os.MkdirAll(storeDir, 0o755)

	s := natssrv.New(&natssrv.Options{
		Port:      port,
		HTTPPort:  httpPort,
		JetStream: true,
		StoreDir:  storeDir,
	})

	s.Start()

	logger.Info("Embedded NATS with JetStream started",
		"port", port,
		"http_port", httpPort,
		"store_dir", storeDir,
	)

	common.WaitForSignal()

	s.Shutdown()
	logger.Info("NATS server shutdown complete")
	return nil
}
