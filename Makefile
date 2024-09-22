# Makefile for timetracker
PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
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
	CGO_ENABLED=1 go build -o $(SERVER_BIN) $(SERVER_DIR)/main.go

# Build client binary
build-client:
	@echo "Building client..."
	@mkdir -p bin
	CGO_ENABLED=1 go build -o $(CLIENT_BIN) $(CLIENT_DIR)/main.go

# Lint code by formatting and vetting
	
lint:
	@echo "Installing golangci-lint" && go get github.com/golangci/golangci-lint/cmd/golangci-lint && go install github.com/golangci/golangci-lint/cmd/golangci-lint
	go mod tidy
	@echo "Running $@"
	${GOPATH}/bin/golangci-lint run -c .golangci.yml

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

