# User Authentication Service

## Objective

This is a Spring Boot REST API service that provides secure user authentication and authorization using JWT (JSON Web Tokens). The service allows users to register accounts and authenticate themselves to access protected resources. It implements stateless authentication, making it suitable for modern microservices architectures and mobile/web applications.

## Authentication Mechanism

This service uses **JWT-based authentication** with Spring Security. Here's how it works step by step:

### 1. User Registration (`POST /api/auth/register`)

- User submits username, email, password, and roles
- System checks if username already exists
- Password is hashed using BCrypt algorithm
- User account is saved to PostgreSQL database
- Returns success message

### 2. User Login (`POST /api/auth/login`)

**Step 1: Credential Verification**
- User submits username and password
- AuthenticationManager validates credentials against database
- UserDetailsServiceImpl loads user from database
- Password is verified using BCrypt encoder

**Step 2: JWT Token Generation**
- Upon successful authentication, JwtUtils generates a JWT token
- Token contains: username (subject), issued date, expiration date
- Token is signed with HMAC-SHA256 using a secret key
- Token is returned to client along with username and roles

### 3. Accessing Protected Resources

**Step 1: Request Interception**
- Client sends request with JWT in Authorization header: `Bearer <token>`
- AuthTokenFilter intercepts every incoming request

**Step 2: Token Extraction**
- Filter extracts JWT from Authorization header (removes "Bearer " prefix)

**Step 3: Token Validation**
- JwtUtils validates token signature and expiration
- If invalid or expired, request proceeds without authentication

**Step 4: User Authentication**
- Username is extracted from valid token
- UserDetailsServiceImpl loads user details from database
- UsernamePasswordAuthenticationToken is created with user authorities
- Authentication is set in SecurityContext

**Step 5: Authorization Check**
- Spring Security checks if authenticated user has access to requested endpoint
- Public endpoints (`/api/auth/**`) are accessible without token
- All other endpoints require valid authentication

### 4. Security Configuration

- **Session Management**: Stateless (no server-side sessions)
- **CSRF Protection**: Disabled (not needed for stateless JWT)
- **Password Encoding**: BCrypt with default strength
- **Token Expiration**: Configurable via `app.jwt.expiration-ms` property
- **Secret Key**: Configurable via `app.jwt.secret` property

## Technology Stack

- Spring Boot 4.0.3
- Spring Security
- JWT (jjwt 0.11.5)
- PostgreSQL
- JPA/Hibernate
- BCrypt password encoding
