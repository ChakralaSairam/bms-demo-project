#!/bin/bash

echo "=========================================="
echo "🧪 Internal API Contract Testing"
echo "=========================================="

BASE_URL="http://localhost:8082"

echo -e "\n✅ TEST 1: Get booked seats for show 1"
curl -s ${BASE_URL}/internal/shows/1/booked-seats | jq

echo -e "\n✅ TEST 2: Get booked seats for show with no bookings"
curl -s ${BASE_URL}/internal/shows/2/booked-seats | jq

echo -e "\n❌ TEST 3: Get booked seats for non-existent show (should return 404)"
curl -s ${BASE_URL}/internal/shows/99999/booked-seats | jq

echo -e "\n📊 TEST 4: Verify database state"
echo "Confirmed bookings in database:"
docker exec spring-auth-db psql -U myuser -d auth_db -c "SELECT show_id, seat_id, status FROM bookings WHERE status = 'CONFIRMED';"

echo -e "\n=========================================="
echo "✅ Testing Complete"
echo "=========================================="
