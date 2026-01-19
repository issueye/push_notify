package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Config struct {
	AccessToken string            `json:"access_token"`
	Headers     map[string]string `json:"headers"`
	Method      string            `json:"method"`
	Secret      string            `json:"secret"`
	WebhookURL  string            `json:"webhook_url"`
}

// 实现Sql序列化和反序列话接口
func (c *Config) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch t := value.(type) {
	case string:
		if t == "" {
			return nil
		}
		return json.Unmarshal([]byte(t), c)
	case []byte:
		if len(t) == 0 {
			return nil
		}
		return json.Unmarshal(t, c)
	default:
		return fmt.Errorf("unsupported type %T", t)
	}
}

func (c *Config) Value() (driver.Value, error) {
	return json.Marshal(c)
}

type Target struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Type      string         `gorm:"size:20;not null" json:"type"`          // dingtalk, email
	Config    *Config        `gorm:"type:text;not null" json:"config"`      // JSON配置
	Scope     string         `gorm:"size:20;default:'global'" json:"scope"` // global, repo
	Status    string         `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Repos       []Repo       `gorm:"many2many:repo_targets;" json:"repos,omitempty"`
	Pushes      []Push       `json:"-"`
	RepoTargets []RepoTarget `json:"-"`
}

// 推送目标类型
const (
	TargetTypeDingTalk = "dingtalk"
	TargetTypeWebhook  = "webhook"
)

// 范围
const (
	TargetScopeGlobal = "global"
	TargetScopeRepo   = "repo"
)
