package middleware

import (
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

		if apiKey != os.Getenv("API_KEY") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{Error: "invalid API key"})
			return
		}

		c.Next()
	}
}
