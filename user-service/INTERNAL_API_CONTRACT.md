# 🤝 Internal API Contract - Seat Booking Integration

## Overview
This document defines the **strict API contract** between:
- **Java Service (user-service)**: Owns booking data in PostgreSQL
- **Go Service (inventory-service)**: Manages seat availability in Redis

---

## 📋 Official Contract

### Endpoint Details

| Property | Value |
|----------|-------|
| **URL** | `http://<RAJESH_IP>:8082/internal/shows/{showId}/booked-seats` |
| **Method** | `GET` |
| **Authentication** | None (Internal network only) |
| **Content-Type** | `application/json` |
| **Purpose** | Retrieve permanently booked (CONFIRMED) seats for a show |

---

## 📥 Request

### Path Parameters
- `showId` (Long, required): The unique identifier of the show

### Example Request
```bash
GET http://localhost:8082/internal/shows/1/booked-seats
Accept: application/json
```

---

## 📤 Response

### Success Response (200 OK)

**When seats are booked:**
```json
["A1", "A5", "C12"]
```

**When no seats are booked:**
```json
[]
```

**Important**: Always returns `200 OK` with an empty array if no confirmed bookings exist. Never returns `404` for this case.

### Error Response (404 Not Found)

**When show doesn't exist:**
```json
{
  "error": "Show not found with ID: 99999"
}
```

---

## 🔍 Business Logic

### What Gets Returned
- Only seats with `status = 'CONFIRMED'` in the bookings table
- Returns seat IDs as strings (e.g., "A1", "B10")
- Excludes seats with status `PENDING` or `FAILED`

### Why This Design
1. **Separation of Concerns**: Java owns "business status" (payment confirmed), Go owns "physical availability" (Redis locks)
2. **Simple Contract**: Just an array of strings - easy to parse in any language
3. **No Postgres Dependency**: Go service doesn't need database access

---

## 🧪 Testing the API

### Test 1: Empty Result
```bash
curl http://localhost:8082/internal/shows/1/booked-seats

Expected: []
Status: 200 OK
```

### Test 2: With Confirmed Bookings
```bash
# First, create confirmed bookings in database
# Then call API
curl http://localhost:8082/internal/shows/1/booked-seats

Expected: ["A1", "A5"]
Status: 200 OK
```

### Test 3: Invalid Show ID
```bash
curl http://localhost:8082/internal/shows/99999/booked-seats

Expected: {"error": "Show not found with ID: 99999"}
Status: 404 Not Found
```

---

## 🔐 Security Considerations

### Current Implementation
- **Public endpoint** (no authentication required)
- Intended for internal network communication only

### Future Enhancement: API Key Protection

To add API key authentication later, modify the controller:

```java
@GetMapping("/{showId}/booked-seats")
public ResponseEntity<List<String>> getBookedSeats(
    @PathVariable Long showId,
    @RequestHeader("X-Internal-API-Key") String apiKey) {
    
    if (!"your-secret-key".equals(apiKey)) {
        return ResponseEntity.status(HttpStatus.UNAUTHORIZED).build();
    }
    
    List<String> bookedSeats = internalShowService.getConfirmedBookedSeats(showId);
    return ResponseEntity.ok(bookedSeats);
}
```

Then Go service would call:
```go
req.Header.Set("X-Internal-API-Key", "your-secret-key")
```

---

## 🔄 Integration Flow

### Complete Booking Flow

```
1. User requests booking
   ↓
2. Java Service → Go Service: Lock seat in Redis (temporary)
   ↓
3. User completes payment
   ↓
4. Java Service: Update booking status to CONFIRMED in Postgres
   ↓
5. Go Service calls: GET /internal/shows/{showId}/booked-seats
   ↓
6. Go Service: Merges confirmed seats with Redis locks
   ↓
7. Go Service: Returns final availability to users
```

---

## 💻 Go Service Integration Code

### Sample Go Client Code for Sairam

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type JavaServiceClient struct {
    BaseURL string
}

func (c *JavaServiceClient) GetConfirmedBookedSeats(showID int64) ([]string, error) {
    url := fmt.Sprintf("%s/internal/shows/%d/booked-seats", c.BaseURL, showID)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to call Java service: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 404 {
        return nil, fmt.Errorf("show not found")
    }
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }
    
    var bookedSeats []string
    if err := json.NewDecoder(resp.Body).Decode(&bookedSeats); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return bookedSeats, nil
}

// Usage in Go service
func main() {
    client := &JavaServiceClient{
        BaseURL: "http://localhost:8082",
    }
    
    bookedSeats, err := client.GetConfirmedBookedSeats(1)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Confirmed booked seats: %v\n", bookedSeats)
    // Now merge with Redis data to get final availability
}
```

---

## 📊 Data Flow Example

### Scenario: Show 1 has 100 seats (A1-A100)

**In Java Service (PostgreSQL):**
```
Bookings table:
- A1: CONFIRMED (paid)
- A5: CONFIRMED (paid)
- B10: PENDING (payment in progress)
```

**API Response:**
```json
["A1", "A5"]
```
Note: B10 is NOT included (still PENDING)

**In Go Service (Redis):**
```
Locked seats (temporary):
- B10: locked for 10 minutes
- C15: locked for 5 minutes
```

**Final Availability (Go Service calculates):**
```
Unavailable seats = ["A1", "A5"] + ["B10", "C15"]
Available seats = All seats - Unavailable seats
                = 96 seats available
```

---

## ✅ Contract Validation Checklist

Before going live, verify:

- [ ] Java service returns `[]` for shows with no confirmed bookings
- [ ] Java service returns `404` for non-existent shows
- [ ] Response is always a JSON array of strings
- [ ] Only CONFIRMED bookings are included
- [ ] PENDING and FAILED bookings are excluded
- [ ] Go service can successfully parse the response
- [ ] Network connectivity between services is established
- [ ] Error handling works on both sides

---

## 🚀 Deployment Notes

### Environment Variables

**Java Service:**
```properties
# application.properties
server.port=8082
```

**Go Service:**
```bash
JAVA_SERVICE_URL=http://user-service:8082
```

### Docker Compose Example
```yaml
services:
  user-service:
    image: user-service:latest
    ports:
      - "8082:8082"
    networks:
      - internal-network
  
  inventory-service:
    image: inventory-service:latest
    environment:
      - JAVA_SERVICE_URL=http://user-service:8082
    networks:
      - internal-network

networks:
  internal-network:
    driver: bridge
```

---

## 📞 Support & Questions

**For Java Service Issues:**
- Check logs: `docker logs spring-auth-app`
- Verify database connection
- Ensure show exists in database

**For Go Service Issues:**
- Verify network connectivity: `curl http://user-service:8082/internal/shows/1/booked-seats`
- Check JSON parsing
- Validate Redis merge logic

---

## 🎯 Success Criteria

The integration is successful when:
1. ✅ Go service can fetch confirmed bookings from Java service
2. ✅ Go service merges data with Redis locks correctly
3. ✅ Users see accurate seat availability
4. ✅ No double bookings occur
5. ✅ System handles errors gracefully

---

**Contract Version:** 1.0  
**Last Updated:** 2024-02-28  
**Status:** ✅ PRODUCTION READY
