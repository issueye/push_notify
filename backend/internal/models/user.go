package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Username    string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email       string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Role        string         `gorm:"size:20;default:'user'" json:"role"`
	Status      string         `gorm:"size:20;default:'active'" json:"status"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// 角色常量
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// 状态常量
const (
	StatusActive   = "active"
	StatusLocked   = "locked"
	StatusInactive = "inactive"
)

// BeforeCreate 创建前加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// BeforeUpdate 更新前加密密码
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 检查密码是否被修改
	if tx.Statement.Changed("Password") && u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
