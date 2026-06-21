package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/services"
	"github.com/lomokwa/mc-manager/types"
	"github.com/lomokwa/mc-manager/utils"
)

func StartServerHandler(c *gin.Context) {
	log.Printf("start request received")

	var req types.StartServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if !utils.FileExists("./minecraft-server/server.jar") {
		log.Printf("server.jar not found, downloading latest")
		err := services.DownloadLatestServerJar("./minecraft-server/server.jar")
		if err != nil {
			log.Printf("failed to download server.jar: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("server.jar downloaded successfully")
	}

	log.Printf("creating server files")
	if err := services.PrepareServerFiles("./minecraft-server", req.CreateLaunchScript, req.ConfigureProperties, req.Properties); err != nil {
		log.Printf("failed to prepare server files: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	log.Printf("starting server process")
	output, err := services.StartServerProcess()
	if err != nil {
		log.Printf("failed to start server process: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("server process started")
	c.JSON(http.StatusOK, gin.H{"output": output})
}
