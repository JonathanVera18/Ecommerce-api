.PHONY: build run test clean migrate seed docker-up docker-down

# Build the application
build:
	go build -o bin/ecommerce-api main.go

# Run the application
run:
	go run main.go

# Run with air for hot reloading (install with: go install github.com/cosmtrek/air@latest)
dev:
	air

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Database migrations
migrate-up:
	go run scripts/migrate.go up

migrate-down:
	go run scripts/migrate.go down

# Seed database
seed:
	go run scripts/seed.go

# Create admin user
create-admin:
	go run scripts/admin_user.go

# Docker commands
docker-build:
	docker build -t ecommerce-api .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (install golangci-lint first)
lint:
	golangci-lint run

# Generate swagger docs (install swag first: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	swag init -g main.go -o ./api

# Security scan (install gosec first: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
security:
	gosec ./...
