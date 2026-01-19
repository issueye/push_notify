package models

import (
	"time"

	"gorm.io/gorm"
)

type Template struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"uniqueIndex:idx_name_type_scene;size:100;not null" json:"name"`
	Type      string         `gorm:"uniqueIndex:idx_name_type_scene;size:20;not null" json:"type"`  // dingtalk, email
	Scene     string         `gorm:"uniqueIndex:idx_name_type_scene;size:50;not null" json:"scene"` // commit_notify, review_notify
	Title     string         `gorm:"size:200;not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	Status    string         `gorm:"size:20;default:'active'" json:"status"`
	Version   int            `gorm:"default:1" json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Pushes []Push `json:"-"`
}

// 模板类型
const (
	TemplateTypeDingTalk = "dingtalk"
)

// 模板场景
const (
	TemplateSceneCommitNotify = "commit_notify"
	TemplateSceneReviewNotify = "review_notify"
)
