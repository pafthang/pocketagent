.PHONY: up down logs build test

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

# Development helpers
run-gateway:
	go run ./services/api-gateway

run-execution:
	go run ./services/execution-service
