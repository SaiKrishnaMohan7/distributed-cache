APP_NAME := distributed-cache

.PHONY: build run test clean

build:
	@echo "Building the application..."
	go build -o bin/$(APP_NAME) .

run:
	@echo "Running the application..."
	./bin/$(APP_NAME)

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning up..."
	rm -rf bin/

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linters..."
	golangci-lint run

run: build
	@echo "Running the application..."
	./bin/distributed-cache

help:
	@echo "Available targets:"
	@echo "  build      Build the application"
	@echo "  run        Build and run the application"
	@echo "  clean      Remove build artifacts"
	@echo "  test       Run tests"
	@echo "  fmt        Format code"
	@echo "  lint       Run code linters"
	@echo "  tidy       Update dependencies"
