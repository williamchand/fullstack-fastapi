#!/bin/bash

# Development environment setup script

set -e

echo "Setting up development environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

# Start PostgreSQL container
echo "Starting PostgreSQL container..."
docker run -d \
    --name myapp-postgres \
    -e POSTGRES_USER=myapp \
    -e POSTGRES_PASSWORD=myapp \
    -e POSTGRES_DB=myapp \
    -p 5432:5432 \
    postgres:15-alpine

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until docker exec myapp-postgres pg_isready -U myapp; do
    sleep 1
done

# Run migrations
echo "Running migrations..."
export DB_URL="postgres://myapp:myapp@localhost:5432/myapp?sslmode=disable"
make migrate-up

echo "Development environment is ready!"
echo "Database is running on localhost:5432"
echo "Run 'make run' to start the application"