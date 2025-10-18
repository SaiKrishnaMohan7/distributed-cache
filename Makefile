APP_NAME := distributed-cache
P2P_NAME := p2p-node
BIN_DIR := bin

.PHONY: build run test clean build-p2p run-p2p fmt lint tidy test-race test-coverage lint-fix check help

build:
	@echo "🔨 Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/cache_server

build-p2p:
	@echo "🔨 Building $(P2P_NAME)..."
	@mkdir -p $(BIN_DIR)
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

test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

clean:
	@echo "🧹 Cleaning up binaries..."
	rm -rf $(BIN_DIR)/*
	rm -f coverage.out coverage.html

fmt:
	@echo "🎨 Formatting code..."
	gofmt -s -w .
	goimports -w .

lint:
	@echo "🔍 Running linters..."
	golangci-lint run

lint-fix:
	@echo "🔧 Running linters with auto-fix..."
	golangci-lint run --fix

tidy:
	@echo "📦 Tidying go.mod/go.sum..."
	go mod tidy

check: fmt lint test
	@echo "✅ All checks passed!"

help:
	@echo "📌 Available targets:"
	@echo "  build         Build the cache-server"
	@echo "  build-p2p     Build the p2p node"
	@echo "  run           Run the cache-server"
	@echo "  run-p2p       Run the p2p node"
	@echo "  test          Run tests"
	@echo "  test-race     Run tests with race detector"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  fmt           Format code"
	@echo "  lint          Run code linters"
	@echo "  lint-fix      Run linters with auto-fix"
	@echo "  tidy          Clean and update go.mod/go.sum"
	@echo "  check         Run fmt, lint, and test"
	@echo "  clean         Remove build artifacts"