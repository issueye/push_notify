package utils

import (
	"net/http"

	"backend/internal/models"
	"backend/utils/logger"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Details   []string    `json:"details,omitempty"`
	RequestID string      `json:"request_id"`
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      200,
		Message:   "success",
		Data:      data,
		RequestID: GetRequestID(c),
	})
}

// 成功响应（带消息）
func SuccessWithMsg(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      200,
		Message:   message,
		Data:      data,
		RequestID: GetRequestID(c),
	})
}

// 分页响应
func SuccessWithPage(c *gin.Context, list interface{}, total, page, size int) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data: map[string]interface{}{
			"list": list,
			"pagination": map[string]interface{}{
				"page":        page,
				"size":        size,
				"total":       total,
				"total_pages": (total + size - 1) / size,
			},
		},
		RequestID: GetRequestID(c),
	})
}

// 失败响应
func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:      code,
		Message:   message,
		RequestID: GetRequestID(c),
	})
}

// 失败响应（带详情）
func FailWithDetails(c *gin.Context, code int, message string, details []string) {
	c.JSON(http.StatusOK, Response{
		Code:      code,
		Message:   message,
		Details:   details,
		RequestID: GetRequestID(c),
	})
}

// 参数错误
func ValidateError(c *gin.Context, details []string) {
	FailWithDetails(c, 400, "参数错误", details)
}

// 未授权
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:      401,
		Message:   message,
		RequestID: GetRequestID(c),
	})
}

// 无权限
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:      403,
		Message:   message,
		RequestID: GetRequestID(c),
	})
}

// 资源不存在
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:      404,
		Message:   message,
		RequestID: GetRequestID(c),
	})
}

// 服务器错误
func ServerError(c *gin.Context, err error) {
	logger.Error("Server error", map[string]interface{}{
		"error": err.Error(),
	})
	c.JSON(http.StatusInternalServerError, Response{
		Code:      500,
		Message:   "服务器内部错误",
		RequestID: GetRequestID(c),
	})
}

// GetRequestID 获取请求ID
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		return id.(string)
	}
	return ""
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) uint {
	if id, exists := c.Get("user_id"); exists {
		return id.(uint)
	}
	return 0
}

// GetUserRole 从上下文获取用户角色
func GetUserRole(c *gin.Context) string {
	if role, exists := c.Get("role"); exists {
		return role.(string)
	}
	return ""
}

// IsAdmin 判断是否为管理员
func IsAdmin(c *gin.Context) bool {
	return GetUserRole(c) == models.RoleAdmin
}
