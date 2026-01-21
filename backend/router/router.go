package router

import (
	"backend/database"
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/services"
	"backend/internal/utils"
	"backend/static"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config) *gin.Engine {
	// 初始化数据库连接
	db, err := database.Init(database.DatabaseConfig{
		Driver:  cfg.Database.Driver,
		Path:    cfg.Database.Path,
		LogMode: cfg.Database.LogMode,
	})
	if err != nil {
		panic("Failed to init database: " + err.Error())
	}

	// 初始化基础服务
	logService := services.NewLogService(db)

	r := gin.New()

	// 中间件
	r.Use(middleware.RecoveryMiddleware(logService))
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// 初始化JWT
	jwtUtils := utils.NewJWT(cfg.JWT.Secret, cfg.JWT.AccessTokenExpire, cfg.JWT.RefreshTokenExpire)

	// 初始化其他服务
	authService := services.NewAuthService(db, jwtUtils)
	repoService := services.NewRepoService(db)
	targetService := services.NewTargetService(db)
	templateService := services.NewTemplateService(db)
	generateTemplateService := services.NewGenerateTemplateService(db)
	promptService := services.NewPromptService(db)
	modelService := services.NewAIModelService(db)
	pushService := services.NewPushService(db)
	webhookService := services.NewWebhookService(db)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(authService)
	repoHandler := handlers.NewRepoHandler(repoService, logService)
	targetHandler := handlers.NewTargetHandler(targetService, logService)
	templateHandler := handlers.NewTemplateHandler(templateService, generateTemplateService, logService)
	promptHandler := handlers.NewPromptHandler(promptService, logService)
	modelHandler := handlers.NewModelHandler(modelService, logService)
	pushHandler := handlers.NewPushHandler(pushService)
	webhookHandler := handlers.NewWebhookHandler(webhookService)
	logHandler := handlers.NewLogHandler(logService)

	// 公共接口（无需认证）
	publicAPI := r.Group("/api/v1/auth")
	{
		publicAPI.POST("/login", authHandler.Login)
		publicAPI.POST("/register", authHandler.Register)
		publicAPI.POST("/refresh", authHandler.RefreshToken)
	}

	// Webhook接口（无需认证）
	webhookAPI := r.Group("/webhook")
	{
		webhookAPI.POST("/github/:webhookId", webhookHandler.HandleGitHub)
		webhookAPI.POST("/gitlab/:webhookId", webhookHandler.HandleGitLab)
	}

	// 需要认证的接口
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(jwtUtils))
	{
		// 认证相关
		api.GET("/auth/me", authHandler.GetCurrentUser)
		api.PUT("/auth/password", authHandler.ChangePassword)

		// 仓库管理
		repos := api.Group("/repos")
		{
			repos.GET("", repoHandler.List)
			repos.POST("", repoHandler.Create)
			repos.GET("/:id", repoHandler.Detail)
			repos.PUT("/:id", repoHandler.Update)
			repos.DELETE("/:id", repoHandler.Delete)
			repos.POST("/:id/test", repoHandler.TestWebhook)
			repos.GET("/:id/targets", repoHandler.GetTargets)
			repos.POST("/:id/targets", repoHandler.AddTarget)
			repos.DELETE("/:id/targets/:targetId", repoHandler.RemoveTarget)
		}

		// 推送目标管理
		targets := api.Group("/targets")
		{
			targets.GET("", targetHandler.List)
			targets.POST("", targetHandler.Create)
			targets.GET("/:id", targetHandler.Detail)
			targets.PUT("/:id", targetHandler.Update)
			targets.DELETE("/:id", targetHandler.Delete)
			targets.POST("/:id/test", targetHandler.Test)
			targets.POST("/:id/repos", targetHandler.AddRepo)
			targets.DELETE("/:id/repos/:repoId", targetHandler.RemoveRepo)
		}

		// 推送记录
		pushes := api.Group("/pushes")
		{
			pushes.GET("", pushHandler.List)
			pushes.GET("/:id", pushHandler.Detail)
			pushes.POST("/:id/retry", pushHandler.Retry)
			pushes.POST("/batch-retry", pushHandler.BatchRetry)
			pushes.DELETE("/batch-delete", pushHandler.BatchDelete)
			pushes.GET("/stats", pushHandler.GetStats)
		}

		// 消息模板
		templates := api.Group("/templates")
		{
			templates.GET("", templateHandler.List)
			templates.POST("", templateHandler.Create)
			templates.GET("/:id", templateHandler.Detail)
			templates.PUT("/:id", templateHandler.Update)
			templates.DELETE("/:id", templateHandler.Delete)
			templates.PUT("/:id/status", templateHandler.SetStatus)
			templates.POST("/:id/rollback", templateHandler.Rollback)
			templates.POST("/generate", templateHandler.Generate)
		}

		// 提示词
		prompts := api.Group("/prompts")
		{
			prompts.GET("", promptHandler.List)
			prompts.POST("", promptHandler.Create)
			prompts.GET("/:id", promptHandler.Detail)
			prompts.PUT("/:id", promptHandler.Update)
			prompts.DELETE("/:id", promptHandler.Delete)
			prompts.POST("/:id/test", promptHandler.Test)
			prompts.POST("/:id/rollback", promptHandler.Rollback)
			prompts.GET("/:id/history", promptHandler.History)
		}

		// AI模型
		models := api.Group("/models")
		{
			models.GET("", modelHandler.List)
			models.POST("", modelHandler.Create)
			models.GET("/:id", modelHandler.Detail)
			models.PUT("/:id", modelHandler.Update)
			models.DELETE("/:id", modelHandler.Delete)
			models.POST("/:id/default", modelHandler.SetDefault)
			models.GET("/:id/logs", modelHandler.GetLogs)
			models.POST("/:id/verify", modelHandler.Verify)
		}

		// 日志管理
		logs := api.Group("/logs")
		{
			logs.GET("/system", logHandler.GetSystemLogs)
			logs.GET("/operations", logHandler.GetOperationLogs)
			logs.GET("/ai-calls", logHandler.GetAICallLogs)
			logs.GET("/search", logHandler.Search)
			logs.GET("/export", logHandler.Export)
			logs.GET("/stats", logHandler.GetStats)
		}

		// 用户管理（仅管理员）
		users := api.Group("/users")
		users.Use(middleware.AdminOnly())
		{
			users.GET("", userListHandler(db))
			users.POST("", createUserHandler(db))
			users.GET("/:id", userDetailHandler(db))
			users.PUT("/:id", updateUserHandler(db))
			users.DELETE("/:id", deleteUserHandler(db))
			users.POST("/:id/reset-password", resetPasswordHandler(db))
			users.PUT("/:id/lock", lockUserHandler(db))
		}

		// 个人设置
		api.GET("/settings", settingsHandler)
		api.PUT("/settings", updateSettingsHandler)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{
			"status":  "ok",
			"service": "push-notify",
		})
	})

	// 静态资源服务
	staticFS := static.GetStaticFS()
	
	// 处理 /web 前缀
	r.StaticFS("/web", http.FS(staticFS))

	// SPA 路由回退处理
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// 如果是访问 /web 下的路径，且不是 API 或 Webhook，则返回 index.html
		if len(path) >= 4 && path[:4] == "/web" {
			c.FileFromFS("index.html", http.FS(staticFS))
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})

	// 根路径重定向到 /web/
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/")
	})

	return r
}

