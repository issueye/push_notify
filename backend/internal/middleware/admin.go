package middleware

import (
	"backend/internal/models"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AdminOnly 仅管理员可访问中间件
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := utils.GetUserRole(c)
		if role != models.RoleAdmin {
			utils.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
