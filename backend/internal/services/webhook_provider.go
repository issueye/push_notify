package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"backend/internal/models"
)

// UnifiedPushPayload 统一的推送负载
type UnifiedPushPayload struct {
	Ref          string
	After        string // CommitID
	Before       string
	CommitMsg    string
	AuthorName   string
	RepoName     string // FullName or PathWithNamespace
	Branch       string
	FileCount    int
	FileList     []string        // Simple list of changed files
	Commits      []UnifiedCommit // For listing in message
	TotalCommits int
}

// UnifiedCommit 统一的提交信息
type UnifiedCommit struct {
	ID      string
	Message string
	Author  string
}

// WebhookProvider Webhook提供者接口
type WebhookProvider interface {
	GetEventType(header http.Header) string
	ParsePushPayload(body []byte) (*UnifiedPushPayload, error)
	BuildMessage(payload *UnifiedPushPayload, template *models.Template) string
}

// GitHubProvider GitHub实现
type GitHubProvider struct{}

func (p *GitHubProvider) GetEventType(header http.Header) string {
	return header.Get("X-GitHub-Event")
}

func (p *GitHubProvider) ParsePushPayload(body []byte) (*UnifiedPushPayload, error) {
	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	// 统计文件
	allFiles := append(append(payload.HeadCommit.Added, payload.HeadCommit.Modified...), payload.HeadCommit.Removed...)
	
	// 转换 Commits
	var commits []UnifiedCommit
	for _, c := range payload.Commits {
		commits = append(commits, UnifiedCommit{
			ID:      c.ID,
			Message: c.Message,
			Author:  c.Author.Name,
		})
	}

	return &UnifiedPushPayload{
		Ref:          payload.Ref,
		After:        payload.After,
		Before:       payload.Before,
		CommitMsg:    payload.HeadCommit.Message,
		AuthorName:   payload.HeadCommit.Author.Name,
		RepoName:     payload.Repository.FullName,
		Branch:       strings.TrimPrefix(payload.Ref, "refs/heads/"),
		FileCount:    len(allFiles),
		FileList:     allFiles,
		Commits:      commits,
		TotalCommits: len(payload.Commits), // GitHub payload usually contains new commits
	}, nil
}

func (p *GitHubProvider) BuildMessage(payload *UnifiedPushPayload, template *models.Template) string {
	// 使用模板
	if template != nil && template.Content != "" {
		return applyTemplate(payload, template)
	}

	// 默认格式
	var content strings.Builder
	content.WriteString("## 代码提交通知\n\n")
	content.WriteString("**仓库**: " + payload.RepoName + "\n")
	content.WriteString("**分支**: " + payload.Branch + "\n")
	content.WriteString("**提交**: " + payload.After[:7] + "\n")
	content.WriteString("**信息**: " + payload.CommitMsg + "\n")
	content.WriteString("**作者**: " + payload.AuthorName + "\n")
	content.WriteString("**文件数**: " + fmt.Sprintf("%d", payload.FileCount) + "\n\n")

	if len(payload.FileList) > 0 {
		content.WriteString("### 变更文件\n")
		for i, file := range payload.FileList {
			if i >= 10 {
				content.WriteString("- ... 还有 " + fmt.Sprintf("%d", len(payload.FileList)-10) + " 个文件\n")
				break
			}
			content.WriteString("- " + file + "\n")
		}
	}

	return content.String()
}

// GitLabProvider GitLab实现
type GitLabProvider struct{}

func (p *GitLabProvider) GetEventType(header http.Header) string {
	return header.Get("X-Gitlab-Event")
}

