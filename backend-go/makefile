# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGENERATE=$(GOCMD) generate
GOMOD=$(GOCMD) mod
BINARY_NAME=myapp
DOCKER_IMAGE=myapp:latest
TOOLS = \
	github.com/bufbuild/buf/cmd/buf \
	google.golang.org/protobuf/cmd/protoc-gen-go \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
	github.com/sqlc-dev/sqlc/cmd/sqlc \
	github.com/golang-migrate/migrate/v4/cmd/migrate \
	github.com/golangci/golangci-lint/cmd/golangci-lint \
	mvdan.cc/gofumpt
# Project structure
PROTO_DIR=api/proto
GEN_PROTO_DIR=gen/proto
DB_QUERIES_DIR=db/queries
DB_MIGRATIONS_DIR=db/migrations
SQLC_OUT_DIR=internal/infrastructure/database/dbgen

# Tools
BUF=$(shell which buf)
SQLC=$(shell which sqlc)
MIGRATE=$(shell which migrate)
GOLANGCI_LINT=$(shell which golangci-lint)
SWAG=$(shell which swag)

.PHONY: all build clean test generate proto sqlc lint migrate help

all: clean build

## Build the application
build: generate
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o bin/$(BINARY_NAME) ./cmd/server

## Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -rf $(GEN_PROTO_DIR)
	rm -rf $(SQLC_OUT_DIR)/*.go

## Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## Generate all code (protobuf, sqlc, etc.)
generate: proto sqlc

## Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	@if [ -z "$(BUF)" ]; then \
		echo "Error: buf is not installed. Run 'make install-tools' first"; \
		exit 1; \
	fi
	cd $(PROTO_DIR) && buf mod update
	cd $(PROTO_DIR) && buf generate

## Generate SQLC code
sqlc:
	@echo "Generating SQLC code..."
	@if [ -z "$(SQLC)" ]; then \
		echo "Error: sqlc is not installed. Run 'make install-tools' first"; \
		exit 1; \
	fi
	$(SQLC) generate

## Run linter
lint:
	@echo "Running linter..."
	@if [ -z "$(GOLANGCI_LINT)" ]; then \
		echo "Error: golangci-lint is not installed. Run 'make install-tools' first"; \
		exit 1; \
	fi
	$(GOLANGCI_LINT) run

## Fix linting issues
lint-fix:
	@echo "Fixing linting issues..."
	$(GOLANGCI_LINT) run --fix

## Create new migration
migrate-create:
	@if [ -z "$(MIGRATE)" ]; then \
		echo "Error: migrate is not installed. Run 'make install-tools' first"; \
		exit 1; \
	fi
	@read -p "Enter migration name: " name; \
	$(MIGRATE) create -ext sql -dir $(DB_MIGRATIONS_DIR) -seq $${name}

## Run migrations up
migrate-up:
	@if [ -z "$(MIGRATE)" ]; then \
		echo "Error: migrate is not installed. Run 'make install-tools' first"; \
		exit 1; \
	fi
	$(MIGRATE) -path $(DB_MIGRATIONS_DIR) -database "$(DB_URL)" up

## Run migrations down
migrate-down:
	@if [ -z "$(MIGRATE)" ]; then \
		echo "Error: migrate is not installed. Run 'make install-tools' first"; \
		exit 1; \
	fi
	$(MIGRATE) -path $(DB_MIGRATIONS_DIR) -database "$(DB_URL)" down

## Rollback one migration
migrate-down-1:
	@if [ -z "$(MIGRATE)" ]; then \
		echo "Error: migrate is not installed. Run 'make install-tools' first"; \
		exit 1; \
	fi
	$(MIGRATE) -path $(DB_MIGRATIONS_DIR) -database "$(DB_URL)" down 1

## Show migration status
migrate-status:
	@if [ -z "$(MIGRATE)" ]; then \
		echo "Error: migrate is not installed. Run 'make install-tools' first"; \
		exit 1; \
	fi
	$(MIGRATE) -path $(DB_MIGRATIONS_DIR) -database "$(DB_URL)" version

install-tools:
	@echo "Installing tools..."
	@for tool in $(TOOLS); do \
		echo "Installing $$tool..."; \
		go install $$tool@latest; \
	done

## Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

## Run the application
run: generate
	@echo "Running application..."
	$(GOCMD) run ./cmd/server

## Format code
fmt:
	@echo "Formatting code..."
	gofumpt -l -w .

## Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

## Vendor dependencies
vendor: tidy
	@echo "Vendoring dependencies..."
	$(GOMOD) vendor

## Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  generate      - Generate all code (protobuf, sqlc)"
	@echo "  proto         - Generate protobuf code"
	@echo "  sqlc          - Generate SQLC code"
	@echo "  lint          - Run linter"
	@echo "  lint-fix      - Fix linting issues"
	@echo "  migrate-create - Create new migration"
	@echo "  migrate-up    - Run migrations up"
	@echo "  migrate-down  - Run migrations down"
	@echo "  install-tools - Install required tools"
	@echo "  docker-build  - Build Docker image"
	@echo "  run           - Run the application"
	@echo "  fmt           - Format code"
	@echo "  tidy          - Tidy dependencies"
	@echo "  help          - Show this help"

# Environment targets
env-example:
	@echo "Copying .env.example to .env..."
	@cp .env.example .env
	@echo "Please edit .env with your configuration values"

generate-jwt-keys:
	@echo "Generating JWT RSA key pair..."
	@mkdir -p config/jwt
	@go run scripts/generate-jwt-keys/main.go

setup-dev: env-example generate-jwt-keys
	@echo "✅ Development environment setup complete!"
	@echo "📝 Please edit .env with your actual configuration values"

# Run with specific environment
run-dev:
	@ENV=development go run cmd/server/main.go

run-prod:
	@ENV=production go run cmd/server/main.go
.DEFAULT_GOAL := help