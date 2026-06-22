package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/handlers"
)

func main() {
	r := gin.Default()

	// Routes
	r.POST("/api/start", handlers.StartServerHandler)
	r.POST("/api/stop", handlers.StopServerHandler)
	r.GET("/api/status", handlers.GetServerStatus)
	// TODO: Add GET /api/status route to check if server is running

	r.Run()
}