func (p *GitLabProvider) ParsePushPayload(body []byte) (*UnifiedPushPayload, error) {
	var payload GitLabPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	// 统计文件
	var allFiles []string
	var commitMsg string
	var authorName string
	
	var commits []UnifiedCommit

	for _, commit := range payload.Commits {
		allFiles = append(allFiles, commit.Added...)
		allFiles = append(allFiles, commit.Modified...)
		allFiles = append(allFiles, commit.Removed...)
		
		commits = append(commits, UnifiedCommit{
			ID:      commit.ID,
			Message: commit.Message,
			Author:  commit.Author.Name,
		})

		if commit.ID == payload.After {
			commitMsg = commit.Message
			authorName = commit.Author.Name
		}
	}
	
	// 如果没找到 HEAD commit
	if commitMsg == "" && len(payload.Commits) > 0 {
		commitMsg = payload.Commits[0].Message
		authorName = payload.Commits[0].Author.Name
	} else if commitMsg == "" {
		commitMsg = "Unknown commit"
		authorName = payload.User.Name
	}

	uniqueFiles := removeDuplicates(allFiles)

	return &UnifiedPushPayload{
		Ref:          payload.Ref,
		After:        payload.After,
		Before:       payload.Before,
		CommitMsg:    commitMsg,
		AuthorName:   authorName,
		RepoName:     payload.Project.PathWithNamespace,
		Branch:       strings.TrimPrefix(payload.Ref, "refs/heads/"),
		FileCount:    len(allFiles), // Total changes
		FileList:     uniqueFiles,   // Unique files for display
		Commits:      commits,
		TotalCommits: payload.TotalCommitsCount,
	}, nil
}

func (p *GitLabProvider) BuildMessage(payload *UnifiedPushPayload, template *models.Template) string {
	// 使用模板
	if template != nil && template.Content != "" {
		return applyTemplate(payload, template)
	}

	// 默认格式
	var content strings.Builder
	content.WriteString("## GitLab 代码提交通知\n\n")
	content.WriteString("**项目**: " + payload.RepoName + "\n")
	content.WriteString("**分支**: " + payload.Branch + "\n")
	content.WriteString("**提交数**: " + fmt.Sprintf("%d", payload.TotalCommits) + "\n")
	content.WriteString("**提交者**: " + payload.AuthorName + "\n\n")

	if len(payload.Commits) > 0 {
		content.WriteString("### 提交记录\n")
		for i, commit := range payload.Commits {
			if i >= 5 {
				content.WriteString("- ... 还有 " + fmt.Sprintf("%d", len(payload.Commits)-5) + " 个提交\n")
				break
			}
			shortID := commit.ID
			if len(shortID) > 7 {
				shortID = shortID[:7]
			}
			content.WriteString(fmt.Sprintf("- `%s` %s\n", shortID, commit.Message))
		}
	}

	if len(payload.FileList) > 0 {
		content.WriteString("\n### 变更文件\n")
		for i, file := range payload.FileList {
			if i >= 10 {
				content.WriteString("- ... 还有 " + fmt.Sprintf("%d", len(payload.FileList)-10) + " 个文件\n")
				break
			}
			content.WriteString("- " + file + "\n")
		}
	}

	return content.String()
}

// 辅助函数：应用模板
func applyTemplate(payload *UnifiedPushPayload, template *models.Template) string {
	content := template.Content
	content = strings.ReplaceAll(content, "{{.RepoName}}", payload.RepoName)
	content = strings.ReplaceAll(content, "{{.Branch}}", payload.Branch)
	content = strings.ReplaceAll(content, "{{.CommitID}}", payload.After)
	content = strings.ReplaceAll(content, "{{.CommitMsg}}", payload.CommitMsg)
	content = strings.ReplaceAll(content, "{{.Author}}", payload.AuthorName)
	content = strings.ReplaceAll(content, "{{.FileCount}}", fmt.Sprintf("%d", payload.FileCount))

	// 构建文件列表
	var fileListBuilder strings.Builder
	for i, file := range payload.FileList {
		if i >= 10 {
			fileListBuilder.WriteString(fmt.Sprintf("- ... 还有 %d 个文件\n", len(payload.FileList)-10))
			break
		}
		fileListBuilder.WriteString("- " + file + "\n")
	}
	content = strings.ReplaceAll(content, "{{.FileList}}", fileListBuilder.String())

	return content
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
