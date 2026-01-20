package services

import (
	"errors"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"backend/utils/logger"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrUserAlreadyExists = errors.New("用户已存在")
	ErrInvalidPassword   = errors.New("密码错误")
	ErrUserLocked        = errors.New("账户已锁定")
)

type AuthService struct {
	userRepo *repository.UserRepo
	jwtUtils *utils.JWT
	logServ  *LogService
}

func NewAuthService(db *gorm.DB, jwtUtils *utils.JWT) *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepo(db),
		jwtUtils: jwtUtils,
		logServ:  NewLogService(db),
	}
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 检查账户状态
	if user.Status == models.StatusLocked {
		return nil, ErrUserLocked
	}

	// 验证密码
	if !user.VerifyPassword(password) {
		return nil, ErrInvalidPassword
	}

	// 生成Token
	accessToken, refreshToken, err := s.jwtUtils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		s.logServ.LogOperation(user.ID, "auth", "登录失败", "auth", 0, map[string]string{"username": username, "error": err.Error()})
		return nil, err
	}

	// 更新最后登录时间
	s.userRepo.UpdateLastLogin(user.ID)

	result := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    s.jwtUtils.AccessTokenExpire,
		"token_type":    "Bearer",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	}

	s.logServ.LogOperation(user.ID, "auth", "登录成功", "auth", user.ID, map[string]string{"username": username})
	logger.Info("User logged in", map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
	})

	return result, nil
}

// Register 用户注册
func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	// 检查用户名是否存在
	_, err := s.userRepo.GetByUsername(username)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 检查邮箱是否存在
	_, err = s.userRepo.GetByEmail(email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
		Role:     models.RoleUser,
		Status:   models.StatusActive,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	s.logServ.LogOperation(user.ID, "auth", "用户注册", "user", user.ID, map[string]string{"username": username, "email": email})
	logger.Info("User registered", map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
	})

	return user, nil
}

// GetUserInfo 获取当前用户信息
func (s *AuthService) GetUserInfo(userID uint) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if !user.VerifyPassword(oldPassword) {
		return ErrInvalidPassword
	}

	user.Password = newPassword
	return s.userRepo.Update(user)
}

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	return s.jwtUtils.RefreshToken(refreshToken)
}
