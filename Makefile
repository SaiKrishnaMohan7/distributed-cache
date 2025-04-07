APP_NAME := distributed-cache
P2P_NAME := p2p-node
BIN_DIR := bin

.PHONY: build run test clean build-p2p run-p2p fmt lint tidy test-race help

build:
	@echo "🔨 Building $(APP_NAME)..."
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/cache_server

build-p2p:
	@echo "🔨 Building $(P2P_NAME)..."
	go build -o $(BIN_DIR)/$(P2P_NAME) ./cmd/p2p-node

run: build
	@echo "🚀 Running $(APP_NAME)..."
	./$(BIN_DIR)/$(APP_NAME)

run-p2p: build-p2p
	@echo "🚀 Running $(P2P_NAME)..."
	./$(BIN_DIR)/$(P2P_NAME)

test:
	@echo "🧪 Running tests..."
	go test ./...

test-race:
	@echo "🧪 Running tests with race detector..."
	go test -race ./...

clean:
	@echo "🧹 Cleaning up binaries..."
	rm -rf $(BIN_DIR)/*

fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

lint:
	@echo "🔍 Running linters..."
	golangci-lint run

tidy:
	@echo "📦 Tidying go.mod/go.sum..."
	go mod tidy

help:
	@echo "📌 Available targets:"
	@echo "  build        Build the cache-server"
	@echo "  build-p2p    Build the p2p node"
	@echo "  run          Run the cache-server"
	@echo "  run-p2p      Run the p2p node"
	@echo "  test         Run tests"
	@echo "  test-race    Run tests with race detector"
	@echo "  fmt          Format code"
	@echo "  lint         Run code linters"
	@echo "  tidy         Clean and update go.mod/go.sum"
	@echo "  clean        Remove build artifacts"
