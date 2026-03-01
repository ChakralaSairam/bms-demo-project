#!/bin/bash

# Quick Manual Test Commands
# Copy and paste these commands to test the service

echo "======================================"
echo "🧪 MANUAL TEST COMMANDS"
echo "======================================"

echo -e "\n1️⃣ Test Catalog - Get All Movies"
echo "curl -s http://localhost:8082/catalog/movies | jq"

echo -e "\n2️⃣ Test Catalog - Get Shows for Movie 1"
echo "curl -s http://localhost:8082/catalog/shows/1 | jq"

echo -e "\n3️⃣ Register New User"
echo 'curl -s -X POST http://localhost:8082/api/auth/register \
  -H "Content-Type: application/json" \
  -d '"'"'{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "SecurePass@123"
  }'"'"

echo -e "\n4️⃣ Login and Get JWT Token"
echo 'curl -s -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '"'"'{
    "username": "john_doe",
    "password": "SecurePass@123"
  }'"'"' | jq'

echo -e "\n5️⃣ Test Booking (Replace TOKEN with actual JWT)"
echo 'TOKEN="your-jwt-token-here"
curl -s -X POST http://localhost:8082/bookings/initiate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '"'"'{
    "showId": 1,
    "seatId": "A15"
  }'"'"' | jq'

echo -e "\n6️⃣ Test Unauthorized Access (No Token)"
echo 'curl -s -X POST http://localhost:8082/bookings/initiate \
  -H "Content-Type: application/json" \
  -d '"'"'{
    "showId": 1,
    "seatId": "A15"
  }'"'"

echo -e "\n7️⃣ Test Validation Error (Missing Fields)"
echo 'TOKEN="your-jwt-token-here"
curl -s -X POST http://localhost:8082/bookings/initiate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '"'"'{}'"'"' | jq'

echo -e "\n8️⃣ Test Show Not Found"
echo 'TOKEN="your-jwt-token-here"
curl -s -X POST http://localhost:8082/bookings/initiate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '"'"'{
    "showId": 99999,
    "seatId": "A15"
  }'"'"' | jq'

echo -e "\n9️⃣ Check Database - View All Bookings"
echo 'docker exec -it spring-auth-db psql -U myuser -d auth_db -c "SELECT * FROM bookings;"'

echo -e "\n🔟 Check Database - View All Users"
echo 'docker exec -it spring-auth-db psql -U myuser -d auth_db -c "SELECT id, username, email FROM users;"'

echo -e "\n======================================"
echo "✅ Copy commands above to test!"
echo "======================================"
