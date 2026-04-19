.PHONY: build run dev test clean docker-build docker-up docker-down docker-dev-up docker-dev-down

# Application
APP_NAME=invitation-api
MAIN_PATH=./cmd/main.go

# Build the application
build:
	go build -o $(APP_NAME).exe $(MAIN_PATH)

# Run the application
run:
	go run $(MAIN_PATH)

# Run with hot-reload (requires air)
dev:
	air -c .air.toml

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f $(APP_NAME).exe
	rm -rf tmp/

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy

# Docker commands - Production
docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

# Docker commands - Development
docker-dev-up:
	docker-compose -f docker-compose.dev.yml up -d

docker-dev-down:
	docker-compose -f docker-compose.dev.yml down

docker-dev-logs:
	docker-compose -f docker-compose.dev.yml logs -f app

docker-dev-rebuild:
	docker-compose -f docker-compose.dev.yml up -d --build

# Database commands
db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down postgres

# Help
help:
	@echo "Available commands:"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application"
	@echo "  dev              - Run with hot-reload"
	@echo "  test             - Run tests"
	@echo "  clean            - Clean build artifacts"
	@echo "  fmt              - Format code"
	@echo "  lint             - Lint code"
	@echo "  deps             - Download dependencies"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-up        - Start production containers"
	@echo "  docker-down      - Stop production containers"
	@echo "  docker-logs      - View production logs"
	@echo "  docker-dev-up    - Start development containers"
	@echo "  docker-dev-down  - Stop development containers"
	@echo "  docker-dev-logs  - View development logs"
	@echo "  db-up            - Start PostgreSQL only"
	@echo "  db-down          - Stop PostgreSQL"
