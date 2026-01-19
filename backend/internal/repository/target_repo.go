package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type TargetRepo struct {
	db *gorm.DB
}

func NewTargetRepo(db *gorm.DB) *TargetRepo {
	return &TargetRepo{db: db}
}

// Create 创建推送目标
func (r *TargetRepo) Create(target *models.Target) error {
	return r.db.Create(target).Error
}

// GetByID 根据ID获取推送目标
func (r *TargetRepo) GetByID(id uint) (*models.Target, error) {
	var target models.Target
	err := r.db.Preload("Repos").First(&target, id).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

// GetList 获取推送目标列表
func (r *TargetRepo) GetList(page, size int, keyword, targetType, scope string) ([]models.Target, int64) {
	var targets []models.Target
	var total int64

	query := r.db.Model(&models.Target{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if targetType != "" {
		query = query.Where("type = ?", targetType)
	}
	if scope != "" {
		query = query.Where("scope = ?", scope)
	}

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&targets)

	return targets, total
}

// Update 更新推送目标
func (r *TargetRepo) Update(target *models.Target) error {
	return r.db.Save(target).Error
}

// Delete 删除推送目标
func (r *TargetRepo) Delete(id uint) error {
	// 先删除关联
	r.db.Where("target_id = ?", id).Delete(&models.RepoTarget{})
	return r.db.Delete(&models.Target{}, id).Error
}

// AddRepo 添加仓库关联
func (r *TargetRepo) AddRepo(targetID, repoID uint) error {
	return r.db.Create(&models.RepoTarget{
		TargetID: targetID,
		RepoID:   repoID,
	}).Error
}

// RemoveRepo 移除仓库关联
func (r *TargetRepo) RemoveRepo(targetID, repoID uint) error {
	return r.db.Where("target_id = ? AND repo_id = ?", targetID, repoID).Delete(&models.RepoTarget{}).Error
}

// GetRepos 获取推送目标关联的仓库
func (r *TargetRepo) GetRepos(targetID uint) ([]models.Repo, error) {
	var repos []models.Repo
	err := r.db.Joins("JOIN repo_targets ON repo_targets.repo_id = repos.id").
		Where("repo_targets.target_id = ?", targetID).
		Find(&repos).Error
	return repos, err
}

// GetByScopeAndRepo 根据范围和仓库获取推送目标
func (r *TargetRepo) GetByScopeAndRepo(repoID uint) ([]models.Target, error) {
	var targets []models.Target

	// 获取全局推送目标
	var globalTargets []models.Target
	r.db.Where("scope = ? AND status = ?", models.TargetScopeGlobal, models.StatusActive).Find(&globalTargets)

	// 获取指定仓库的推送目标
	var repoTargets []models.Target
	r.db.Joins("JOIN repo_targets ON repo_targets.target_id = targets.id").
		Where("repo_targets.repo_id = ? AND targets.status = ?", repoID, models.StatusActive).
		Find(&repoTargets)

	targets = append(globalTargets, repoTargets...)
	return targets, nil
}
