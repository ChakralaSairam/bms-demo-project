package cache

// Redis client setup for inventory-service.
//
// Java equivalent (Spring Boot):
//   @Configuration
//   public class RedisConfig {
//       @Bean
//       public RedisTemplate<String, Object> redisTemplate(RedisConnectionFactory factory) {
//           RedisTemplate<String, Object> template = new RedisTemplate<>();
//           template.setConnectionFactory(factory);
//           return template;
//       }
//   }
//   # application.properties:
//   spring.data.redis.host=localhost
//   spring.data.redis.port=6379
//
// In Go, we create a single shared *redis.Client and pass it to handlers
// that need it — like a Spring @Autowired RedisTemplate.

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates and validates a Redis connection.
// Returns a *redis.Client to be shared across the application.
//
// Java equivalent:
//   new RedisStandaloneConfiguration("localhost", 6379)
func NewRedisClient(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr, // e.g., "localhost:6379" or "bms-redis:6379" inside Docker
		DB:   0,    // Redis has 16 DBs (0–15). DB 0 is the default — like the "default" schema.
	})

	// Ping Redis to confirm the connection is alive — like a JDBC connection test.
	// context.Background() is Go's way of saying "no timeout, no cancellation" — like
	// passing null for the executor context in Java.
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("could not connect to Redis at %s: %v", addr, err)
	}

	fmt.Printf("Connected to Redis at %s\n", addr)
	return client
}
