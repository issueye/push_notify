package services

import (
	"fmt"
	"io"
	"strings"
	"sync"
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
	codeReviewQ  *CodeReviewQueue
	pushNotifyQ  *PushNotifyQueue
	baseURL      string
}

func NewWebhookService(db *gorm.DB, baseURL string) *WebhookService {
	s := &WebhookService{
		db:           db,
		repoRepo:     repository.NewRepoRepo(db),
		targetRepo:   repository.NewTargetRepo(db),
		pushRepo:     repository.NewPushRepo(db),
		templateRepo: repository.NewTemplateRepo(db),
		promptRepo:   repository.NewPromptRepo(db),
		modelRepo:    repository.NewAIModelRepo(db),
		codeviewServ: NewCodeViewService(db),
		baseURL:      baseURL,
	}
	s.codeReviewQ = NewCodeReviewQueue(200, 2, s.processCodeReviewJob)
	s.pushNotifyQ = NewPushNotifyQueue(500, 5, s.processPushNotifyJob)
	return s
}

func (s *WebhookService) processPushNotifyJob(job PushNotifyJob) {
	s.sendUnifiedPushNotification(job.Repo, job.Target, job.Payload, job.Template, job.Provider)
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
		t := target // 局部变量，防止闭包问题
		s.pushNotifyQ.Enqueue(PushNotifyJob{
			Repo:     repo,
			Target:   &t,
			Payload:  payload,
			Template: template,
			Provider: provider,
		})
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
	err = s.sendDingTalk(target, content)

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
		s.codeReviewQ.Enqueue(CodeReviewJob{
			RepoID:   repo.ID,
			PushID:   push.ID,
			CommitID: payload.After,
			Branch:   payload.Branch,
		})
	}
}

func (s *WebhookService) processCodeReviewJob(job CodeReviewJob) {
	repo, err := s.repoRepo.GetByID(job.RepoID)
	if err != nil {
		logger.Error("Repo not found for codeview", map[string]interface{}{
			"repo_id": job.RepoID,
		})
		return
	}

	push, err := s.pushRepo.GetByID(job.PushID)
	if err != nil {
		logger.Error("Push not found for codeview", map[string]interface{}{
			"push_id": job.PushID,
		})
		return
	}

	logger.Info("Starting code review", map[string]interface{}{
		"repo_id":   repo.ID,
		"commit_id": job.CommitID,
	})

	s.pushRepo.UpdateCodeview(repo.ID, job.CommitID, models.CodeviewStatusPending, nil)

	// 默认使用 go-git
	gitClient := git.NewGoGitClient(repo.URL, repo.AccessToken)
	if gitClient == nil {
		logger.Warn("Unsupported repo type for codeview", map[string]interface{}{
			"repo_id": repo.ID,
			"type":    repo.Type,
		})
		resultText := "不支持的仓库类型，已跳过"
		s.pushRepo.UpdateCodeview(repo.ID, job.CommitID, models.CodeviewStatusSkipped, &resultText)
		return
	}

	// 获取差异文件
	files, err := gitClient.GetSingleCommitDiff(job.CommitID)
	if err != nil {
		logger.Error("Failed to get diff", map[string]interface{}{
			"repo_id":   repo.ID,
			"commit_id": job.CommitID,
			"error":     err.Error(),
		})
		resultText := "获取差异失败: " + err.Error()
		s.pushRepo.UpdateCodeview(repo.ID, job.CommitID, models.CodeviewStatusFailed, &resultText)
		return
	}

	// 过滤需要审查的文件 (代码文件)
	codeFiles := filterCodeFiles(files)
	if len(codeFiles) == 0 {
		logger.Info("No code files to review", map[string]interface{}{
			"repo_id":   repo.ID,
			"commit_id": job.CommitID,
			"files":     len(files),
		})
		resultText := "无代码文件，已跳过"
		s.pushRepo.UpdateCodeview(repo.ID, job.CommitID, models.CodeviewStatusSkipped, &resultText)
		s.sendReviewNotification(repo, push, codeFiles, resultText)
		return
	}

	type fileTask struct {
		file git.DiffFile
	}
	type fileResult struct {
		fileName string
		summary  string
		err      error
	}

	fileCh := make(chan fileTask, len(codeFiles))
	resultCh := make(chan fileResult, len(codeFiles))

	workerCount := 3
	if workerCount > len(codeFiles) {
		workerCount = len(codeFiles)
	}

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for task := range fileCh {
				input := CodeViewInput{
					FileName:    task.file.Filename,
					DiffContent: task.file.Patch,
					RepoName:    repo.Name,
					Branch:      job.Branch,
					CommitMsg:   push.CommitMsg,
					Language:    detectLanguage(task.file.Filename),
				}

				res, err := s.codeviewServ.Review(repo.ID, input)
				if err != nil {
					resultCh <- fileResult{fileName: task.file.Filename, err: err}
					continue
				}

				if res != nil {
					resultCh <- fileResult{fileName: task.file.Filename, summary: strings.TrimSpace(res.Summary)}
				} else {
					resultCh <- fileResult{fileName: task.file.Filename}
				}
			}
		}()
	}

	for _, f := range codeFiles {
		fileCh <- fileTask{file: f}
	}
	close(fileCh)

	wg.Wait()
	close(resultCh)

	var allIssues strings.Builder
	for r := range resultCh {
		if r.err != nil {
			logger.Error("Failed to review file", map[string]interface{}{
				"file_name": r.fileName,
				"error":     r.err.Error(),
			})
			allIssues.WriteString(fmt.Sprintf("### %s\n审查失败: %s\n\n", r.fileName, r.err.Error()))
			continue
		}
		if r.summary != "" {
			allIssues.WriteString(fmt.Sprintf("### %s\n%s\n\n", r.fileName, r.summary))
		}
	}

	resultText := strings.TrimSpace(allIssues.String())
	if strings.TrimSpace(resultText) == "" {
		resultText = "未发现明显问题"
	}
	s.pushRepo.UpdateCodeview(repo.ID, job.CommitID, models.CodeviewStatusSuccess, &resultText)

	// 发送审查结果通知
	s.sendReviewNotification(repo, push, codeFiles, resultText)

	logger.Info("Code review completed", map[string]interface{}{
		"repo_id":   repo.ID,
		"commit_id": job.CommitID,
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
	return s.sendDingTalkMarkdown(target, "代码提交通知", content)
}

func (s *WebhookService) sendDingTalkMarkdown(target *models.Target, title string, content string) error {
	if target.Config == nil {
		return fmt.Errorf("config is required for DingTalk target")
	}

	accessToken := target.Config.AccessToken
	secret := target.Config.Secret

	if accessToken == "" {
		return fmt.Errorf("access_token is required for DingTalk target")
	}

	client := dingtalk.NewClient(accessToken, secret)

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
				s.sendDingTalkMarkdown(&target, "代码审查报告", content)
			}
		}
	}
}

