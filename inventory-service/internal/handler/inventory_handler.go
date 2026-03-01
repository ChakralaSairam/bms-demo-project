package handler

// InventoryHandler manages seat locking for high-concurrency flash sale scenarios.
//
// Java Spring equivalent:
//   @RestController
//   @RequestMapping("/inventory")
//   public class InventoryController {
//       @Autowired private RedisTemplate<String, String> redisTemplate;
//
//       @PostMapping("/lock")
//       public ResponseEntity<?> lockSeat(@RequestBody SeatLockRequest req) { ... }
//
//       @DeleteMapping("/lock")
//       public ResponseEntity<?> releaseSeat(@RequestParam String showId,
//                                            @RequestParam String seatId) { ... }
//
//       @GetMapping("/status")
//       public ResponseEntity<?> getLockStatus() { ... }   // public, no auth
//   }

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"inventory-service/internal/client"
	"inventory-service/internal/model"
)

// LockEntry is what each lock looks like in the GET /inventory/status response.
type LockEntry struct {
	Key     string `json:"key"`      // full Redis key e.g. lock:show:SH123:seat:A10
	LockedBy string `json:"locked_by"` // user_id who holds the lock
	TTLSecs int64  `json:"ttl_seconds"` // seconds remaining before auto-expiry
}

// seatLockTTL is how long a seat is held when locked.
// After 5 minutes with no confirmation, the lock auto-expires.
// Java equivalent: redisTemplate.expire(key, 5, TimeUnit.MINUTES)
const seatLockTTL = 5 * time.Minute

// InventoryHandler holds Redis + ShowClient dependencies.
// Java equivalent:
//   @Autowired private RedisTemplate<String, String> redisTemplate;
//   @Autowired private ShowClient showClient;
type InventoryHandler struct {
	redis      *redis.Client
	showClient *client.ShowClient
}

// NewInventoryHandler is the constructor.
func NewInventoryHandler(redisClient *redis.Client, showClient *client.ShowClient) *InventoryHandler {
	return &InventoryHandler{redis: redisClient, showClient: showClient}
}

// lockKey builds the Redis key for a seat lock.
//
// Format: lock:show:{show_id}:seat:{seat_id}
// Example: lock:show:101:seat:A10
//
// Java equivalent:
//   String key = String.format("lock:show:%d:seat:%s", showId, seatId);
func lockKey(showID int64, seatID string) string {
	return fmt.Sprintf("lock:show:%d:seat:%s", showID, seatID)
}

// LockSeat handles POST /inventory/lock
//
// This uses Redis SETNX (Set if Not Exists) — the heart of the flash sale logic.
//
// SETNX is ATOMIC — even with 10,000 concurrent users hitting this endpoint
// at the same time, Redis guarantees only ONE of them will get OK=true.
// Everyone else gets false and receives 409 Conflict.
//
// Java equivalent (using Lua script or RedisTemplate):
//   Boolean success = redisTemplate.opsForValue()
//       .setIfAbsent("lock:show:SH123:seat:A10", "user_99", 5, TimeUnit.MINUTES);
//   if (!success) return ResponseEntity.status(HttpStatus.CONFLICT).build();
func (h *InventoryHandler) LockSeat(c *gin.Context) {
	// Step 1: Parse and validate the request body.
	// ShouldBindJSON reads the JSON body and maps it to SeatLockRequest.
	// If any required field is missing or malformed, it returns an error.
	// Java equivalent: @RequestBody SeatLockRequest req (with @Valid)
	var req model.SeatLockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[LockSeat] ERROR - failed to parse request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	log.Printf("[LockSeat] Received request — show_id=%d seat_id=%q user_id=%d", req.ShowID, req.SeatID, req.UserID)

	// Step 2: Validate all required fields are present.
	if req.ShowID == 0 || req.SeatID == "" || req.UserID == 0 {
		log.Printf("[LockSeat] ERROR - validation failed — show_id=%d seat_id=%q user_id=%d", req.ShowID, req.SeatID, req.UserID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "show_id (Long), seat_id and user_id (Long) are required"})
		return
	}

	// Step 3: Build the Redis key.
	key := lockKey(req.ShowID, req.SeatID)
	log.Printf("[LockSeat] Redis key to set: %q", key)

	// Step 4: Attempt atomic SETNX with TTL.
	//
	// redis.SetNX = SET key value NX EX ttl
	//   NX  → only set if key does NOT exist (the "if not exists" part)
	//   EX  → set expiry at the same time (atomic — no separate EXPIRE call needed)
	//
	// Returns:
	//   true  → key was SET   = seat is now locked for this user  ✅
	//   false → key existed   = seat already taken by someone else ❌
	//
	// context.Background() — no timeout/cancellation context (like passing null in Java)
	// Store user_id as string in Redis (Redis values are always strings).
	// fmt.Sprintf converts int64 → string: 99 → "99"
	log.Printf("[LockSeat] Attempting SETNX — key=%q value=%d ttl=%s", key, req.UserID, seatLockTTL)
	ok, err := h.redis.SetNX(context.Background(), key, fmt.Sprintf("%d", req.UserID), seatLockTTL).Result()

	if err != nil {
		log.Printf("[LockSeat] ERROR - Redis SETNX failed for key=%q: %v", key, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not reach cache"})
		return
	}

	log.Printf("[LockSeat] SETNX result for key=%q: acquired=%v", key, ok)

	if !ok {
		// Key already existed — check who holds the lock for debugging
		existingVal, getErr := h.redis.Get(context.Background(), key).Result()
		if getErr == nil {
			log.Printf("[LockSeat] CONFLICT - seat already locked — key=%q held_by_user_id=%s", key, existingVal)
		} else {
			log.Printf("[LockSeat] CONFLICT - seat already locked — key=%q (could not read existing value: %v)", key, getErr)
		}
		c.JSON(http.StatusConflict, gin.H{
			"error":   "seat already locked",
			"show_id": req.ShowID,
			"seat_id": req.SeatID,
		})
		return
	}

	log.Printf("[LockSeat] SUCCESS - seat locked — key=%q user_id=%d ttl=%s", key, req.UserID, seatLockTTL)
	c.JSON(http.StatusOK, gin.H{
		"message":  "seat locked successfully",
		"show_id":  req.ShowID,
		"seat_id":  req.SeatID,
		"user_id":  req.UserID,
		"held_for": seatLockTTL.String(),
	})
}

