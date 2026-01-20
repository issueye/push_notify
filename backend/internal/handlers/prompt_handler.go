package handlers

import (
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type PromptHandler struct {
	promptService *services.PromptService
	logService    *services.LogService
}

func NewPromptHandler(promptService *services.PromptService, logService *services.LogService) *PromptHandler {
	return &PromptHandler{
		promptService: promptService,
		logService:    logService,
	}
}

// List 获取提示词列表
func (h *PromptHandler) List(c *gin.Context) {
	page := utils.GetPage(c)
	size := utils.GetSize(c)
	keyword := c.Query("keyword")
	promptType := c.Query("type")
	scene := c.Query("scene")

	prompts, total, err := h.promptService.GetList(page, size, keyword, promptType, scene)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithPage(c, prompts, int(total), page, size)
}

// Create 创建提示词
func (h *PromptHandler) Create(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	prompt, err := h.promptService.Create(data)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	h.logService.LogOperation(utils.GetUserID(c), "prompt", "创建提示词", "prompt", prompt.ID, data)
	utils.SuccessWithMsg(c, "创建成功", prompt)
}

// Detail 获取提示词详情
func (h *PromptHandler) Detail(c *gin.Context) {
	id := utils.GetID(c)
	prompt, err := h.promptService.GetByID(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, prompt)
}

// Update 更新提示词
func (h *PromptHandler) Update(c *gin.Context) {
	id := utils.GetID(c)
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	err := h.promptService.Update(id, data)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	h.logService.LogOperation(utils.GetUserID(c), "prompt", "更新提示词", "prompt", id, data)
	utils.SuccessWithMsg(c, "更新成功", nil)
}

// Delete 删除提示词
func (h *PromptHandler) Delete(c *gin.Context) {
	id := utils.GetID(c)
	err := h.promptService.Delete(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	h.logService.LogOperation(utils.GetUserID(c), "prompt", "删除提示词", "prompt", id, nil)
	utils.SuccessWithMsg(c, "删除成功", nil)
}

// Test 测试提示词
func (h *PromptHandler) Test(c *gin.Context) {
	id := utils.GetID(c)
	var req struct {
		TestData map[string]interface{} `json:"test_data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	result, err := h.promptService.Test(id, req.TestData)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, map[string]interface{}{
		"result": result,
	})
}

// Rollback 回滚版本
func (h *PromptHandler) Rollback(c *gin.Context) {
	id := utils.GetID(c)
	var req struct {
		Version int `json:"version" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	err := h.promptService.Rollback(id, req.Version)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "回滚成功", nil)
}

// History 获取历史列表
func (h *PromptHandler) History(c *gin.Context) {
	id := utils.GetID(c)
	history, err := h.promptService.GetHistoryList(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, history)
}
