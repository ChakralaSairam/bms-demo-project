package handler

// MovieHandler serves the /health and /movies endpoints.
//
// Java Spring equivalent:
//   @RestController
//   public class MovieController {
//       @Autowired private RedisTemplate<String, String> redisTemplate;
//
//       @GetMapping("/health")
//       public ResponseEntity<Map<String,String>> health() { ... }
//
//       @GetMapping("/movies")
//       public ResponseEntity<List<Movie>> listMovies() { ... }
//   }
//
// Redis caching pattern used here:
//   1. Check Redis for cached response (Cache-Aside / Lazy Loading)
//   2. Cache HIT  → return immediately (fast path)
//   3. Cache MISS → fetch data, store in Redis with TTL, return response
// This is identical to Spring's @Cacheable annotation behaviour.

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"inventory-service/internal/model"
)

const moviesCacheKey = "movies:all"
const moviesCacheTTL = 5 * time.Minute

// MovieHandler holds dependencies for movie-related routes.
type MovieHandler struct {
	// redis is injected — like @Autowired RedisTemplate in Spring.
	redis *redis.Client
}

// NewMovieHandler is the constructor — accepts a Redis client.
//
// Java equivalent:
//   public MovieController(RedisTemplate<String,String> redisTemplate) {
//       this.redisTemplate = redisTemplate;
//   }
func NewMovieHandler(redisClient *redis.Client) *MovieHandler {
	return &MovieHandler{redis: redisClient}
}

// Health handles GET /health — public, no auth required.
func (h *MovieHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "inventory-service",
	})
}

// ListMovies handles GET /movies — protected by JWT middleware.
// Uses Redis to cache the movie list for 5 minutes.
//
// Java @Cacheable equivalent:
//   @Cacheable(value = "movies", key = "'all'")
//   public List<Movie> listMovies() { ... }
func (h *MovieHandler) ListMovies(c *gin.Context) {
	ctx := context.Background()

	// --- Cache HIT check ---
	// Like: String cached = redisTemplate.opsForValue().get("movies:all");
	cached, err := h.redis.Get(ctx, moviesCacheKey).Result()
	if err == nil {
		// Cache HIT — deserialize and return immediately
		var movies []model.Movie
		if jsonErr := json.Unmarshal([]byte(cached), &movies); jsonErr == nil {
			c.Header("X-Cache", "HIT") // useful header to confirm cache is working in Postman
			c.JSON(http.StatusOK, movies)
			return
		}
	}

	// --- Cache MISS — build the data ---
	// In a real service this would be a DB query.
	movies := []model.Movie{
		{ID: "1", Title: "Inception", Genre: "Sci-Fi", ReleaseYear: 2010, Rating: 8.8},
		{ID: "2", Title: "The Dark Knight", Genre: "Action", ReleaseYear: 2008, Rating: 9.0},
		{ID: "3", Title: "Interstellar", Genre: "Sci-Fi", ReleaseYear: 2014, Rating: 8.6},
		{ID: "4", Title: "Parasite", Genre: "Thriller", ReleaseYear: 2019, Rating: 8.5},
		{ID: "5", Title: "The Shawshank Redemption", Genre: "Drama", ReleaseYear: 1994, Rating: 9.3},
	}

	// Serialize to JSON and store in Redis with a 5-minute TTL.
	// Like: redisTemplate.opsForValue().set("movies:all", json, 5, TimeUnit.MINUTES);
	if data, jsonErr := json.Marshal(movies); jsonErr == nil {
		h.redis.Set(ctx, moviesCacheKey, data, moviesCacheTTL)
	}

	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, movies)
}
