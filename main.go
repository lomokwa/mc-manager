package main

//go:generate go run github.com/swaggo/swag/cmd/swag@latest init

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lomokwa/mc-manager/handlers"
	"github.com/lomokwa/mc-manager/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lomokwa/mc-manager/docs"
)

// @title MC Manager API
// @version 1.0
// @description API for managing a Minecraft server
// @host localhost:8080
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment")
	}

	r := gin.Default()

	// Rate limiter: 10 requests/sec, burst of 20
	limiter := middleware.NewRateLimiter(10, 20)
	r.Use(limiter.Middleware())
	r.Use(middleware.ValidateAPIKey())

	// Routes
	api := r.Group("/api", middleware.ValidateAPIKey())
	api.POST("/start", handlers.StartServerHandler)
	api.POST("/stop", handlers.StopServerHandler)

	// Console WebSocket
	api.GET("/console", handlers.ConsoleHandler)

	// Server Health check
	api.GET("/status", handlers.StatusHandler)

	// Serve API Docs
	r.GET("/api/docs/*any", func(c *gin.Context) {
		if c.Param("any") == "/" || c.Param("any") == "" {
			c.Redirect(http.StatusMovedPermanently, "/api/docs/index.html")
			return
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1), ginSwagger.URL("/api/docs/doc.json"))(c)
	})

	r.Run()
}
