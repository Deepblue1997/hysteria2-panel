package handlers

import (
	"net/http"

	"hysteria2-panel/models"
	"hysteria2-panel/services"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	configManager *services.ConfigManagerService
}

func NewConfigHandler(configManager *services.ConfigManagerService) *ConfigHandler {
	return &ConfigHandler{configManager: configManager}
}

func (h *ConfigHandler) GetUserConfig(c *gin.Context) {
	id := c.Param("id")

	config, err := h.configManager.GetUserConfig(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"config": config})
}

func (h *ConfigHandler) UpdateUserConfig(c *gin.Context) {
	id := c.Param("id")
	var config models.UserConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.configManager.UpdateUserConfig(id, &config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "配置更新成功", "config": config})
}
