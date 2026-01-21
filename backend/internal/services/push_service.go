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

	// 使用事务确保先删后加的原子性
	var newID uint
	err = s.pushRepo.Transaction(func(txRepo *repository.PushRepo) error {
		// 1. 准备新数据
		newPushData := map[string]interface{}{
			"repo_id":    push.RepoID,
			"target_id":  push.TargetID,
			"commit_id":  push.CommitID,
			"commit_msg": push.CommitMsg,
			"content":    push.Content,
		}
		if push.TemplateID != nil {
			newPushData["template_id"] = *push.TemplateID
		}

		// 2. 物理删除旧记录
		if err := txRepo.Delete(id); err != nil {
			return err
		}

		// 3. 创建新记录
		// 这里直接构造模型对象以避免依赖 Service 的 Create 方法（可能涉及不同的 Repo 实例）
		newPush := &models.Push{
			RepoID:     newPushData["repo_id"].(uint),
			TargetID:   newPushData["target_id"].(uint),
			CommitID:   newPushData["commit_id"].(string),
			CommitMsg:  newPushData["commit_msg"].(string),
			Content:    newPushData["content"].(string),
			Status:     models.PushStatusPending,
			RetryCount: 0,
		}
		if tid, ok := newPushData["template_id"].(uint); ok {
			newPush.TemplateID = &tid
		}

		if err := txRepo.Create(newPush); err != nil {
			return err
		}
		newID = newPush.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

	logger.Info("Push retried and old record physically removed", map[string]interface{}{
		"original_id": id,
		"new_id":      newID,
	})

	return newID, nil
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

// Delete 删除推送记录
func (s *PushService) Delete(id uint) error {
	return s.pushRepo.Delete(id)
}

// GetStats 获取统计
func (s *PushService) GetStats(startDate, endDate string) (map[string]interface{}, error) {
	return s.pushRepo.GetStats(startDate, endDate)
}

// UpdateStatus 更新推送状态
func (s *PushService) UpdateStatus(id uint, status, errorMsg string) error {
	return s.pushRepo.UpdateStatus(id, status, errorMsg)
}