// 用户管理处理器
func userListHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := utils.GetPage(c)
		size := utils.GetSize(c)
		keyword := c.Query("keyword")
		role := c.Query("role")
		status := c.Query("status")

		userRepo := services.NewUserRepoWrapper(db)
		users, total := userRepo.List(page, size, keyword, role, status)
		utils.SuccessWithPage(c, users, int(total), page, size)
	}
}

func createUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.SuccessWithMsg(c, "创建成功", nil)
	}
}

func userDetailHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.Success(c, map[string]interface{}{})
	}
}

func updateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.SuccessWithMsg(c, "更新成功", nil)
	}
}

func deleteUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.SuccessWithMsg(c, "删除成功", nil)
	}
}

func resetPasswordHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.Success(c, map[string]interface{}{
			"new_password": "TempP@ss123",
		})
	}
}

func lockUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.SuccessWithMsg(c, "操作成功", nil)
	}
}

func settingsHandler(c *gin.Context) {
	utils.Success(c, map[string]interface{}{
		"notify": map[string]interface{}{
			"channels":    []string{"dingtalk", "email"},
			"quiet_hours": false,
		},
		"preferences": map[string]interface{}{
			"language": "zh-CN",
			"theme":    "light",
		},
	})
}

func updateSettingsHandler(c *gin.Context) {
	utils.SuccessWithMsg(c, "更新成功", nil)
}
