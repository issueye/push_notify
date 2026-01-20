package services

import (
	"encoding/json"
	"strings"
	"text/template"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/ai"
	"backend/utils/logger"

	"gorm.io/gorm"
)

type CodeViewService struct {
	db         *gorm.DB
	promptRepo *repository.PromptRepo
	modelRepo  *repository.AIModelRepo
	modelServ  *AIModelService
}

func NewCodeViewService(db *gorm.DB) *CodeViewService {
	return &CodeViewService{
		db:         db,
		promptRepo: repository.NewPromptRepo(db),
		modelRepo:  repository.NewAIModelRepo(db),
	}
}

// CodeViewInput CODEVIEW输入
type CodeViewInput struct {
	FileName    string
	FileContent string
	DiffContent string
	Language    string
	RepoName    string
	Branch      string
	CommitMsg   string
}

// CodeViewResult CODEVIEW结果
type CodeViewResult struct {
	Result  string  `json:"result"` // 通过/有建议/有问题
	Issues  []Issue `json:"issues"`
	Summary string  `json:"summary"`
}

// Issue 问题
type Issue struct {
	Line     int    `json:"line"`
	Severity string `json:"severity"` // info/warning/error
	Message  string `json:"message"`
	Code     string `json:"code"`
}

// Review 执行代码审查
func (s *CodeViewService) Review(repoID uint, input CodeViewInput) (*CodeViewResult, error) {
	// 获取仓库信息
	repoRepo := repository.NewRepoRepo(s.db)
	repo, err := repoRepo.GetByID(repoID)
	if err != nil {
		return nil, err
	}

	// 获取提示词
	prompts, err := s.promptRepo.GetByTypeAndScene(models.PromptTypeCodeView, "")
	if err != nil || len(prompts) == 0 {
		// 使用默认提示词 (优先使用diff审查)
		prompts = []models.Prompt{
			{
				Content: `请作为专业的代码审查助手，审查以下代码变更（diff格式）。

审查要点：
1. 代码规范和最佳实践
2. 潜在的bug和安全漏洞
3. 性能问题
4. 逻辑错误
5. 是否有重复代码

请以简洁的方式输出审查结果，直接列出问题，无需过多铺垫。

--- 变更信息 ---
文件名：{{.FileName}}
语言：{{.Language}}
仓库：{{.RepoName}}
分支：{{.Branch}}
提交信息：{{.CommitMsg}}

--- 代码差异 ---
{{if .DiffContent}}
{{.DiffContent}}
{{else}}
{{.FileContent}}
{{end}}`,
			},
		}
	}

	// 获取模型
	var aiClient *ai.Client
	if repo.ModelID != nil {
		model, err := s.modelRepo.GetByID(*repo.ModelID)
		if err == nil {
			var params map[string]interface{}
			json.Unmarshal([]byte(model.Params), &params)
			modelName := strings.TrimSpace(model.Name)
			if modelName == "" {
				modelName = strings.TrimSpace(model.Type)
			}
			aiClient = ai.NewClientWithConfig(ai.Config{
				APIURL: model.APIURL,
				APIKey: model.APIKey,
				Model:  modelName,
				Params: params,
			})
		}
	}

	// 如果没有指定模型，使用默认模型
	if aiClient == nil {
		model, err := s.modelRepo.GetDefault()
		if err == nil {
			var params map[string]interface{}
			json.Unmarshal([]byte(model.Params), &params)
			modelName := strings.TrimSpace(model.Name)
			if modelName == "" {
				modelName = strings.TrimSpace(model.Type)
			}
			aiClient = ai.NewClientWithConfig(ai.Config{
				APIURL: model.APIURL,
				APIKey: model.APIKey,
				Model:  modelName,
				Params: params,
			})
		}
	}

	// 如果还是没有AI客户端，返回空结果
	if aiClient == nil {
		return &CodeViewResult{
			Result:  "跳过",
			Summary: "未配置AI模型",
		}, nil
	}

	// 使用第一个提示词
	prompt := prompts[0]

	// 渲染提示词
	promptText, err := s.renderPrompt(prompt.Content, input)
	if err != nil {
		return nil, err
	}


	// 调用AI
	// 提示词模板中已经包含代码内容，这里传空字符串避免重复
	result, err := aiClient.CodeReview(promptText, "")
	if err != nil {
		return nil, err
	}

	// 解析结果
	codeViewResult := s.parseResult(result)

	// 更新模型调用次数
	if repo.ModelID != nil {
		s.modelRepo.IncrementCallCount(*repo.ModelID)
	}

	logger.Info("CodeView completed", map[string]interface{}{
		"repo_id":   repoID,
		"file_name": input.FileName,
		"result":    codeViewResult.Result,
	})

	return codeViewResult, nil
}

// renderPrompt 渲染提示词模板
func (s *CodeViewService) renderPrompt(promptTpl string, input CodeViewInput) (string, error) {
	funcMap := template.FuncMap{}

	tmpl, err := template.New("prompt").Funcs(funcMap).Parse(promptTpl)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, input)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// parseResult 解析AI返回结果
func (s *CodeViewService) parseResult(result string) *CodeViewResult {
	// 简化解析逻辑
	// 实际项目中可以使用更复杂的解析方式

	codeViewResult := &CodeViewResult{
		Result:  "有建议",
		Summary: result,
		Issues:  []Issue{},
	}

	// 检查结果中是否包含关键信息
	lowerResult := strings.ToLower(result)
	if strings.Contains(lowerResult, "通过") || strings.Contains(lowerResult, "pass") {
		codeViewResult.Result = "通过"
	}
	if strings.Contains(lowerResult, "错误") || strings.Contains(lowerResult, "error") {
		codeViewResult.Result = "有问题"
	}

	return codeViewResult
}

// BatchReview 批量审查
func (s *CodeViewService) BatchReview(repoID uint, files []CodeViewInput) ([]CodeViewResult, error) {
	results := make([]CodeViewResult, 0, len(files))

	for _, file := range files {
		result, err := s.Review(repoID, file)
		if err != nil {
			logger.Error("Failed to review file", map[string]interface{}{
				"file_name": file.FileName,
				"error":     err.Error(),
			})
			continue
		}
		results = append(results, *result)
	}

	return results, nil
}
