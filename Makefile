.PHONY: up down logs build

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

build:
	go build ./...
