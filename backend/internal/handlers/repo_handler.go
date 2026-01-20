package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type RepoHandler struct {
	repoService *services.RepoService
}

func NewRepoHandler(repoService *services.RepoService) *RepoHandler {
	return &RepoHandler{repoService: repoService}
}

// List 获取仓库列表
func (h *RepoHandler) List(c *gin.Context) {
	page := utils.GetPage(c)
	size := utils.GetSize(c)
	keyword := c.Query("keyword")

	repos, total, err := h.repoService.GetList(page, size, keyword)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithPage(c, repos, int(total), page, size)
}

// Create 创建仓库
func (h *RepoHandler) Create(c *gin.Context) {
	var data models.CreateRepo
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	repo, err := h.repoService.Create(data)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "创建成功", repo)
}

// Detail 获取仓库详情
func (h *RepoHandler) Detail(c *gin.Context) {
	id := utils.GetID(c)
	repo, err := h.repoService.GetByID(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, repo)
}

// Update 更新仓库
func (h *RepoHandler) Update(c *gin.Context) {
	id := utils.GetID(c)
	var data models.UpdateRepo
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	err := h.repoService.Update(id, data)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "更新成功", nil)
}

// Delete 删除仓库
func (h *RepoHandler) Delete(c *gin.Context) {
	id := utils.GetID(c)
	err := h.repoService.Delete(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}

// TestWebhook 测试Webhook
func (h *RepoHandler) TestWebhook(c *gin.Context) {
	id := utils.GetID(c)
	result, err := h.repoService.TestWebhook(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, result)
}

// GetTargets 获取仓库关联的推送目标
func (h *RepoHandler) GetTargets(c *gin.Context) {
	id := utils.GetID(c)
	targets, err := h.repoService.GetTargets(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, targets)
}

// AddTarget 关联推送目标
func (h *RepoHandler) AddTarget(c *gin.Context) {
	id := utils.GetID(c)
	var req struct {
		TargetID uint `json:"target_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	err := h.repoService.AddTarget(id, req.TargetID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "关联成功", nil)
}

// RemoveTarget 取消关联推送目标
func (h *RepoHandler) RemoveTarget(c *gin.Context) {
	id := utils.GetID(c)
	targetID := utils.GetIDParam(c, "targetId")

	err := h.repoService.RemoveTarget(id, targetID)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "取消关联成功", nil)
}
