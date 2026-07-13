package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/pafthang/pocketagent/internal/common"
)

func main() {
	var port = flag.Int("port", 4222, "NATS port")
	flag.Parse()

	logger := common.NewSlogLogger("nats-server")

	s := server.New(&server.Options{
		Port:      *port,
		JetStream: true,
		StoreDir:  "/data",
	})

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	logger.Info("Embedded NATS with JetStream started", "port", *port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.Shutdown()
	logger.Info("NATS server shutdown complete")
}
