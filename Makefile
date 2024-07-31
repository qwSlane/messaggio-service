.PHONY: build run migrate_up migrate_down

migrate_up:
		migrate -database postgres://postgres:postgres@localhost:5432/messages_db?sslmode=disable -path migrations up 1

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5432/messages_db?sslmode=disable -path migrations down 1

build:
	go build .\cmd\main.go

local:
	@echo Starting local docker compose
	docker-compose up -d

.DEFAULT_GOAL := build

