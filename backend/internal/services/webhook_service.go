package services

import (
	"fmt"
	"io"
	"strings"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/dingtalk"
	"backend/pkg/git"
	"backend/utils/logger"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebhookService struct {
	db           *gorm.DB
	repoRepo     *repository.RepoRepo
	targetRepo   *repository.TargetRepo
	pushRepo     *repository.PushRepo
	templateRepo *repository.TemplateRepo
	promptRepo   *repository.PromptRepo
	modelRepo    *repository.AIModelRepo
	codeviewServ *CodeViewService
}

func NewWebhookService(db *gorm.DB) *WebhookService {
	return &WebhookService{
		db:           db,
		repoRepo:     repository.NewRepoRepo(db),
		targetRepo:   repository.NewTargetRepo(db),
		pushRepo:     repository.NewPushRepo(db),
		templateRepo: repository.NewTemplateRepo(db),
		promptRepo:   repository.NewPromptRepo(db),
		modelRepo:    repository.NewAIModelRepo(db),
		codeviewServ: NewCodeViewService(db),
	}
}

// HandleGitHubWebhook 处理GitHub Webhook
func (s *WebhookService) HandleGitHubWebhook(c *gin.Context) {
	s.handleWebhook(c, &GitHubProvider{})
}

// HandleGitLabWebhook 处理GitLab Webhook
func (s *WebhookService) HandleGitLabWebhook(c *gin.Context) {
	s.handleWebhook(c, &GitLabProvider{})
}

// handleWebhook 通用Webhook处理逻辑
func (s *WebhookService) handleWebhook(c *gin.Context, provider WebhookProvider) {
	webhookID := c.Param("webhookId")

	// 获取仓库
	repo, err := s.repoRepo.GetByWebhookID(webhookID)
	if err != nil {
		logger.Error("Repo not found for webhook", map[string]interface{}{
			"webhook_id": webhookID,
		})
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "仓库不存在",
		})
		return
	}

	// 读取请求体
	body, _ := io.ReadAll(c.Request.Body)

	// 解析事件类型
	eventType := provider.GetEventType(c.Request.Header)
	logger.Info("Received webhook", map[string]interface{}{
		"event":     eventType,
		"repo_id":   repo.ID,
		"repo_name": repo.Name,
	})

	// 统一处理逻辑
	// 注意：GitHub的push事件是 "push"，GitLab的push事件是 "Push Hook"
	if eventType == "push" || eventType == "Push Hook" {
		payload, err := provider.ParsePushPayload(body)
		if err != nil {
			logger.Error("Failed to parse webhook payload", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		s.dispatchPushNotification(repo, payload, provider)
	} else if eventType == "ping" || eventType == "Event Hook" {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
		return
	} else {
		logger.Info("Unsupported event type", map[string]interface{}{
			"event": eventType,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Webhook处理成功",
		"data": map[string]interface{}{
			"repo_id":   repo.ID,
			"repo_name": repo.Name,
			"status":    "processing",
		},
	})
}

// dispatchPushNotification 分发推送通知
func (s *WebhookService) dispatchPushNotification(repo *models.Repo, payload *UnifiedPushPayload, provider WebhookProvider) {
	// 获取推送目标
	targets, err := s.targetRepo.GetByScopeAndRepo(repo.ID)
	if err != nil {
		logger.Error("Failed to get targets", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if len(targets) == 0 {
		logger.Info("No targets configured", map[string]interface{}{
			"repo_id": repo.ID,
		})
		return
	}

	// 获取默认模板
	var template *models.Template
	if repo.CommitTemplate != nil {
		template = repo.CommitTemplate
	} else {
		template, _ = s.templateRepo.GetByTypeAndScene(models.TemplateTypeDingTalk, models.TemplateSceneCommitNotify)
	}

	// 为每个推送目标发送通知
	for _, target := range targets {
		go s.sendUnifiedPushNotification(repo, &target, payload, template, provider)
	}
}

// sendUnifiedPushNotification 发送统一推送通知
func (s *WebhookService) sendUnifiedPushNotification(repo *models.Repo, target *models.Target, payload *UnifiedPushPayload, template *models.Template, provider WebhookProvider) {
	// 去重检查：同一提交同一目标不重复推送
	if s.pushRepo.ExistsByCommitAndTarget(payload.After, target.ID) {
		logger.Info("Duplicate push detected, skipping", map[string]interface{}{
			"commit_id": payload.After[:7],
			"target_id": target.ID,
			"repo_name": repo.Name,
		})
		return
	}

	// 构建消息内容
	content := provider.BuildMessage(payload, template)

	// 创建推送记录
	push := &models.Push{
		RepoID:    repo.ID,
		TargetID:  target.ID,
		CommitID:  payload.After,
		CommitMsg: payload.CommitMsg,
		Status:    models.PushStatusPending,
		Content:   content,
	}

	if template != nil {
		push.TemplateID = &template.ID
	}

	if err := s.pushRepo.Create(push); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
			logger.Info("Duplicate push detected (DB constraint), skipping", map[string]interface{}{
				"commit_id": payload.After[:7],
				"target_id": target.ID,
				"repo_name": repo.Name,
			})
			return
		}
		logger.Error("Failed to create push record", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// 发送通知
	var err error
	if target.Type == models.TargetTypeDingTalk {
		err = s.sendDingTalk(target, content)
	} else if target.Type == models.TargetTypeWebhook {
		err = s.sendWebhook(target, content)
	}

	// 更新推送状态
	if err != nil {
		push.Status = models.PushStatusFailed
		push.ErrorMsg = err.Error()
		logger.Error("Push failed", map[string]interface{}{
			"push_id":   push.ID,
			"target_id": target.ID,
			"error":     err.Error(),
		})
	} else {
		push.Status = models.PushStatusSuccess
		logger.Info("Push succeeded", map[string]interface{}{
			"push_id":   push.ID,
			"target_id": target.ID,
		})
	}

	s.pushRepo.Update(push)

	// 执行代码审查 (异步)
	if push.Status == models.PushStatusSuccess && repo.ModelID != nil {
		go func() {
			s.doCodeReview(repo, payload, push)
		}()
	}
}

// doCodeReview 执行代码审查
func (s *WebhookService) doCodeReview(repo *models.Repo, payload *UnifiedPushPayload, push *models.Push) {
	logger.Info("Starting code review", map[string]interface{}{
		"repo_id":   repo.ID,
		"commit_id": payload.After,
	})

	// 更新状态为进行中
	push.CodeviewStatus = models.CodeviewStatusPending
	s.pushRepo.Update(push)

	// 根据仓库类型创建Git客户端
	var gitClient git.GitClient
	if repo.Type == models.RepoTypeGitHub {
		gitClient = git.NewClient(
			repo.URL,
			repo.AccessToken,
			extractOwner(repo.URL),
			extractRepoName(repo.URL),
		)
	} else if repo.Type == models.RepoTypeGitLab {
		// 提取BaseURL
		baseURL := extractGitLabBaseURL(repo.URL)
		// 对于GitLab，Project ID可以是数字ID或URL编码的path_with_namespace
		// 这里假设repo.Name是path_with_namespace (e.g. group/project)
		projectID := url.PathEscape(repo.Name)
		gitClient = git.NewGitLabClient(
			baseURL,
			repo.AccessToken,
			projectID,
		)
	}

	if gitClient == nil {
		logger.Warn("Unsupported repo type for codeview", map[string]interface{}{
			"repo_id": repo.ID,
			"type":    repo.Type,
		})
		push.CodeviewStatus = models.CodeviewStatusSkipped
		s.pushRepo.Update(push)
		return
	}

	// 获取差异文件
	files, err := gitClient.GetSingleCommitDiff(payload.After)
	if err != nil {
		logger.Error("Failed to get diff", map[string]interface{}{
			"repo_id":   repo.ID,
			"commit_id": payload.After,
			"error":     err.Error(),
		})
		push.CodeviewStatus = models.CodeviewStatusFailed
		s.pushRepo.Update(push)
		return
	}

	// 过滤需要审查的文件 (代码文件)
	codeFiles := filterCodeFiles(files)
	if len(codeFiles) == 0 {
		logger.Info("No code files to review", map[string]interface{}{
			"repo_id":   repo.ID,
			"commit_id": payload.After,
			"files":     len(files),
		})
		push.CodeviewStatus = models.CodeviewStatusSkipped
		s.pushRepo.Update(push)
		return
	}

	// 获取分支
	branch := payload.Branch

	// 批量审查
	var allIssues strings.Builder
	for _, file := range codeFiles {
		input := CodeViewInput{
			FileName:    file.Filename,
			DiffContent: file.Patch,
			RepoName:    repo.Name,
			Branch:      branch,
			CommitMsg:   payload.CommitMsg,
			Language:    detectLanguage(file.Filename),
		}

		result, err := s.codeviewServ.Review(repo.ID, input)
		if err != nil {
			logger.Error("Failed to review file", map[string]interface{}{
				"file_name": file.Filename,
				"error":     err.Error(),
			})
			continue
		}

		if result != nil && result.Summary != "" {
			allIssues.WriteString(fmt.Sprintf("### %s\n%s\n\n", file.Filename, result.Summary))
		}
	}

	// 更新推送记录
	if allIssues.Len() > 0 {
		result := allIssues.String()
		push.CodeviewResult = &result
	}
	push.CodeviewStatus = models.CodeviewStatusSuccess
	s.pushRepo.Update(push)

	// 发送审查结果通知
	if push.CodeviewStatus == models.CodeviewStatusSuccess && allIssues.Len() > 0 {
		s.sendReviewNotification(repo, push, codeFiles, allIssues.String())
	}

	logger.Info("Code review completed", map[string]interface{}{
		"repo_id":   repo.ID,
		"commit_id": payload.After,
		"files":     len(codeFiles),
	})
}

// extractOwner 从URL提取owner
func extractOwner(urlStr string) string {
	parts := strings.Split(strings.TrimSuffix(urlStr, ".git"), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return ""
}

// extractRepoName 从URL提取仓库名
func extractRepoName(urlStr string) string {
	parts := strings.Split(strings.TrimSuffix(urlStr, ".git"), "/")
	if len(parts) >= 1 {
		return parts[len(parts)-1]
	}
	return ""
}

// extractGitLabBaseURL 从URL提取GitLab BaseURL
func extractGitLabBaseURL(repoURL string) string {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "https://gitlab.com" // Default fallback
	}
	return u.Scheme + "://" + u.Host
}

// filterCodeFiles 过滤需要审查的代码文件
func filterCodeFiles(files []git.DiffFile) []git.DiffFile {
	codeExtensions := map[string]bool{
		".go": true, ".java": true, ".py": true, ".js": true, ".ts": true,
		".cpp": true, ".c": true, ".h": true, ".cs": true, ".rb": true,
		".php": true, ".swift": true, ".kt": true, ".rs": true, ".vue": true,
		".tsx": true, ".jsx": true, ".yaml": true, ".yml": true,
	}

	var result []git.DiffFile
	for _, f := range files {
		ext := getFileExtension(f.Filename)
		if _, ok := codeExtensions[ext]; ok && f.Status != "deleted" {
			result = append(result, f)
		}
	}
	return result
}

// getFileExtension 获取文件扩展名
func getFileExtension(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx > 0 {
		return strings.ToLower(filename[idx:])
	}
	return ""
}

// detectLanguage 检测编程语言
func detectLanguage(filename string) string {
	ext := getFileExtension(filename)
	langMap := map[string]string{
		".go":   "Go",
		".java": "Java",
		".py":   "Python",
		".js":   "JavaScript",
		".ts":   "TypeScript",
		".vue":  "Vue",
		".rb":   "Ruby",
		".php":  "PHP",
		".rs":   "Rust",
		".cpp":  "C++",
		".c":    "C",
		".cs":   "C#",
	}
	if lang, ok := langMap[ext]; ok {
		return lang
	}
	return "Unknown"
}

// sendDingTalk 发送钉钉通知
func (s *WebhookService) sendDingTalk(target *models.Target, content string) error {
	if target.Config == nil {
		return fmt.Errorf("config is required for DingTalk target")
	}

	accessToken := target.Config.AccessToken
	secret := target.Config.Secret

	if accessToken == "" {
		return fmt.Errorf("access_token is required for DingTalk target")
	}

	client := dingtalk.NewClient(accessToken, secret)

	// 解析内容生成Markdown
	title := "代码提交通知"
	return client.SendMarkdown(target.Config.WebhookURL, title, content)
}

// sendWebhook 发送Webhook通知
func (s *WebhookService) sendWebhook(target *models.Target, content string) error {
	if target.Config == nil {
		return fmt.Errorf("config is required for webhook target")
	}

	if target.Config.WebhookURL == "" {
		return fmt.Errorf("webhook_url is required for webhook target")
	}

	method := target.Config.Method

	// 准备请求体
	var body io.Reader
	if method == "POST" {
		body = strings.NewReader(content)
	}

	req, err := http.NewRequest(method, target.Config.WebhookURL, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned error status: %d", resp.StatusCode)
	}

	return nil
}

// sendReviewNotification 发送审查结果通知
func (s *WebhookService) sendReviewNotification(repo *models.Repo, push *models.Push, codeFiles []git.DiffFile, issues string) {
	// 获取推送目标
	targets, err := s.targetRepo.GetByScopeAndRepo(repo.ID)
	if err != nil || len(targets) == 0 {
		return
	}

	// 确定涉及的语言
	languages := make(map[string]bool)
	for _, file := range codeFiles {
		lang := detectLanguage(file.Filename)
		if lang != "Unknown" {
			languages[lang] = true
		}
	}

	// 匹配模板
	templatesToSend := make(map[uint]*models.Template)

	if len(repo.ReviewTemplates) == 0 {
		// 使用系统默认
		tpl, err := s.templateRepo.GetByTypeAndScene(models.TemplateTypeDingTalk, models.TemplateSceneReviewNotify)
		if err == nil {
			templatesToSend[tpl.ID] = tpl
		}
	} else {
		matched := false
		// 尝试匹配特定语言
		for lang := range languages {
			for _, rt := range repo.ReviewTemplates {
				if strings.EqualFold(rt.Language, lang) {
					templatesToSend[rt.TemplateID] = &rt.Template
					matched = true
				}
			}
		}

		// 如果没有特定语言匹配，且没有已选模板，尝试 default
		if !matched && len(templatesToSend) == 0 {
			for _, rt := range repo.ReviewTemplates {
				if rt.Language == "default" {
					templatesToSend[rt.TemplateID] = &rt.Template
				}
			}
		}
	}

	// 如果没有找到任何模板，不发送
	if len(templatesToSend) == 0 {
		return
	}

	// 发送通知
	for _, tpl := range templatesToSend {
		content := s.buildReviewMessageContent(repo, push, issues, tpl)
		for _, target := range targets {
			if target.Type == models.TargetTypeDingTalk {
				s.sendDingTalk(&target, content)
			} else if target.Type == models.TargetTypeWebhook {
				s.sendWebhook(&target, content)
			}
		}
	}
}

// buildReviewMessageContent 构建审查结果消息内容
func (s *WebhookService) buildReviewMessageContent(repo *models.Repo, push *models.Push, issues string, template *models.Template) string {
	if template == nil || template.Content == "" {
		var content strings.Builder
		content.WriteString("## 代码审查报告\n\n")
		content.WriteString("**仓库**: " + repo.Name + "\n")
		content.WriteString("**提交**: " + push.CommitID[:7] + "\n")
		content.WriteString("**信息**: " + push.CommitMsg + "\n\n")
		content.WriteString(issues)
		return content.String()
	}

	// 简单的模板替换
	content := template.Content
	content = strings.ReplaceAll(content, "{{.RepoName}}", repo.Name)
	content = strings.ReplaceAll(content, "{{.CommitID}}", push.CommitID)
	content = strings.ReplaceAll(content, "{{.CommitMsg}}", push.CommitMsg)
	content = strings.ReplaceAll(content, "{{.Issues}}", issues)
	
	return content
}
