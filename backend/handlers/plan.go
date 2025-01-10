package handlers

import (
	"net/http"
	"strconv"

	"hysteria2-panel/models"
	"hysteria2-panel/services"

	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	planService *services.PlanService
}

func NewPlanHandler(planService *services.PlanService) *PlanHandler {
	return &PlanHandler{planService: planService}
}

// 创建套餐
func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var plan models.Plan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.planService.CreatePlan(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "套餐创建成功", "plan": plan})
}

// 获取套餐列表
func (h *PlanHandler) GetPlans(c *gin.Context) {
	plans, err := h.planService.GetPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

// 更新套餐
func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的套餐ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.planService.UpdatePlan(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "套餐更新成功"})
}

// 订阅套餐
func (h *PlanHandler) Subscribe(c *gin.Context) {
	planID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的套餐ID"})
		return
	}

	// 从JWT中获取用户ID
	userID := uint(1) // TODO: 从JWT中获取

	if err := h.planService.Subscribe(userID, uint(planID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "订阅成功"})
}

// 创建订单
func (h *PlanHandler) CreateOrder(c *gin.Context) {
	planID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的套餐ID"})
		return
	}

	// 从JWT中获取用户ID
	userID := uint(1) // TODO: 从JWT中获取

	order, err := h.planService.CreateOrder(userID, uint(planID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order": order})
}
