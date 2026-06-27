package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Check Authorization header first
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{Error: "invalid Authorization format"})
				return
			}
		} else if t := c.Query("token"); t != "" {
			// Fallback to query param (for WebSocket connections)
			tokenString = t
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{Error: "missing Authorization header"})
			return
		}

		// Parse and validate
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{Error: "invalid or expired token"})
			return
		}

		// Extract claims and set in context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{Error: "invalid token claims"})
			return
		}

		c.Set("userID", claims["user_id"])
		c.Set("username", claims["username"])
		c.Next()
	}
}
