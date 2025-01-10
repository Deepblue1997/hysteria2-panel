package handlers

import (
	"net/http"

	"hysteria2-panel/services"

	"github.com/gin-gonic/gin"
)

type Hysteria2Handler struct {
	configManager *services.ConfigManagerService
	hy2Service    *services.Hysteria2Service
}

func NewHysteria2Handler(configManager *services.ConfigManagerService, hy2Service *services.Hysteria2Service) *Hysteria2Handler {
	return &Hysteria2Handler{
		configManager: configManager,
		hy2Service:    hy2Service,
	}
}

func (h *Hysteria2Handler) GenerateServerConfig(c *gin.Context) {
	id := c.Param("id")

	config, err := h.configManager.GetUserConfig(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.hy2Service.GenerateServerConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "服务器配置生成成功"})
}

func (h *Hysteria2Handler) GetClientConfig(c *gin.Context) {
	id := c.Param("id")
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain参数不能为空"})
		return
	}

	config, err := h.configManager.GetUserConfig(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientConfig, err := h.hy2Service.GenerateClientConfig(domain, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"config": clientConfig})
}
