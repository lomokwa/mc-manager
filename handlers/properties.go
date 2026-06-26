package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lomokwa/mc-manager/services"
	"github.com/lomokwa/mc-manager/types"
)

// @Summary Update server properties
// @Description Merges the provided key/values into server.properties (falling back to defaults for keys not set yet) and writes the file. Takes effect on the next server start.
// @Tags server
// @Accept json
// @Produce json
// @Param request body types.UpdateServerPropertiesRequest true "Properties to update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/properties [patch]
func UpdateServerPropertiesHandler(c *gin.Context) {
	var req types.UpdateServerPropertiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponse{Error: "invalid request body"})
		return
	}

	if len(req.Properties) == 0 {
		c.JSON(http.StatusBadRequest, types.APIResponse{Error: "no properties provided"})
		return
	}

	if err := services.UpdateServerProperties(req.Properties); err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{Success: true, Data: gin.H{"updated": len(req.Properties)}})
}