// GetLockStatus handles GET /inventory/status — PUBLIC, no JWT required.
// Scans Redis for all active seat locks matching pattern "lock:show:*"
// and returns them with their TTL remaining.
//
// Java equivalent:
//   @GetMapping("/status")
//   public ResponseEntity<List<LockEntry>> getLockStatus() {
//       Set<String> keys = redisTemplate.keys("lock:show:*");
//       // build response from keys...
//   }
func (h *InventoryHandler) GetLockStatus(c *gin.Context) {
	ctx := context.Background()

	// KEYS pattern scan — find all seat lock keys.
	// "lock:show:*" matches any key starting with "lock:show:"
	// Java: redisTemplate.keys("lock:show:*") returns Set<String>
	//
	// NOTE: KEYS is fine for debugging/dev. In production with millions of
	// keys, use SCAN instead (non-blocking). Like using a cursor-based
	// paginated query vs a SELECT * in SQL.
	keys, err := h.redis.Keys(ctx, "lock:show:*").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not reach cache"})
		return
	}

	// Build a response slice — like new ArrayList<LockEntry>() in Java.
	locks := make([]LockEntry, 0, len(keys))

	for _, key := range keys {
		// GET the value (user_id) stored at this key
		// Java: String userId = redisTemplate.opsForValue().get(key);
		val, err := h.redis.Get(ctx, key).Result()
		if err != nil {
			continue // key may have expired between KEYS and GET — skip it
		}

		// TTL returns remaining lifetime in time.Duration.
		// -1 means no expiry set, -2 means key doesn't exist.
		// Java: Long ttl = redisTemplate.getExpire(key, TimeUnit.SECONDS);
		ttl, err := h.redis.TTL(ctx, key).Result()
		ttlSecs := int64(0)
		if err == nil {
			ttlSecs = int64(ttl.Seconds())
		}

		locks = append(locks, LockEntry{
			Key:      key,
			LockedBy: val,
			TTLSecs:  ttlSecs,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total_locks": len(locks),
		"locks":       locks,
	})
}

// SeatStatus represents one seat entry in the availability response.
// Status can be: "locked" (held in Redis) or "confirmed" (booked in Postgres).
type SeatStatus struct {
	SeatID   string `json:"seat_id"`
	Status   string `json:"status"`              // "locked" | "confirmed"
	LockedBy string `json:"locked_by,omitempty"` // user_id — only present for locked seats
	TTLSecs  int64  `json:"ttl_seconds,omitempty"` // only present for locked seats
}

