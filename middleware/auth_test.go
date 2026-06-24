package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/types"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestValidateAPIKey_MissingKey(t *testing.T) {
	r := gin.New()
	r.Use(ValidateAPIKey())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, types.APIResponse{Success: true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	var resp types.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error != "missing API key" {
		t.Errorf("expected 'missing API key', got %q", resp.Error)
	}
}

func TestValidateAPIKey_InvalidKey(t *testing.T) {
	os.Setenv("API_KEY", "correct-key")
	defer os.Unsetenv("API_KEY")

	r := gin.New()
	r.Use(ValidateAPIKey())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, types.APIResponse{Success: true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	var resp types.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error != "invalid API key" {
		t.Errorf("expected 'invalid API key', got %q", resp.Error)
	}
}

func TestValidateAPIKey_ValidKey(t *testing.T) {
	os.Setenv("API_KEY", "correct-key")
	defer os.Unsetenv("API_KEY")

	r := gin.New()
	r.Use(ValidateAPIKey())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, types.APIResponse{Success: true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "correct-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp types.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if !resp.Success {
		t.Error("expected success to be true")
	}
}
