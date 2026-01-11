.PHONY: all build run dev clean test lint fmt tidy docker-up docker-down docker-logs help

APP_NAME := fiozap
BINARY := bin/$(APP_NAME)
MAIN := ./cmd/server
GO := /snap/bin/go

all: build

## Build
build:
	@echo "Building $(APP_NAME)..."
	@$(GO) build -o $(BINARY) $(MAIN)

build-linux:
	@echo "Building $(APP_NAME) for Linux..."
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY)-linux $(MAIN)

## Run
run: build
	@./$(BINARY)

dev:
	@$(GO) run $(MAIN)

## Dependencies
tidy:
	@$(GO) mod tidy

download:
	@$(GO) mod download

## Code Quality
fmt:
	@$(GO) fmt ./...

lint:
	@golangci-lint run ./...

vet:
	@$(GO) vet ./...

## Testing
test:
	@$(GO) test -v ./...

test-cover:
	@$(GO) test -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html

## Docker
docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

docker-logs:
	@docker-compose logs -f

docker-ps:
	@docker-compose ps

## Database
db-reset:
	@docker-compose down -v
	@docker-compose up -d postgres
	@echo "Waiting for PostgreSQL..."
	@sleep 3

## Clean
clean:
	@rm -rf bin/
	@rm -f coverage.out coverage.html

## Setup
setup: docker-up
	@cp -n .env.example .env 2>/dev/null || true
	@echo "Setup complete. Edit .env as needed."

## Help
help:
	@echo "FioZap Makefile Commands:"
	@echo ""
	@echo "  make build        - Build the application"
	@echo "  make build-linux  - Build for Linux amd64"
	@echo "  make run          - Build and run the application"
	@echo "  make dev          - Run without building (go run)"
	@echo ""
	@echo "  make tidy         - Run go mod tidy"
	@echo "  make download     - Download dependencies"
	@echo ""
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"
	@echo "  make vet          - Run go vet"
	@echo ""
	@echo "  make test         - Run tests"
	@echo "  make test-cover   - Run tests with coverage"
	@echo ""
	@echo "  make docker-up    - Start docker services"
	@echo "  make docker-down  - Stop docker services"
	@echo "  make docker-logs  - View docker logs"
	@echo "  make docker-ps    - List docker services"
	@echo ""
	@echo "  make db-reset     - Reset database (destroys data)"
	@echo "  make setup        - Initial setup (docker + .env)"
	@echo "  make clean        - Remove build artifacts"
