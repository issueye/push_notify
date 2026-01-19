package handlers

import (
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	result, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "登录成功", result)
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required,min=3,max=20"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=8"`
		ConfirmPass string `json:"confirm_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	if req.Password != req.ConfirmPass {
		utils.Fail(c, 400, "两次密码不一致")
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "注册成功", map[string]interface{}{
		"user_id": user.ID,
	})
}

// GetCurrentUser 获取当前用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := utils.GetUserID(c)
	user, err := h.authService.GetUserInfo(userID)
	if err != nil {
		utils.Fail(c, 400, "获取用户信息失败")
		return
	}

	utils.Success(c, map[string]interface{}{
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"role":          user.Role,
		"status":        user.Status,
		"last_login_at": user.LastLoginAt,
		"created_at":    user.CreatedAt,
	})
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	userID := utils.GetUserID(c)
	err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "密码修改成功", nil)
}

// RefreshToken 刷新Token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateError(c, []string{err.Error()})
		return
	}

	token, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.Fail(c, 401, "无效的刷新令牌")
		return
	}

	utils.Success(c, map[string]interface{}{
		"access_token": token,
		"token_type":   "Bearer",
	})
}
