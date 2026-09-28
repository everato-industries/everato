#!/bin/bash

set -e
source .env

# Check if .env is loaded correctly
if [ -z "$DB_URL" ]; then
  echo "Error: DB_URL not found in .env"
  echo "Please create a .env file with DB_URL=..."
  exit 1
fi

echo "Database URL: $DB_URL"
echo "Migrations path: ./internal/db/migrations"

# Run migrations
migrate -database $DB_URL -path ./internal/db/migrations up

# Check the exit code
if [ $? -eq 0 ]; then
  echo "✅ Migrations applied successfully"
else
  echo "❌ Migration failed"
  exit 1
fi
