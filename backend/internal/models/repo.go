package models

import (
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Name          string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
	URL           string         `gorm:"size:500;not null" json:"url"`
	Type          string         `gorm:"size:50;not null" json:"type"` // github, gitlab, gitee
	AccessToken   string         `gorm:"size:255" json:"access_token"`
	WebhookID     string         `gorm:"uniqueIndex;size:100" json:"webhook_id"`
	WebhookURL    string         `gorm:"size:500;not null" json:"webhook_url"`
	WebhookSecret string         `gorm:"size:100" json:"webhook_secret"`
	ModelID       *uint          `json:"model_id"`
	Model         *AIModel       `gorm:"foreignKey:ModelID" json:"model,omitempty"`
	
	// 模板关联
	CommitTemplateID *uint          `json:"commit_template_id"`
	CommitTemplate   *Template      `gorm:"foreignKey:CommitTemplateID" json:"commit_template,omitempty"`
	ReviewTemplates  []RepoTemplate `gorm:"foreignKey:RepoID" json:"review_templates,omitempty"`

	Status        string         `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Targets     []Target     `gorm:"many2many:repo_targets;" json:"targets,omitempty"`
	Pushes      []Push       `json:"-"`
	RepoTargets []RepoTarget `json:"repo_targets,omitempty"`
}

// 仓库类型
const (
	RepoTypeGitHub = "github"
	RepoTypeGitLab = "gitlab"
	RepoTypeGitee  = "gitee"
)

// 状态常量
const (
	RepoStatusActive   = "active"
	RepoStatusInactive = "inactive"
)

type UpdateRepo struct {
	ID               uint                `json:"id"`
	Name             string              `json:"name"`
	URL              string              `json:"url"`
	Type             string              `json:"type"` // github, gitlab, gitee
	Status           string              `json:"status"`
	ModelID          *uint               `json:"model_id"`
	TargetIds        []uint              `json:"target_ids"`
	CommitTemplateID *uint               `json:"commit_template_id"`
	ReviewTemplates  []RepoTemplateConfig `json:"review_templates"`
	AccessToken      string              `json:"access_token"`
}

type CreateRepo struct {
	Name             string              `json:"name"`
	URL              string              `json:"url"`
	Type             string              `json:"type"` // github, gitlab, gitee
	Status           string              `json:"status"`
	ModelID          *uint               `json:"model_id"`
	TargetIds        []uint              `json:"target_ids"`
	CommitTemplateID *uint               `json:"commit_template_id"`
	ReviewTemplates  []RepoTemplateConfig `json:"review_templates"`
	AccessToken      string              `json:"access_token"`
}

type RepoTemplateConfig struct {
	TemplateID uint   `json:"template_id"`
	Language   string `json:"language"`
}
