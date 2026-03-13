package server

import (
	"testing"
	"time"
)

func TestTokenManager_GenerateAccessToken(t *testing.T) {
	config := DefaultOAuthConfig()
	tm := NewTokenManager(config)

	clientID := "test_client"
	userID := uint(1)
	clientSecret := "test_secret_key"

	token, err := tm.GenerateAccessToken(clientID, userID, clientSecret)
	if err != nil {
		t.Fatalf("Failed to generate access token: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}

	t.Logf("Generated token: %s", token)
}

func TestTokenManager_ValidateAccessToken(t *testing.T) {
	config := DefaultOAuthConfig()
	tm := NewTokenManager(config)

	clientID := "test_client"
	userID := uint(1)
	clientSecret := "test_secret_key"

	// 生成 Token
	token, err := tm.GenerateAccessToken(clientID, userID, clientSecret)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 验证 Token
	claims, err := tm.ValidateAccessToken(token, clientSecret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.ClientID != clientID {
		t.Errorf("Expected ClientID %s, got %s", clientID, claims.ClientID)
	}

	t.Logf("Token claims: %+v", claims)
}

func TestTokenManager_ValidateAccessToken_WrongSecret(t *testing.T) {
	config := DefaultOAuthConfig()
	tm := NewTokenManager(config)

	clientID := "test_client"
	userID := uint(1)
	clientSecret := "test_secret_key"

	// 生成 Token
	token, _ := tm.GenerateAccessToken(clientID, userID, clientSecret)

	// 使用错误的密钥验证
	_, err := tm.ValidateAccessToken(token, "wrong_secret")
	if err == nil {
		t.Error("Expected error for wrong secret, but got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestTokenManager_GenerateAuthCode(t *testing.T) {
	config := DefaultOAuthConfig()
	tm := NewTokenManager(config)

	clientID := "test_client"
	userID := uint(1)
	clientSecret := "test_secret_key"

	code, err := tm.GenerateAuthCode(clientID, userID, clientSecret)
	if err != nil {
		t.Fatalf("Failed to generate auth code: %v", err)
	}

	if code == "" {
		t.Error("Code should not be empty")
	}

	t.Logf("Generated auth code: %s", code)
}

func TestTokenManager_ValidateAuthCode(t *testing.T) {
	config := DefaultOAuthConfig()
	tm := NewTokenManager(config)

	clientID := "test_client"
	userID := uint(1)
	clientSecret := "test_secret_key"

	// 生成授权码
	code, err := tm.GenerateAuthCode(clientID, userID, clientSecret)
	if err != nil {
		t.Fatalf("Failed to generate auth code: %v", err)
	}

	// 验证授权码
	claims, err := tm.ValidateAuthCode(code, clientSecret)
	if err != nil {
		t.Fatalf("Failed to validate auth code: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.ClientID != clientID {
		t.Errorf("Expected ClientID %s, got %s", clientID, claims.ClientID)
	}

	t.Logf("Auth code claims: %+v", claims)
}

func TestTokenManager_GenerateRefreshToken(t *testing.T) {
	config := DefaultOAuthConfig()
	tm := NewTokenManager(config)

	refreshToken := tm.GenerateRefreshToken()

	if refreshToken == "" {
		t.Error("Refresh token should not be empty")
	}

	t.Logf("Generated refresh token: %s", refreshToken)
}

func TestTokenManager_GetExpireTime(t *testing.T) {
	config := &OAuthConfig{
		AccessTokenExpireTime:  7200,
		RefreshTokenExpireTime: 604800,
		AuthCodeExpireTime:     600,
	}
	tm := NewTokenManager(config)

	expireTime := tm.GetExpireTime()
	if expireTime != 7200 {
		t.Errorf("Expected expire time 7200, got %d", expireTime)
	}

	refreshExpireTime := tm.GetRefreshTokenExpireTime()
	if refreshExpireTime != 604800 {
		t.Errorf("Expected refresh expire time 604800, got %d", refreshExpireTime)
	}
}

func TestDefaultOAuthConfig(t *testing.T) {
	config := DefaultOAuthConfig()

	if config.AccessTokenExpireTime != 7200 {
		t.Errorf("Expected default access token expire time 7200, got %d", config.AccessTokenExpireTime)
	}

	if config.RefreshTokenExpireTime != 604800 {
		t.Errorf("Expected default refresh token expire time 604800, got %d", config.RefreshTokenExpireTime)
	}

	if config.AuthCodeExpireTime != 600 {
		t.Errorf("Expected default auth code expire time 600, got %d", config.AuthCodeExpireTime)
	}

	if !config.EnableSSOLogout {
		t.Error("Expected default enable SSO logout to be true")
	}

	t.Logf("Default config: %+v", config)
}

func TestTokenManager_TokenExpiration(t *testing.T) {
	// 创建一个极短过期时间的配置
	config := &OAuthConfig{
		AccessTokenExpireTime:  1, // 1 秒
		RefreshTokenExpireTime: 1,
		AuthCodeExpireTime:     1,
	}
	tm := NewTokenManager(config)

	clientID := "test_client"
	userID := uint(1)
	clientSecret := "test_secret_key"

	// 生成 Token
	token, _ := tm.GenerateAccessToken(clientID, userID, clientSecret)

	// 立即验证应该成功
	_, err := tm.ValidateAccessToken(token, clientSecret)
	if err != nil {
		t.Errorf("Token should be valid immediately: %v", err)
	}

	// 等待过期
	time.Sleep(2 * time.Second)

	// 过期后验证应该失败
	_, err = tm.ValidateAccessToken(token, clientSecret)
	if err == nil {
		t.Error("Expected error for expired token")
	} else {
		t.Logf("Expected expiration error: %v", err)
	}
}
