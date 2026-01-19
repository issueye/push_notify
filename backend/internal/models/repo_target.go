package models

import (
	"time"
)

// RepoTarget 仓库-推送目标关联表
type RepoTarget struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	RepoID    uint      `gorm:"not null;index" json:"repo_id"`
	TargetID  uint      `gorm:"not null;index" json:"target_id"`
	CreatedAt time.Time `json:"created_at"`

	Repo   Repo   `gorm:"foreignKey:RepoID" json:"repo,omitempty"`
	Target Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}
