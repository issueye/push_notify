package git

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GitLabCommit GitLab提交
type GitLabCommit struct {
	ID      string `json:"id"`
	ShortID string `json:"short_id"`
	Message string `json:"message"`
	Author  struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
	ParentIDs []string `json:"parent_ids"`
}

// GitLabDiff 差异
type GitLabDiff struct {
	OldPath     string `json:"old_path"`
	NewPath     string `json:"new_path"`
	NewFile     string `json:"new_file"`
	RenamedFile string `json:"renamed_file"`
	DeletedFile string `json:"deleted_file"`
	ModeChange  string `json:"mode_change"`
	Diff        string `json:"diff"`
}

// GitLabCompare 比较结果
type GitLabCompare struct {
	Commit         *GitLabCommit `json:"commit"`
	Diffs          []GitLabDiff  `json:"diffs"`
	CompareTimeout bool          `json:"compare_timeout"`
	CompareCache   bool          `json:"compare_cache"`
}

// GitLabClient GitLab客户端
type GitLabClient struct {
	baseURL   string
	token     string
	projectID string
	client    *http.Client
}

// NewGitLabClient 创建GitLab客户端
func NewGitLabClient(baseURL, token, projectID string) *GitLabClient {
	return &GitLabClient{
		baseURL:   strings.TrimSuffix(baseURL, "/"),
		token:     token,
		projectID: projectID,
		client:    &http.Client{},
	}
}

// GetDiff 获取两次提交之间的差异
func (c *GitLabClient) GetDiff(from, to string) ([]DiffFile, error) {
	// URL编码commit SHA
	fromEncoded := url.PathEscape(from)
	toEncoded := url.PathEscape(to)

	url := fmt.Sprintf("%s/projects/%s/repository/compare?from=%s&to=%s", c.baseURL, c.projectID, fromEncoded, toEncoded)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("PRIVATE-TOKEN", c.token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error: %s", string(body))
	}

	var result GitLabCompare
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 转换GitLab差异为通用格式
	files := make([]DiffFile, 0, len(result.Diffs))
	for _, d := range result.Diffs {
		status := "modified"
		if d.NewFile == "true" {
			status = "added"
		} else if d.DeletedFile == "true" {
			status = "deleted"
		}

		files = append(files, DiffFile{
			Filename: d.NewPath,
			Status:   status,
			Patch:    d.Diff,
		})
	}

	return files, nil
}

// GetSingleCommitDiff 获取单次提交的差异
func (c *GitLabClient) GetSingleCommitDiff(commitSHA string) ([]DiffFile, error) {
	url := fmt.Sprintf("%s/projects/%s/repository/commits/%s/diff", c.baseURL, c.projectID, commitSHA)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("PRIVATE-TOKEN", c.token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error: %s", string(body))
	}

	var diffs []GitLabDiff
	if err := json.NewDecoder(resp.Body).Decode(&diffs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 转换GitLab差异为通用格式
	files := make([]DiffFile, 0, len(diffs))
	for _, d := range diffs {
		status := "modified"
		if d.NewFile == "true" {
			status = "added"
		} else if d.DeletedFile == "true" {
			status = "deleted"
		}

		files = append(files, DiffFile{
			Filename: d.NewPath,
			Status:   status,
			Patch:    d.Diff,
		})
	}

	return files, nil
}
