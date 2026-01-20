package models

import (
	"time"
)

type RepoTemplate struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	RepoID     uint      `gorm:"not null;index;uniqueIndex:idx_repo_lang" json:"repo_id"`
	TemplateID uint      `gorm:"not null;index" json:"template_id"`
	Repo       Repo      `gorm:"foreignKey:RepoID" json:"repo,omitempty"`
	Template   Template  `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Language   string    `gorm:"size:50;default:'default';uniqueIndex:idx_repo_lang" json:"language"` // 适用语言，default表示默认
	CreatedAt  time.Time `json:"created_at"`
}
