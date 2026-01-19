package handlers

import (
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type ModelHandler struct {
	modelService *services.AIModelService
}

func NewModelHandler(modelService *services.AIModelService) *ModelHandler {
	return &ModelHandler{modelService: modelService}
}

// List 获取模型列表
func (h *ModelHandler) List(c *gin.Context) {
	page := utils.GetPage(c)
	size := utils.GetSize(c)
	keyword := c.Query("keyword")
	provider := c.Query("provider")

	models, total, err := h.modelService.GetList(page, size, keyword, provider)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithPage(c, models, int(total), page, size)
}

// Create 创建模型
func (h *ModelHandler) Create(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	model, err := h.modelService.Create(data)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "创建成功", model)
}

// Detail 获取模型详情
func (h *ModelHandler) Detail(c *gin.Context) {
	id := utils.GetID(c)
	model, err := h.modelService.GetByID(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, model)
}

// Update 更新模型
func (h *ModelHandler) Update(c *gin.Context) {
	id := utils.GetID(c)
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	err := h.modelService.Update(id, data)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "更新成功", nil)
}

// Delete 删除模型
func (h *ModelHandler) Delete(c *gin.Context) {
	id := utils.GetID(c)
	err := h.modelService.Delete(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}

// SetDefault 设置默认模型
func (h *ModelHandler) SetDefault(c *gin.Context) {
	id := utils.GetID(c)
	err := h.modelService.SetDefault(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "设置成功", nil)
}

// GetLogs 获取调用日志
func (h *ModelHandler) GetLogs(c *gin.Context) {
	page := utils.GetPage(c)
	size := utils.GetSize(c)

	// TODO: 从日志服务获取
	utils.SuccessWithPage(c, []interface{}{}, 0, page, size)
}

// Verify 验证模型配置
func (h *ModelHandler) Verify(c *gin.Context) {
	id := utils.GetID(c)
	model, err := h.modelService.GetByID(id)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	// TODO: 验证API配置
	utils.Success(c, map[string]interface{}{
		"valid": true,
		"model": model.Name,
	})
}
