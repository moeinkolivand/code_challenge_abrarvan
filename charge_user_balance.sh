#!/bin/bash

# Simple PostgreSQL Balance Update Script
# Usage: ./update_balance.sh <phone_number> <new_balance>

# Configuration - Edit these variables for your setup
CONTAINER_NAME="postgres"
DB_NAME="sms_gateway"
DB_USER="postgres"
DB_PASSWORD="postgres"
DB_HOST="localhost"

# Check arguments
if [ $# -ne 2 ]; then
    echo "Usage: $0 <phone_number> <new_balance>"
    echo "Example: $0 09332823692 500.50"
    exit 1
fi

PHONE_NUMBER="$1"
NEW_BALANCE="$2"

echo "Updating balance for phone: $PHONE_NUMBER to $NEW_BALANCE"

# Set password for psql
export PGPASSWORD="$DB_PASSWORD"

# Update balance using phone number filter
docker exec "$CONTAINER_NAME" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c \
"UPDATE users SET balance = balance + $NEW_BALANCE WHERE phone_number = '$PHONE_NUMBER';"

if [ $? -eq 0 ]; then
    echo "✓ Balance updated successfully!"


    docker exec "$CONTAINER_NAME" psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c \
    "SELECT id, phone_number, balance FROM users WHERE phone_number = '$PHONE_NUMBER';"
else
    echo "✗ Update failed!"
fi