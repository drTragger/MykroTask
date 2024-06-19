# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif

# Variables
DB_DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Commands
.PHONY: all setup run migrate-up migrate-down

all: setup run

setup:
	@echo "Setting up the environment..."
	@go mod tidy

run:
	@echo "Running the server..."
	@go run main.go

migrate-up:
	@echo "Running migrations up..."
	migrate -database $(DB_DSN) -path db/migrations up

migrate-down:
	@echo "Running migrations down..."
	migrate -database $(DB_DSN) -path db/migrations down

migrate-force:
	@echo "Forcing migrations version..."
	migrate -database $(DB_DSN) -path db/migrations force $(VERSION)

migrate-create:
	@echo "Creating new migration..."
	migrate create -ext sql -dir db/migrations -seq $(NAME)