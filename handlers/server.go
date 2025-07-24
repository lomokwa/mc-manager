package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/services"
	"github.com/lomokwa/mc-manager/utils"
)

func StartServerHandler(c *gin.Context) {
	if !utils.FileExists("./minecraft-server/server.jar") {
		err := services.DownloadLatestServerJar("./minecraft-server/server.jar")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Server jar downloaded successfully"})
	}

	output, err := services.StartServerProcess()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"output": output})
}