// buildReviewMessageContent 构建审查结果消息内容
func (s *WebhookService) buildReviewMessageContent(repo *models.Repo, push *models.Push, issues string, template *models.Template) string {
	if strings.TrimSpace(issues) == "" {
		issues = "未发现明显问题"
	}

	if template == nil || template.Content == "" {
		// 生成审查链接
		reviewURL := ""
		if s.baseURL != "" {
			reviewURL = fmt.Sprintf("%s/web/#/pushes/review?id=%d", s.baseURL, push.ID)
		}

		var content strings.Builder
		content.WriteString("### 🔍 代码审查结果\n\n")
		content.WriteString("**仓库名称：** " + repo.Name + "\n")
		content.WriteString("**提交ID：** `" + push.CommitID + "`\n")
		content.WriteString("**提交信息：** " + push.CommitMsg + "\n\n")
		content.WriteString("---\n")
		if reviewURL != "" {
			content.WriteString("[查看审查详情](" + reviewURL + ")")
		} else {
			content.WriteString(issues)
		}
		return content.String()
	}

	// 简单的模板替换
	content := template.Content
	content = strings.ReplaceAll(content, "{{.RepoName}}", repo.Name)
	content = strings.ReplaceAll(content, "{{.CommitID}}", push.CommitID)
	content = strings.ReplaceAll(content, "{{.CommitMsg}}", push.CommitMsg)
	content = strings.ReplaceAll(content, "{{.Issues}}", issues)

	// 生成审查链接
	reviewURL := ""
	if s.baseURL != "" {
		reviewURL = fmt.Sprintf("%s/web/#/pushes/review?id=%d", s.baseURL, push.ID)
	}
	content = strings.ReplaceAll(content, "{{.ReviewURL}}", reviewURL)

	return content
}
