package models

import (
	"time"

	"gorm.io/gorm"
)

type Prompt struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Type      string         `gorm:"size:20;not null" json:"type"` // codeview, message
	Scene     string         `gorm:"size:50" json:"scene,omitempty"`
	Language  string         `gorm:"size:50" json:"language,omitempty"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	ModelID   *uint          `json:"model_id,omitempty"`
	Model     *AIModel       `gorm:"foreignKey:ModelID" json:"model,omitempty"`
	Version   int            `gorm:"default:1" json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type PromptHistory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	PromptID  uint      `gorm:"index" json:"prompt_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

// 提示词类型
const (
	PromptTypeCodeView = "codeview"
	PromptTypeMessage  = "message"
)
