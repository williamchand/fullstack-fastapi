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
    --name db \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=9Rdl_tmcO298f97J0FKzf1O1QFePclhGwLJwH-0A390 \
    -e POSTGRES_DB=app \
    -p 5433:5432 \
    postgres:18

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until docker exec db pg_isready -U app; do
    sleep 1
done

# Run migrations
echo "Running migrations..."
export DB_URL="postgres://postgres:9Rdl_tmcO298f97J0FKzf1O1QFePclhGwLJwH-0A390@localhost:5433/app?sslmode=disable"
make migrate-up

echo "Development environment is ready!"
echo "Database is running on localhost:5432"
echo "Run 'make run' to start the application"