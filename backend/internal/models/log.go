package models

import (
	"time"

	"gorm.io/gorm"
)

type Log struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Type      string         `gorm:"size:20;not null;index" json:"type"`  // system, operation, ai_call
	Level     string         `gorm:"size:20;not null;index" json:"level"` // debug, info, warn, error
	Module    string         `gorm:"size:50" json:"module,omitempty"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Detail    string         `gorm:"type:text" json:"detail,omitempty"`
	UserID    *uint          `json:"user_id,omitempty"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RequestID string         `gorm:"size:100" json:"request_id,omitempty"`
	CreatedAt time.Time      `gorm:"index" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 日志类型
const (
	LogTypeSystem    = "system"
	LogTypeOperation = "operation"
	LogTypeAICall    = "ai_call"
)

// 日志级别
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)
