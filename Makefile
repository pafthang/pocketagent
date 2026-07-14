.PHONY: up down logs build test run run-ctrl

up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f

build:
	go build ./...

test:
	go test ./...

# Local dev — start all services (reads ./configs)
run run-ctrl:
	go run ./cmd/ctrl -config-dir=configs

# Individual services
run-gate:
	go run ./cmd/gate

run-exec:
	go run ./cmd/exec