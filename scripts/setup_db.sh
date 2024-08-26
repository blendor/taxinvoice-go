#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Database configuration
DB_USER="your_username"
DB_PASSWORD="your_password"
DB_NAME="tax_invoice_app"

# Check if psql is installed
if ! command -v psql &> /dev/null
then
    echo "psql could not be found. Please install PostgreSQL."
    exit 1
fi

# Create database
echo "Creating database..."
createdb -U "$DB_USER" "$DB_NAME"

# Run migrations
echo "Running migrations..."
migrate -path ./internal/db/migrations -database "postgresql://$DB_USER:$DB_PASSWORD@localhost/$DB_NAME?sslmode=disable" up

echo "Database setup completed successfully."