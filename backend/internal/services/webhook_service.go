package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/dingtalk"
	"backend/pkg/git"
	"backend/utils/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WebhookPayload GitHub Webhook通用负载
type WebhookPayload struct {
	Ref        string     `json:"ref"`
	Before     string     `json:"before"`
	After      string     `json:"after"`
	Repository Repository `json:"repository"`
	Pusher     Pusher     `json:"pusher"`
	Sender     Sender     `json:"sender"`
	Commits    []Commit   `json:"commits"`
	HeadCommit Commit     `json:"head_commit"`
}

// GitLabPayload GitLab Webhook负载
type GitLabPayload struct {
	ObjectKind        string         `json:"object_kind"`
	EventName         string         `json:"event_name"`
	Before            string         `json:"before"`
	After             string         `json:"after"`
	Ref               string         `json:"ref"`
	CheckoutSHA       string         `json:"checkout_sha"`
	User              GitLabUser     `json:"user"`
	Project           GitLabProject  `json:"project"`
	Commits           []GitLabCommit `json:"commits"`
	TotalCommitsCount int            `json:"total_commits_count"`
}

// GitLabUser GitLab用户
type GitLabUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GitLabProject GitLab项目
type GitLabProject struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebURL            string `json:"web_url"`
	GitHTTPURL        string `json:"git_http_url"`
}

// GitLabCommit GitLab提交
type GitLabCommit struct {
	ID        string       `json:"id"`
	Message   string       `json:"message"`
	Timestamp string       `json:"timestamp"`
	Author    GitLabAuthor `json:"author"`
	URL       string       `json:"url"`
	Added     []string     `json:"added"`
	Modified  []string     `json:"modified"`
	Removed   []string     `json:"removed"`
}

// GitLabAuthor GitLab作者
type GitLabAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Repository 仓库信息
type Repository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
	CloneURL string `json:"clone_url"`
}

// Pusher 推送者
type Pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Sender 发送者
type Sender struct {
	Login string `json:"login"`
}

