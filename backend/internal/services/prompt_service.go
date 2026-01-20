package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/ai"
	"backend/utils/logger"

	"gorm.io/gorm"
)

var (
	ErrPromptNotFound      = errors.New("提示词不存在")
	ErrPromptAlreadyExists = errors.New("提示词名称已存在")
)

type PromptService struct {
	db         *gorm.DB
	promptRepo *repository.PromptRepo
	modelRepo  *repository.AIModelRepo
	logServ    *LogService
}

func NewPromptService(db *gorm.DB) *PromptService {
	return &PromptService{
		db:         db,
		promptRepo: repository.NewPromptRepo(db),
		modelRepo:  repository.NewAIModelRepo(db),
		logServ:    NewLogService(db),
	}
}

// Create 创建提示词
func (s *PromptService) Create(data map[string]interface{}) (*models.Prompt, error) {
	name := data["name"].(string)
	promptType := data["type"].(string)

	prompt := &models.Prompt{
		Name:     name,
		Type:     promptType,
		Scene:    getStringWithDefault(data, "scene", ""),
		Language: getStringWithDefault(data, "language", ""),
		Content:  data["content"].(string),
		Version:  1,
	}

	if modelID, ok := data["model_id"].(uint); ok && modelID > 0 {
		prompt.ModelID = &modelID
	}

	if err := s.promptRepo.Create(prompt); err != nil {
		return nil, err
	}

	logger.Info("Prompt created", map[string]interface{}{
		"prompt_id": prompt.ID,
		"name":      prompt.Name,
		"type":      prompt.Type,
	})

	return prompt, nil
}

// GetByID 获取提示词详情
func (s *PromptService) GetByID(id uint) (*models.Prompt, error) {
	return s.promptRepo.GetByID(id)
}

// GetList 获取提示词列表
func (s *PromptService) GetList(page, size int, keyword, promptType, scene string) ([]models.Prompt, int64, error) {
	prompts, total := s.promptRepo.GetList(page, size, keyword, promptType, scene)
	return prompts, total, nil
}

// Update 更新提示词
func (s *PromptService) Update(id uint, data map[string]interface{}) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.promptRepo.WithTx(tx)
		prompt, err := txRepo.GetByID(id)
		if err != nil {
			return err
		}

		// 保存历史记录
		history := &models.PromptHistory{
			PromptID:  prompt.ID,
			Content:   prompt.Content,
			Version:   prompt.Version,
			CreatedAt: prompt.UpdatedAt,
		}
		if err := txRepo.CreateHistory(history); err != nil {
			return err
		}

		if name, ok := data["name"].(string); ok && name != "" {
			prompt.Name = name
		}
		if scene, ok := data["scene"].(string); ok {
			prompt.Scene = scene
		}
		if language, ok := data["language"].(string); ok {
			prompt.Language = language
		}
		if content, ok := data["content"].(string); ok && content != "" {
			prompt.Content = content
		}
		if modelID, ok := data["model_id"].(uint); ok {
			prompt.ModelID = &modelID
		}

		prompt.Version++
		return txRepo.Update(prompt)
	})
}

// Rollback 回滚版本
func (s *PromptService) Rollback(id uint, version int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.promptRepo.WithTx(tx)
		prompt, err := txRepo.GetByID(id)
		if err != nil {
			return err
		}

		history, err := txRepo.GetHistory(id, version)
		if err != nil {
			return fmt.Errorf("未找到指定版本的历史记录: %w", err)
		}

		// 保存当前版本为历史记录
		currentHistory := &models.PromptHistory{
			PromptID:  prompt.ID,
			Content:   prompt.Content,
			Version:   prompt.Version,
			CreatedAt: prompt.UpdatedAt,
		}
		if err := txRepo.CreateHistory(currentHistory); err != nil {
			return err
		}

		// 回滚
		prompt.Content = history.Content
		prompt.Version++ // 回滚也视为一次新版本，或者保持版本号？这里选择增加版本号。
		return txRepo.Update(prompt)
	})
}

// Test 测试提示词
func (s *PromptService) Test(id uint, testData map[string]interface{}) (string, error) {
	prompt, err := s.promptRepo.GetByID(id)
	if err != nil {
		return "", err
	}

	modelID := prompt.ModelID
	if modelID == nil {
		return "", errors.New("提示词未配置模型")
	}

	model, err := s.modelRepo.GetByID(*modelID)
	if err != nil {
		return "", ErrModelNotFound
	}

	// 创建AI客户端
	client := ai.NewClient(model.APIURL, model.APIKey, 60)
	
	// 构建测试内容
	content := prompt.Content
	for k, v := range testData {
		placeholder := fmt.Sprintf("{{.%s}}", k)
		content = strings.ReplaceAll(content, placeholder, fmt.Sprintf("%v", v))
	}

	startTime := time.Now()
	res, err := client.Chat([]ai.Message{{Role: "user", Content: content}}, "")
	duration := int(time.Since(startTime).Milliseconds())

	if err != nil {
		s.logServ.LogAICall(*modelID, content, err.Error(), duration, false)
		return "", err
	}

	s.logServ.LogAICall(*modelID, content, res, duration, true)
	return res, nil
}

// GetHistoryList 获取历史列表
func (s *PromptService) GetHistoryList(id uint) ([]models.PromptHistory, error) {
	return s.promptRepo.GetHistoryList(id)
}

// Delete 删除提示词
func (s *PromptService) Delete(id uint) error {
	return s.promptRepo.Delete(id)
}

// IncrementVersion 增加版本号
func (s *PromptService) IncrementVersion(id uint) error {
	return s.promptRepo.IncrementVersion(id)
}

// GetByTypeAndScene 根据类型和场景获取提示词
func (s *PromptService) GetByTypeAndScene(promptType, scene string) ([]models.Prompt, error) {
	return s.promptRepo.GetByTypeAndScene(promptType, scene)
}
