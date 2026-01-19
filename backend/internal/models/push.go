package models

import (
	"time"

	"gorm.io/gorm"
)

type Push struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	RepoID         uint       `gorm:"not null;index" json:"repo_id"`
	Repo           Repo       `gorm:"foreignKey:RepoID" json:"repo,omitempty"`
	TargetID       uint       `gorm:"not null;index" json:"target_id"`
	Target         Target     `gorm:"foreignKey:TargetID" json:"target,omitempty"`
	TemplateID     *uint      `json:"template_id"`
	Template       Template   `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	CommitID       string     `gorm:"size:50;not null" json:"commit_id"`
	CommitMsg      string     `gorm:"size:500;not null" json:"commit_msg"`
	Status         string     `gorm:"size:20;default:'pending'" json:"status"`
	Content        string     `gorm:"type:text;not null" json:"content"`
	ErrorMsg       string     `gorm:"type:text" json:"error_msg,omitempty"`
	CodeviewResult *string    `gorm:"type:text" json:"codeview_result,omitempty"`
	CodeviewStatus string     `gorm:"size:20;default:'pending'" json:"codeview_status"`
	RetryCount     int        `gorm:"default:0" json:"retry_count"`
	PushedAt       *time.Time `json:"pushed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 推送状态
const (
	PushStatusPending = "pending"
	PushStatusSuccess = "success"
	PushStatusFailed  = "failed"
)

// Codeview 状态
const (
	CodeviewStatusPending = "pending"
	CodeviewStatusSuccess = "success"
	CodeviewStatusFailed  = "failed"
	CodeviewStatusSkipped = "skipped"
)
