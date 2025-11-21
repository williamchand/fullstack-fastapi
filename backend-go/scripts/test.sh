#!/bin/bash

# Test script with coverage and reporting

set -e

echo "Running tests..."

# Run tests with coverage
make test-coverage

# Show coverage report
echo "Coverage report:"
go tool cover -func=coverage.out | grep total:

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
echo "Detailed coverage report: coverage.html"

# Run linting
echo "Running linting..."
make lint

echo "All checks passed!"