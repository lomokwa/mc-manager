package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/types"
)

func UpdateServerPropertiesHandler(c *gin.Context) {
	var req types.UpdateServerPropertiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, types.APIResponse{Error: "invalid request body"})
		return
	}
}
