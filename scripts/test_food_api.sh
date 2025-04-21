#!/bin/bash
set -e

# Base URL for the API
BASE_URL="http://localhost:8080/api/v1"

# Test the food endpoints
echo "Testing food endpoints..."

# Test listing foods
echo -e "\nListing foods:"
curl -s "${BASE_URL}/foods?page=1&limit=5" | jq

# Test searching foods
echo -e "\nSearching for 'apple':"
curl -s "${BASE_URL}/foods?q=apple&page=1&limit=5" | jq

# Test getting food by category
echo -e "\nGetting foods by category 'fruit':"
curl -s "${BASE_URL}/foods/category/fruit?page=1&limit=5" | jq

# Test getting a specific food (replace with a valid ID from your database)
echo -e "\nGetting a specific food:"
curl -s "${BASE_URL}/foods/fd_123456" | jq

echo -e "\nTests completed!"
