package handlers

import (
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	webhookService *services.WebhookService
}

func NewWebhookHandler(webhookService *services.WebhookService) *WebhookHandler {
	return &WebhookHandler{webhookService: webhookService}
}

// HandleGitHub GitHub Webhook回调
func (h *WebhookHandler) HandleGitHub(c *gin.Context) {
	h.webhookService.HandleGitHubWebhook(c)
}

// HandleGitLab GitLab Webhook回调
func (h *WebhookHandler) HandleGitLab(c *gin.Context) {
	h.webhookService.HandleGitLabWebhook(c)
}

// HandleGitee Gitee Webhook回调
func (h *WebhookHandler) HandleGitee(c *gin.Context) {
	// TODO: 实现Gitee Webhook处理
	utils.Success(c, map[string]interface{}{
		"status": "received",
		"source": "gitee",
	})
}
