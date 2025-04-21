#!/bin/bash
set -e

# Database configuration
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-nutrimatch}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}

# Create database if it doesn't exist
echo "Creating database $DB_NAME if it doesn't exist..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME"

# Run migrations
echo "Running migrations..."
migrate -path ./migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up

echo "Database setup completed successfully!"
