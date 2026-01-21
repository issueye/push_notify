package repository

import (
	"time"

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
	query.Preload("Repo").Preload("Target").Preload("Template").
		Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&pushes)

	return pushes, total
}

// Update 更新推送记录
func (r *PushRepo) Update(push *models.Push) error {
	return r.db.Save(push).Error
}

func (r *PushRepo) UpdateCodeview(repoID uint, commitID string, status string, result *string) error {
	updates := map[string]interface{}{
		"codeview_status": status,
	}
	if result != nil {
		updates["codeview_result"] = *result
	}
	return r.db.Model(&models.Push{}).
		Where("repo_id = ? AND commit_id = ?", repoID, commitID).
		Updates(updates).Error
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
		updates["pushed_at"] = time.Now()
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
	now := time.Now()

	// 统计函数辅助
	getCount := func(start time.Time) (total, success, failed int64) {
		r.db.Model(&models.Push{}).Where("created_at >= ?", start).Count(&total)
		r.db.Model(&models.Push{}).Where("created_at >= ? AND status = ?", start, models.PushStatusSuccess).Count(&success)
		r.db.Model(&models.Push{}).Where("created_at >= ? AND status = ?", start, models.PushStatusFailed).Count(&failed)
		return
	}

	// 今日统计
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tTotal, tSuccess, tFailed := getCount(todayStart)
	stats["today"] = map[string]interface{}{
		"total":   tTotal,
		"success": tSuccess,
		"failed":  tFailed,
	}

	// 本周统计
	offset := int(now.Weekday()) - 1
	if offset < 0 {
		offset = 6
	}
	weekStart := todayStart.AddDate(0, 0, -offset)
	wTotal, wSuccess, wFailed := getCount(weekStart)
	stats["this_week"] = map[string]interface{}{
		"total":   wTotal,
		"success": wSuccess,
		"failed":  wFailed,
	}

	// 本月统计
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	mTotal, mSuccess, mFailed := getCount(monthStart)
	stats["this_month"] = map[string]interface{}{
		"total":   mTotal,
		"success": mSuccess,
		"failed":  mFailed,
	}

	// 趋势统计 (近 7 天)
	var trend []map[string]interface{}
	for i := 6; i >= 0; i-- {
		day := todayStart.AddDate(0, 0, -i)
		nextDay := day.AddDate(0, 0, 1)
		var dTotal, dSuccess, dFailed int64
		r.db.Model(&models.Push{}).Where("created_at >= ? AND created_at < ?", day, nextDay).Count(&dTotal)
		r.db.Model(&models.Push{}).Where("created_at >= ? AND created_at < ? AND status = ?", day, nextDay, models.PushStatusSuccess).Count(&dSuccess)
		r.db.Model(&models.Push{}).Where("created_at >= ? AND created_at < ? AND status = ?", day, nextDay, models.PushStatusFailed).Count(&dFailed)
		trend = append(trend, map[string]interface{}{
			"date":    day.Format("2006-01-02"),
			"total":   dTotal,
			"success": dSuccess,
			"failed":  dFailed,
		})
	}
	stats["trend"] = trend

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

// Delete 删除推送记录（物理删除，避免唯一索引冲突）
func (r *PushRepo) Delete(id uint) error {
	return r.db.Unscoped().Delete(&models.Push{}, id).Error
}

// Transaction 执行事务
func (r *PushRepo) Transaction(fn func(txRepo *PushRepo) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &PushRepo{db: tx}
		return fn(txRepo)
	})
}

// BatchDelete 批量删除（物理删除）
func (r *PushRepo) BatchDelete(ids []uint) error {
	return r.db.Unscoped().Where("id IN ?", ids).Delete(&models.Push{}).Error
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
