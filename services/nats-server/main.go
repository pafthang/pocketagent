package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats-server/v2/server"
)

func main() {
	var port = flag.Int("port", 4222, "NATS port")
	flag.Parse()

	s := server.New(&server.Options{
		Port: *port,
		JetStream: true,
		StoreDir: "/data",
	})

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Embedded NATS with JetStream started on port %d", *port)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	s.Shutdown()
}
