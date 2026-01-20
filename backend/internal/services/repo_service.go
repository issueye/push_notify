package services

import (
	"errors"
	"fmt"
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
	db         *gorm.DB
	repoRepo   *repository.RepoRepo
	targetRepo *repository.TargetRepo
}

func NewRepoService(db *gorm.DB) *RepoService {
	return &RepoService{
		db:         db,
		repoRepo:   repository.NewRepoRepo(db),
		targetRepo: repository.NewTargetRepo(db),
	}
}

// Create 创建仓库
func (s *RepoService) Create(data models.CreateRepo) (*models.Repo, error) {
	name := data.Name

	// 检查名称是否重复
	_, err := s.repoRepo.GetByName(name)
	if err == nil {
		return nil, ErrRepoAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 验证URL格式
	url := data.URL
	if !isValidGitURL(url) {
		return nil, ErrInvalidRepoURL
	}

	// 生成Webhook URL
	webhookID := uuid.New().String()
	webhookURL := fmt.Sprintf("/webhook/%s/%s", data.Type, webhookID)

	// 生成Webhook Secret
	webhookSecret, _ := utils.GenerateWebhookSecret()

	repo := &models.Repo{
		Name:          name,
		URL:           url,
		Type:          data.Type,
		WebhookID:     webhookID,
		WebhookURL:    webhookURL,
		WebhookSecret: webhookSecret,
		Status:        models.RepoStatusActive,
		AccessToken:   data.AccessToken,
	}

	repo.ModelID = data.ModelID
	repo.CommitTemplateID = data.CommitTemplateID

	if err := s.repoRepo.Create(repo); err != nil {
		return nil, err
	}

	// 绑定推送目标
	if len(data.TargetIds) > 0 {
		if err := s.repoRepo.InsertTargets(repo.ID, data.TargetIds); err != nil {
			return nil, err
		}
	}

	// 绑定审查模板
	if len(data.ReviewTemplates) > 0 {
		// 去重
		uniqueTemplates := deduplicateReviewTemplates(data.ReviewTemplates)
		var configs []models.RepoTemplateConfig
		for _, rt := range uniqueTemplates {
			templateID := rt.TemplateID
			language := rt.Language
			if templateID > 0 {
				configs = append(configs, models.RepoTemplateConfig{
					TemplateID: uint(templateID),
					Language:   language,
				})
			}
		}
		if len(configs) > 0 {
			s.repoRepo.InsertReviewTemplates(repo.ID, configs)
		}
	}

	logger.Info("Repo created", map[string]interface{}{
		"repo_id": repo.ID,
		"name":    repo.Name,
		"targets": data.TargetIds,
		"review_templates": data.ReviewTemplates,
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
func (s *RepoService) Update(id uint, data models.UpdateRepo) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.repoRepo.WithTx(tx)

		repo, err := txRepo.GetByID(id)
		if err != nil {
			return err
		}

		repo.Name = data.Name
		repo.URL = data.URL
		repo.Type = data.Type
		repo.Status = data.Status
		repo.ModelID = data.ModelID
		repo.CommitTemplateID = data.CommitTemplateID

		if data.AccessToken != "" {
			repo.AccessToken = data.AccessToken
		}

		// 处理推送目标
		if err := txRepo.DeleteTargets(repo.ID); err != nil {
			return err
		}
		if len(data.TargetIds) > 0 {
			if err := txRepo.InsertTargets(repo.ID, data.TargetIds); err != nil {
				return err
			}
		}

		// 清除 Gorm 加载的关联，避免 Save 时重复插入
		repo.Targets = nil
		repo.RepoTargets = nil

		// 更新审查模板
		if err := txRepo.DeleteReviewTemplates(repo.ID); err != nil {
			return err
		}
		if len(data.ReviewTemplates) > 0 {
			// 去重
			uniqueTemplates := deduplicateReviewTemplates(data.ReviewTemplates)
			if err := txRepo.InsertReviewTemplates(repo.ID, uniqueTemplates); err != nil {
				return err
			}
		}

		return txRepo.Update(repo)
	})
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

// deduplicateReviewTemplates 去重审查模板配置
func deduplicateReviewTemplates(templates []models.RepoTemplateConfig) []models.RepoTemplateConfig {
	seen := make(map[string]struct{})
	var result []models.RepoTemplateConfig
	for _, t := range templates {
		// 统一转换为小写进行比较，避免大小写造成的重复
		langKey := strings.ToLower(t.Language)
		if _, ok := seen[langKey]; !ok {
			seen[langKey] = struct{}{}
			result = append(result, t)
		}
	}
	return result
}
