package handlers

import (
	"net/http"

	"hysteria2-panel/models"
	"hysteria2-panel/services"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	settingService *services.SettingService
}

func NewSettingHandler(settingService *services.SettingService) *SettingHandler {
	return &SettingHandler{settingService: settingService}
}

// 获取TLS配置
func (h *SettingHandler) GetTLSConfig(c *gin.Context) {
	config, err := h.settingService.GetTLSConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"config": config})
}

// 更新TLS配置
func (h *SettingHandler) UpdateTLSConfig(c *gin.Context) {
	var config models.TLSConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.settingService.UpdateTLSConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "TLS配置更新成功"})
}

// 获取SMTP配置
func (h *SettingHandler) GetSMTPConfig(c *gin.Context) {
	config, err := h.settingService.GetSMTPConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"config": config})
}

// 更新SMTP配置
func (h *SettingHandler) UpdateSMTPConfig(c *gin.Context) {
	var config models.SMTPConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.settingService.UpdateSMTPConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMTP配置更新成功"})
}

// 获取系统公告
func (h *SettingHandler) GetAnnouncement(c *gin.Context) {
	announcement, err := h.settingService.GetAnnouncement()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"announcement": announcement})
}

// 更新系统公告
func (h *SettingHandler) UpdateAnnouncement(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.settingService.UpdateAnnouncement(req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "系统公告更新成功"})
}
