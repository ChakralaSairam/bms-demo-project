#!/bin/bash

# Mock Inventory Service
# Run as: bash mock-inventory-service.sh

PORT=9090

echo "🚀 Starting Mock Inventory Service on port $PORT"

while true; do
  response=$(echo -e "HTTP/1.1 201 Created\r\nContent-Type: application/json\r\nContent-Length: 27\r\n\r\n{\"status\":\"seat_locked\"}" | nc -l $PORT)
  echo "[$(date)] Received request - Responding with 201 Created"
done
