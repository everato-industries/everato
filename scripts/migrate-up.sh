#!/bin/bash

set -e
source .env

migrate -database $DB_URL -path ./internal/db/migrations up