// Commit 提交信息
type Commit struct {
	ID        string   `json:"id"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	Author    Author   `json:"author"`
	Added     []string `json:"added"`
	Modified  []string `json:"modified"`
	Removed   []string `json:"removed"`
}

// Author 作者信息
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

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
	eventType := c.GetHeader("X-GitHub-Event")
	logger.Info("Received GitHub webhook", map[string]interface{}{
		"event":     eventType,
		"repo_id":   repo.ID,
		"repo_name": repo.Name,
	})

	// 处理不同事件
	switch eventType {
	case "push":
		s.handlePushEvent(repo, body)
	case "ping":
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
		return
	default:
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

// HandleGitLabWebhook 处理GitLab Webhook
func (s *WebhookService) HandleGitLabWebhook(c *gin.Context) {
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
	eventType := c.GetHeader("X-Gitlab-Event")
	logger.Info("Received GitLab webhook", map[string]interface{}{
		"event":     eventType,
		"repo_id":   repo.ID,
		"repo_name": repo.Name,
	})

	// 处理不同事件
	switch eventType {
	case "Push Hook":
		s.handleGitLabPushEvent(repo, body)
	case "Event Hook":
		// GitLab ping event
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
		return
	default:
		logger.Info("Unsupported GitLab event type", map[string]interface{}{
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

// handlePushEvent 处理push事件
func (s *WebhookService) handlePushEvent(repo *models.Repo, body []byte) {
	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		logger.Error("Failed to parse webhook payload", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

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
	template, _ := s.templateRepo.GetByTypeAndScene(models.TemplateTypeDingTalk, models.TemplateSceneCommitNotify)

	// 为每个推送目标发送通知
	for _, target := range targets {
		go s.sendPushNotification(repo, &target, &payload, template)
	}
}

// handleGitLabPushEvent 处理GitLab push事件
func (s *WebhookService) handleGitLabPushEvent(repo *models.Repo, body []byte) {
	var payload GitLabPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		logger.Error("Failed to parse GitLab webhook payload", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

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
	template, _ := s.templateRepo.GetByTypeAndScene(models.TemplateTypeDingTalk, models.TemplateSceneCommitNotify)

	// 为每个推送目标发送通知
	for _, target := range targets {
		go s.sendGitLabPushNotification(repo, &target, &payload, template)
	}
}

// sendPushNotification 发送推送通知
func (s *WebhookService) sendPushNotification(repo *models.Repo, target *models.Target, payload *WebhookPayload, template *models.Template) {
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
	content := s.buildMessageContent(repo, payload)

	// 创建推送记录
	push := &models.Push{
		RepoID:    repo.ID,
		TargetID:  target.ID,
		CommitID:  payload.After,
		CommitMsg: payload.HeadCommit.Message,
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
func (s *WebhookService) doCodeReview(repo *models.Repo, payload *WebhookPayload, push *models.Push) {
	logger.Info("Starting code review", map[string]interface{}{
		"repo_id":   repo.ID,
		"commit_id": payload.After,
	})

	// 更新状态为进行中
	push.CodeviewStatus = models.CodeviewStatusPending
	s.pushRepo.Update(push)

	// 根据仓库类型创建Git客户端
	var gitClient *git.Client
	if repo.Type == models.RepoTypeGitHub {
		gitClient = git.NewClient(
			repo.URL,
			repo.AccessToken,
			extractOwner(repo.URL),
			extractRepoName(repo.URL),
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
		})
		push.CodeviewStatus = models.CodeviewStatusSkipped
		s.pushRepo.Update(push)
		return
	}

	// 获取分支
	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

	// 批量审查
	var allIssues strings.Builder
	for _, file := range codeFiles {
		input := CodeViewInput{
			FileName:    file.Filename,
			DiffContent: file.Patch,
			RepoName:    repo.Name,
			Branch:      branch,
			CommitMsg:   payload.HeadCommit.Message,
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

// buildMessageContent 构建消息内容
func (s *WebhookService) buildMessageContent(repo *models.Repo, payload *WebhookPayload) string {
	var content strings.Builder

	// 统计变更文件
	allFiles := append(append(payload.HeadCommit.Added, payload.HeadCommit.Modified...), payload.HeadCommit.Removed...)
	fileCount := len(allFiles)

	content.WriteString("## 代码提交通知\n\n")
	content.WriteString("**仓库**: " + payload.Repository.FullName + "\n")
	content.WriteString("**分支**: " + strings.TrimPrefix(payload.Ref, "refs/heads/") + "\n")
	content.WriteString("**提交**: " + payload.HeadCommit.ID[:7] + "\n")
	content.WriteString("**信息**: " + payload.HeadCommit.Message + "\n")
	content.WriteString("**作者**: " + payload.HeadCommit.Author.Name + "\n")
	content.WriteString("**文件数**: " + string(rune(fileCount)) + "\n\n")

	if len(allFiles) > 0 {
		content.WriteString("### 变更文件\n")
		for i, file := range allFiles {
			if i >= 10 {
				content.WriteString("- ... 还有 " + string(rune(len(allFiles)-10)) + " 个文件\n")
				break
			}
			content.WriteString("- " + file + "\n")
		}
	}

	return content.String()
}

// buildGitLabMessageContent 构建GitLab消息内容
func (s *WebhookService) buildGitLabMessageContent(repo *models.Repo, payload *GitLabPayload) string {
	var content strings.Builder

	// 统计变更文件
	var allFiles []string
	for _, commit := range payload.Commits {
		allFiles = append(allFiles, commit.Added...)
		allFiles = append(allFiles, commit.Modified...)
		allFiles = append(allFiles, commit.Removed...)
	}
	fileCount := len(allFiles)

	content.WriteString("## GitLab 代码提交通知\n\n")
	content.WriteString("**项目**: " + payload.Project.PathWithNamespace + "\n")
	content.WriteString("**分支**: " + strings.TrimPrefix(payload.Ref, "refs/heads/") + "\n")
	content.WriteString("**提交数**: " + fmt.Sprintf("%d", payload.TotalCommitsCount) + "\n")
	content.WriteString("**提交者**: " + payload.User.Name + "\n\n")

	if len(payload.Commits) > 0 {
		content.WriteString("### 提交记录\n")
		for i, commit := range payload.Commits {
			if i >= 5 {
				content.WriteString("- ... 还有 " + fmt.Sprintf("%d", len(payload.Commits)-5) + " 个提交\n")
				break
			}
			shortID := commit.ID[:7]
			content.WriteString(fmt.Sprintf("- `%s` %s\n", shortID, commit.Message))
		}
	}

	if fileCount > 0 {
		content.WriteString("\n### 变更文件\n")
		uniqueFiles := removeDuplicates(allFiles)
		for i, file := range uniqueFiles {
			if i >= 10 {
				content.WriteString("- ... 还有 " + fmt.Sprintf("%d", len(uniqueFiles)-10) + " 个文件\n")
				break
			}
			content.WriteString("- " + file + "\n")
		}
	}

	return content.String()
}

// removeDuplicates 移除切片中的重复元素
func removeDuplicates(slice []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, s := range slice {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
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

// sendGitLabPushNotification 发送GitLab推送通知
func (s *WebhookService) sendGitLabPushNotification(repo *models.Repo, target *models.Target, payload *GitLabPayload, template *models.Template) {
	// 去重检查：同一提交同一目标不重复推送
	if s.pushRepo.ExistsByCommitAndTarget(payload.After, target.ID) {
		logger.Info("Duplicate GitLab push detected, skipping", map[string]interface{}{
			"commit_id": payload.After[:7],
			"target_id": target.ID,
			"repo_name": repo.Name,
		})
		return
	}

	// 构建消息内容
	content := s.buildGitLabMessageContent(repo, payload)

	// 创建推送记录
	push := &models.Push{
		RepoID:    repo.ID,
		TargetID:  target.ID,
		CommitID:  payload.After,
		CommitMsg: payload.Commits[0].Message,
		Status:    models.PushStatusPending,
		Content:   content,
	}

	if template != nil {
		push.TemplateID = &template.ID
	}

	if err := s.pushRepo.Create(push); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "Duplicate entry") {
			logger.Info("Duplicate GitLab push detected (DB constraint), skipping", map[string]interface{}{
				"commit_id": payload.After[:7],
				"target_id": target.ID,
				"repo_name": repo.Name,
			})
			return
		}
		logger.Error("Failed to create GitLab push record", map[string]interface{}{
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
		logger.Error("GitLab push failed", map[string]interface{}{
			"push_id":   push.ID,
			"target_id": target.ID,
			"error":     err.Error(),
		})
	} else {
		push.Status = models.PushStatusSuccess
		logger.Info("GitLab push succeeded", map[string]interface{}{
			"push_id":   push.ID,
			"target_id": target.ID,
		})
	}

	s.pushRepo.Update(push)
}
