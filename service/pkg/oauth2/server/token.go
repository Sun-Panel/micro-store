package server

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenManager Token 管理器
type TokenManager struct {
	config *OAuthConfig
}

// NewTokenManager 创建 Token 管理器
func NewTokenManager(config *OAuthConfig) *TokenManager {
	if config == nil {
		config = DefaultOAuthConfig()
	}
	return &TokenManager{config: config}
}

// GenerateAccessToken 生成 Access Token
func (tm *TokenManager) GenerateAccessToken(clientID string, userID uint, clientSecret string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(tm.config.AccessTokenExpireTime) * time.Second)

	claims := &OAuthClaims{
		UserID:   userID,
		ClientID: clientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken 生成 Refresh Token
func (tm *TokenManager) GenerateRefreshToken() string {
	// 简单实现：生成随机字符串
	// 实际生产环境中应该使用更安全的方式
	return fmt.Sprintf("refresh_%d_%d", time.Now().UnixNano(), tm.config.RefreshTokenExpireTime)
}

// GenerateAuthCode 生成授权码
func (tm *TokenManager) GenerateAuthCode(clientID string, userID uint, clientSecret string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(tm.config.AuthCodeExpireTime) * time.Second)

	claims := &AuthCodeClaims{
		UserID:   userID,
		ClientID: clientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign auth code: %w", err)
	}

	return tokenString, nil
}

// ValidateAccessToken 验证 Access Token
func (tm *TokenManager) ValidateAccessToken(tokenString, clientSecret string) (*OAuthClaims, error) {
	claims := &OAuthClaims{}
	err := tm.validateJWT(tokenString, clientSecret, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// ValidateAuthCode 验证授权码
func (tm *TokenManager) ValidateAuthCode(tokenString, clientSecret string) (*AuthCodeClaims, error) {
	claims := &AuthCodeClaims{}
	err := tm.validateJWT(tokenString, clientSecret, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// validateJWT 验证 JWT Token
func (tm *TokenManager) validateJWT(tokenString, clientSecret string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(clientSecret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return fmt.Errorf("malformed token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return fmt.Errorf("expired token")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return fmt.Errorf("token not valid yet")
			}
		}
		return fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// GetExpireTime 获取过期时间（秒）
func (tm *TokenManager) GetExpireTime() int {
	return tm.config.AccessTokenExpireTime
}

// GetRefreshTokenExpireTime 获取刷新令牌过期时间（秒）
func (tm *TokenManager) GetRefreshTokenExpireTime() int {
	return tm.config.RefreshTokenExpireTime
}
