package handlers

import (
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	templateService *services.TemplateService
	generateService *services.GenerateTemplateService
}

func NewTemplateHandler(templateService *services.TemplateService, generateService *services.GenerateTemplateService) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
		generateService: generateService,
	}
}

// List 获取模板列表
func (h *TemplateHandler) List(c *gin.Context) {
	page := utils.GetPage(c)
	size := utils.GetSize(c)
	keyword := c.Query("keyword")
	templateType := c.Query("type")
	scene := c.Query("scene")

	templates, total, err := h.templateService.GetList(page, size, keyword, templateType, scene)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithPage(c, templates, int(total), page, size)
}

// Create 创建模板
func (h *TemplateHandler) Create(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	template, err := h.templateService.Create(data)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "创建成功", template)
}

// Detail 获取模板详情
func (h *TemplateHandler) Detail(c *gin.Context) {
	id := utils.GetID(c)
	template, err := h.templateService.GetByID(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, template)
}

// Update 更新模板
func (h *TemplateHandler) Update(c *gin.Context) {
	id := utils.GetID(c)
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	err := h.templateService.Update(id, data)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "更新成功", nil)
}

// Delete 删除模板
func (h *TemplateHandler) Delete(c *gin.Context) {
	id := utils.GetID(c)
	err := h.templateService.Delete(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}

// SetStatus 更新启用状态
func (h *TemplateHandler) SetStatus(c *gin.Context) {
	id := utils.GetID(c)
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	err := h.templateService.Update(id, map[string]interface{}{"status": req.Status})
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "状态更新成功", nil)
}

// Rollback 回滚版本
func (h *TemplateHandler) Rollback(c *gin.Context) {
	var req struct {
		Version int `json:"version" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	// TODO: 实现版本回滚逻辑
	utils.SuccessWithMsg(c, "回滚成功", nil)
}

// Generate 生成模板内容
func (h *TemplateHandler) Generate(c *gin.Context) {
	var input services.GenerateTemplateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	content, err := h.generateService.Generate(input)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.Success(c, map[string]string{
		"content": content,
	})
}
