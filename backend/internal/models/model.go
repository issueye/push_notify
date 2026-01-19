package models

import (
	"time"

	"gorm.io/gorm"
)

type AIModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Type      string         `gorm:"size:50;not null" json:"type"` // gpt-4, claude, etc.
	Provider  string         `gorm:"size:50" json:"provider,omitempty"`
	APIURL    string         `gorm:"size:500;not null" json:"api_url"`
	APIKey    string         `gorm:"size:255;not null" json:"-"`
	Params    string         `gorm:"type:text" json:"params,omitempty"` // JSON参数
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	CallCount int            `gorm:"default:0" json:"call_count"`
	Status    string         `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
