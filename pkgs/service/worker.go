package service

import "github.com/pafthang/pocketagent/pkgs/service/worker"

// Worker is a base NATS JetStream worker lifecycle.
type Worker = worker.Worker

// NewWorker connects to NATS and prepares a consumer.
func NewWorker(name, natsURL, logLevel string) (*Worker, error) {
	return worker.New(name, natsURL, logLevel)
}