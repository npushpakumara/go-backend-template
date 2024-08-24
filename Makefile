# Variables
APP_NAME := server
GO_FILES := $(shell find . -name '*.go' -type f)
BUILD_DIR := build
BINARY := $(BUILD_DIR)/$(APP_NAME)
MODULE := $(shell go list -m)

# Default target
.PHONY: all
all: build

# Help target
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all           Build the application (default)"
	@echo "  clean         Clean up the build directory"
	@echo "  verify        Run dependency verification"
	@echo "  fmt           Run go fmt on all source files"
	@echo "  vet           Run go vet on all source files"
	@echo "  lint          Run go lint on all source files"
	@echo "  test          Run tests"
	@echo "  build         Build the application"
	@echo "  run           Run the application"
	@echo "  validate      Run fmt, vet, lint, test, and build"
	@echo "  docker        Build a Docker image for the application"
	@echo "  docker-run    Run the application in a Docker container"
	@echo "  help          Show this help message"

# Clean up the build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Run go fmt on all source files
.PHONY: fmt
fmt:
	go fmt ./...

# Run go vet on all source files
.PHONY: vet
vet:
	go vet ./...

# Run go dependency verification
.PHONY: verify
verify:
	go mod verify

# Run go lint (you need to install golint first)
.PHONY: lint
lint:
	golangci-lint run ./...

# Run tests
.PHONY: test
test:
	go test -v ./...

# Build the application
.PHONY: build
build: clean
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY) $(MODULE)/cmd/$(APP_NAME)

# Run the application
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	@$(BINARY)

# Run all (fmt, vet, lint, test, build)
.PHONY: validate
validate: fmt vet verify lint test build