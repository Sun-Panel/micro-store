# OAuth2 模块

通用的 OAuth2 服务端和客户端实现，支持多种授权模式。

## 功能特性

### 服务端
- ✅ 授权码模式（Authorization Code）
- ✅ 密码模式（Password）
- ✅ 客户端凭证模式（Client Credentials）
- ✅ 刷新令牌（Refresh Token）
- ✅ 单点登出（SSO Logout）
- ✅ JWT Token 生成和验证
- ✅ 灵活的存储接口

### 客户端
- ✅ 授权码模式
- ✅ 密码模式
- ✅ 客户端凭证模式
- ✅ 刷新令牌
- ✅ API 调用封装

## 目录结构

```
pkg/oauth2/
├── common/          # 公共类型和常量
│   └── types.go     # 公共类型定义
├── server/          # OAuth2 服务端实现
│   ├── types.go     # 服务端类型定义
│   ├── token.go     # Token 生成和验证
│   └── handler.go   # HTTP 处理器
├── client/          # OAuth2 客户端实现
│   ├── types.go     # 客户端类型定义
│   └── client.go    # 客户端实现
└── README.md        # 文档
```

## 快速开始

### 服务端使用

#### 1. 定义存储实现

```go
package main

import (
    "time"
    "sync"
    "sun-panel-micro-store/pkg/oauth2/server"
)

// 实现第三方应用存储接口
type ThirdAppStoreImpl struct {
    apps map[string]ThirdAppInfo
}

func (s *ThirdAppStoreImpl) GetByClientID(clientID string) (server.ThirdAppInfo, error) {
    app, ok := s.apps[clientID]
    if !ok {
        return nil, errors.New("not found")
    }
    return app, nil
}

func (s *ThirdAppStoreImpl) GetByClientIDAndSecret(clientID, clientSecret string) (server.ThirdAppInfo, error) {
    app, ok := s.apps[clientID]
    if !ok || app.GetClientSecret() != clientSecret {
        return nil, errors.New("invalid credentials")
    }
    return app, nil
}

// 实现其他存储接口...
```

#### 2. 创建 OAuth2 Handler

```go
package main

import (
    "github.com/gin-gonic/gin"
    "sun-panel-micro-store/pkg/oauth2/server"
)

func main() {
    // 创建配置
    config := &server.OAuthConfig{
        AccessTokenExpireTime:  7200,        // 2小时
        RefreshTokenExpireTime: 604800,     // 7天
        AuthCodeExpireTime:     600,        // 10分钟
        EnableSSOLogout:        true,
    }

    // 创建 handler
    handler := server.NewOAuthHandler(config)

    // 设置存储实现
    handler.SetStores(
        &ThirdAppStoreImpl{},
        &UserStoreImpl{},
        &TokenStoreImpl{},
        &AuthCodeStoreImpl{},
        &RefreshTokenStoreImpl{},
    )

    // 注册路由
    r := gin.Default()
    handler.RegisterRoutes(r.Group("/api"))

    r.Run(":8080")
}
```

### 客户端使用

