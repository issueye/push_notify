package main

import (
	"backend/database"
	"backend/internal/config"
	"backend/router"
	"backend/utils/logger"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(cfg.Logging.Level, cfg.Logging.Format); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting application", map[string]interface{}{
		"env":  cfg.App.Env,
		"port": cfg.App.Port,
	})

	// 初始化数据库
	db, err := database.Init(database.DatabaseConfig{
		Driver:  cfg.Database.Driver,
		Path:    cfg.Database.Path,
		LogMode: cfg.Database.LogMode,
	})
	if err != nil {
		logger.Fatal("Failed to init database", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer database.Close(db)

	// 自动迁移
	if err := database.Migrate(db); err != nil {
		logger.Fatal("Failed to migrate database", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 初始化数据
	if err := database.InitData(db); err != nil {
		logger.Fatal("Failed to init data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 设置Gin模式
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := router.Setup(cfg)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	logger.Info("Server starting", map[string]interface{}{
		"addr": addr,
	})

	if err := r.Run(addr); err != nil {
		logger.Fatal("Server failed to start", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
