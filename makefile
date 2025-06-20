.PHONY: build run test tidy clean

env:
	@cat .env || echo "No .env file found"

# Build CLI scraper
build-cli:
	go build -o bin/udccli ./cmd/udccli

# Build REST API server
build-server:
	go build -o bin/server ./cmd/server

build-Autopipeline:
	go build -o bin/Autopipeline ./cmd/Autopipeline

build-Webserver:
	go build -o bin/Webserver ./cmd/Webserver

# Build both
build: build-cli build-server build-Autopipeline build-Webserver

# Run server
run-server: build-server
	DB_PATH=tags.db DATA_DIR =./data SERVER_PORT=8080 ./bin/server

# Run scraper
run-cli: build-cli
	./bin/udccli

# Run tests
test:
	go test ./pkg/...

# Tidy dependencies
tidy:
	go mod tidy

# Clean all builds
clean:
	rm -rf bin/

bootstrap:
	go run scripts/bootstrap_import.go