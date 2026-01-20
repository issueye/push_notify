package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) int {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	return page
}

func GetSize(c *gin.Context) int {
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return size
}

func GetID(c *gin.Context) uint {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	return uint(id)
}

func GetIDParam(c *gin.Context, name string) uint {
	id, _ := strconv.ParseUint(c.Param(name), 10, 32)
	return uint(id)
}

func GetUserID(c *gin.Context) uint {
	userID, _ := c.Get("user_id")
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

func GetRequestID(c *gin.Context) string {
	requestID, _ := c.Get("request_id")
	if id, ok := requestID.(string); ok {
		return id
	}
	return ""
}
