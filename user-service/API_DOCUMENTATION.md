# User Service - Movie Booking System

## Architecture Overview

This is a Spring Boot 3.x microservice that handles:
- Movie catalog management
- Show listings
- Booking coordination with external inventory-service

## Tech Stack

- Java 17+
- Spring Boot 3.4.2
- PostgreSQL
- Spring Data JPA
- Spring Security with JWT
- OpenFeign for inter-service communication
- Lombok
- Maven

## API Endpoints

### Catalog Module

#### GET `/catalog/movies`
Returns list of all available movies.

**Response:**
```json
[
  {
    "id": 1,
    "title": "Movie Title",
    "genre": "Action",
    "duration": 120,
    "language": "English"
  }
]
```

#### GET `/catalog/shows/{movieId}`
Returns all shows for a specific movie.

**Response:**
```json
[
  {
    "id": 1,
    "movieId": 1,
    "theaterId": 1,
    "startTime": "2024-01-15T18:00:00",
    "price": 250.00
  }
]
```

### Booking Module

#### POST `/bookings/initiate`
Initiates a booking by locking a seat via inventory-service.

**Request:**
```json
{
  "showId": 1,
  "seatId": "A12"
}
```

**Success Response (201):**
```json
{
  "bookingId": "550e8400-e29b-41d4-a716-446655440000",
  "status": "PENDING",
  "message": "Booking initiated successfully"
}
```

**Error Responses:**
- `404` - Show not found
- `409` - Seat already locked
- `503` - Inventory service unavailable

## Configuration

### application.properties

```properties
# Inventory Service
inventory.service.url=http://localhost:9090

# Database
spring.datasource.url=jdbc:postgresql://localhost:5432/auth_db
spring.datasource.username=myuser
spring.datasource.password=mypassword

# JWT
app.jwt.secret=YourSuperSecretKeyForSigningTokensMustBeLongEnough12345
app.jwt.expiration-ms=86400000
```

## Database Schema

### Movies
- id (BIGINT, PK)
- title (VARCHAR)
- genre (VARCHAR)
- duration (INTEGER)
- language (VARCHAR)

### Theaters
- id (BIGINT, PK)
- name (VARCHAR)
- city (VARCHAR)

### Shows
- id (BIGINT, PK)
- movie_id (BIGINT, FK, indexed)
- theater_id (BIGINT, FK, indexed)
- start_time (TIMESTAMP)
- price (DECIMAL)

### Bookings
- id (UUID, PK)
- user_id (BIGINT, indexed)
- show_id (BIGINT, indexed)
- seat_id (VARCHAR)
- status (VARCHAR: PENDING, CONFIRMED, FAILED)
- amount (DECIMAL)
- created_at (TIMESTAMP)

## Error Handling

Global exception handler provides consistent error responses:

- `ShowNotFoundException` → 404
- `SeatAlreadyLockedException` → 409
- `InventoryServiceException` → 503
- Validation errors → 400

## Inter-Service Communication

Uses OpenFeign to communicate with inventory-service:

**Endpoint:** `POST /inventory/lock`

**Request:**
```json
{
  "showId": "1",
  "seatId": "A12"
}
```

**Expected Responses:**
- `201` - Seat locked successfully
- `409` - Seat already locked
- `500/503` - Service unavailable

## Security

- JWT-based authentication
- User ID extracted from SecurityContext
- All booking endpoints require authentication

## Build & Run

```bash
mvn clean install
mvn spring-boot:run
```

## Project Structure

```
com.test.user
├── client/              # Feign clients
├── config/              # Configuration classes
├── controller/          # REST controllers
├── dto/                 # Data Transfer Objects
├── exception/           # Custom exceptions
├── models/              # JPA entities
├── repository/          # JPA repositories
├── security/            # Security configuration
└── service/             # Business logic
```
