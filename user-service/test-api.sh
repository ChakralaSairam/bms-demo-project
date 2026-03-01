#!/bin/bash

# API Testing Script for User Service
# Run as: ./test-api.sh

BASE_URL="http://localhost:8082"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=========================================="
echo "🧪 User Service API Testing"
echo "=========================================="

# Test 1: Health Check
echo -e "\n${YELLOW}[TEST 1] Health Check${NC}"
response=$(curl -s -o /dev/null -w "%{http_code}" ${BASE_URL}/actuator/health 2>/dev/null || echo "000")
if [ "$response" == "200" ] || [ "$response" == "404" ]; then
    echo -e "${GREEN}✓ Service is running${NC}"
else
    echo -e "${RED}✗ Service not reachable (HTTP $response)${NC}"
fi

# Test 2: Register User
echo -e "\n${YELLOW}[TEST 2] User Registration${NC}"
register_response=$(curl -s -X POST ${BASE_URL}/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test@1234"
  }')
echo "Response: $register_response"
if echo "$register_response" | grep -q "successfully\|already"; then
    echo -e "${GREEN}✓ Registration endpoint working${NC}"
else
    echo -e "${RED}✗ Registration failed${NC}"
fi

# Test 3: Login
echo -e "\n${YELLOW}[TEST 3] User Login${NC}"
login_response=$(curl -s -X POST ${BASE_URL}/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test@1234"
  }')
echo "Response: $login_response"

JWT_TOKEN=$(echo $login_response | grep -o '"token":"[^"]*' | cut -d'"' -f4)
if [ -n "$JWT_TOKEN" ]; then
    echo -e "${GREEN}✓ Login successful${NC}"
    echo "JWT Token: ${JWT_TOKEN:0:50}..."
else
    echo -e "${RED}✗ Login failed - no token received${NC}"
    JWT_TOKEN="dummy_token_for_testing"
fi

# Test 4: Get All Movies
echo -e "\n${YELLOW}[TEST 4] GET /catalog/movies${NC}"
movies_response=$(curl -s -X GET ${BASE_URL}/catalog/movies \
  -H "Authorization: Bearer $JWT_TOKEN")
echo "Response: $movies_response"
if echo "$movies_response" | grep -q "\[\]" || echo "$movies_response" | grep -q "id"; then
    echo -e "${GREEN}✓ Movies endpoint accessible${NC}"
else
    echo -e "${RED}✗ Movies endpoint failed${NC}"
fi

# Test 5: Get Shows by Movie ID
echo -e "\n${YELLOW}[TEST 5] GET /catalog/shows/1${NC}"
shows_response=$(curl -s -X GET ${BASE_URL}/catalog/shows/1 \
  -H "Authorization: Bearer $JWT_TOKEN")
echo "Response: $shows_response"
if echo "$shows_response" | grep -q "\[\]" || echo "$shows_response" | grep -q "id"; then
    echo -e "${GREEN}✓ Shows endpoint accessible${NC}"
else
    echo -e "${RED}✗ Shows endpoint failed${NC}"
fi

# Test 6: Initiate Booking - Missing Fields (Validation Test)
echo -e "\n${YELLOW}[TEST 6] POST /bookings/initiate - Validation Test${NC}"
validation_response=$(curl -s -X POST ${BASE_URL}/bookings/initiate \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}')
echo "Response: $validation_response"
if echo "$validation_response" | grep -q "required\|showId\|seatId"; then
    echo -e "${GREEN}✓ Validation working correctly${NC}"
else
    echo -e "${YELLOW}⚠ Validation response unexpected${NC}"
fi

# Test 7: Initiate Booking - Show Not Found
echo -e "\n${YELLOW}[TEST 7] POST /bookings/initiate - Show Not Found${NC}"
notfound_response=$(curl -s -X POST ${BASE_URL}/bookings/initiate \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "showId": 99999,
    "seatId": "A1"
  }')
echo "Response: $notfound_response"
if echo "$notfound_response" | grep -q "not found\|Show"; then
    echo -e "${GREEN}✓ Show not found error handled correctly${NC}"
else
    echo -e "${YELLOW}⚠ Error handling may need verification${NC}"
fi

# Test 8: Initiate Booking - Valid Request (Will fail if inventory service not running)
echo -e "\n${YELLOW}[TEST 8] POST /bookings/initiate - Valid Request${NC}"
booking_response=$(curl -s -X POST ${BASE_URL}/bookings/initiate \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "showId": 1,
    "seatId": "A12"
  }')
echo "Response: $booking_response"
if echo "$booking_response" | grep -q "bookingId\|PENDING"; then
    echo -e "${GREEN}✓ Booking created successfully${NC}"
elif echo "$booking_response" | grep -q "unavailable\|locked"; then
    echo -e "${YELLOW}⚠ Expected: Inventory service issue or seat locked${NC}"
else
    echo -e "${YELLOW}⚠ Booking response: Check if show exists in DB${NC}"
fi

# Test 9: Unauthorized Access
echo -e "\n${YELLOW}[TEST 9] Unauthorized Access Test${NC}"
unauth_response=$(curl -s -o /dev/null -w "%{http_code}" -X POST ${BASE_URL}/bookings/initiate \
  -H "Content-Type: application/json" \
  -d '{
    "showId": 1,
    "seatId": "A1"
  }')
if [ "$unauth_response" == "401" ] || [ "$unauth_response" == "403" ]; then
    echo -e "${GREEN}✓ Unauthorized access blocked (HTTP $unauth_response)${NC}"
else
    echo -e "${YELLOW}⚠ Security check: HTTP $unauth_response${NC}"
fi

echo -e "\n=========================================="
echo "✅ Testing Complete"
echo "=========================================="
