package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type LogRepo struct {
	db *gorm.DB
}

func NewLogRepo(db *gorm.DB) *LogRepo {
	return &LogRepo{db: db}
}

// Create 创建日志
func (r *LogRepo) Create(log *models.Log) error {
	return r.db.Create(log).Error
}

// GetByID 根据ID获取日志
func (r *LogRepo) GetByID(id uint) (*models.Log, error) {
	var log models.Log
	err := r.db.Preload("User").First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetSystemLogs 获取系统日志
func (r *LogRepo) GetSystemLogs(page, size int, level, keyword, startTime, endTime string) ([]models.Log, int64) {
	var logs []models.Log
	var total int64

	query := r.db.Model(&models.Log{}).Where("type = ?", models.LogTypeSystem)
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if keyword != "" {
		query = query.Where("message LIKE ?", "%"+keyword+"%")
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&logs)

	return logs, total
}

// GetOperationLogs 获取操作日志
func (r *LogRepo) GetOperationLogs(page, size int, userID uint, action, startTime, endTime string) ([]models.Log, int64) {
	var logs []models.Log
	var total int64

	query := r.db.Model(&models.Log{}).Where("type = ?", models.LogTypeOperation)
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if action != "" {
		query = query.Where("module = ?", action)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&logs)

	return logs, total
}

// GetAICallLogs 获取AI调用日志
func (r *LogRepo) GetAICallLogs(page, size int, startTime, endTime string) ([]models.Log, int64) {
	var logs []models.Log
	var total int64

	query := r.db.Model(&models.Log{}).Where("type = ?", models.LogTypeAICall)
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&logs)

	return logs, total
}

// GetStats 获取日志统计
func (r *LogRepo) GetStats(startDate, endDate string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 按级别统计
	var debugCount, infoCount, warnCount, errorCount int64
	r.db.Model(&models.Log{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("level").
		Pluck("count", &debugCount)
	r.db.Model(&models.Log{}).
		Where("created_at >= ? AND created_at <= ? AND level = ?", startDate, endDate, "info").
		Count(&infoCount)
	r.db.Model(&models.Log{}).
		Where("created_at >= ? AND created_at <= ? AND level = ?", startDate, endDate, "warn").
		Count(&warnCount)
	r.db.Model(&models.Log{}).
		Where("created_at >= ? AND created_at <= ? AND level = ?", startDate, endDate, "error").
		Count(&errorCount)

	stats["level_distribution"] = map[string]interface{}{
		"debug": debugCount,
		"info":  infoCount,
		"warn":  warnCount,
		"error": errorCount,
	}

	return stats, nil
}

// Delete 清理旧日志
func (r *LogRepo) Delete(beforeDays int) error {
	return r.db.Where("created_at < DATE_SUB(NOW(), INTERVAL ? DAY)", beforeDays).Delete(&models.Log{}).Error
}
