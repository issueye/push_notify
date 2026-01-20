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
