#!/bin/bash

set -e
source .env

read -p ">>> Enter the migration version to force: " version
migrate -database $DB_URL -path $MIGRATIONS_DIR force $version
