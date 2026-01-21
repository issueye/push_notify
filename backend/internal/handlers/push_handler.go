package handlers

import (
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type PushHandler struct {
	pushService *services.PushService
}

func NewPushHandler(pushService *services.PushService) *PushHandler {
	return &PushHandler{pushService: pushService}
}

// List 获取推送记录列表
func (h *PushHandler) List(c *gin.Context) {
	page := utils.GetPage(c)
	size := utils.GetSize(c)
	repoID := utils.GetIDParam(c, "repoId")
	targetID := utils.GetIDParam(c, "targetId")
	status := c.Query("status")
	keyword := c.Query("keyword")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	pushes, total, err := h.pushService.GetList(page, size, repoID, targetID, status, keyword, startTime, endTime)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithPage(c, pushes, int(total), page, size)
}

// Detail 获取推送详情
func (h *PushHandler) Detail(c *gin.Context) {
	id := utils.GetID(c)
	push, err := h.pushService.GetByID(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, push)
}

// Delete 删除推送记录
func (h *PushHandler) Delete(c *gin.Context) {
	id := utils.GetID(c)
	err := h.pushService.Delete(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}

// Retry 重试推送
func (h *PushHandler) Retry(c *gin.Context) {
	id := utils.GetID(c)
	newID, err := h.pushService.Retry(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, map[string]interface{}{
		"new_push_id":      newID,
		"original_push_id": id,
	})
}

// BatchRetry 批量重试
func (h *PushHandler) BatchRetry(c *gin.Context) {
	var req struct {
		PushIDs []uint `json:"push_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	total, err := h.pushService.BatchRetry(req.PushIDs)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, map[string]interface{}{
		"total":   total,
		"success": total,
	})
}

// BatchDelete 批量删除
func (h *PushHandler) BatchDelete(c *gin.Context) {
	var req struct {
		PushIDs    []uint `json:"push_ids"`
		BeforeDate string `json:"before_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	err := h.pushService.BatchDelete(req.PushIDs, req.BeforeDate)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}

// GetStats 获取统计
func (h *PushHandler) GetStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.pushService.GetStats(startDate, endDate)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, stats)
}