#### 1. 创建 OAuth2 客户端

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
        APIServerURL:  "http://localhost:8080",
        ClientID:      "your_client_id",
        ClientSecret:  "your_client_secret",
        RedirectURI:   "http://localhost:3000/callback",
        Timeout:       30,
    }

    oauthClient, err := client.NewOAuth2Client(config)
    if err != nil {
        panic(err)
    }

    // 获取授权 URL
    authURL := oauthClient.GetAuthorizationURL(config.RedirectURI, "random_state")
    fmt.Println("授权 URL:", authURL)
}
```

#### 2. 使用授权码获取 Token

```go
func handleCallback(code string) {
    ctx := context.Background()
    
    tokenResp, err := oauthClient.GetAccessTokenByCode(ctx, code)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Access Token: %s\n", tokenResp.AccessToken)
    fmt.Printf("Refresh Token: %s\n", tokenResp.RefreshToken)
    fmt.Printf("Expires In: %d\n", tokenResp.ExpiresIn)
}
```

#### 3. 使用客户端凭证模式

```go
func getClientToken() {
    ctx := context.Background()
    
    tokenResp, err := oauthClient.GetClientCredentialsToken(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Access Token: %s\n", tokenResp.AccessToken)
}
```

#### 4. 调用 API

```go
func callAPI(accessToken string) {
    apiClient := client.NewAPIClient("http://localhost:8080", 30)
    
    ctx := context.Background()
    userInfo, err := apiClient.GetUserInfo(ctx, accessToken)
    if err != nil {
        panic(err)
    }

    fmt.Printf("User: %s (%s)\n", userInfo.Name, userInfo.Email)
}
```

## API 端点

### 服务端端点

| 端点 | 方法 | 描述 |
|------|------|------|
| `/oauth2/v1/authorize` | GET | 授权端点 |
| `/oauth2/v1/token` | POST | Token 端点（授权码模式、密码模式） |
| `/oauth2/v1/clientCredentials/token` | POST | 客户端凭证模式 Token 端点 |
| `/oauth2/v1/sso/logout` | POST | 单点登出 |

### 支持的授权模式

#### 1. 授权码模式（Authorization Code）

```
1. 用户访问授权 URL
2. 用户授权后获得 code
3. 使用 code 换取 access_token
```

#### 2. 密码模式（Password）

```
1. 用户提供用户名和密码
2. 直接获取 access_token
```

#### 3. 客户端凭证模式（Client Credentials）

```
1. 客户端提供 client_id 和 client_secret
2. 获取 access_token（无用户上下文）
```

#### 4. 刷新令牌（Refresh Token）

```
1. 使用 refresh_token 获取新的 access_token
```

## 存储接口

### ThirdAppStore - 第三方应用存储
```go
type ThirdAppStore interface {
    GetByClientID(clientID string) (ThirdAppInfo, error)
    GetByClientIDAndSecret(clientID, clientSecret string) (ThirdAppInfo, error)
}
```

### UserStore - 用户存储
```go
type UserStore interface {
    GetByUsernameAndPassword(username, password string) (UserInfo, error)
    GetByID(userID uint) (UserInfo, error)
}
```

### TokenStore - Token 存储
```go
type TokenStore interface {
    SetAccessToken(token string, data AccessTokenData) error
    GetAccessToken(token string) (AccessTokenData, error)
    DeleteAccessToken(token string) error
}
```

### AuthCodeStore - 授权码存储
```go
type AuthCodeStore interface {
    SetAuthCode(code string, data OAuthCodeData) error
    GetAuthCode(code string) (OAuthCodeData, error)
    DeleteAuthCode(code string) error
}
```

### RefreshTokenStore - 刷新令牌存储
```go
type RefreshTokenStore interface {
    SetRefreshToken(token string, data RefreshTokenData) error
    GetRefreshToken(token string) (RefreshTokenData, error)
    DeleteRefreshToken(token string) error
}
```

## 自定义存储实现

可以使用 Redis、数据库或其他存储方式实现存储接口。

### Redis 示例

```go
import "github.com/go-redis/redis/v8"

type RedisTokenStore struct {
    client *redis.Client
}

func (s *RedisTokenStore) SetAccessToken(token string, data server.AccessTokenData) error {
    ctx := context.Background()
    jsonData, _ := json.Marshal(data)
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

## 已知问题和兼容性说明

### ⚠️ Token 响应格式兼容性（临时）

**问题描述：**

当前授权中心（192.168.3.101:3088）的 Token 端点返回的是包装格式：
```json
{
  "code": 0,
  "data": {
    "access_token": "...",
    "token_type": "Bearer",
    "expires_in": 3600
  },
  "msg": "OK"
}
```

而标准 OAuth2 协议要求 Token 端点直接返回：
```json
{
  "access_token": "...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

**临时解决方案：**

客户端已实现兼容性代码，支持同时解析两种格式：
- 优先尝试解析包装格式（检查 `code=0` 和 `data` 字段）
- 如果不是包装格式，则按标准 OAuth2 格式解析

**兼容代码位置：**
- 文件：`pkg/oauth2/client/client.go`
- 方法：`requestToken()`
- 标记：搜索 `兼容性代码开始` 和 `兼容性代码结束` 注释

**待办事项：**

- [ ] 修改授权中心 Token 端点，直接返回标准 OAuth2 格式
- [ ] 删除客户端兼容性代码（标记为"兼容性代码"的部分）
- [ ] 更新此文档

**影响范围：**
- 仅影响 Token 端点的响应解析
- 不影响其他 OAuth2 功能
- 其他标准 OAuth2 客户端可能无法连接当前授权中心

**修复日期：** 2026-03-13

## 安全建议

1. **使用 HTTPS**: 所有 OAuth2 端点必须使用 HTTPS
2. **State 参数**: 使用 state 参数防止 CSRF 攻击
3. **Token 过期**: 设置合理的 Token 过期时间
4. **安全存储**: Client Secret 必须安全存储
5. **授权码一次性**: 授权码使用后立即删除
6. **限制 Scope**: 实现最小权限原则

## 许可证

MIT License
