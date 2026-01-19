package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type PushRepo struct {
	db *gorm.DB
}

func NewPushRepo(db *gorm.DB) *PushRepo {
	return &PushRepo{db: db}
}

// Create 创建推送记录
func (r *PushRepo) Create(push *models.Push) error {
	return r.db.Create(push).Error
}

// GetByID 根据ID获取推送记录
func (r *PushRepo) GetByID(id uint) (*models.Push, error) {
	var push models.Push
	err := r.db.Preload("Repo").Preload("Target").Preload("Template").First(&push, id).Error
	if err != nil {
		return nil, err
	}
	return &push, nil
}

// GetList 获取推送记录列表
func (r *PushRepo) GetList(page, size int, repoID, targetID uint, status, keyword, startTime, endTime string) ([]models.Push, int64) {
	var pushes []models.Push
	var total int64

	query := r.db.Model(&models.Push{})
	if repoID > 0 {
		query = query.Where("repo_id = ?", repoID)
	}
	if targetID > 0 {
		query = query.Where("target_id = ?", targetID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("commit_msg LIKE ?", "%"+keyword+"%")
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&pushes)

	return pushes, total
}

// Update 更新推送记录
func (r *PushRepo) Update(push *models.Push) error {
	return r.db.Save(push).Error
}

// UpdateStatus 更新推送状态
func (r *PushRepo) UpdateStatus(id uint, status, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}
	if status == models.PushStatusSuccess {
		updates["pushed_at"] = gorm.Expr("datetime('now')")
	}

	return r.db.Model(&models.Push{}).Where("id = ?", id).Updates(updates).Error
}

// IncrementRetryCount 增加重试次数
func (r *PushRepo) IncrementRetryCount(id uint) error {
	return r.db.Model(&models.Push{}).Where("id = ?", id).Update("retry_count", gorm.Expr("retry_count + 1")).Error
}

// GetStats 获取统计数据
func (r *PushRepo) GetStats(startDate, endDate string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 今日统计
	var todayTotal, todaySuccess, todayFailed int64
	r.db.Model(&models.Push{}).
		Where("date(created_at) = date('now')").
		Count(&todayTotal)
	r.db.Model(&models.Push{}).
		Where("date(created_at) = date('now') AND status = ?", models.PushStatusSuccess).
		Count(&todaySuccess)
	r.db.Model(&models.Push{}).
		Where("date(created_at) = date('now') AND status = ?", models.PushStatusFailed).
		Count(&todayFailed)

	stats["today"] = map[string]interface{}{
		"total":   todayTotal,
		"success": todaySuccess,
		"failed":  todayFailed,
	}

	// 本周统计
	var weekTotal, weekSuccess, weekFailed int64
	r.db.Model(&models.Push{}).
		Where("strftime('%Y-%W', created_at) = strftime('%Y-%W', 'now')").
		Count(&weekTotal)
	r.db.Model(&models.Push{}).
		Where("strftime('%Y-%W', created_at) = strftime('%Y-%W', 'now') AND status = ?", models.PushStatusSuccess).
		Count(&weekSuccess)
	r.db.Model(&models.Push{}).
		Where("strftime('%Y-%W', created_at) = strftime('%Y-%W', 'now') AND status = ?", models.PushStatusFailed).
		Count(&weekFailed)

	stats["this_week"] = map[string]interface{}{
		"total":   weekTotal,
		"success": weekSuccess,
		"failed":  weekFailed,
	}

	// 本月统计
	var monthTotal, monthSuccess, monthFailed int64
	r.db.Model(&models.Push{}).
		Where("strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now')").
		Count(&monthTotal)
	r.db.Model(&models.Push{}).
		Where("strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now') AND status = ?", models.PushStatusSuccess).
		Count(&monthSuccess)
	r.db.Model(&models.Push{}).
		Where("strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now') AND status = ?", models.PushStatusFailed).
		Count(&monthFailed)

	stats["this_month"] = map[string]interface{}{
		"total":   monthTotal,
		"success": monthSuccess,
		"failed":  monthFailed,
	}

	return stats, nil
}

// GetPendingPushes 获取待推送的记录
func (r *PushRepo) GetPendingPushes(limit int) ([]models.Push, error) {
	var pushes []models.Push
	err := r.db.Where("status = ? AND retry_count < 3", models.PushStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&pushes).Error
	return pushes, err
}

// Delete 删除推送记录
func (r *PushRepo) Delete(id uint) error {
	return r.db.Delete(&models.Push{}, id).Error
}

// BatchDelete 批量删除
func (r *PushRepo) BatchDelete(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&models.Push{}).Error
}

// ExistsByCommitAndTarget 检查是否存在相同提交和目标的推送记录
func (r *PushRepo) ExistsByCommitAndTarget(commitID string, targetID uint) bool {
	var count int64
	r.db.Model(&models.Push{}).
		Where("commit_id = ? AND target_id = ?", commitID, targetID).
		Count(&count)
	return count > 0
}

// GetByCommitAndTarget 根据commit_id和target_id获取推送记录
func (r *PushRepo) GetByCommitAndTarget(commitID string, targetID uint) (*models.Push, error) {
	var push models.Push
	err := r.db.Where("commit_id = ? AND target_id = ?", commitID, targetID).First(&push).Error
	if err != nil {
		return nil, err
	}
	return &push, nil
}
