#!/bin/bash
set -e

echo "Installing Swagger dependencies..."

# Install swag CLI tool
go install github.com/swaggo/swag/cmd/swag@latest

# Install required packages
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/files
go get -u github.com/swaggo/swag

echo "Swagger dependencies installed successfully!"
