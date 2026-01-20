package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"
	"backend/utils/logger"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(jwtUtils *utils.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		// 验证Token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证Token
		claims, err := jwtUtils.ValidateToken(tokenString)
		if err != nil {
			if err == utils.ErrTokenExpired {
				utils.Unauthorized(c, "认证令牌已过期")
			} else {
				utils.Unauthorized(c, "无效的认证令牌")
			}
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := utils.GenerateUUID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// 记录请求开始
		logger.Info("Request started", map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"request_id": requestID,
			"client_ip":  c.ClientIP(),
		})

		// 处理请求
		c.Next()

		// 记录请求完成
		logger.Info("Request completed", map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"request_id": requestID,
		})
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware(logService *services.LogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := string(debug.Stack())
				
				// 记录系统错误日志
				message := fmt.Sprintf("System Panic: %v", err)
				logService.LogSystem(models.LogLevelError, "system", message, stack)
				
				// 记录到文件日志
				logger.Error("System panic recovered", map[string]interface{}{
					"error": err,
					"stack": stack,
					"path":  c.Request.URL.Path,
				})

				// 返回错误信息
				utils.Fail(c, 500, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RateLimitMiddleware 限流中间件（简化版）
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现Redis限流
		c.Next()
	}
}
