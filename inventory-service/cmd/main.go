package main

// In Java Spring Boot, your entry point looks like:
//
//   @SpringBootApplication
//   public class InventoryServiceApplication {
//       public static void main(String[] args) {
//           SpringApplication.run(InventoryServiceApplication.class, args);
//       }
//   }
//
// In Go with Gin:
// - gin.Default() creates a router with Logger + Recovery middleware pre-attached.
//   Like Spring Boot's auto-configured DispatcherServlet.
// - r.Group("/path") → like @RequestMapping("/path") on a controller class.
// - r.Use(middleware) → like adding a filter to Spring Security filter chain.
// - r.Run(":8081")   → like server.port=8081 in application.properties.

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"inventory-service/internal/cache"
	"inventory-service/internal/client"
	"inventory-service/internal/handler"
	"inventory-service/internal/middleware"
)

func main() {
	// gin.Default() — like SpringApplication.run() but returns the router/engine.
	// It includes two built-in middlewares:
	//   - Logger: logs every request (like Spring's CommonsRequestLoggingFilter)
	//   - Recovery: catches panics and returns 500 (like Spring's @ExceptionHandler)
	r := gin.Default()

	// Read Redis address from environment variable — like @Value("${spring.data.redis.host}")
	// Falls back to "localhost:6379" for local development.
	// Inside Docker Compose, REDIS_ADDR will be set to "bms-redis:6379".
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient := cache.NewRedisClient(redisAddr)

	// Base URL of Rajesh's Java service — read from env variable.
	// Defaults to his local IP for development. Override via SHOW_SERVICE_URL in docker-compose.
	showServiceURL := os.Getenv("SHOW_SERVICE_URL")
	if showServiceURL == "" {
		showServiceURL = "http://192.168.1.7:8082"
	}
	showClient := client.NewShowClient(showServiceURL)

	// --- Handlers (like @Autowired controllers in Spring) ---
	movieHandler := handler.NewMovieHandler(redisClient)
	inventoryHandler := handler.NewInventoryHandler(redisClient, showClient)

	// --- Public routes (no auth required) ---
	// Java equivalent: permitAll() in Spring Security config
	r.GET("/health", movieHandler.Health)
	r.GET("/inventory/status", inventoryHandler.GetLockStatus) // debug view of Redis locks

	// --- Protected routes (JWT required) ---
	// Java equivalent:
	//   http.authorizeHttpRequests(auth -> auth.requestMatchers("/inventory/**").authenticated())
	//
	// r.Group() creates a sub-router — like @RequestMapping("/inventory") on a class.
	// .Use(middleware.JWTAuth()) attaches the auth filter to this group only.
	protected := r.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		// POST /inventory/lock  — lock a seat (SETNX atomic operation)
		protected.POST("/inventory/lock", inventoryHandler.LockSeat)

		// DELETE /inventory/lock — release a seat lock manually
		protected.DELETE("/inventory/lock", inventoryHandler.ReleaseSeat)

		// GET /inventory/show/:id/availability — locked (Redis) + booked seats from Java service merged
		protected.GET("/inventory/show/:id/availability", inventoryHandler.GetShowAvailability)
	}

	port := ":8081"
	log.Printf("Inventory service running on http://localhost%s", port)

	// r.Run() starts the HTTP server — Gin's equivalent of SpringApplication.run()
	// embedding a Tomcat/Netty server.
	if err := r.Run(port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
