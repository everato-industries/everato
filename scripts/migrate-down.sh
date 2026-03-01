#!/bin/bash

set -e
source .env

migrate -database $DB_URL -path $MIGRATIONS_DIR down 1
