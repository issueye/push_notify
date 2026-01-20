package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type PromptRepo struct {
	db *gorm.DB
}

func NewPromptRepo(db *gorm.DB) *PromptRepo {
	return &PromptRepo{db: db}
}

// WithTx 返回一个使用指定事务 DB 的 PromptRepo
func (r *PromptRepo) WithTx(tx *gorm.DB) *PromptRepo {
	return &PromptRepo{db: tx}
}

// Create 创建提示词
func (r *PromptRepo) Create(prompt *models.Prompt) error {
	return r.db.Create(prompt).Error
}

// GetByID 根据ID获取提示词
func (r *PromptRepo) GetByID(id uint) (*models.Prompt, error) {
	var prompt models.Prompt
	err := r.db.Preload("Model").First(&prompt, id).Error
	if err != nil {
		return nil, err
	}
	return &prompt, nil
}

// GetList 获取提示词列表
func (r *PromptRepo) GetList(page, size int, keyword, promptType, scene string) ([]models.Prompt, int64) {
	var prompts []models.Prompt
	var total int64

	query := r.db.Model(&models.Prompt{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if promptType != "" {
		query = query.Where("type = ?", promptType)
	}
	if scene != "" {
		query = query.Where("scene = ?", scene)
	}

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&prompts)

	return prompts, total
}

// Update 更新提示词
func (r *PromptRepo) Update(prompt *models.Prompt) error {
	return r.db.Save(prompt).Error
}

// IncrementVersion 增加版本号
func (r *PromptRepo) IncrementVersion(id uint) error {
	return r.db.Model(&models.Prompt{}).Where("id = ?", id).Update("version", gorm.Expr("version + 1")).Error
}

// Delete 删除提示词
func (r *PromptRepo) Delete(id uint) error {
	return r.db.Delete(&models.Prompt{}, id).Error
}

// GetByTypeAndScene 根据类型和场景获取提示词
func (r *PromptRepo) GetByTypeAndScene(promptType, scene string) ([]models.Prompt, error) {
	var prompts []models.Prompt
	err := r.db.Where("type = ? AND scene = ?", promptType, scene).Find(&prompts).Error
	return prompts, err
}

// CreateHistory 创建历史记录
func (r *PromptRepo) CreateHistory(history *models.PromptHistory) error {
	return r.db.Create(history).Error
}

// GetHistory 获取历史记录
func (r *PromptRepo) GetHistory(promptID uint, version int) (*models.PromptHistory, error) {
	var history models.PromptHistory
	err := r.db.Where("prompt_id = ? AND version = ?", promptID, version).First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// GetHistoryList 获取历史列表
func (r *PromptRepo) GetHistoryList(promptID uint) ([]models.PromptHistory, error) {
	var history []models.PromptHistory
	err := r.db.Where("prompt_id = ?", promptID).Order("version DESC").Find(&history).Error
	return history, err
}
