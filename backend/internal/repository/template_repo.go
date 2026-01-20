package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type TemplateRepo struct {
	db *gorm.DB
}

func NewTemplateRepo(db *gorm.DB) *TemplateRepo {
	return &TemplateRepo{db: db}
}

// Create 创建模板
func (r *TemplateRepo) Create(template *models.Template) error {
	return r.db.Create(template).Error
}

// GetByID 根据ID获取模板
func (r *TemplateRepo) GetByID(id uint) (*models.Template, error) {
	var template models.Template
	err := r.db.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByNameTypeScene 根据名称、类型和场景获取模板
func (r *TemplateRepo) GetByNameTypeScene(name, templateType, scene string) (*models.Template, error) {
	var template models.Template
	err := r.db.Where("name = ? AND type = ? AND scene = ?", name, templateType, scene).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByTypeAndScene 根据类型和场景获取默认模板
func (r *TemplateRepo) GetByTypeAndScene(templateType, scene string) (*models.Template, error) {
	var template models.Template
	err := r.db.Where("type = ? AND scene = ? AND status = ?", templateType, scene, models.StatusActive).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetList 获取模板列表
func (r *TemplateRepo) GetList(page, size int, keyword, templateType, scene string) ([]models.Template, int64) {
	var templates []models.Template
	var total int64

	query := r.db.Model(&models.Template{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if templateType != "" {
		query = query.Where("type = ?", templateType)
	}
	if scene != "" {
		query = query.Where("scene = ?", scene)
	}

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&templates)

	return templates, total
}

// Update 更新模板
func (r *TemplateRepo) Update(template *models.Template) error {
	return r.db.Save(template).Error
}

// IncrementVersion 增加版本号
func (r *TemplateRepo) IncrementVersion(id uint) error {
	return r.db.Model(&models.Template{}).Where("id = ?", id).Update("version", gorm.Expr("version + 1")).Error
}

// Delete 删除模板
func (r *TemplateRepo) Delete(id uint) error {
	return r.db.Delete(&models.Template{}, id).Error
}

// SetDefault 设置默认模板
func (r *TemplateRepo) SetDefault(id uint, templateType, scene string) error {
	// 先取消同类型场景的默认状态
	r.db.Model(&models.Template{}).
		Where("type = ? AND scene = ?", templateType, scene).
		Update("is_default", false)

	// 设置新的默认模板
	return r.db.Model(&models.Template{}).Where("id = ?", id).Update("is_default", true).Error
}
