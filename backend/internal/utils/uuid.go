package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateUUID 生成UUID
func GenerateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateWebhookSecret 生成Webhook密钥
func GenerateWebhookSecret() (string, error) {
	return GenerateRandomString(32)
}
