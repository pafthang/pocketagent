module github.com/pafthang/pocketagent/services/execution-service

go 1.23

require (
	github.com/nats-io/nats.go v1.38.0
	github.com/pafthang/pocketagent/internal/models v0.0.0
	github.com/pafthang/pocketagent/internal/nats v0.0.0
)

replace (
	github.com/pafthang/pocketagent/internal/models => ../../internal/models
	github.com/pafthang/pocketagent/internal/nats => ../../internal/nats
)
