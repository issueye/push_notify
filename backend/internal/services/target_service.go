package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/dingtalk"
	"backend/utils/logger"

	"gorm.io/gorm"
)

var (
	ErrTargetNotFound       = errors.New("推送目标不存在")
	ErrTargetAlreadyExists  = errors.New("推送目标名称已存在")
	ErrInvalidDingTalkToken = errors.New("无效的钉钉AccessToken")
)

type TargetService struct {
	targetRepo *repository.TargetRepo
	repoRepo   *repository.RepoRepo
}

func NewTargetService(db *gorm.DB) *TargetService {
	return &TargetService{
		targetRepo: repository.NewTargetRepo(db),
		repoRepo:   repository.NewRepoRepo(db),
	}
}

// Create 创建推送目标
func (s *TargetService) Create(data map[string]interface{}) (*models.Target, error) {
	name := data["name"].(string)
	targetType := data["type"].(string)

	// 验证配置
	configStr, err := json.Marshal(data["config"])
	if err != nil {
		return nil, err
	}

	// 验证配置有效性
	if targetType == models.TargetTypeDingTalk {
		if err := validateDingTalkConfig(string(configStr)); err != nil {
			return nil, err
		}
	} else if targetType == models.TargetTypeWebhook {
		if err := validateWebhookConfig(string(configStr)); err != nil {
			return nil, err
		}
	}

	cfg := &models.Config{}
	if err := json.Unmarshal(configStr, cfg); err != nil {
		return nil, err
	}

	target := &models.Target{
		Name:   name,
		Type:   targetType,
		Config: cfg,
		Scope:  getStringWithDefault(data, "scope", models.TargetScopeGlobal),
		Status: models.StatusActive,
	}

	if err := s.targetRepo.Create(target); err != nil {
		return nil, err
	}

	// 如果是指定仓库模式，关联仓库
	if scope, ok := data["scope"].(string); ok && scope == models.TargetScopeRepo {
		if repoIDs, ok := data["repo_ids"].([]interface{}); ok {
			for _, id := range repoIDs {
				if repoID, ok := id.(uint); ok {
					s.targetRepo.AddRepo(target.ID, repoID)
				}
			}
		}
	}

	logger.Info("Target created", map[string]interface{}{
		"target_id": target.ID,
		"name":      target.Name,
		"type":      target.Type,
	})

	return target, nil
}

// GetByID 获取推送目标详情
func (s *TargetService) GetByID(id uint) (*models.Target, error) {
	return s.targetRepo.GetByID(id)
}

// GetList 获取推送目标列表
func (s *TargetService) GetList(page, size int, keyword, targetType, scope string) ([]models.Target, int64, error) {
	targets, total := s.targetRepo.GetList(page, size, keyword, targetType, scope)
	return targets, total, nil
}

// Update 更新推送目标
func (s *TargetService) Update(id uint, data map[string]interface{}) error {
	target, err := s.targetRepo.GetByID(id)
	if err != nil {
		return err
	}

	if name, ok := data["name"].(string); ok && name != "" {
		target.Name = name
	}
	if config, ok := data["config"].(map[string]interface{}); ok {
		configStr, _ := json.Marshal(config)
		cfg := &models.Config{}
		if err := json.Unmarshal(configStr, cfg); err != nil {
			return err
		}
		target.Config = cfg
	}
	if scope, ok := data["scope"].(string); ok && scope != "" {
		target.Scope = scope
	}
	if status, ok := data["status"].(string); ok && status != "" {
		target.Status = status
	}

	return s.targetRepo.Update(target)
}

// Delete 删除推送目标
func (s *TargetService) Delete(id uint) error {
	return s.targetRepo.Delete(id)
}

