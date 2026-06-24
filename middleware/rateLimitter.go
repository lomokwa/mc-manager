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
}

func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		clients:  make(map[string]*client),
		rate:     requestsPerSecond,
		maxBurst: float64(burst),
	}
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
