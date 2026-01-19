package main

import (
	"fmt"
	"log"

	"backend/database"
	"backend/internal/models"
	"backend/utils/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// 初始化日志
	if err := logger.Init("info", "json"); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	// 连接数据库
	db, err := database.Init(database.DatabaseConfig{
		Driver:  "sqlite",
		Path:    "data.db",
		LogMode: "false",
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer database.Close(db)

	// 创建 admin 用户
	if err := createAdminUser(db); err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Println("Admin user created successfully!")
}

func createAdminUser(db *gorm.DB) error {
	// 检查是否已存在
	var count int64
	db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)
	if count > 0 {
		fmt.Println("Admin user already exists")
		return nil
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: string(hashedPassword),
		Role:     models.RoleAdmin,
		Status:   models.StatusActive,
	}

	if err := db.Create(user).Error; err != nil {
		return err
	}

	fmt.Printf("Created admin user: %s (ID: %d)\n", user.Username, user.ID)
	return nil
}
