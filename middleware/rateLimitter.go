package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/types"
)

type client struct {
	tokens    float64
	lastCheck time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*client
	rate     float64 // tokens per second
	maxBurst float64 // max tokens
	ttl      time.Duration
}

const (
	// defaultClientTTL is how long an idle client entry is kept before the
	// janitor evicts it.
	defaultClientTTL = 10 * time.Minute
	// defaultCleanupInterval is how often the janitor scans for stale clients.
	defaultCleanupInterval = time.Minute
)

func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*client),
		rate:     requestsPerSecond,
		maxBurst: float64(burst),
		ttl:      defaultClientTTL,
	}
	go rl.cleanupLoop(defaultCleanupInterval)
	return rl
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[ip]
	if !exists {
		rl.clients[ip] = &client{tokens: rl.maxBurst - 1, lastCheck: time.Now()}
		return true
	}

	elapsed := time.Since(c.lastCheck).Seconds()
	c.tokens += elapsed * rl.rate
	if c.tokens > rl.maxBurst {
		c.tokens = rl.maxBurst
	}
	c.lastCheck = time.Now()

	if c.tokens < 1 {
		return false
	}

	c.tokens--
	return true
}

// cleanupLoop periodically evicts clients that have been idle for longer than
// the configured TTL, so the clients map doesn't grow without bound.
func (rl *RateLimiter) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		rl.cleanup(time.Now().Add(-rl.ttl))
	}
}

// cleanup removes every client whose last activity is older than cutoff. An
// idle client has already replenished to maxBurst, so evicting it is
// equivalent to letting a fresh entry be created on its next request.
func (rl *RateLimiter) cleanup(cutoff time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for ip, c := range rl.clients {
		if c.lastCheck.Before(cutoff) {
			delete(rl.clients, ip)
		}
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, types.APIResponse{Error: "rate limit exceeded"})
			return
		}

		c.Next()
	}
}
