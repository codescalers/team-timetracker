# Makefile for timetracker

# Directories
SERVER_DIR=cmd/server
CLIENT_DIR=cmd/client

# Binaries
SERVER_BIN=bin/timetrackerd
CLIENT_BIN=bin/timetracker-cli

# Default target
.PHONY: all build lint test clean run-server

all: build

# Build both server and client
build: build-server build-client

# Build server binary
build-server:
	@echo "Building server..."
	@mkdir -p bin
	go build -o $(SERVER_BIN) $(SERVER_DIR)/main.go

# Build client binary
build-client:
	@echo "Building client..."
	@mkdir -p bin
	go build -o $(CLIENT_BIN) $(CLIENT_DIR)/main.go

# Format code using go fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code using go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Lint code by formatting and vetting
lint: fmt vet

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin

# Run server
run-server: build-server
	@echo "Running server..."
	$(SERVER_BIN) -config=./team-timetracker-server.sample.json

