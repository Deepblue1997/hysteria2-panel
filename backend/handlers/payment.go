package handlers

import (
	"net/http"

	"hysteria2-panel/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// 创建支付
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	orderNo := c.Query("order_no")
	method := c.Query("method")
	if orderNo == "" || method == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要参数"})
		return
	}

	paymentURL, err := h.paymentService.CreatePayment(orderNo, method)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_url": paymentURL})
}

// 支付回调
func (h *PaymentHandler) HandleCallback(c *gin.Context) {
	method := c.Param("method")
	params := make(map[string]string)

	// 根据不同支付方式处理参数
	switch method {
	case "alipay":
		for k, v := range c.Request.Form {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	case "wechat":
		// 微信支付使用XML格式
		// TODO: 处理XML数据
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的支付方式"})
		return
	}

	if err := h.paymentService.HandleCallback(method, params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	switch method {
	case "alipay":
		c.String(http.StatusOK, "success")
	case "wechat":
		c.XML(http.StatusOK, gin.H{"return_code": "SUCCESS"})
	}
}

// 查询支付状态
func (h *PaymentHandler) QueryPaymentStatus(c *gin.Context) {
	orderNo := c.Query("order_no")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少订单号"})
		return
	}

	paid, err := h.paymentService.QueryPaymentStatus(orderNo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"paid": paid})
}
