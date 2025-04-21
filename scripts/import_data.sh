#!/bin/bash
set -e

# Database configuration
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-nutrimatch}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_CONN_STRING="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Data file path
DATA_FILE=${1:-"./opennutrition_foods.tsv"}

if [ ! -f "$DATA_FILE" ]; then
    echo "Error: Data file not found at $DATA_FILE"
    echo "Usage: $0 [path/to/opennutrition_foods.tsv]"
    exit 1
fi

echo "Importing data from $DATA_FILE into database..."
# Using the correct import script path
go run ./scripts/import_food_data.go "$DATA_FILE"

echo "Data import completed!"
