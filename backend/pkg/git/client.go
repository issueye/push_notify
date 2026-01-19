package git

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client Git客户端
type Client struct {
	baseURL   string
	token     string
	repoOwner string
	repoName  string
	client    *http.Client
}

// NewClient 创建Git客户端
func NewClient(baseURL, token, repoOwner, repoName string) *Client {
	return &Client{
		baseURL:   strings.TrimSuffix(baseURL, "/"),
		token:     token,
		repoOwner: repoOwner,
		repoName:  repoName,
		client:    &http.Client{},
	}
}

// DiffFile 差异文件
type DiffFile struct {
	Filename    string `json:"filename"`
	Status      string `json:"status"` // added, modified, deleted
	Patch       string `json:"patch"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	Changes     int    `json:"changes"`
	SHA         string `json:"sha"`
	PreviousSHA string `json:"previous_sha"`
}

// CompareResult 比较结果
type CompareResult struct {
	Status       string     `json:"status"`
	AheadBy      int        `json:"ahead_by"`
	BehindBy     int        `json:"behind_by"`
	TotalCommits int        `json:"total_commits"`
	Files        []DiffFile `json:"files"`
}

// GetDiff 获取两次提交之间的差异
func (c *Client) GetDiff(base, head string) ([]DiffFile, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/compare/%s...%s", c.baseURL, c.repoOwner, c.repoName, base, head)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error: %s", string(body))
	}

	var result CompareResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Files, nil
}

// GetSingleCommitDiff 获取单次提交的差异
func (c *Client) GetSingleCommitDiff(commitSHA string) ([]DiffFile, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits/%s", c.baseURL, c.repoOwner, c.repoName, commitSHA)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error: %s", string(body))
	}

	// 解析响应获取parent
	var commitData struct {
		SHA     string `json:"sha"`
		Parents []struct {
			SHA string `json:"sha"`
		} `json:"parents"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&commitData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 如果有parent，获取与parent的差异
	if len(commitData.Parents) > 0 {
		return c.GetDiff(commitData.Parents[0].SHA, commitSHA)
	}

	// 首次提交，返回空差异
	return []DiffFile{}, nil
}

// GetFileContent 获取文件内容
func (c *Client) GetFileContent(filePath, ref string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s", c.baseURL, c.repoOwner, c.repoName, filePath, ref)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3.raw")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("api error: %s", string(body))
	}

	content, _ := io.ReadAll(resp.Body)
	return string(content), nil
}
