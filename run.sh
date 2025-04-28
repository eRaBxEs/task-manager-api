#!/bin/bash

# Load environment variables from .env
set -o allexport
source .env
set +o allexport

# Run database migrations
goose -dir migrations postgres "host=$DB_HOST port=$DB_PORT dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD sslmode=disable" up

# Build the Go application
go build -o myapp cmd/myapp/main.go  # Adjust the path if needed

# Run the application
./myapp
