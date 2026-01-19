package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Message 钉钉消息结构
type Message struct {
	MsgType  string          `json:"msgtype"`
	Text     TextContent     `json:"text"`
	Markdown MarkdownContent `json:"markdown,omitempty"`
}

// TextContent 文本消息
type TextContent struct {
	Content string `json:"content"`
}

// MarkdownContent Markdown消息
type MarkdownContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// Config 钉钉配置
type Config struct {
	AccessToken string
	Secret      string
}

// Client 钉钉客户端
type Client struct {
	config Config
	client *http.Client
}

// NewClient 创建钉钉客户端
func NewClient(accessToken, secret string) *Client {
	return &Client{
		config: Config{
			AccessToken: accessToken,
			Secret:      secret,
		},
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send 发送消息
func (c *Client) Send(hookUrl string, msg Message) error {
	apiURL := hookUrl

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	fmt.Println("data ->", string(data))

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode != 200 {
		return fmt.Errorf("dingtalk api error: %s", string(body))
	}

	if errCode, ok := result["errcode"].(float64); ok && errCode != 0 {
		return fmt.Errorf("dingtalk error: %v", result["errmsg"])
	}

	return nil
}

// SendText 发送文本消息
func (c *Client) SendText(hookUrl string, content string) error {
	msg := Message{
		MsgType: "text",
		Text: TextContent{
			Content: content,
		},
	}
	return c.Send(hookUrl, msg)
}

// SendMarkdown 发送Markdown消息
func (c *Client) SendMarkdown(hookUrl string, title, content string) error {
	msg := Message{
		MsgType: "markdown",
		Markdown: MarkdownContent{
			Title: title,
			Text:  content,
		},
	}
	return c.Send(hookUrl, msg)
}

// sign 生成签名
func (c *Client) sign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
