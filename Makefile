# Makefile for Tech Docs AI Application

.PHONY: help build build-server build-worker clean up down test scrape

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build both server and worker binaries"
	@echo "  build-server - Build only the server binary"
	@echo "  build-worker - Build only the worker binary"
	@echo "  up           - Start all services with Docker Compose"
	@echo "  down         - Stop all services"
	@echo "  clean        - Clean up build artifacts"
	@echo "  test         - Run tests"
	@echo "  scrape       - Scrape a W3Schools URL"

# Build both server and worker
build: build-server build-worker

# Build server binary
build-server:
	@echo "Building server binary..."
	@go build -o server cmd/server/main.go

# Build worker binary
build-worker:
	@echo "Building worker binary..."
	@go build -o worker cmd/worker/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f server worker
	@go clean

# Start all services
up:
	@echo "Starting Tech Docs AI services..."
	@echo "Starting Docker Compose services..."
	@docker compose up -d
	@echo "Waiting for Ollama to be ready..."
	@sleep 10
	@echo "Pulling AI Model..."
	@docker compose exec ollama ollama pull tinyllama:latest
	@docker compose exec ollama ollama pull nomic-embed-text
	@docker compose exec -d ollama ollama serve
	@echo "Services started! Access the application at http://localhost"

# Stop all services
down:
	@echo "Stopping services..."
	@docker compose down

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Example scrape command
scrape:
	@echo "Example: curl -X POST http://localhost/api/v1/scrape \\"
	@echo "  -H 'Content-Type: application/json' \\"
	@echo "  -d '{\"url\": \"https://www.w3schools.com/html/html_intro.asp\", \"category\": \"HTML\", \"tags\": [\"tutorial\", \"basics\"]}'"

# Example chat command
chat:
	@echo "Example: curl -X POST http://localhost/api/v1/chat \\"
	@echo "  -H 'Content-Type: application/json' \\"
	@echo "  -d '{\"message\": \"What is HTML?\"}'"

# Example search command
search:
	@echo "Example: curl 'http://localhost/api/v1/documents/search?q=HTML&limit=5'"
