package main

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/start", func(c *gin.Context) {
		cmd := exec.Command("echo", "Starting server...")

		output, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"output": string(output),
		})
	})

	r.Run()
}
