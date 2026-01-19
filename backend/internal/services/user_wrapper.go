package services

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

// UserRepoWrapper 用户仓库包装器
type UserRepoWrapper struct {
	db *gorm.DB
}

func NewUserRepoWrapper(db *gorm.DB) *UserRepoWrapper {
	return &UserRepoWrapper{db: db}
}

// List 获取用户列表
func (w *UserRepoWrapper) List(page, size int, keyword, role, status string) ([]models.User, int64) {
	var users []models.User
	var total int64

	query := w.db.Model(&models.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&users)

	return users, total
}
