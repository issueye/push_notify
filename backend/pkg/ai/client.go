package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Config AI模型配置
type Config struct {
	APIURL  string
	APIKey  string
	Timeout int // 秒
	Params  map[string]interface{}
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request 请求结构
type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

// Response 响应结构
type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 选择
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Client AI客户端
type Client struct {
	config Config
	client *http.Client
}

// NewClient 创建AI客户端
func NewClient(apiURL, apiKey string, timeout int) *Client {
	if timeout <= 0 {
		timeout = 60
	}
	return &Client{
		config: Config{
			APIURL:  apiURL,
			APIKey:  apiKey,
			Timeout: timeout,
		},
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// NewClientWithConfig 使用配置创建客户端
func NewClientWithConfig(config Config) *Client {
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 60
	}
	return &Client{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// Chat 发送聊天请求
func (c *Client) Chat(messages []Message, systemPrompt string) (string, error) {
	// 构建消息
	var allMessages []Message
	if systemPrompt != "" {
		allMessages = append(allMessages, Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}
	allMessages = append(allMessages, messages...)

	// 构建请求
	req := Request{
		Model:    "gpt-4", // 可配置
		Messages: allMessages,
	}

	// 应用参数
	if temp, ok := c.config.Params["temperature"].(float64); ok {
		req.Temperature = temp
	}
	if maxTokens, ok := c.config.Params["max_tokens"].(float64); ok {
		req.MaxTokens = int(maxTokens)
	}
	if topP, ok := c.config.Params["top_p"].(float64); ok {
		req.TopP = topP
	}

	// 发送请求
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.config.APIURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("api error: %s", string(respBody))
	}

	var response Response
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}

// CodeReview 进行代码审查
func (c *Client) CodeReview(prompt, code string) (string, error) {
	messages := []Message{
		{
			Role:    "user",
			Content: prompt + "\n\n代码内容:\n" + code,
		},
	}
	fmt.Println("code", code)
	return c.Chat(messages, "")
}

// ValidateConfig 验证配置
func (c *Client) ValidateConfig() error {
	// 发送一个简单的请求验证配置
	messages := []Message{
		{
			Role:    "user",
			Content: "Hello",
		},
	}

	_, err := c.Chat(messages, "")
	return err
}

// ParseResponse 解析流式响应（简化版）
func ParseResponse(responseBody []byte) (string, error) {
	var response Response
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", err
	}
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}
	return "", nil
}

// IsValidURL 验证URL格式
func IsValidURL(apiURL string) bool {
	parsed, err := url.Parse(apiURL)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}

// NormalizeURL 标准化API URL
func NormalizeURL(apiURL string) string {
	apiURL = strings.TrimSpace(apiURL)
	if !strings.HasPrefix(apiURL, "http://") && !strings.HasPrefix(apiURL, "https://") {
		apiURL = "https://" + apiURL
	}
	return apiURL
}
