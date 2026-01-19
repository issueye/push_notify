package services

import (
	"encoding/json"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/utils/logger"

	"gorm.io/gorm"
)

type LogService struct {
	db *gorm.DB
}

func NewLogService(db *gorm.DB) *LogService {
	return &LogService{db: db}
}

// CreateLog 创建日志
func (s *LogService) CreateLog(logType, level, module, message, detail string, userID uint, requestID string) error {
	log := &models.Log{
		Type:      logType,
		Level:     level,
		Module:    module,
		Message:   message,
		Detail:    detail,
		UserID:    &userID,
		RequestID: requestID,
	}

	db := repository.NewLogRepo(s.db)
	return db.Create(log)
}

// LogSystem 系统日志
func (s *LogService) LogSystem(level, module, message string, detail interface{}) {
	detailStr := ""
	if detail != nil {
		if d, err := json.Marshal(detail); err == nil {
			detailStr = string(d)
		}
	}
	s.CreateLog(models.LogTypeSystem, level, module, message, detailStr, 0, "")
	logger.Info(message, map[string]interface{}{
		"level":  level,
		"module": module,
	})
}

// LogOperation 操作日志
func (s *LogService) LogOperation(userID uint, module, action, objectType string, objectID uint, detail interface{}) {
	detailStr := ""
	if detail != nil {
		if d, err := json.Marshal(detail); err == nil {
			detailStr = string(d)
		}
	}
	s.CreateLog(models.LogTypeOperation, models.LogLevelInfo, module,
		action,
		detailStr, userID, "")
}

// LogAICall AI调用日志
func (s *LogService) LogAICall(modelID uint, input, output string, duration int, success bool) {
	detail := map[string]interface{}{
		"model_id": modelID,
		"duration": duration,
		"success":  success,
	}
	detailBytes, _ := json.Marshal(detail)
	detailStr := string(detailBytes)

	status := models.LogLevelInfo
	if !success {
		status = models.LogLevelError
	}

	s.CreateLog(models.LogTypeAICall, status, "ai",
		"AI模型调用",
		detailStr, 0, "")
}

// GetSystemLogs 获取系统日志
func (s *LogService) GetSystemLogs(page, size int, level, keyword, startTime, endTime string) ([]models.Log, int64, error) {
	db := repository.NewLogRepo(s.db)
	logs, total := db.GetSystemLogs(page, size, level, keyword, startTime, endTime)
	return logs, total, nil
}

// GetOperationLogs 获取操作日志
func (s *LogService) GetOperationLogs(page, size int, userID uint, action, startTime, endTime string) ([]models.Log, int64, error) {
	db := repository.NewLogRepo(s.db)
	logs, total := db.GetOperationLogs(page, size, userID, action, startTime, endTime)
	return logs, total, nil
}

// GetAICallLogs 获取AI调用日志
func (s *LogService) GetAICallLogs(page, size int, startTime, endTime string) ([]models.Log, int64, error) {
	db := repository.NewLogRepo(s.db)
	logs, total := db.GetAICallLogs(page, size, startTime, endTime)
	return logs, total, nil
}

// GetStats 获取日志统计
func (s *LogService) GetStats(startDate, endDate string) (map[string]interface{}, error) {
	db := repository.NewLogRepo(s.db)
	return db.GetStats(startDate, endDate)
}
