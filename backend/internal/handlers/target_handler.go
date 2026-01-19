package handlers

import (
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type TargetHandler struct {
	targetService *services.TargetService
}

func NewTargetHandler(targetService *services.TargetService) *TargetHandler {
	return &TargetHandler{targetService: targetService}
}

// List 获取推送目标列表
func (h *TargetHandler) List(c *gin.Context) {
	page := utils.GetPage(c)
	size := utils.GetSize(c)
	keyword := c.Query("keyword")
	targetType := c.Query("type")
	scope := c.Query("scope")

	targets, total, err := h.targetService.GetList(page, size, keyword, targetType, scope)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithPage(c, targets, int(total), page, size)
}

// Create 创建推送目标
func (h *TargetHandler) Create(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	target, err := h.targetService.Create(data)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "创建成功", target)
}

// Detail 获取推送目标详情
func (h *TargetHandler) Detail(c *gin.Context) {
	id := utils.GetID(c)
	target, err := h.targetService.GetByID(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, target)
}

// Update 更新推送目标
func (h *TargetHandler) Update(c *gin.Context) {
	id := utils.GetID(c)
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	err := h.targetService.Update(id, data)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "更新成功", nil)
}

// Delete 删除推送目标
func (h *TargetHandler) Delete(c *gin.Context) {
	id := utils.GetID(c)
	err := h.targetService.Delete(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}

// Test 发送测试消息
func (h *TargetHandler) Test(c *gin.Context) {
	id := utils.GetID(c)
	result, err := h.targetService.Test(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "测试消息已发送", result)
}

// AddRepo 关联仓库
func (h *TargetHandler) AddRepo(c *gin.Context) {
	id := utils.GetID(c)
	var req struct {
		RepoIDs []uint `json:"repo_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	for _, repoID := range req.RepoIDs {
		h.targetService.AddRepo(id, repoID)
	}

	utils.SuccessWithMsg(c, "关联成功", nil)
}

// RemoveRepo 取消关联仓库
func (h *TargetHandler) RemoveRepo(c *gin.Context) {
	id := utils.GetID(c)
	repoID := utils.GetIDParam(c, "repoId")

	err := h.targetService.RemoveRepo(id, repoID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "取消关联成功", nil)
}
