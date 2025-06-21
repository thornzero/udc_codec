.PHONY: build run test tidy clean help

env: ## Print the contents of the .env file if it exists
	@cat .env || echo "No .env file found"

# Build CLI scraper
build-cli: ## Build the CLI scraper binary
	go build -o bin/udccli ./cmd/udccli

# Build REST API server
build-server: ## Build the REST API server binary
	go build -o bin/server ./cmd/server

build-autopipeline: ## Build the autopipeline tool binary
	go build -o bin/autopipeline ./cmd/autopipeline

build-webserver: ## Build the webserver binary
	go build -o bin/webserver ./cmd/webserver

# Build both
build: build-cli build-server build-autopipeline build-webserver ## Build all binaries (CLI, server, autopipeline, webserver)

# Run server
run-server: build-server ## Run the REST API server (requires tags.db and data directory)
	DB_PATH=tags.db DATA_DIR =./data SERVER_PORT=8080 ./bin/server

# Run scraper
run-cli: build-cli ## Run the CLI tool (udccli)
	./bin/udccli

# Run tests
test: build-cli ## Run all Go tests in the pkg directory
	go test ./pkg/...

test-verbose: build-cli ## Run all Go tests with verbose output
	go test -v ./pkg/...

test-udc: build-cli	## Run tests specifically for the UDC module
	go test -v ./pkg/udc/...

test-coverage: build-cli ## Run tests with coverage reporting
	go test -cover ./pkg/...

test-coverage-html: build-cli ## Run tests with HTML coverage report
	go test -coverprofile=coverage.out ./pkg/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-all: test-verbose test-coverage ## Run all tests with verbose output and coverage

# Tidy dependencies
tidy: ## Clean up and verify Go module dependencies
	go mod tidy

# Clean all builds
clean: ## Remove all built binaries from the bin directory
	rm -rf bin/

bootstrap: ## Run the bootstrap import script to initialize data
	go run scripts/bootstrap_import.go

update-udc: build-cli ## Update the UDC data by scraping the latest hierarchy into data/udc_full.yaml
	./bin/udccli scrape

help: ## Show this help message listing available targets
	@echo "Available targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Build process generates the following binaries in ./bin/:"
	@echo "  udccli         - Command-line tool for scraping and processing"
	@echo "  server         - REST API server for tag management"
	@echo "  autopipeline   - Automation pipeline tool"
	@echo "  webserver      - Web frontend server"
	@echo ""
	@echo "Next steps after building:"
	@echo "  - To run the REST API server:      ./bin/server (ensure tags.db and ./data exist)"
	@echo "  - To use the CLI tool:             ./bin/udccli"
	@echo "  - To run the web frontend:         ./bin/webserver"
	@echo "  - To initialize data:              make bootstrap"
	@echo ""
	@echo "See Makefile targets for more options, or run 'make help'."