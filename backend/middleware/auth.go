package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 移除 "Bearer " 前缀
		token = strings.TrimPrefix(token, "Bearer ")

		// TODO: 验证 JWT token
		// userID, err := utils.ValidateToken(token)
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
		// 	c.Abort()
		// 	return
		// }

		// c.Set("userID", userID)
		c.Next()
	}
}
