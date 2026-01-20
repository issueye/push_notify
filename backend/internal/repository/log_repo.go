package repository

import (
	"time"

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
func (r *LogRepo) GetOperationLogs(page, size int, userID uint, module, keyword, startTime, endTime string) ([]models.Log, int64) {
	var logs []models.Log
	var total int64

	query := r.db.Model(&models.Log{}).Preload("User").Where("type = ?", models.LogTypeOperation)
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if keyword != "" {
		query = query.Where("(message LIKE ? OR detail LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
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
func (r *LogRepo) GetAICallLogs(page, size int, keyword, startTime, endTime string) ([]models.Log, int64) {
	var logs []models.Log
	var total int64

	query := r.db.Model(&models.Log{}).Where("type = ?", models.LogTypeAICall)
	if keyword != "" {
		query = query.Where("(message LIKE ? OR detail LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
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

// GetStats 获取日志统计
func (r *LogRepo) GetStats(startDate, endDate string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 按级别统计
	var counts []struct {
		Level string
		Count int64
	}
	
	query := r.db.Model(&models.Log{})
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}
	
	query.Select("level, count(*) as count").Group("level").Scan(&counts)

	levelDist := map[string]int64{
		"debug": 0,
		"info":  0,
		"warn":  0,
		"error": 0,
	}
	for _, c := range counts {
		levelDist[c.Level] = c.Count
	}

	stats["level_distribution"] = levelDist

	// 按类型统计
	var typeCounts []struct {
		Type  string
		Count int64
	}
	r.db.Model(&models.Log{}).Select("type, count(*) as count").Group("type").Scan(&typeCounts)
	typeDist := map[string]int64{
		"system":    0,
		"operation": 0,
		"ai_call":   0,
	}
	for _, c := range typeCounts {
		typeDist[c.Type] = c.Count
	}
	stats["type_distribution"] = typeDist

	return stats, nil
}

// Delete 清理旧日志
func (r *LogRepo) Delete(beforeDays int) error {
	beforeTime := time.Now().AddDate(0, 0, -beforeDays)
	return r.db.Where("created_at < ?", beforeTime).Delete(&models.Log{}).Error
}
