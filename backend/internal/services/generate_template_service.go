package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/ai"
	"backend/utils/logger"

	"gorm.io/gorm"
)

type GenerateTemplateService struct {
	modelRepo *repository.AIModelRepo
}

func NewGenerateTemplateService(db *gorm.DB) *GenerateTemplateService {
	return &GenerateTemplateService{
		modelRepo: repository.NewAIModelRepo(db),
	}
}

// GenerateTemplateInput 生成模板输入
type GenerateTemplateInput struct {
	Name    string `json:"name"`
	Type    string `json:"type"`  // dingtalk, email
	Scene   string `json:"scene"` // commit_notify, review_notify
	Title   string `json:"title"`
	ModelID uint   `json:"model_id"` // 可选，指定模型
}

// GenerateTemplate 生成模板内容
func (s *GenerateTemplateService) Generate(input GenerateTemplateInput) (string, error) {
	// 获取模型
	var aiModel *models.AIModel
	var err error

	if input.ModelID > 0 {
		aiModel, err = s.modelRepo.GetByID(input.ModelID)
	} else {
		aiModel, err = s.modelRepo.GetDefault()
	}

	if err != nil {
		return "", fmt.Errorf("failed to get AI model: %v", err)
	}

	// 构建提示词
	prompt := s.buildPrompt(input)

	// 初始化AI客户端
	var params map[string]interface{}
	if aiModel.Params != "" {
		json.Unmarshal([]byte(aiModel.Params), &params)
	}

	aiClient := ai.NewClientWithConfig(ai.Config{
		APIURL: aiModel.APIURL,
		APIKey: aiModel.APIKey,
		Model:  aiModel.Name,
		Params: params,
	})

	// 调用AI
	messages := []ai.Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	result, err := aiClient.Chat(messages, "你是一个专业的消息模板生成助手。")
	if err != nil {
		logger.Error("Failed to generate template", map[string]interface{}{
			"error":    err.Error(),
			"model_id": aiModel.ID,
			"api_url":  aiModel.APIURL,
		})
		return "", err
	}

	// 更新调用次数
	s.modelRepo.IncrementCallCount(aiModel.ID)

	return strings.TrimSpace(result), nil
}

func (s *GenerateTemplateService) buildPrompt(input GenerateTemplateInput) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("请为我生成一个 %s 消息通知模板。\n", input.Type))
	sb.WriteString(fmt.Sprintf("场景：%s\n", input.Scene))
	sb.WriteString(fmt.Sprintf("模板名称：%s\n", input.Name))
	sb.WriteString(fmt.Sprintf("消息标题：%s\n", input.Title))
	sb.WriteString("\n")

	if input.Type == models.TemplateTypeDingTalk {
		sb.WriteString("请生成 Markdown 格式的模板内容。\n")
	} else {
		sb.WriteString("请生成适合该渠道的模板内容。\n")
	}

	sb.WriteString("\n可用变量：\n")
	if input.Scene == models.TemplateSceneCommitNotify {
		sb.WriteString("- {{.RepoName}}: 仓库名称\n")
		sb.WriteString("- {{.Branch}}: 分支名称\n")
		sb.WriteString("- {{.CommitID}}: 提交ID\n")
		sb.WriteString("- {{.CommitMsg}}: 提交信息\n")
		sb.WriteString("- {{.Author}}: 作者\n")
		sb.WriteString("- {{.FileCount}}: 变更文件数\n")
		sb.WriteString("- {{.FileList}}: 变更文件列表\n")
	} else if input.Scene == models.TemplateSceneReviewNotify {
		sb.WriteString("- {{.RepoName}}: 仓库名称\n")
		sb.WriteString("- {{.CommitID}}: 提交ID\n")
		sb.WriteString("- {{.CommitMsg}}: 提交信息\n")
		sb.WriteString("- {{.Issues}}: 审查问题列表\n")
	}

	sb.WriteString("\n要求：\n")
	sb.WriteString("1. 只返回模板内容，不要包含任何解释性文字。\n")
	sb.WriteString("2. 确保变量格式正确，使用 Go Template 语法 (例如 {{.RepoName}})。\n")
	sb.WriteString("3. 排版美观，突出重点信息。\n")

	return sb.String()
}
