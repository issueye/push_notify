package services

import (
	"encoding/json"
	"errors"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/utils/logger"

	"gorm.io/gorm"
)

var (
	ErrModelNotFound      = errors.New("AI模型不存在")
	ErrModelAlreadyExists = errors.New("模型名称已存在")
)

type AIModelService struct {
	modelRepo *repository.AIModelRepo
}

func NewAIModelService(db *gorm.DB) *AIModelService {
	return &AIModelService{
		modelRepo: repository.NewAIModelRepo(db),
	}
}

// Create 创建AI模型
func (s *AIModelService) Create(data map[string]interface{}) (*models.AIModel, error) {
	name := data["name"].(string)

	model := &models.AIModel{
		Name:      name,
		Type:      getStringWithDefault(data, "type", ""),
		Provider:  getStringWithDefault(data, "provider", ""),
		APIURL:    data["api_url"].(string),
		APIKey:    data["api_key"].(string),
		IsDefault: false,
		Status:    models.StatusActive,
		CallCount: 0,
	}

	if params, ok := data["params"].(map[string]interface{}); ok {
		if paramsStr, err := json.Marshal(params); err == nil {
			model.Params = string(paramsStr)
		}
	}

	if err := s.modelRepo.Create(model); err != nil {
		return nil, err
	}

	logger.Info("AIModel created", map[string]interface{}{
		"model_id": model.ID,
		"name":     model.Name,
	})

	return model, nil
}

// GetByID 获取模型详情
func (s *AIModelService) GetByID(id uint) (*models.AIModel, error) {
	return s.modelRepo.GetByID(id)
}

// GetList 获取模型列表
func (s *AIModelService) GetList(page, size int, keyword, provider string) ([]models.AIModel, int64, error) {
	models, total := s.modelRepo.GetList(page, size, keyword, provider)
	return models, total, nil
}

// Update 更新模型
func (s *AIModelService) Update(id uint, data map[string]interface{}) error {
	model, err := s.modelRepo.GetByID(id)
	if err != nil {
		return err
	}

	if name, ok := data["name"].(string); ok && name != "" {
		model.Name = name
	}
	if apiURL, ok := data["api_url"].(string); ok && apiURL != "" {
		model.APIURL = apiURL
	}
	if apiKey, ok := data["api_key"].(string); ok && apiKey != "" {
		model.APIKey = apiKey
	}
	if status, ok := data["status"].(string); ok && status != "" {
		model.Status = status
	}
	if params, ok := data["params"].(map[string]interface{}); ok {
		if paramsStr, err := json.Marshal(params); err == nil {
			model.Params = string(paramsStr)
		}
	}

	return s.modelRepo.Update(model)
}

// Delete 删除模型
func (s *AIModelService) Delete(id uint) error {
	return s.modelRepo.Delete(id)
}

// SetDefault 设置默认模型
func (s *AIModelService) SetDefault(id uint) error {
	return s.modelRepo.SetDefault(id)
}

// GetDefault 获取默认模型
func (s *AIModelService) GetDefault() (*models.AIModel, error) {
	return s.modelRepo.GetDefault()
}

// IncrementCallCount 增加调用次数
func (s *AIModelService) IncrementCallCount(id uint) error {
	return s.modelRepo.IncrementCallCount(id)
}

func getStringWithDefault(data map[string]interface{}, key, defaultVal string) string {
	if val, ok := data[key].(string); ok && val != "" {
		return val
	}
	return defaultVal
}
