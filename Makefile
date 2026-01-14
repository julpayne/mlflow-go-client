.PHONY: help install-mlflow run-mlflow stop-mlflow deps build test clean fmt vet lint example

# Default target
.DEFAULT_GOAL := help

# Variables
MLFLOW_HOST ?= 127.0.0.1
MLFLOW_PORT ?= 5000
MLFLOW_BACKEND_STORE_URI ?= sqlite:///mlflow.db
MLFLOW_DEFAULT_ARTIFACT_ROOT ?= ./mlruns

## help: Show this help message
help:
	@echo "Available targets:"
	@echo ""
	@echo "  make install-mlflow    - Download and install MLflow server"
	@echo "  make run-mlflow         - Run MLflow server locally"
	@echo "  make stop-mlflow        - Stop MLflow server (if running)"
	@echo "  make deps               - Download all Go dependencies"
	@echo "  make build              - Build the Go client library"
	@echo "  make test               - Run tests"
	@echo "  make test-godog         - Run godog BDD tests (requires MLflow server)"
	@echo "  make test-godog-server  - Run godog BDD tests with automatic server"
	@echo "  make fmt                - Format Go code"
	@echo "  make vet                - Run go vet"
	@echo "  make lint               - Run golangci-lint (if installed)"
	@echo "  make example            - Build and run the example"
	@echo "  make clean              - Clean build artifacts"
	@echo ""
	@echo "MLflow server configuration (use as environment variables):"
	@echo "  MLFLOW_HOST                    - Server host (default: 127.0.0.1)"
	@echo "  MLFLOW_PORT                    - Server port (default: 5000)"
	@echo "  MLFLOW_BACKEND_STORE_URI       - Backend store URI (default: sqlite:///mlflow.db)"
	@echo "  MLFLOW_DEFAULT_ARTIFACT_ROOT    - Artifact root (default: ./mlruns)"
	@echo ""
	@echo "Examples:"
	@echo "  make install-mlflow"
	@echo "  make run-mlflow"
	@echo "  MLFLOW_PORT=8080 make run-mlflow"

## install-mlflow: Download and install MLflow server
install-mlflow:
	@echo "📥 Installing MLflow server..."
	@./scripts/download_mlflow.sh

## run-mlflow: Run MLflow server locally in the background
run-mlflow:
	@echo "🛑 Stopping MLflow server..."
	- @make stop-mlflow
	@echo "🚀 Starting MLflow server..."
	@echo "   Host: $(MLFLOW_HOST)"
	@echo "   Port: $(MLFLOW_PORT)"
	@echo "   Backend: $(MLFLOW_BACKEND_STORE_URI)"
	@echo "   Artifacts: $(MLFLOW_DEFAULT_ARTIFACT_ROOT)"
	@echo ""
	@MLFLOW_HOST=$(MLFLOW_HOST) \
	 MLFLOW_PORT=$(MLFLOW_PORT) \
	 MLFLOW_BACKEND_STORE_URI=$(MLFLOW_BACKEND_STORE_URI) \
	 MLFLOW_DEFAULT_ARTIFACT_ROOT=$(MLFLOW_DEFAULT_ARTIFACT_ROOT) \
	 ./scripts/run_mlflow.sh
	@echo "   MLflow server started in background. Use 'make stop-mlflow' to stop it."
	@echo "   You can use 'make test-godog-server' to run BDD tests with the server automatically."
	@echo ""
	@echo "   To use the server in your tests:"
	@echo "   - Set MLFLOW_TEST_URL=http://localhost:5000 in your test environment"
	@echo "   - Run tests with 'make test-godog'"
	@echo ""
	@echo "   To stop the server:"
	@echo "   make stop-mlflow"

## stop-mlflow: Stop MLflow server (if running)
stop-mlflow:
	./scripts/stop_mlflow.sh

## deps: Download all Go dependencies
deps:
	@echo "📦 Downloading Go dependencies..."
	@go mod download
	@echo "🔧 Tidying go.mod and go.sum..."
	@go mod tidy
	@echo "✅ Dependencies downloaded and synced"

## build: Build the Go client library
build:
	@echo "🔨 Building Go client library..."
	@go build ./...

## test: Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

## test-godog: Run godog BDD tests
test-godog:
	@echo "🧪 Running godog BDD tests..."
	@echo "   Make sure MLflow server is running or set MLFLOW_TEST_URL"
	@cd tests/features && go test -v -tags=godog

MLFLOW_TEST_URL ?= http://$(MLFLOW_HOST):$(MLFLOW_PORT)

## test-godog-server: Run godog tests with automatic server management
test-godog-server: run-mlflow ; $(info The MLflow server was successfully started)
	@echo "🧪 Running godog BDD tests with test server..."
	@rm -f mlflow.db tests/features/test_mlflow_${MLFLOW_PORT}.db
	@cd tests/features && MLFLOW_TEST_URL=${MLFLOW_TEST_URL} go test -v -tags=godog
	@make stop-mlflow

## fmt: Format Go code
fmt:
	@echo "📝 Formatting Go code..."
	@go fmt ./...

## vet: Run go vet
vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

## lint: Run golangci-lint (if installed)
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## example: Build and run the example
example:
	@echo "🏃 Running example..."
	@go run example/main.go

## clean: Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@go clean ./...
	@rm -f mlflow.db
	@rm -rf mlruns
	@rm -rf *.db

.PHONY: cls
cls:
	printf "\33c\e[3J"
