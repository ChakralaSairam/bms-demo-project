# inventory-service

A Go microservice for a movie catalog inventory system, built with [Gin](https://github.com/gin-gonic/gin).

## Tech Stack

| Concern | Technology |
|---|---|
| Language | Go 1.21 |
| Web Framework | Gin (`github.com/gin-gonic/gin`) |
| Auth | JWT (`github.com/golang-jwt/jwt/v5`) |
| Containerization | Docker (multi-stage) + Docker Compose |

---

## Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/health` | None | Health check |
| GET | `/movies` | JWT | List all movies |
| GET | `/items` | JWT | List inventory items |
| POST | `/items` | JWT | Create inventory item |
| GET | `/items/:id` | JWT | Get item by ID |

---

## Run Locally

```bash
cd inventory-service
go run ./cmd/main.go
# Server starts on http://localhost:8081
```

### Test health (no token needed):
```bash
curl http://localhost:8081/health
```

### Test protected route (token required):
```bash
# First, get a token from the Java user-service (Day 3):
TOKEN="<paste JWT from Java /register or /login here>"

curl http://localhost:8081/movies \
  -H "Authorization: Bearer $TOKEN"
```

---

## Run with Docker

```bash
# From the Demo Project root:
docker-compose up --build
```

This starts the Go service on port `8081`. When the Java `user-service` is ready, uncomment its block in `docker-compose.yml`.

---

## Architectural Decision: JWT for Stateless Inter-Service Communication

### Why JWT?

In a distributed microservice architecture, services must verify caller identity **without sharing session state**. We chose **JWT (JSON Web Tokens)** over alternatives like session cookies or API keys for the following reasons:

| Concern | JWT Approach |
|---|---|
| **Stateless** | No shared session store (Redis/DB) needed between Java and Go services. Each service validates the token independently using the shared secret. |
| **Self-contained** | The token carries claims (user ID, roles, expiry) — no DB lookup on every request. Like a signed passport. |
| **Language-agnostic** | Java uses `io.jsonwebtoken` (JJWT); Go uses `golang-jwt/jwt`. Both implement the same RFC 7519 standard — they can verify each other's tokens. |
| **Scalable** | Horizontal scaling is trivial — any instance of any service can verify a token without coordination. |

### The Shared Secret Pattern (Current: Week 1)

Both services use the **same HMAC-SHA256 secret key** (`JWTSecret` in Go, `jwt.secret` property in Java).

```
Java user-service          Go inventory-service
─────────────────          ────────────────────
User registers/logs in  →  Signs JWT with secret
                           ↓
Client sends JWT        →  Go validates with same secret
                           → 200 OK or 401 Unauthorized
```

> **Week 11 upgrade:** The hardcoded secret moves to **AWS Secrets Manager**. Both services fetch it at startup — no secrets in source code.

### Trade-offs Accepted

- **Token revocation** is not instant (must wait for expiry). Mitigation: short expiry (15 min) + refresh tokens.
- **Secret rotation** requires coordinated redeployment of both services until AWS Secrets Manager is integrated.
