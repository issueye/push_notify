package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrTokenExpired  = errors.New("token expired")
	ErrInvalidClaims = errors.New("invalid claims")
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type JWT struct {
	SecretKey          string
	AccessTokenExpire  int // 秒
	RefreshTokenExpire int // 秒
}

func NewJWT(secretKey string, accessExpire, refreshExpire int) *JWT {
	return &JWT{
		SecretKey:          secretKey,
		AccessTokenExpire:  accessExpire,
		RefreshTokenExpire: refreshExpire,
	}
}

// GenerateToken 生成访问令牌和刷新令牌
func (j *JWT) GenerateToken(userID uint, username, role string) (accessToken, refreshToken string, err error) {
	now := time.Now()

	// 生成访问令牌
	accessClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.AccessTokenExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "push-notify",
			Subject:   username,
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.RefreshTokenExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "push-notify",
			Subject:   username,
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken 验证令牌并返回Claims
func (j *JWT) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidClaims
}

// RefreshToken 使用刷新令牌获取新的访问令牌
func (j *JWT) RefreshToken(refreshToken string) (string, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	now := time.Now()
	newAccessClaims := JWTClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.AccessTokenExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "push-notify",
			Subject:   claims.Username,
		},
	}

	newTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims)
	return newTokenObj.SignedString([]byte(j.SecretKey))
}
