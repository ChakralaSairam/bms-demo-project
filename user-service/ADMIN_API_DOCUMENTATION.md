# Admin API Documentation

## Overview
Complete CRUD operations for Movies, Theaters, and Shows with role-based access control.

**Access**: ADMIN role required for all endpoints

---

## Authentication

### Register Admin User
```bash
POST /api/auth/register
Content-Type: application/json

{
  "username": "admin",
  "email": "admin@example.com",
  "password": "Admin@123",
  "roles": ["ADMIN"]
}
```

### Register Regular User (Default: USER role)
```bash
POST /api/auth/register
Content-Type: application/json

{
  "username": "user",
  "email": "user@example.com",
  "password": "User@123"
}
```

### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "Admin@123"
}

Response:
{
  "username": "admin",
  "roles": ["ROLE_ADMIN"],
  "token": "eyJhbGc..."
}
```

---

## Movie Management (Admin Only)

### Create Movie
```bash
POST /api/admin/movies
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "title": "Avatar",
  "genre": "Sci-Fi",
  "duration": 162,
  "language": "English"
}

Response: 201 Created
{
  "id": 4,
  "title": "Avatar",
  "genre": "Sci-Fi",
  "duration": 162,
  "language": "English"
}
```

### Update Movie
```bash
PUT /api/admin/movies/{id}
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "title": "Avatar: The Way of Water",
  "genre": "Sci-Fi",
  "duration": 192,
  "language": "English"
}

Response: 200 OK
```

### Delete Movie
```bash
DELETE /api/admin/movies/{id}
Authorization: Bearer <admin-token>

Response: 204 No Content
```

---

## Theater Management (Admin Only)

### Get All Theaters
```bash
GET /api/admin/theaters
Authorization: Bearer <admin-token>

Response: 200 OK
[
  {
    "id": 1,
    "name": "PVR Cinemas",
    "city": "Mumbai"
  }
]
```

### Create Theater
```bash
POST /api/admin/theaters
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Carnival Cinemas",
  "city": "Hyderabad"
}

Response: 201 Created
{
  "id": 4,
  "name": "Carnival Cinemas",
  "city": "Hyderabad"
}
```

### Update Theater
```bash
PUT /api/admin/theaters/{id}
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Carnival Cinemas IMAX",
  "city": "Hyderabad"
}

Response: 200 OK
```

### Delete Theater
```bash
DELETE /api/admin/theaters/{id}
Authorization: Bearer <admin-token>

Response: 204 No Content
```

---

## Show Management (Admin Only)

### Get All Shows
```bash
GET /api/admin/shows
Authorization: Bearer <admin-token>

Response: 200 OK
[
  {
    "id": 1,
    "movieId": 1,
    "theaterId": 1,
    "startTime": "2024-01-20T18:00:00",
    "price": 250.00
  }
]
```

### Create Show
```bash
POST /api/admin/shows
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "movieId": 4,
  "theaterId": 4,
  "startTime": "2024-02-01T20:00:00",
  "price": 400.00
}

Response: 201 Created
{
  "id": 6,
  "movieId": 4,
  "theaterId": 4,
  "startTime": "2024-02-01T20:00:00",
  "price": 400.00
}
```

### Update Show
```bash
PUT /api/admin/shows/{id}
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "movieId": 4,
  "theaterId": 4,
  "startTime": "2024-02-01T21:00:00",
  "price": 450.00
}

Response: 200 OK
```

### Delete Show
```bash
DELETE /api/admin/shows/{id}
Authorization: Bearer <admin-token>

Response: 204 No Content
```

---

## Public Endpoints (No Auth Required)

### Get All Movies
```bash
GET /api/catalog/movies

Response: 200 OK
[
  {
    "id": 1,
    "title": "Inception",
    "genre": "Sci-Fi",
    "duration": 148,
    "language": "English"
  }
]
```

### Get Shows by Movie
```bash
GET /api/catalog/shows/{movieId}

Response: 200 OK
[
  {
    "id": 1,
    "movieId": 1,
    "theaterId": 1,
    "startTime": "2024-01-20T18:00:00",
    "price": 250.00
  }
]
```

---

## Booking Endpoints (Authenticated Users)

### Initiate Booking
```bash
POST /api/bookings/initiate
Authorization: Bearer <user-token>
Content-Type: application/json

{
  "showId": 1,
  "seatId": "A12"
}

Response: 201 Created
{
  "bookingId": "550e8400-e29b-41d4-a716-446655440000",
  "status": "PENDING",
  "message": "Booking initiated successfully"
}
```

---

## Error Responses

### 400 Bad Request - Validation Error
```json
{
  "title": "Title is required",
  "duration": "Duration is required"
}
```

### 401 Unauthorized
```json
{
  "timestamp": "2024-01-20T10:00:00.000+00:00",
  "status": 401,
  "error": "Unauthorized",
  "path": "/api/admin/movies"
}
```

### 403 Forbidden - Insufficient Permissions
```json
{
  "timestamp": "2024-01-20T10:00:00.000+00:00",
  "status": 403,
  "error": "Forbidden",
  "path": "/api/admin/movies"
}
```

### 404 Not Found
```json
{
  "error": "Movie not found"
}
```

---

## Security Features

### Role-Based Access Control
- **ADMIN**: Full CRUD access to movies, theaters, and shows
- **USER**: Can browse catalog and create bookings
- **Public**: Can view catalog only

### JWT Authentication
- Token expires in 24 hours
- Token includes user roles
- Stateless authentication

### Password Security
- BCrypt encryption
- Minimum 6 characters required
- Stored securely in database

---

## Testing

### Test Admin Access
```bash
# 1. Register admin
curl -X POST http://localhost:8082/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","email":"admin@test.com","password":"Admin@123","roles":["ADMIN"]}'

# 2. Login
TOKEN=$(curl -s -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin@123"}' | jq -r '.token')

# 3. Create movie
curl -X POST http://localhost:8082/api/admin/movies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Movie","genre":"Action","duration":120,"language":"English"}'
```

### Test Regular User Cannot Access Admin
```bash
# 1. Register user
curl -X POST http://localhost:8082/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user","email":"user@test.com","password":"User@123"}'

# 2. Login
USER_TOKEN=$(curl -s -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"User@123"}' | jq -r '.token')

# 3. Try to create movie (should fail with 403)
curl -X POST http://localhost:8082/api/admin/movies \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","genre":"Action","duration":120,"language":"English"}'
```

---

## Complete API Endpoints Summary

| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| POST | /api/auth/register | Public | Register new user |
| POST | /api/auth/login | Public | Login and get JWT |
| GET | /api/catalog/movies | Public | Get all movies |
| GET | /api/catalog/shows/{movieId} | Public | Get shows by movie |
| POST | /api/bookings/initiate | USER | Create booking |
| POST | /api/admin/movies | ADMIN | Create movie |
| PUT | /api/admin/movies/{id} | ADMIN | Update movie |
| DELETE | /api/admin/movies/{id} | ADMIN | Delete movie |
| GET | /api/admin/theaters | ADMIN | Get all theaters |
| POST | /api/admin/theaters | ADMIN | Create theater |
| PUT | /api/admin/theaters/{id} | ADMIN | Update theater |
| DELETE | /api/admin/theaters/{id} | ADMIN | Delete theater |
| GET | /api/admin/shows | ADMIN | Get all shows |
| POST | /api/admin/shows | ADMIN | Create show |
| PUT | /api/admin/shows/{id} | ADMIN | Update show |
| DELETE | /api/admin/shows/{id} | ADMIN | Delete show |
