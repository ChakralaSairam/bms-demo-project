package middleware

// JWTAuthMiddleware validates a JWT token on every protected request.
//
// Java Spring Security equivalent:
//   public class JwtAuthFilter extends OncePerRequestFilter {
//       @Override
//       protected void doFilterInternal(HttpServletRequest request,
//                                       HttpServletResponse response,
//                                       FilterChain filterChain) {
//           String token = request.getHeader("Authorization");
//           // validate token...
//           filterChain.doFilter(request, response);
//       }
//   }
//
// In Gin, middleware is a function that receives a *gin.Context.
// - c.Next()  → like filterChain.doFilter() — proceed to the next handler
// - c.Abort() → like response.sendError() — stop the chain

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"inventory-service/internal/config"
)

// JWTAuth returns a Gin middleware function that validates Bearer JWT tokens.
// Apply it to any route group you want to protect.
//
// Java equivalent: @PreAuthorize or Spring Security filter chain protecting specific routes.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the Authorization header — like request.getHeader("Authorization")
		authHeader := c.GetHeader("Authorization")

		// Expect format: "Bearer <token>"
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// Like: return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body("missing token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or malformed Authorization header",
			})
			return
		}

		// Strip "Bearer " prefix to get the raw token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the JWT using the shared secret.
		// golang-jwt equivalent of io.jsonwebtoken (JJWT) in Java:
		//   Jwts.parserBuilder().setSigningKey(secret).build().parseClaimsJws(token)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing algorithm is HMAC (HS256/HS384/HS512).
			// This prevents the "algorithm confusion" attack — always validate the alg.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// Token is valid — store claims in context for downstream handlers.
		// Like: SecurityContextHolder.getContext().setAuthentication(auth)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("claims", claims)
		}

		// Proceed to the next handler in the chain
		c.Next()
	}
}
