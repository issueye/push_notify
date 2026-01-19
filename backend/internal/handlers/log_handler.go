package handlers

import (
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logService *services.LogService
}

func NewLogHandler(logService *services.LogService) *LogHandler {
	return &LogHandler{logService: logService}
}

// GetSystemLogs 获取系统日志
func (h *LogHandler) GetSystemLogs(c *gin.Context) {
	page := utils.GetPage(c)
	size := utils.GetSize(c)
	level := c.Query("level")
	keyword := c.Query("keyword")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	logs, total, err := h.logService.GetSystemLogs(page, size, level, keyword, startTime, endTime)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithPage(c, logs, int(total), page, size)
}

// GetOperationLogs 获取操作日志
func (h *LogHandler) GetOperationLogs(c *gin.Context) {
	page := utils.GetPage(c)
	size := utils.GetSize(c)
	userID := utils.GetIDParam(c, "userId")
	action := c.Query("action")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	logs, total, err := h.logService.GetOperationLogs(page, size, userID, action, startTime, endTime)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithPage(c, logs, int(total), page, size)
}

// GetAICallLogs 获取AI调用日志
func (h *LogHandler) GetAICallLogs(c *gin.Context) {
	page := utils.GetPage(c)
	size := utils.GetSize(c)
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	logs, total, err := h.logService.GetAICallLogs(page, size, startTime, endTime)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithPage(c, logs, int(total), page, size)
}

// Search 日志搜索
func (h *LogHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	page := utils.GetPage(c)
	size := utils.GetSize(c)

	// 搜索所有日志类型
	logs, total, _ := h.logService.GetSystemLogs(page, size, "", keyword, "", "")
	utils.SuccessWithPage(c, logs, int(total), page, size)
}

// Export 导出日志
func (h *LogHandler) Export(c *gin.Context) {
	logType := c.Query("type")

	// TODO: 实现日志导出
	utils.Success(c, map[string]interface{}{
		"download_url": "/api/v1/logs/export/file",
		"type":         logType,
	})
}

// GetStats 获取统计
func (h *LogHandler) GetStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.logService.GetStats(startDate, endDate)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, stats)
}
