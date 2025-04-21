.PHONY: setup migrate migrate-down generate build run run-swagger test clean swagger swagger-deps

# Default target
all: build

# Setup the project
setup:
	go mod tidy

# Run database migrations
migrate:
	go run scripts/migrate.go up

# Rollback database migrations
migrate-down:
	go run scripts/migrate_simple.go down

# Generate code
generate:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

# Build the application
build:
	go build -o bin/nutrimatch ./cmd/api

# Run the application
run:
	go run cmd/api/main.go

# Run the application with Swagger
run-swagger: swagger
	go run cmd/api/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf internal/repository/postgres/db/

# Import OpenNutrition dataset
import:
	@if [ -z "$(file)" ]; then \
		echo "Usage: make import file=<path-to-tsv-file>"; \
		exit 1; \
	fi
	go run scripts/import.go $(file)

# Generate Swagger documentation
swagger:
	./scripts/generate_swagger.sh

# Install Swagger dependencies
swagger-deps:
	./scripts/install_swagger_deps.sh
