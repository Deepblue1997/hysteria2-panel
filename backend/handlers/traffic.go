package handlers

import (
	"net/http"

	"hysteria2-panel/services"

	"github.com/gin-gonic/gin"
)

type TrafficHandler struct {
	trafficService *services.TrafficService
}

func NewTrafficHandler(trafficService *services.TrafficService) *TrafficHandler {
	return &TrafficHandler{trafficService: trafficService}
}

// 记录流量使用
func (h *TrafficHandler) RecordTraffic(c *gin.Context) {
	type TrafficRecord struct {
		UserID   uint  `json:"user_id" binding:"required"`
		Upload   int64 `json:"upload" binding:"required"`
		Download int64 `json:"download" binding:"required"`
	}

	var record TrafficRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.trafficService.RecordTraffic(record.UserID, record.Upload, record.Download); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "流量记录成功"})
}

// 检查用户流量限制
func (h *TrafficHandler) CheckTrafficLimit(c *gin.Context) {
	userID := c.Param("id")

	allowed, err := h.trafficService.CheckTrafficLimit(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"allowed": allowed})
}
