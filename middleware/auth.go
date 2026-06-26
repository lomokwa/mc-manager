package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/types"
)

func ValidateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("key")
		}

		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{Error: "missing API key"})
			return
		}

		// Constant-time comparison to avoid leaking the key through response
		// timing. An empty/unset API_KEY fails closed (no key can be valid).
		expected := os.Getenv("API_KEY")
		if expected == "" || subtle.ConstantTimeCompare([]byte(apiKey), []byte(expected)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{Error: "invalid API key"})
			return
		}

		c.Next()
	}
}
