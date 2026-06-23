package main

//go:generate go run github.com/swaggo/swag/cmd/swag@latest init

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/handlers"
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
	r := gin.Default()

	// Routes
	r.POST("/api/start", handlers.StartServerHandler)
	r.POST("/api/stop", handlers.StopServerHandler)
	r.GET("/api/status", handlers.StatusHandler)

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
