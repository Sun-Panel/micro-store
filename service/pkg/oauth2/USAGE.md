# OAuth2 模块使用指南

本文档提供 OAuth2 模块的详细使用说明和最佳实践。

## 目录

1. [快速开始](#快速开始)
2. [服务端实现](#服务端实现)
3. [客户端实现](#客户端实现)
4. [存储实现](#存储实现)
5. [安全最佳实践](#安全最佳实践)
6. [常见问题](#常见问题)

## 快速开始

### 安装

```bash
# 将模块添加到你的项目中
import "sun-panel-micro-store/pkg/oauth2"
```

### 最小服务端示例

```go
package main

import (
    "github.com/gin-gonic/gin"
    "sun-panel-micro-store/pkg/oauth2/server"
)

func main() {
    // 创建配置
    config := server.DefaultOAuthConfig()
    
    // 创建 handler
    handler := server.NewOAuthHandler(config)
    
    // 设置存储（需要自己实现）
    handler.SetStores(
        &MyThirdAppStore{},
        &MyUserStore{},
        &MyTokenStore{},
        &MyAuthCodeStore{},
        &MyRefreshTokenStore{},
    )
    
    // 注册路由
    r := gin.Default()
    handler.RegisterRoutes(r.Group("/api"))
    
    r.Run(":8080")
}
```

### 最小客户端示例

```go
package main

import (
    "context"
    "fmt"
    "sun-panel-micro-store/pkg/oauth2/client"
)

func main() {
    config := &client.Config{
        AuthServerURL: "http://localhost:8080",
        ClientID:      "your_client_id",
        ClientSecret:  "your_client_secret",
    }
    
    oauthClient, _ := client.NewOAuth2Client(config)
    
    // 获取客户端凭证 Token
    token, _ := oauthClient.GetClientCredentialsToken(context.Background())
    fmt.Println("Access Token:", token.AccessToken)
}
```

## 服务端实现

### 1. 存储接口实现

#### Redis 实现（推荐生产环境使用）

```go
package storage

import (
    "context"
    "encoding/json"
    "time"
    "github.com/go-redis/redis/v8"
    "sun-panel-micro-store/pkg/oauth2/server"
)

type RedisTokenStore struct {
    client *redis.Client
}

func NewRedisTokenStore(addr string) *RedisTokenStore {
    return &RedisTokenStore{
        client: redis.NewClient(&redis.Options{
            Addr: addr,
        }),
    }
}

func (s *RedisTokenStore) SetAccessToken(token string, data server.AccessTokenData) error {
    ctx := context.Background()
    jsonData, _ := json.Marshal(data)
    // 设置过期时间为 2 小时
    return s.client.Set(ctx, "token:"+token, jsonData, 2*time.Hour).Err()
}

func (s *RedisTokenStore) GetAccessToken(token string) (server.AccessTokenData, error) {
    ctx := context.Background()
    data, err := s.client.Get(ctx, "token:"+token).Bytes()
    if err != nil {
        return server.AccessTokenData{}, err
    }
    var tokenData server.AccessTokenData
    json.Unmarshal(data, &tokenData)
    return tokenData, nil
}

func (s *RedisTokenStore) DeleteAccessToken(token string) error {
    ctx := context.Background()
    return s.client.Del(ctx, "token:"+token).Err()
}
```

#### 数据库实现

```go
package storage

import (
    "gorm.io/gorm"
    "sun-panel-micro-store/pkg/oauth2/server"
)

type DBTokenStore struct {
    db *gorm.DB
}

type TokenModel struct {
    Token    string `gorm:"primaryKey"`
    UserID   uint
    ClientID string
    CToken   string
    OpenID   string
    ExpiredAt time.Time
}

func (TokenModel) TableName() string {
    return "oauth_tokens"
}

func NewDBTokenStore(db *gorm.DB) *DBTokenStore {
    return &DBTokenStore{db: db}
}

func (s *DBTokenStore) SetAccessToken(token string, data server.AccessTokenData) error {
    model := TokenModel{
        Token:     token,
        UserID:    data.UserID,
        ClientID:  data.ClientID,
        CToken:    data.CToken,
        OpenID:    data.OpenID,
        ExpiredAt: time.Now().Add(2 * time.Hour),
    }
    return s.db.Create(&model).Error
}

func (s *DBTokenStore) GetAccessToken(token string) (server.AccessTokenData, error) {
    var model TokenModel
    err := s.db.Where("token = ? AND expired_at > ?", token, time.Now()).First(&model).Error
    if err != nil {
        return server.AccessTokenData{}, err
    }
    return server.AccessTokenData{
        UserID:   model.UserID,
        ClientID: model.ClientID,
        CToken:   model.CToken,
        OpenID:   model.OpenID,
    }, nil
}

func (s *DBTokenStore) DeleteAccessToken(token string) error {
    return s.db.Where("token = ?", token).Delete(&TokenModel{}).Error
}
```

### 2. 用户认证中间件

```go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenStore server.TokenStore) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取 Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        
        // 解析 Bearer token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
            c.Abort()
            return
        }
        
        accessToken := parts[1]
        
        // 验证 token
        tokenData, err := tokenStore.GetAccessToken(accessToken)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }
        
        // 注入用户信息到上下文
        c.Set("user_id", tokenData.UserID)
        c.Set("client_id", tokenData.ClientID)
        
        c.Next()
    }
}
```

## 客户端实现

### 1. Token 管理

```go
package auth

import (
    "context"
    "sync"
    "time"
    "sun-panel-micro-store/pkg/oauth2/client"
)

type TokenManager struct {
    client       *client.OAuth2Client
    currentToken *client.TokenResponse
    mu           sync.RWMutex
}

func NewTokenManager(config *client.Config) (*TokenManager, error) {
    oauthClient, err := client.NewOAuth2Client(config)
    if err != nil {
        return nil, err
    }
    
    return &TokenManager{
        client: oauthClient,
    }, nil
}

// GetValidToken 获取有效的 Token（自动刷新）
func (tm *TokenManager) GetValidToken(ctx context.Context) (string, error) {
    tm.mu.RLock()
    
    // 如果当前 Token 有效，直接返回
    if tm.currentToken != nil && tm.isValid(tm.currentToken) {
        token := tm.currentToken.AccessToken
        tm.mu.RUnlock()
        return token, nil
    }
    
    tm.mu.RUnlock()
    
    // Token 无效或过期，获取新的 Token
    tm.mu.Lock()
    defer tm.mu.Unlock()
    
    // 双重检查
    if tm.currentToken != nil && tm.isValid(tm.currentToken) {
        return tm.currentToken.AccessToken, nil
    }
    
    // 尝试刷新 Token
    if tm.currentToken != nil && tm.currentToken.RefreshToken != "" {
        newToken, err := tm.client.RefreshAccessToken(ctx, tm.currentToken.RefreshToken)
        if err == nil {
            tm.currentToken = newToken
            return newToken.AccessToken, nil
        }
    }
    
    // 获取新的 Token
    newToken, err := tm.client.GetClientCredentialsToken(ctx)
    if err != nil {
        return "", err
    }
    
    tm.currentToken = newToken
    return newToken.AccessToken, nil
}

func (tm *TokenManager) isValid(token *client.TokenResponse) bool {
    // 简单判断：Token 获取后 1 小时内有效
    // 实际应该记录 Token 的获取时间
    return true
}
```

### 2. 自动重试的 API 客户端

```go
package api

import (
    "context"
    "sun-panel-micro-store/pkg/oauth2/client"
)

type AutoRefreshAPIClient struct {
    apiClient    *client.APIClient
    tokenManager *TokenManager
}

func NewAutoRefreshAPIClient(baseURL string, tokenManager *TokenManager) *AutoRefreshAPIClient {
    return &AutoRefreshAPIClient{
        apiClient:    client.NewAPIClient(baseURL, 30),
        tokenManager: tokenManager,
    }
}

func (c *AutoRefreshAPIClient) Call(ctx context.Context, method, endpoint string, requestData interface{}, responseData interface{}) error {
    // 获取有效的 Token
    token, err := c.tokenManager.GetValidToken(ctx)
    if err != nil {
        return err
    }
    
    // 调用 API
    err = c.apiClient.Call(ctx, method, endpoint, token, requestData, responseData)
    if err != nil {
        // 如果是 Token 过期错误，尝试刷新后重试
        // 这里需要根据实际错误判断
        return err
    }
    
    return nil
}
```

## 存储实现

### 授权码存储（必须实现过期时间）

```go
type AuthCodeModel struct {
    Code        string    `gorm:"primaryKey"`
    AccessToken string
    CToken      string
    ClientID    string
    UserID      uint
    CreatedAt   time.Time
    ExpiredAt   time.Time
}

func (s *DBAuthCodeStore) SetAuthCode(code string, data server.OAuthCodeData) error {
    model := AuthCodeModel{
        Code:        code,
        AccessToken: data.AccessToken,
        CToken:      data.CToken,
        ClientID:    data.ClientID,
        UserID:      data.UserID,
        CreatedAt:   time.Now(),
        ExpiredAt:   time.Now().Add(10 * time.Minute),
    }
    return s.db.Create(&model).Error
}

func (s *DBAuthCodeStore) GetAuthCode(code string) (server.OAuthCodeData, error) {
    var model AuthCodeModel
    err := s.db.Where("code = ? AND expired_at > ?", code, time.Now()).First(&model).Error
    if err != nil {
        return server.OAuthCodeData{}, err
    }
    return server.OAuthCodeData{
        Code:        model.Code,
        AccessToken: model.AccessToken,
        CToken:      model.CToken,
        ClientID:    model.ClientID,
        UserID:      model.UserID,
    }, nil
}
```

## 安全最佳实践

### 1. State 参数验证

```go
func (h *OAuthHandler) Authorize(c *gin.Context) {
    state := c.Query("state")
    if state == "" {
        // 生成随机 state
        state = generateRandomState()
        // 存储到 session
        session := sessions.Default(c)
        session.Set("oauth_state", state)
        session.Save()
    }
    
    // ... 重定向到授权页面
}

func (h *OAuthHandler) handleCallback(c *gin.Context) {
    receivedState := c.Query("state")
    
    // 验证 state
    session := sessions.Default(c)
    expectedState := session.Get("oauth_state")
    
    if receivedState != expectedState {
        c.JSON(400, gin.H{"error": "invalid state"})
        return
    }
    
    // ... 处理授权码
}
```

### 2. PKCE 支持（可选增强）

```go
import "crypto/sha256"

func generateCodeVerifier() string {
    // 生成随机字符串
    return generateRandomString(128)
}

func generateCodeChallenge(verifier string) string {
    h := sha256.New()
    h.Write([]byte(verifier))
    return base64URLEncode(h.Sum(nil))
}

// 在授权请求中使用
func (c *OAuth2Client) GetAuthorizationURLWithPKCE(redirectURI, state string) (string, string) {
    verifier := generateCodeVerifier()
    challenge := generateCodeChallenge(verifier)
    
    params := url.Values{}
    params.Set("code_challenge", challenge)
    params.Set("code_challenge_method", "S256")
    // ... 其他参数
    
    return authURL, verifier
}
```

### 3. Token 安全存储

```go
// 前端
// 使用 HttpOnly Cookie 存储 Token
func handleTokenResponse(tokenResp *client.TokenResponse, w http.ResponseWriter) {
    // Access Token 存储在内存中
    // Refresh Token 存储在 HttpOnly Cookie 中
    
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    tokenResp.RefreshToken,
        HttpOnly: true,
        Secure:   true, // HTTPS
        SameSite: http.SameSiteStrictMode,
        MaxAge:   7 * 24 * 3600, // 7天
    })
}
```

## 常见问题

### Q1: 如何实现多租户？

```go
// 在存储接口中增加租户 ID
type MultiTenantTokenStore struct {
    db *gorm.DB
}

func (s *MultiTenantTokenStore) GetAccessToken(tenantID string, token string) (server.AccessTokenData, error) {
    var model TokenModel
    err := s.db.Where("tenant_id = ? AND token = ?", tenantID, token).First(&model).Error
    // ...
}
```

### Q2: 如何实现 Token 黑名单？

```go
type BlacklistTokenStore struct {
    store      server.TokenStore
    blacklist  *redis.Client
}

func (s *BlacklistTokenStore) GetAccessToken(token string) (server.AccessTokenData, error) {
    // 检查黑名单
    if s.isBlacklisted(token) {
        return server.AccessTokenData{}, errors.New("token is blacklisted")
    }
    
    return s.store.GetAccessToken(token)
}

func (s *BlacklistTokenStore) isBlacklisted(token string) bool {
    ctx := context.Background()
    return s.blacklist.Exists(ctx, "blacklist:"+token).Val() > 0
}
```

### Q3: 如何实现单点登录（SSO）？

```go
// 在用户登录时，生成 SSO Session
func (h *OAuthHandler) SSOLogin(c *gin.Context) {
    // 验证用户登录
    
    // 生成 SSO Session
    ssoSession := generateSSOSession(userID)
    
    // 存储到 Redis
    h.ssoStore.Set(ssoSession, userID)
    
    // 设置 Cookie
    c.SetCookie("sso_session", ssoSession, ...)
}

// 其他应用验证 SSO Session
func (h *OAuthHandler) ValidateSSO(c *gin.Context) {
    ssoSession := c.GetHeader("SSO-Session")
    userID := h.ssoStore.Get(ssoSession)
    // ...
}
```

## 性能优化

### 1. Token 缓存

```go
type CachedTokenStore struct {
    store server.TokenStore
    cache *lru.Cache
}

func (s *CachedTokenStore) GetAccessToken(token string) (server.AccessTokenData, error) {
    // 先从缓存读取
    if data, ok := s.cache.Get(token); ok {
        return data.(server.AccessTokenData), nil
    }
    
    // 从存储读取
    data, err := s.store.GetAccessToken(token)
    if err != nil {
        return data, err
    }
    
    // 写入缓存
    s.cache.Add(token, data)
    return data, nil
}
```

### 2. 批量删除过期 Token

```go
// 定时任务清理过期 Token
func CleanExpiredTokens(db *gorm.DB) {
    ticker := time.NewTicker(1 * time.Hour)
    for range ticker.C {
        db.Where("expired_at < ?", time.Now()).Delete(&TokenModel{})
    }
}
```

## 测试

### 单元测试示例

```go
package server_test

import (
    "testing"
    "sun-panel-micro-store/pkg/oauth2/server"
)

func TestTokenGeneration(t *testing.T) {
    config := server.DefaultOAuthConfig()
    tm := server.NewTokenManager(config)
    
    token, err := tm.GenerateAccessToken("client1", 1, "secret")
    if err != nil {
        t.Fatalf("Failed to generate token: %v", err)
    }
    
    if token == "" {
        t.Error("Token should not be empty")
    }
    
    claims, err := tm.ValidateAccessToken(token, "secret")
    if err != nil {
        t.Fatalf("Failed to validate token: %v", err)
    }
    
    if claims.UserID != 1 {
        t.Errorf("Expected UserID 1, got %d", claims.UserID)
    }
}
```

## 监控和日志

```go
import "github.com/sirupsen/logrus"

type LoggingTokenStore struct {
    store server.TokenStore
    log   *logrus.Logger
}

func (s *LoggingTokenStore) SetAccessToken(token string, data server.AccessTokenData) error {
    s.log.WithFields(logrus.Fields{
        "client_id": data.ClientID,
        "user_id":   data.UserID,
    }).Info("Setting access token")
    
    return s.store.SetAccessToken(token, data)
}
```
