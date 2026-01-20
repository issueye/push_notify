package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type RepoRepo struct {
	db *gorm.DB
}

func NewRepoRepo(db *gorm.DB) *RepoRepo {
	return &RepoRepo{db: db}
}

// WithTx 返回一个使用指定事务 DB 的 RepoRepo
func (r *RepoRepo) WithTx(tx *gorm.DB) *RepoRepo {
	return &RepoRepo{db: tx}
}

// Create 创建仓库
func (r *RepoRepo) Create(repo *models.Repo) error {
	return r.db.Create(repo).Error
}

// GetByID 根据ID获取仓库
func (r *RepoRepo) GetByID(id uint) (*models.Repo, error) {
	var repo models.Repo
	err := r.db.Preload("Model").
		Preload("Targets").
		Preload("CommitTemplate").
		Preload("ReviewTemplates").
		Preload("ReviewTemplates.Template").
		First(&repo, id).Error
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// GetByWebhookID 根据WebhookID获取仓库
func (r *RepoRepo) GetByWebhookID(webhookID string) (*models.Repo, error) {
	var repo models.Repo
	err := r.db.Where("webhook_id = ?", webhookID).
		Preload("CommitTemplate").
		Preload("ReviewTemplates").
		Preload("ReviewTemplates.Template").
		First(&repo).Error
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// GetByName 根据名称获取仓库
func (r *RepoRepo) GetByName(name string) (*models.Repo, error) {
	var repo models.Repo
	err := r.db.Where("name = ?", name).First(&repo).Error
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// GetList 获取仓库列表
func (r *RepoRepo) GetList(page, size int, keyword string) ([]models.Repo, int64) {
	var repos []models.Repo
	var total int64

	query := r.db.Model(&models.Repo{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 加载关联数据
	query.Preload("RepoTargets").
		Preload("CommitTemplate").
		Preload("ReviewTemplates").
		Preload("ReviewTemplates.Template")

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&repos)

	return repos, total
}

// Update 更新仓库
func (r *RepoRepo) Update(repo *models.Repo) error {
	updates := map[string]interface{}{
		"name":               repo.Name,
		"url":                repo.URL,
		"type":               repo.Type,
		"status":             repo.Status,
		"model_id":           repo.ModelID,
		"commit_template_id": repo.CommitTemplateID,
	}
	if repo.AccessToken != "" {
		updates["access_token"] = repo.AccessToken
	}
	return r.db.Model(&models.Repo{}).Where("id = ?", repo.ID).Updates(updates).Error
}

// Delete 删除仓库
func (r *RepoRepo) Delete(id uint) error {
	// 先删除关联
	r.db.Where("repo_id = ?", id).Delete(&models.RepoTarget{})
	r.db.Where("repo_id = ?", id).Delete(&models.RepoTemplate{})
	return r.db.Delete(&models.Repo{}, id).Error
}

// AddTarget 添加推送目标关联
func (r *RepoRepo) AddTarget(repoID, targetID uint) error {
	return r.db.Create(&models.RepoTarget{
		RepoID:   repoID,
		TargetID: targetID,
	}).Error
}

// RemoveTarget 移除推送目标关联
func (r *RepoRepo) RemoveTarget(repoID, targetID uint) error {
	return r.db.Where("repo_id = ? AND target_id = ?", repoID, targetID).Delete(&models.RepoTarget{}).Error
}

func (r *RepoRepo) DeleteTargets(repoID uint) error {
	return r.db.Where("repo_id = ?", repoID).Delete(&models.RepoTarget{}).Error
}

// InsertTargets
func (r *RepoRepo) InsertTargets(repoID uint, targetIDs []uint) error {
	insertData := make([]models.RepoTarget, len(targetIDs))
	for i, targetID := range targetIDs {
		insertData[i] = models.RepoTarget{
			RepoID:   repoID,
			TargetID: targetID,
		}
	}
	return r.db.Create(&insertData).Error
}

func (r *RepoRepo) DeleteReviewTemplates(repoID uint) error {
	return r.db.Where("repo_id = ?", repoID).Delete(&models.RepoTemplate{}).Error
}

func (r *RepoRepo) InsertReviewTemplates(repoID uint, configs []models.RepoTemplateConfig) error {
	insertData := make([]models.RepoTemplate, len(configs))
	for i, config := range configs {
		insertData[i] = models.RepoTemplate{
			RepoID:     repoID,
			TemplateID: config.TemplateID,
			Language:   config.Language,
		}
	}
	return r.db.Create(&insertData).Error
}

// GetTargets 获取仓库关联的推送目标
func (r *RepoRepo) GetTargets(repoID uint) ([]models.Target, error) {
	var targets []models.Target
	err := r.db.Joins("JOIN repo_targets ON repo_targets.target_id = targets.id").
		Where("repo_targets.repo_id = ?", repoID).
		Find(&targets).Error
	return targets, err
}
