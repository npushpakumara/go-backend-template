# Go application variables
APP_NAME := server
MODULE_NAME := $(shell go list -m)
BUILD_DIR := build
SRC_DIR := .
GO_FILES := $(shell find $(SRC_DIR) -name '*.go')

# Docker variables
DOCKER_IMAGE := backend-api:latest

# Default target: build the Go application
all: build

# Build the Go application
build: $(BIN_DIR)/$(APP_NAME)

$(BIN_DIR)/$(APP_NAME): $(GO_FILES)
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME) $(MODULE_NAME)/cmd/server

# Clean the binary and other build artifacts
clean:
	@rm -rf $(BIN_DIR)

# Docker build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Docker Compose up
docker-compose-up:
	@docker compose up -d

# Docker Compose down
docker-compose-down:
	@docker compose down

# Lint the Go code
lint:
	@golangci-lint run

# Format the Go code
fmt:
	@go fmt ./...

# Detect race conditions
race:
	@go test -race ./...

# Help message
help:
	@echo "Usage:"
	@echo "  make                   Build the Go application (default target)"
	@echo "  make build             Build the Go application"
	@echo "  make clean             Clean the binary and other build artifacts"
	@echo "  make docker-build      Build the Docker image"
	@echo "  make docker-compose-up Start the application using Docker Compose"
	@echo "  make docker-compose-down Stop the application using Docker Compose"
	@echo "  make lint              Lint the Go code"
	@echo "  make fmt               Format the Go code"
	@echo "  make race              Detect race conditions"
	@echo "  make help              Show this help message"

.PHONY: all build clean docker-build docker-compose-up docker-compose-down lint fmt race help
