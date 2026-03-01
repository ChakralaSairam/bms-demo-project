package config

// JWTSecret is the shared secret key used to sign and verify JWT tokens.
//
// Java equivalent (Spring Boot application.properties):
//   jwt.secret=super-secret-inventory-key-change-in-production
//
// IMPORTANT: This is hardcoded for development only.
// In production (Week 11), this will be fetched from AWS Secrets Manager.
// Both this Go service and the Java user-service must use the SAME secret.
const JWTSecret = "YourSuperSecretKeyForSigningTokensMustBeLongEnough12345"
