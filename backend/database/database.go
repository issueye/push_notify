package database

import (
	"fmt"

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
	return db.AutoMigrate(
		&models.User{},
		&models.Repo{},
		&models.RepoTemplate{},
		&models.Target{},
		&models.RepoTarget{},
		&models.Push{},
		&models.Template{},
		&models.Prompt{},
		&models.AIModel{},
		&models.Log{},
	)
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
