package services

import (
	"errors"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/utils/logger"

	"gorm.io/gorm"
)

var (
	ErrPromptNotFound      = errors.New("提示词不存在")
	ErrPromptAlreadyExists = errors.New("提示词名称已存在")
)

type PromptService struct {
	promptRepo *repository.PromptRepo
}

func NewPromptService(db *gorm.DB) *PromptService {
	return &PromptService{
		promptRepo: repository.NewPromptRepo(db),
	}
}

// Create 创建提示词
func (s *PromptService) Create(data map[string]interface{}) (*models.Prompt, error) {
	name := data["name"].(string)
	promptType := data["type"].(string)

	prompt := &models.Prompt{
		Name:     name,
		Type:     promptType,
		Scene:    getStringWithDefault(data, "scene", ""),
		Language: getStringWithDefault(data, "language", ""),
		Content:  data["content"].(string),
		Version:  1,
	}

	if modelID, ok := data["model_id"].(uint); ok && modelID > 0 {
		prompt.ModelID = &modelID
	}

	if err := s.promptRepo.Create(prompt); err != nil {
		return nil, err
	}

	logger.Info("Prompt created", map[string]interface{}{
		"prompt_id": prompt.ID,
		"name":      prompt.Name,
		"type":      prompt.Type,
	})

	return prompt, nil
}

// GetByID 获取提示词详情
func (s *PromptService) GetByID(id uint) (*models.Prompt, error) {
	return s.promptRepo.GetByID(id)
}

// GetList 获取提示词列表
func (s *PromptService) GetList(page, size int, keyword, promptType, scene string) ([]models.Prompt, int64, error) {
	prompts, total := s.promptRepo.GetList(page, size, keyword, promptType, scene)
	return prompts, total, nil
}

// Update 更新提示词
func (s *PromptService) Update(id uint, data map[string]interface{}) error {
	prompt, err := s.promptRepo.GetByID(id)
	if err != nil {
		return err
	}

	if name, ok := data["name"].(string); ok && name != "" {
		prompt.Name = name
	}
	if scene, ok := data["scene"].(string); ok {
		prompt.Scene = scene
	}
	if language, ok := data["language"].(string); ok {
		prompt.Language = language
	}
	if content, ok := data["content"].(string); ok && content != "" {
		prompt.Content = content
	}
	if modelID, ok := data["model_id"].(uint); ok {
		prompt.ModelID = &modelID
	}

	return s.promptRepo.Update(prompt)
}

// Delete 删除提示词
func (s *PromptService) Delete(id uint) error {
	return s.promptRepo.Delete(id)
}

// IncrementVersion 增加版本号
func (s *PromptService) IncrementVersion(id uint) error {
	return s.promptRepo.IncrementVersion(id)
}

// GetByTypeAndScene 根据类型和场景获取提示词
func (s *PromptService) GetByTypeAndScene(promptType, scene string) ([]models.Prompt, error) {
	return s.promptRepo.GetByTypeAndScene(promptType, scene)
}
