package client

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

// Connect opens a NATS connection, JetStream context, and ensures required streams.
func Connect(url string) (*nats.Conn, nats.JetStreamContext, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, nil, fmt.Errorf("connect nats: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, nil, fmt.Errorf("jetstream: %w", err)
	}

	if err := EnsureStreams(js); err != nil {
		nc.Close()
		return nil, nil, err
	}

	return nc, js, nil
}
