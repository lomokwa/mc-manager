package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/handlers"
)

func main() {
	r := gin.Default()

	// Routes
	r.POST("/api/start", handlers.StartServerHandler)

	r.Run()
}