// Test 发送测试消息
func (s *TargetService) Test(id uint) (map[string]interface{}, error) {
	target, err := s.targetRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	var sendErr error

	switch target.Type {
	case models.TargetTypeDingTalk:
		result, sendErr = s.testDingTalk(target)
	case models.TargetTypeWebhook:
		result, sendErr = s.testWebhook(target)
	default:
		return nil, errors.New("不支持的推送类型: " + target.Type)
	}

	if sendErr != nil {
		logger.Error("Target test failed", map[string]interface{}{
			"target_id": id,
			"name":      target.Name,
			"type":      target.Type,
			"error":     sendErr.Error(),
		})
		return nil, sendErr
	}

	logger.Info("Target test success", map[string]interface{}{
		"target_id": id,
		"name":      target.Name,
		"type":      target.Type,
	})

	return result, nil
}

// testDingTalk 测试钉钉通知
func (s *TargetService) testDingTalk(target *models.Target) (map[string]interface{}, error) {
	if target.Config == nil {
		return nil, errors.New("钉钉配置无效")
	}

	accessToken := target.Config.AccessToken
	secret := target.Config.Secret

	if accessToken == "" {
		return nil, errors.New("access_token不能为空")
	}

	client := dingtalk.NewClient(accessToken, secret)

	// 构建测试消息
	content := `## 代码提交通知

**这是一条测试消息**

- 推送目标: ` + target.Name + `
- 测试时间: ` + time.Now().Format("2006-01-02 15:04:05") + `
- 状态: 正常

如果收到此消息，说明钉钉配置正确。`

	title := "推送通知测试"
	if err := client.SendMarkdown(target.Config.WebhookURL, title, content); err != nil {
		return nil, errors.New("发送失败: " + err.Error())
	}

	return map[string]interface{}{
		"status":  "success",
		"message": "测试消息已发送",
		"type":    "dingtalk",
	}, nil
}

// testWebhook 测试Webhook通知
func (s *TargetService) testWebhook(target *models.Target) (map[string]interface{}, error) {
	if target.Config == nil {
		return nil, errors.New("Webhook配置无效")
	}

	webhookURL := target.Config.WebhookURL
	if webhookURL == "" {
		return nil, errors.New("webhook_url不能为空")
	}

	// 构建测试消息
	testPayload := map[string]interface{}{
		"event":       "test",
		"message":     "这是一条测试消息",
		"target_name": target.Name,
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	body, _ := json.Marshal(testPayload)

	method := "POST"
	if target.Config.Method != "" {
		method = target.Config.Method
	}

	req, err := http.NewRequest(method, webhookURL, bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("创建请求失败: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range target.Config.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("请求失败: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New("返回状态码: " + resp.Status)
	}

	// 读取响应内容（限制1000字符）
	respBody := make([]byte, 1000)
	n, _ := resp.Body.Read(respBody)
	respStr := strings.TrimSpace(string(respBody[:n]))

	return map[string]interface{}{
		"status":      "success",
		"message":     "测试消息已发送",
		"type":        "webhook",
		"status_code": resp.StatusCode,
		"response":    respStr,
	}, nil
}

// AddRepo 关联仓库
func (s *TargetService) AddRepo(targetID, repoID uint) error {
	return s.targetRepo.AddRepo(targetID, repoID)
}

// RemoveRepo 取消关联仓库
func (s *TargetService) RemoveRepo(targetID, repoID uint) error {
	return s.targetRepo.RemoveRepo(targetID, repoID)
}

// GetRepos 获取推送目标关联的仓库
func (s *TargetService) GetRepos(targetID uint) ([]models.Repo, error) {
	return s.targetRepo.GetRepos(targetID)
}

// GetByScope 获取指定范围的推送目标
func (s *TargetService) GetByScope(repoID uint) ([]models.Target, error) {
	return s.targetRepo.GetByScopeAndRepo(repoID)
}

func validateDingTalkConfig(configStr string) error {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return ErrInvalidDingTalkToken
	}
	if _, ok := config["access_token"]; !ok {
		return ErrInvalidDingTalkToken
	}
	return nil
}

func validateWebhookConfig(configStr string) error {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return errors.New("无效的Webhook配置")
	}
	if _, ok := config["webhook_url"]; !ok || config["webhook_url"] == "" {
		return errors.New("Webhook URL不能为空")
	}
	return nil
}
