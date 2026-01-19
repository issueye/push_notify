package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type AIModelRepo struct {
	db *gorm.DB
}

func NewAIModelRepo(db *gorm.DB) *AIModelRepo {
	return &AIModelRepo{db: db}
}

// Create 创建AI模型
func (r *AIModelRepo) Create(model *models.AIModel) error {
	return r.db.Create(model).Error
}

// GetByID 根据ID获取AI模型
func (r *AIModelRepo) GetByID(id uint) (*models.AIModel, error) {
	var model models.AIModel
	err := r.db.First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// GetDefault 获取默认模型
func (r *AIModelRepo) GetDefault() (*models.AIModel, error) {
	var model models.AIModel
	err := r.db.Where("is_default = ? AND status = ?", true, models.StatusActive).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// GetList 获取AI模型列表
func (r *AIModelRepo) GetList(page, size int, keyword, provider string) ([]models.AIModel, int64) {
	var modelList []models.AIModel
	var total int64

	query := r.db.Model(&models.AIModel{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR type LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if provider != "" {
		query = query.Where("provider = ?", provider)
	}

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&modelList)

	return modelList, total
}

// Update 更新AI模型
func (r *AIModelRepo) Update(model *models.AIModel) error {
	return r.db.Save(model).Error
}

// IncrementCallCount 增加调用次数
func (r *AIModelRepo) IncrementCallCount(id uint) error {
	return r.db.Model(&models.AIModel{}).Where("id = ?", id).Update("call_count", gorm.Expr("call_count + 1")).Error
}

// Delete 删除AI模型
func (r *AIModelRepo) Delete(id uint) error {
	return r.db.Delete(&models.AIModel{}, id).Error
}

// SetDefault 设置默认模型
func (r *AIModelRepo) SetDefault(id uint) error {
	// 先取消当前默认
	r.db.Model(&models.AIModel{}).Where("is_default = ?", true).Update("is_default", false)
	// 设置新的默认模型
	return r.db.Model(&models.AIModel{}).Where("id = ?", id).Update("is_default", true).Error
}
