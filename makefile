include .env

.PHONY: run build test lint clean env dev 
.ALL: run

APP_NAME=rinha-backend-2025
CMD_DIR=cmd/server
ENV_FILE=.env
OUT_DIR=bin

run:
	@echo ">> Running $(APP_NAME)..."
	@ENV_FILE=$(ENV_FILE) go run $(CMD_DIR)/main.go

build:
	@echo ">> Building $(APP_NAME)..."
	@go build -o $(OUT_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

docs:
	@echo ">> Generating docs..."
	@swag init -g $(CMD_DIR)/main.go -o internal/docs

test:
	@echo ">> Running tests..."
	@go test ./...

lint:
	@echo ">> Linting code..."
	@go vet ./...

clean:
	@echo ">> Cleaning build artifacts..."
	@rm -f $(APP_NAME)

env:
	@echo ">> Showing environment variables from $(ENV_FILE)"
	@cat $(ENV_FILE)

dev:
	@echo ">> Starting development environment with Docker Compose..."
	@docker-compose -f docker-compose.local.yml up --build --remove-orphans --force-recreate -d
