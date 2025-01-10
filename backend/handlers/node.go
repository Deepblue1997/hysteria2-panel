package handlers

import (
	"net/http"
	"strconv"

	"hysteria2-panel/models"
	"hysteria2-panel/services"

	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	nodeService *services.NodeService
}

func NewNodeHandler(nodeService *services.NodeService) *NodeHandler {
	return &NodeHandler{nodeService: nodeService}
}

// 创建节点
func (h *NodeHandler) CreateNode(c *gin.Context) {
	var node models.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.nodeService.CreateNode(&node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "节点创建成功", "node": node})
}

// 获取节点列表
func (h *NodeHandler) GetNodes(c *gin.Context) {
	nodes, err := h.nodeService.GetNodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nodes": nodes})
}

// 更新节点状态
func (h *NodeHandler) UpdateNodeStatus(c *gin.Context) {
	nodeID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的节点ID"})
		return
	}

	var status models.NodeStatus
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.nodeService.UpdateNodeStatus(uint(nodeID), &status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "节点状态更新成功"})
}

// 获取节点状态
func (h *NodeHandler) GetNodeStatus(c *gin.Context) {
	nodeID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的节点ID"})
		return
	}

	status, err := h.nodeService.GetNodeStatus(uint(nodeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}
