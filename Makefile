.PHONY: help run build test fmt tidy clean docker-build docker-push install

# Variables
BINARY_NAME=axinova-mcp-server
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the MCP server
	@echo "Starting MCP server..."
	go run ./cmd/server

build: ## Build the binary
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/server
	@echo "Binary built: bin/${BINARY_NAME}"

install: build ## Install the binary to /usr/local/bin
	@echo "Installing ${BINARY_NAME}..."
	cp bin/${BINARY_NAME} /usr/local/bin/
	@echo "Installed to /usr/local/bin/${BINARY_NAME}"

test: ## Run tests
	go test -v ./...

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

tidy: ## Tidy and verify go modules
	go mod tidy
	go mod verify

clean: ## Remove build artifacts
	rm -rf bin/
	go clean

docker-build: ## Build Docker image
	docker build -t axinova-mcp-server:${VERSION} -t axinova-mcp-server:latest .

docker-push: ## Push Docker image to registry
	docker tag axinova-mcp-server:${VERSION} ghcr.io/axinova-ai/axinova-mcp-server:${VERSION}
	docker tag axinova-mcp-server:latest ghcr.io/axinova-ai/axinova-mcp-server:latest
	docker push ghcr.io/axinova-ai/axinova-mcp-server:${VERSION}
	docker push ghcr.io/axinova-ai/axinova-mcp-server:latest

# Development helpers
dev: ## Run in development mode with hot reload (requires air)
	air

.DEFAULT_GOAL := help
