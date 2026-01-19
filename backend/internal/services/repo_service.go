package services

import (
	"errors"
	"strings"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"backend/utils/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRepoNotFound      = errors.New("仓库不存在")
	ErrRepoAlreadyExists = errors.New("仓库名称已存在")
	ErrInvalidRepoURL    = errors.New("无效的仓库地址")
)

type RepoService struct {
	repoRepo   *repository.RepoRepo
	targetRepo *repository.TargetRepo
}

func NewRepoService(db *gorm.DB) *RepoService {
	return &RepoService{
		repoRepo:   repository.NewRepoRepo(db),
		targetRepo: repository.NewTargetRepo(db),
	}
}

// Create 创建仓库
func (s *RepoService) Create(data map[string]interface{}) (*models.Repo, error) {
	name := data["name"].(string)

	// 检查名称是否重复
	_, err := s.repoRepo.GetByName(name)
	if err == nil {
		return nil, ErrRepoAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 验证URL格式
	url := data["url"].(string)
	if !isValidGitURL(url) {
		return nil, ErrInvalidRepoURL
	}

	// 生成Webhook URL
	webhookID := uuid.New().String()
	webhookURL := "/webhook/" + webhookID

	// 生成Webhook Secret
	webhookSecret, _ := utils.GenerateWebhookSecret()

	repo := &models.Repo{
		Name:          name,
		URL:           url,
		Type:          data["type"].(string),
		AccessToken:   getString(data, "access_token"),
		WebhookURL:    webhookURL,
		WebhookSecret: webhookSecret,
		Status:        models.RepoStatusActive,
	}

	if modelID, ok := data["model_id"].(uint); ok && modelID > 0 {
		repo.ModelID = &modelID
	}

	if err := s.repoRepo.Create(repo); err != nil {
		return nil, err
	}

	// 绑定推送目标
	if targetIDs, ok := data["target_ids"].([]interface{}); ok {
		for _, id := range targetIDs {
			if targetID, ok := id.(float64); ok {
				s.repoRepo.AddTarget(repo.ID, uint(targetID))
			}
		}
	}

	logger.Info("Repo created", map[string]interface{}{
		"repo_id": repo.ID,
		"name":    repo.Name,
		"targets": data["target_ids"],
	})

	return repo, nil
}

// GetByID 获取仓库详情
func (s *RepoService) GetByID(id uint) (*models.Repo, error) {
	return s.repoRepo.GetByID(id)
}

// GetList 获取仓库列表
func (s *RepoService) GetList(page, size int, keyword string) ([]models.Repo, int64, error) {
	repos, total := s.repoRepo.GetList(page, size, keyword)
	return repos, total, nil
}

// Update 更新仓库
// TODO 添加事务
func (s *RepoService) Update(id uint, data models.UpdateRepo) error {
	repo, err := s.repoRepo.GetByID(id)
	if err != nil {
		return err
	}

	repo.Name = data.Name
	repo.URL = data.URL
	repo.Type = data.Type
	repo.Status = data.Status
	repo.ModelID = data.ModelID

	if len(data.TargetIds) > 0 {
		// 删除再写入
		if err := s.repoRepo.DeleteTargets(repo.ID); err != nil {
			return err
		}

		if err := s.repoRepo.InsertTargets(repo.ID, data.TargetIds); err != nil {
			return err
		}
	}

	return s.repoRepo.Update(repo)
}

// DeleteTargets 删除推送目标
func (s *RepoService) DeleteTargets(repoID uint) error {
	return s.repoRepo.DeleteTargets(repoID)
}

// Delete 删除仓库
func (s *RepoService) Delete(id uint) error {
	_, err := s.repoRepo.GetByID(id)
	if err != nil {
		return err
	}
	return s.repoRepo.Delete(id)
}

// TestWebhook 测试Webhook
func (s *RepoService) TestWebhook(id uint) (map[string]interface{}, error) {
	repo, err := s.repoRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 返回Webhook配置信息用于测试
	result := map[string]interface{}{
		"webhook_url":    repo.WebhookURL,
		"webhook_secret": repo.WebhookSecret,
		"status":         "configured",
	}

	logger.Info("Webhook tested", map[string]interface{}{
		"repo_id": repo.ID,
		"name":    repo.Name,
	})

	return result, nil
}

// AddTarget 关联推送目标
func (s *RepoService) AddTarget(repoID, targetID uint) error {
	_, err := s.repoRepo.GetByID(repoID)
	if err != nil {
		return err
	}
	return s.repoRepo.AddTarget(repoID, targetID)
}

// RemoveTarget 取消关联推送目标
func (s *RepoService) RemoveTarget(repoID, targetID uint) error {
	return s.repoRepo.RemoveTarget(repoID, targetID)
}

// GetTargets 获取仓库关联的推送目标
func (s *RepoService) GetTargets(repoID uint) ([]models.Target, error) {
	return s.repoRepo.GetTargets(repoID)
}

// GetByWebhookID 根据Webhook ID获取仓库
func (s *RepoService) GetByWebhookID(webhookID string) (*models.Repo, error) {
	return s.repoRepo.GetByWebhookID(webhookID)
}

// 辅助函数
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func isValidGitURL(url string) bool {
	// 简单的URL验证
	return strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://") ||
		strings.HasPrefix(url, "git@") ||
		strings.HasPrefix(url, "ssh://")
}
