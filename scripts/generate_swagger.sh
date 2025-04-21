#!/bin/bash
set -e

echo "Generating Swagger documentation..."

# Install swag if not already installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate Swagger docs
cd "$(dirname "$0")/.." # Navigate to project root

# Clean up any existing docs
rm -f docs/docs.go docs/swagger.json docs/swagger.yaml

# Find the swag binary
SWAG_BIN="$HOME/go/bin/swag"
if [ ! -f "$SWAG_BIN" ]; then
    SWAG_BIN=$(which swag 2>/dev/null)
    if [ -z "$SWAG_BIN" ]; then
        echo "Error: swag command not found. Please install it with: go install github.com/swaggo/swag/cmd/swag@latest"
        exit 1
    fi
fi

# Generate new docs
"$SWAG_BIN" init -g cmd/api/main.go -o docs --parseDependency --parseInternal --parseDepth 10

echo "Swagger documentation generated successfully!"
echo "You can access the Swagger UI at http://localhost:8080/swagger/ui/ when the server is running."
