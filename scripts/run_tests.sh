#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Go to the root directory of the project
cd "$(dirname "$0")/.."

# Run go fmt
echo "Running go fmt..."
go fmt ./...

# Run go vet
echo "Running go vet..."
go vet ./...

# Run tests
echo "Running tests..."
go test -v -race -cover ./...

echo "All tests completed successfully."