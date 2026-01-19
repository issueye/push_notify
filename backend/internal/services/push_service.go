package services

import (
	"errors"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/utils/logger"

	"gorm.io/gorm"
)

var (
	ErrPushNotFound = errors.New("推送记录不存在")
)

type PushService struct {
	pushRepo *repository.PushRepo
}

func NewPushService(db *gorm.DB) *PushService {
	return &PushService{
		pushRepo: repository.NewPushRepo(db),
	}
}

// Create 创建推送记录
func (s *PushService) Create(data map[string]interface{}) (*models.Push, error) {
	push := &models.Push{
		RepoID:     data["repo_id"].(uint),
		TargetID:   data["target_id"].(uint),
		CommitID:   data["commit_id"].(string),
		CommitMsg:  data["commit_msg"].(string),
		Status:     models.PushStatusPending,
		Content:    data["content"].(string),
		RetryCount: 0,
	}

	if templateID, ok := data["template_id"].(uint); ok && templateID > 0 {
		push.TemplateID = &templateID
	}

	if err := s.pushRepo.Create(push); err != nil {
		return nil, err
	}

	return push, nil
}

// GetByID 获取推送记录
func (s *PushService) GetByID(id uint) (*models.Push, error) {
	return s.pushRepo.GetByID(id)
}

// GetList 获取推送记录列表
func (s *PushService) GetList(page, size int, repoID, targetID uint, status, keyword, startTime, endTime string) ([]models.Push, int64, error) {
	pushes, total := s.pushRepo.GetList(page, size, repoID, targetID, status, keyword, startTime, endTime)
	return pushes, total, nil
}

// Retry 重试推送
func (s *PushService) Retry(id uint) (uint, error) {
	push, err := s.pushRepo.GetByID(id)
	if err != nil {
		return 0, err
	}

	// 创建新的推送记录
	newPush := &models.Push{
		RepoID:     push.RepoID,
		TargetID:   push.TargetID,
		TemplateID: push.TemplateID,
		CommitID:   push.CommitID,
		CommitMsg:  push.CommitMsg,
		Status:     models.PushStatusPending,
		Content:    push.Content,
		RetryCount: 0,
	}

	if err := s.pushRepo.Create(newPush); err != nil {
		return 0, err
	}

	// 更新原记录的重试次数
	s.pushRepo.IncrementRetryCount(id)

	logger.Info("Push retried", map[string]interface{}{
		"original_id": id,
		"new_id":      newPush.ID,
	})

	return newPush.ID, nil
}

// BatchRetry 批量重试
func (s *PushService) BatchRetry(ids []uint) (int, error) {
	count := 0
	for _, id := range ids {
		_, err := s.Retry(id)
		if err == nil {
			count++
		}
	}
	return count, nil
}

// BatchDelete 批量删除
func (s *PushService) BatchDelete(ids []uint, beforeDate string) error {
	if len(ids) > 0 {
		return s.pushRepo.BatchDelete(ids)
	}
	return nil
}

// GetStats 获取统计
func (s *PushService) GetStats(startDate, endDate string) (map[string]interface{}, error) {
	return s.pushRepo.GetStats(startDate, endDate)
}

// UpdateStatus 更新推送状态
func (s *PushService) UpdateStatus(id uint, status, errorMsg string) error {
	return s.pushRepo.UpdateStatus(id, status, errorMsg)
}