// GetShowAvailability handles GET /inventory/show/:id/availability — JWT protected.
//
// Merges two data sources for a complete picture of unavailable seats:
//   1. Redis  → seats currently locked (held for 5 min, not yet confirmed)
//   2. Postgres → seats already confirmed (permanent bookings from Rajesh's DB)
//
// Java equivalent:
//   List<SeatStatus> locked    = redisTemplate.keys("lock:show:SH123:seat:*")...;
//   List<SeatStatus> confirmed = seatRepo.findByShowIdAndStatus(id, "CONFIRMED");
//   locked.addAll(confirmed);
func (h *InventoryHandler) GetShowAvailability(c *gin.Context) {
	ctx := context.Background()

	showID := c.Param("id")
	// Validate show_id is a valid Long (int64) — matches the Java service type.
	// Java equivalent: @PathVariable Long id (Spring rejects non-numeric automatically)
	if _, err := strconv.ParseInt(showID, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "show_id must be a valid numeric ID (Long)"})
		return
	}

	// ── Step 1: Fetch locked seats from Redis ────────────────────────────────
	// Scan for all keys matching this show: lock:show:SH123:seat:*
	pattern := fmt.Sprintf("lock:show:%s:seat:*", showID)
	keys, err := h.redis.Keys(ctx, pattern).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not reach cache"})
		return
	}

	// Use a map keyed by seat_id to hold results — lets us deduplicate if a
	// seat appears in both Redis and Postgres (confirmed wins over locked).
	// Java equivalent: Map<String, SeatStatus> seatMap = new LinkedHashMap<>();
	seatMap := make(map[string]SeatStatus)

	prefix := fmt.Sprintf("lock:show:%s:seat:", showID)
	for _, key := range keys {
		seatID := key[len(prefix):]

		val, err := h.redis.Get(ctx, key).Result()
		if err != nil {
			continue // expired between KEYS and GET — skip
		}

		ttl, _ := h.redis.TTL(ctx, key).Result()

		seatMap[seatID] = SeatStatus{
			SeatID:   seatID,
			Status:   "locked",
			LockedBy: val,
			TTLSecs:  int64(ttl.Seconds()),
		}
	}

	// ── Step 2: Fetch booked seats from Rajesh's Java show-service ───────────
	// Calls GET /internal/shows/{showId}/booked-seats — returns []string of seat IDs.
	// If the Java service is unreachable, we still return Redis data with a warning
	// header — partial data is better than a hard 500. Like @Fallback in Resilience4j.
	//
	// Java RestTemplate equivalent:
	//   List<String> booked = showClient.getBookedSeats(showId);
	if h.showClient != nil {
		booked, err := h.showClient.GetBookedSeats(c.Request.Context(), showID)
		if err != nil {
			// Don't fail the whole request — log a warning header and continue.
			c.Header("X-ShowService-Warning", fmt.Sprintf("show-service call failed: %v", err))
		} else {
			for _, seatID := range booked {
				// Confirmed booking always wins over a Redis lock for the same seat.
				seatMap[seatID] = SeatStatus{
					SeatID: seatID,
					Status: "confirmed",
				}
			}
		}
	}

	// ── Step 3: Flatten map into a slice for the JSON response ───────────────
	// Maps in Go have no guaranteed order — like HashMap in Java.
	// We collect into a slice so the response is a clean JSON array.
	unavailable := make([]SeatStatus, 0, len(seatMap))
	for _, seat := range seatMap {
		unavailable = append(unavailable, seat)
	}

	c.JSON(http.StatusOK, gin.H{
		"show_id":          showID,
		"unavailable_count": len(unavailable),
		"unavailable_seats": unavailable,
	})
}

// ReleaseSeat handles DELETE /inventory/lock
//
// Manually removes the Redis lock so the seat becomes available immediately,
// without waiting for the 5-minute TTL to expire.
//
// Use case: User deselects a seat → release it immediately so others can grab it.
//
// Java equivalent:
//   @DeleteMapping("/lock")
//   public ResponseEntity<?> releaseSeat(@RequestParam String showId,
//                                         @RequestParam String seatId) {
//       redisTemplate.delete("lock:show:" + showId + ":seat:" + seatId);
//       return ResponseEntity.ok("released");
//   }
func (h *InventoryHandler) ReleaseSeat(c *gin.Context) {
	// Step 1: Read query parameters from the URL.
	// Example URL: DELETE /inventory/lock?show_id=101&seat_id=A10
	//
	// c.Query returns a string — must parse to int64 to match Java Long.
	showIDStr := c.Query("show_id")
	seatID := c.Query("seat_id")

	// Step 2: Validate required params are present and show_id is a valid Long.
	if showIDStr == "" || seatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "show_id and seat_id query params are required"})
		return
	}
	showID, err := strconv.ParseInt(showIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "show_id must be a valid numeric ID (Long)"})
		return
	}

	// Step 3: Build the same key format used during lock.
	key := lockKey(showID, seatID)

	// Step 4: Delete the key from Redis.
	// Del returns the NUMBER of keys deleted (0 or 1 here).
	// 0 means the key didn't exist — seat was never locked or already expired.
	// Java: Long deleted = redisTemplate.delete("lock:show:SH123:seat:A10");
	deleted, err := h.redis.Del(context.Background(), key).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not reach cache"})
		return
	}

	// Step 5: Check if there was actually a lock to delete.
	if deleted == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "no active lock found for this seat",
			"show_id": showID,
			"seat_id": seatID,
		})
		return
	}

	// Step 6: Lock released — seat is now available for others.
	c.JSON(http.StatusOK, gin.H{
		"message": "seat lock released",
		"show_id": showID,
		"seat_id": seatID,
	})
}
