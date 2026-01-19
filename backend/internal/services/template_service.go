package services

import (
	"errors"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/utils/logger"

	"gorm.io/gorm"
)

var (
	ErrTemplateNotFound      = errors.New("模板不存在")
	ErrTemplateAlreadyExists = errors.New("模板名称已存在")
)

type TemplateService struct {
	templateRepo *repository.TemplateRepo
}

func NewTemplateService(db *gorm.DB) *TemplateService {
	return &TemplateService{
		templateRepo: repository.NewTemplateRepo(db),
	}
}

// Create 创建模板
func (s *TemplateService) Create(data map[string]interface{}) (*models.Template, error) {
	name := data["name"].(string)
	templateType := data["type"].(string)
	scene := data["scene"].(string)

	template := &models.Template{
		Name:    name,
		Type:    templateType,
		Scene:   scene,
		Title:   data["title"].(string),
		Content: data["content"].(string),
		Status:  models.StatusActive,
		Version: 1,
	}

	if isDefault, ok := data["is_default"].(bool); ok && isDefault {
		template.IsDefault = true
		s.templateRepo.SetDefault(0, templateType, scene)
	}

	if err := s.templateRepo.Create(template); err != nil {
		return nil, err
	}

	logger.Info("Template created", map[string]interface{}{
		"template_id": template.ID,
		"name":        template.Name,
	})

	return template, nil
}

// GetByID 获取模板详情
func (s *TemplateService) GetByID(id uint) (*models.Template, error) {
	return s.templateRepo.GetByID(id)
}

// GetList 获取模板列表
func (s *TemplateService) GetList(page, size int, keyword, templateType, scene string) ([]models.Template, int64, error) {
	templates, total := s.templateRepo.GetList(page, size, keyword, templateType, scene)
	return templates, total, nil
}

// Update 更新模板
func (s *TemplateService) Update(id uint, data map[string]interface{}) error {
	template, err := s.templateRepo.GetByID(id)
	if err != nil {
		return err
	}

	if name, ok := data["name"].(string); ok && name != "" {
		template.Name = name
	}
	if title, ok := data["title"].(string); ok && title != "" {
		template.Title = title
	}
	if content, ok := data["content"].(string); ok && content != "" {
		template.Content = content
	}
	if status, ok := data["status"].(string); ok && status != "" {
		template.Status = status
	}

	return s.templateRepo.Update(template)
}

// Delete 删除模板
func (s *TemplateService) Delete(id uint) error {
	_, err := s.templateRepo.GetByID(id)
	if err != nil {
		return err
	}
	return s.templateRepo.Delete(id)
}

// SetDefault 设置默认模板
func (s *TemplateService) SetDefault(id uint, templateType, scene string) error {
	return s.templateRepo.SetDefault(id, templateType, scene)
}

// IncrementVersion 增加版本号
func (s *TemplateService) IncrementVersion(id uint) error {
	return s.templateRepo.IncrementVersion(id)
}

// GetDefault 获取默认模板
func (s *TemplateService) GetDefault(templateType, scene string) (*models.Template, error) {
	return s.templateRepo.GetByTypeAndScene(templateType, scene)
}
