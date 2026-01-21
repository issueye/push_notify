package database

import (
	"fmt"
	"strings"

	"backend/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DatabaseConfig struct {
	Driver  string `mapstructure:"driver"`
	Path    string `mapstructure:"path"`
	LogMode string `mapstructure:"log_mode"`
}

func Init(cfg DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.Path

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
		},
	}

	if cfg.LogMode == "debug" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(sqlite.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)

	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Repo{},
		&models.RepoTemplate{},
		&models.Target{},
		&models.RepoTarget{},
		&models.Push{},
		&models.Template{},
		&models.Prompt{},
		&models.PromptHistory{},
		&models.AIModel{},
		&models.Log{},
	)
	if err != nil {
		return err
	}

	// 数据迁移：填充 WebhookID
	var repos []models.Repo
	db.Where("webhook_id = '' OR webhook_id IS NULL").Find(&repos)
	for _, repo := range repos {
		// 从 /webhook/{id} 中提取 id
		parts := strings.Split(repo.WebhookURL, "/")
		if len(parts) >= 3 {
			webhookID := parts[len(parts)-1]
			db.Model(&models.Repo{}).Where("id = ?", repo.ID).Update("webhook_id", webhookID)
		}
	}

	return nil
}

// InitData 初始化基础数据
func InitData(db *gorm.DB) error {
	// 1. 初始化 admin 用户
	var userCount int64
	db.Model(&models.User{}).Where("username = ?", "admin").Count(&userCount)
	if userCount == 0 {
		admin := &models.User{
			Username: "admin",
			Password: "admin123", // 初始密码
			Email:    "admin@example.com",
			Role:     models.RoleAdmin,
			Status:   models.StatusActive,
		}
		if err := db.Create(admin).Error; err != nil {
			return fmt.Errorf("failed to init admin user: %w", err)
		}
		fmt.Println("Initial admin user created: admin / admin123")
	}

	// 2. 初始化默认模板（如果没有任何模板）
	var templateCount int64
	db.Model(&models.Template{}).Count(&templateCount)
	if templateCount == 0 {
		defaultTemplates := []models.Template{
			{
				Name:      "默认提交通知",
				Type:      "dingtalk",
				Scene:     "commit_notify",
				Title:     "代码提交通知",
				Content:   "### 代码提交通知\n- **仓库**: {{.RepoName}}\n- **提交人**: {{.Author}}\n- **提交信息**: {{.CommitMsg}}\n- **提交ID**: {{.CommitID}}\n- **文件变更**: {{.FileCount}} 个文件\n\n[查看详情]({{.RepoURL}})",
				IsDefault: true,
				Status:    models.StatusActive,
				Version:   1,
			},
			{
				Name:      "默认审查通知",
				Type:      "dingtalk",
				Scene:     "review_notify",
				Title:     "代码审查结果",
				Content:   "### 🔍 代码审查结果\n\n**仓库名称：** {{.RepoName}}\n**提交ID：** `{{.CommitID}}` \n**提交信息：** {{.CommitMsg}}\n\n---\n[查看审查详情]({{.ReviewURL}})",
				IsDefault: true,
				Status:    models.StatusActive,
				Version:   1,
			},
		}
		if err := db.Create(&defaultTemplates).Error; err != nil {
			return fmt.Errorf("failed to init default templates: %w", err)
		}
		fmt.Println("Initial default templates created")
	}

	return nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB 获取数据库实例（用于外部访问）
func GetDB(db *gorm.DB) *gorm.DB {
	return db
}
