#!/bin/bash

set -e
source .env

# Check if .env is loaded correctly
if [ -z "$DB_URL" ]; then
  echo "Error: DB_URL not found in .env"
  echo "Please create a .env file with DB_URL=..."
  exit 1
fi

read -p ">>> Enter the migration version to force: " version

# Run migrations
migrate -database $DB_URL -path ./internal/db/migrations force $version

# Check the exit code
if [ $? -eq 0 ]; then
  echo "✅ Migration forced successfully to version: $version"
else
  echo "❌ Migration force failed"
  exit 1
fi
