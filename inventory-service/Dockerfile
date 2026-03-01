# Multi-stage Dockerfile for the Go inventory-service
#
# Java equivalent: Like a Maven multi-stage build where you compile in one
# container and run in a minimal JRE container to reduce image size.
#
# Stage 1 — Builder
# Uses the full Go SDK image to compile the binary.
# Like: FROM maven:3.9-eclipse-temurin AS build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy all source including the vendor directory.
# Using vendored deps avoids any network calls inside the container,
# which fixes TLS/proxy issues in corporate/Docker Desktop environments.
# Like: Maven offline build using a pre-populated .m2 cache.
COPY . .

# Build a statically linked binary — CGO_ENABLED=0 ensures no C dependencies.
# -mod=vendor tells Go to use the local vendor/ folder instead of downloading.
# Like: mvn package -DskipTests -o (offline mode)
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o bin/inventory-service ./cmd/main.go

# ─────────────────────────────────────────────────────────────────────────────
# Stage 2 — Runtime
# Minimal Alpine image — only the compiled binary is copied over.
# Like: FROM eclipse-temurin:21-jre (just the JRE, not the full JDK)
# Final image size: ~10–15 MB vs ~300 MB for a full Go SDK image.
FROM alpine:latest

WORKDIR /app

# Copy ONLY the compiled binary from the builder stage
COPY --from=builder /app/bin/inventory-service .

# Expose the service port
EXPOSE 8081

# Run the binary — like ENTRYPOINT ["java", "-jar", "app.jar"]
ENTRYPOINT ["./inventory-service"]
