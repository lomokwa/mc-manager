package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/types"
)

func TestRateLimiter_AllowsRequests(t *testing.T) {
	r := gin.New()
	limiter := NewRateLimiter(10, 5)
	r.Use(limiter.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, types.APIResponse{Success: true})
	})

	// First 5 requests should pass (burst of 5)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}
}

func TestRateLimiter_BlocksExcessRequests(t *testing.T) {
	r := gin.New()
	limiter := NewRateLimiter(1, 2) // 1 req/sec, burst of 2
	r.Use(limiter.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, types.APIResponse{Success: true})
	})

	// Use up the burst
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	// Next request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}

	var resp types.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error != "rate limit exceeded" {
		t.Errorf("expected 'rate limit exceeded', got %q", resp.Error)
	}
}

func TestRateLimiter_DifferentIPsIndependent(t *testing.T) {
	r := gin.New()
	limiter := NewRateLimiter(1, 1) // very restrictive
	r.Use(limiter.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, types.APIResponse{Success: true})
	})

	// First IP uses its token
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "1.2.3.4:1234"
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("first IP: expected 200, got %d", w1.Code)
	}

	// Second IP should still be allowed
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "5.6.7.8:5678"
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("second IP: expected 200, got %d", w2.Code)
	}
}

func TestRateLimiter_EvictsStaleClients(t *testing.T) {
	rl := NewRateLimiter(10, 5)

	// A fresh, active client.
	rl.allow("fresh")

	// An old client that hasn't been seen in longer than the TTL.
	rl.mu.Lock()
	rl.clients["stale"] = &client{tokens: rl.maxBurst, lastCheck: time.Now().Add(-rl.ttl - time.Minute)}
	rl.mu.Unlock()

	rl.cleanup(time.Now().Add(-rl.ttl))

	rl.mu.Lock()
	_, freshKept := rl.clients["fresh"]
	_, staleKept := rl.clients["stale"]
	remaining := len(rl.clients)
	rl.mu.Unlock()

	if !freshKept {
		t.Error("fresh client should be kept")
	}
	if staleKept {
		t.Error("stale client should be evicted")
	}
	if remaining != 1 {
		t.Errorf("expected 1 client remaining, got %d", remaining)
	}
}
