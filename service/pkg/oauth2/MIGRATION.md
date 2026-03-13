# 迁移指南

本文档帮助你从现有的 OAuth2 实现迁移到新的统一模块。

## 迁移概述

### 现有代码结构

**授权中心（服务端）：**
- `sun-auth-center/service/apiOAuth2/v1.go` - OAuth2 API 实现
- `sun-auth-center/service/biz/thirdApp.go` - 第三方应用业务逻辑

**微商城（客户端）：**
- `sun-panel-micro-store/service/lib/sunStore/oAuth2.go` - OAuth2 客户端
- `sun-panel-micro-store/service/api/api_v1/oAuth2/oauth2.go` - OAuth2 处理器
- `sun-panel-micro-store/service/api/api_v1/system/login.go` - 登录处理

### 新模块结构

```
pkg/oauth2/
├── common/          # 公共类型
├── server/          # 服务端实现
├── client/          # 客户端实现
├── example/         # 使用示例
└── README.md        # 文档
```

## 迁移步骤

### 第一步：迁移服务端（授权中心）

#### 1. 创建存储适配器

```go
package adapter

import (
    "sun-panel-micro-store/pkg/oauth2/server"
    "sun-panel/models"
    "sun-panel/global"
    "gorm.io/gorm"
)

// ThirdAppStoreAdapter 适配现有的第三方应用模型
type ThirdAppStoreAdapter struct{}

func (a *ThirdAppStoreAdapter) GetByClientID(clientID string) (server.ThirdAppInfo, error) {
    info := models.ThirdApp{}
    err := global.Db.First(&info, "client_id=?", clientID).Error
    if err != nil {
        return nil, err
    }
    return &ThirdAppAdapter{app: info}, nil
}

func (a *ThirdAppStoreAdapter) GetByClientIDAndSecret(clientID, clientSecret string) (server.ThirdAppInfo, error) {
    info := models.ThirdApp{}
    err := global.Db.First(&info, "client_id=? AND client_secret=?", clientID, clientSecret).Error
    if err != nil {
        return nil, err
    }
    return &ThirdAppAdapter{app: info}, nil
}

// ThirdAppAdapter 适配第三方应用模型
type ThirdAppAdapter struct {
    app models.ThirdApp
}

func (a *ThirdAppAdapter) GetClientID() string     { return a.app.ClientId }
func (a *ThirdAppAdapter) GetClientSecret() string { return a.app.ClientSecret }
func (a *ThirdAppAdapter) IsEnabled() bool         { return a.app.IsEnabled }
func (a *ThirdAppAdapter) IsSSOLogout() bool       { return a.app.IsSsoLogout }
func (a *ThirdAppAdapter) GetSSOLogoutURL() string { return a.app.SsoLogoutUrl }

// UserStoreAdapter 适配现有的用户模型
type UserStoreAdapter struct{}

func (a *UserStoreAdapter) GetByUsernameAndPassword(username, password string) (server.UserInfo, error) {
    mUser := models.User{}
    user, err := mUser.GetUserInfoByUsernameAndPassword(username, password)
    if err != nil {
        return nil, err
    }
    return &UserAdapter{user: user}, nil
}

func (a *UserStoreAdapter) GetByID(userID uint) (server.UserInfo, error) {
    mUser := models.User{}
    user, err := mUser.GetUserInfoByUid(userID)
    if err != nil {
        return nil, err
    }
    return &UserAdapter{user: user}, nil
}

// UserAdapter 适配用户模型
type UserAdapter struct {
    user models.User
}

func (a *UserAdapter) GetUserID() uint      { return a.user.ID }
func (a *UserAdapter) GetUsername() string  { return a.user.Username }
func (a *UserAdapter) GetPassword() string  { return a.user.Password }
func (a *UserAdapter) IsActive() bool       { return a.user.Status == 1 }
```

#### 2. 创建 Token 存储适配器

```go
package adapter

import (
    "sun-panel-micro-store/pkg/oauth2/server"
    "sun-panel/global"
)

// TokenStoreAdapter 适配现有的 Token 缓存
type TokenStoreAdapter struct{}

func (a *TokenStoreAdapter) SetAccessToken(token string, data server.AccessTokenData) error {
    global.CUserToken.SetDefault(token, data.CToken)
    return nil
}

func (a *TokenStoreAdapter) GetAccessToken(token string) (server.AccessTokenData, error) {
    // 实现根据你的缓存逻辑
    // 这里使用现有的缓存机制
    data := server.AccessTokenData{}
    // ...
    return data, nil
}

func (a *TokenStoreAdapter) DeleteAccessToken(token string) error {
    global.CUserToken.Delete(token)
    return nil
}

// AuthCodeStoreAdapter 适配授权码缓存
type AuthCodeStoreAdapter struct{}

func (a *AuthCodeStoreAdapter) SetAuthCode(code string, data server.OAuthCodeData) error {
    // 使用现有的缓存机制
    // biz.ThirdApp.OAuthCodeCache.SetDefault(code, data)
    return nil
}

func (a *AuthCodeStoreAdapter) GetAuthCode(code string) (server.OAuthCodeData, error) {
    // 使用现有的缓存机制
    // data, ok := biz.ThirdApp.OAuthCodeCache.Get(code)
    // ...
    return server.OAuthCodeData{}, nil
}

func (a *AuthCodeStoreAdapter) DeleteAuthCode(code string) error {
    // biz.ThirdApp.OAuthCodeCache.Delete(code)
    return nil
}
```

#### 3. 替换路由注册

**旧代码：**
```go
// service/router/oAuth2/oAuth2.go
func InitOAuth2(router *gin.RouterGroup) {
    api := apiOAuth2.OAuth2{}
    r := router.Group("v1")
    r.GET("authorize", api.Auth)
    r.POST("token", api.Token)
    r.POST("clientCredentials/token", api.ClientCredentialsToken)
}
```

**新代码：**
```go
// service/router/oAuth2/oAuth2.go
package oAuth2

import (
    "sun-panel-micro-store/pkg/oauth2/server"
    "sun-panel/adapter"
    "github.com/gin-gonic/gin"
)

func InitOAuth2(router *gin.RouterGroup) {
    // 创建配置
    config := &server.OAuthConfig{
        AccessTokenExpireTime:  7200,
        RefreshTokenExpireTime: 604800,
        AuthCodeExpireTime:     600,
        EnableSSOLogout:        true,
    }
    
    // 创建 handler
    handler := server.NewOAuthHandler(config)
    
    // 设置存储适配器
    handler.SetStores(
        &adapter.ThirdAppStoreAdapter{},
        &adapter.UserStoreAdapter{},
        &adapter.TokenStoreAdapter{},
        &adapter.AuthCodeStoreAdapter{},
        &adapter.RefreshTokenStoreAdapter{},
    )
    
    // 注册路由
    handler.RegisterRoutes(router)
    
    // 添加用户授权登录接口
    router.POST("auth/login", func(c *gin.Context) {
        // 从上下文获取当前用户
        userInfo, _ := base.GetCurrentUserInfo(c)
        c.Set("user_id", userInfo.ID)
        handler.AuthLogin(c)
    })
}
```

### 第二步：迁移客户端（微商城）

#### 1. 创建客户端实例

**旧代码：**
```go
// service/biz/sunStore.go
func (s *SunStoreType) GetClientIdAndSecret() (clientId string, clientSecret string) {
    clientId = global.Config.GetValueString("sun_store", "client_id")
    clientSecret = global.Config.GetValueString("sun_store", "client_secret")
    return
}

// service/api/api_v1/system/login.go
sunapi := sunStore.NewSunStoreApi(apiHost, clientId, clientSecret)
accessToken, err := sunapi.GetAccessToken(sunStore.AccessTokenResquest{
    Code: req.Code,
})
```

**新代码：**
```go
// service/biz/sunStore.go
package biz

import (
    "sun-panel-micro-store/pkg/oauth2/client"
)

type SunStoreType struct {
    oauthClient *client.OAuth2Client
}

func (s *SunStoreType) InitOAuthClient() error {
    config := &client.Config{
        AuthServerURL: s.ApiHost(),
        APIServerURL:  s.ApiHost(),
        ClientID:      global.Config.GetValueString("sun_store", "client_id"),
        ClientSecret:  global.Config.GetValueString("sun_store", "client_secret"),
        RedirectURI:   "", // 根据实际情况设置
        Timeout:       30,
    }
    
    var err error
    s.oauthClient, err = client.NewOAuth2Client(config)
    return err
}

func (s *SunStoreType) GetAccessToken(code string) (*client.TokenResponse, error) {
    ctx := context.Background()
    return s.oauthClient.GetAccessTokenByCode(ctx, code)
}
```

#### 2. 替换登录处理

**旧代码：**
```go
// service/api/api_v1/system/login.go
sunapi := sunStore.NewSunStoreApi(apiHost, clientId, clientSecret)
accessToken, err := sunapi.GetAccessToken(sunStore.AccessTokenResquest{
    Code: req.Code,
})
```

**新代码：**
```go
// service/api/api_v1/system/login.go
import "sun-panel-micro-store/pkg/oauth2/client"

func (l *LoginApi) OAuth2CodeLogin(c *gin.Context) {
    req := OAuth2CodeLoginReq{}
    if err := c.ShouldBindJSON(&req); err != nil {
        apiReturn.ErrorParamFomat(c, err.Error())
        return
    }
    
    // 使用新的 OAuth2 客户端
    ctx := context.Background()
    tokenResp, err := biz.SunStore.GetAccessToken(req.Code)
    if err != nil {
        global.Logger.Errorln("获取access_token失败", err.Error())
        apiReturn.ErrorByCodeAndMsg(c, -2, err.Error())
        return
    }
    
    // 使用 Access Token 获取用户信息
    apiClient := client.NewAPIClient(biz.SunStore.ApiHost(), 30)
    userInfo, err := apiClient.GetUserInfo(ctx, tokenResp.AccessToken)
    
    // ... 后续处理逻辑保持不变
}
```

### 第三步：逐步迁移

建议按以下顺序迁移：

1. **测试环境验证**
   - 在测试环境部署新模块
   - 运行完整的 OAuth2 流程测试
   - 验证所有授权模式

2. **灰度发布**
   - 先迁移客户端凭证模式
   - 再迁移授权码模式
   - 最后迁移密码模式

3. **监控和回滚**
   - 监控 Token 生成和验证的成功率
   - 准备回滚方案
   - 记录所有异常情况

## 迁移对照表

| 功能 | 旧实现 | 新实现 | 备注 |
|------|--------|--------|------|
| 授权端点 | `apiOAuth2.OAuth2.Auth` | `server.OAuthHandler.Authorize` | 路径不变 |
| Token 端点 | `apiOAuth2.OAuth2.Token` | `server.OAuthHandler.Token` | 路径不变 |
| 客户端凭证 | `sunapi.ClientCredentialsAuth` | `client.GetClientCredentialsToken` | 接口更简洁 |
| 授权码换 Token | `sunapi.GetAccessToken` | `client.GetAccessTokenByCode` | 接口更简洁 |
| 密码模式 | `sunapi.PasswordAuth` | `client.GetAccessTokenByPassword` | 接口更简洁 |
| API 调用 | 自定义实现 | `client.APIClient.Call` | 统一封装 |

## 兼容性说明

### 路由兼容

新模块保持与旧实现相同的路由结构：

```
GET  /oauth2/v1/authorize          - 授权端点
POST /oauth2/v1/token              - Token 端点
POST /oauth2/v1/clientCredentials/token - 客户端凭证模式
POST /oauth2/v1/sso/logout         - 单点登出
```

### Token 格式兼容

新模块使用相同的 JWT 格式，确保现有 Token 可以继续使用：

```go
type OAuthClaims struct {
    UserID   uint   `json:"user_id"`
    ClientID string `json:"client_id"`
    OpenID   string `json:"openid,omitempty"`
    jwt.RegisteredClaims
}
```

### 缓存兼容

可以使用现有的缓存机制，只需实现适配器接口即可。

## 注意事项

1. **不要同时使用新旧模块**：避免在同一个服务中同时使用新旧实现
2. **保持 Token 密钥一致**：迁移时确保 JWT 签名密钥不变
3. **测试所有场景**：迁移后需要测试所有授权模式
4. **监控异常**：密切监控 Token 相关的错误日志
5. **保留旧代码**：迁移成功前保留旧代码，以便回滚

## 迁移检查清单

- [ ] 创建存储适配器
- [ ] 替换路由注册
- [ ] 创建 OAuth2 客户端实例
- [ ] 替换 Token 获取逻辑
- [ ] 替换 API 调用逻辑
- [ ] 测试授权码模式
- [ ] 测试密码模式
- [ ] 测试客户端凭证模式
- [ ] 测试刷新 Token
- [ ] 测试单点登出
- [ ] 验证 Token 格式兼容性
- [ ] 性能测试
- [ ] 灰度发布
- [ ] 全量发布
- [ ] 清理旧代码

## 常见问题

### Q1: 迁移后 Token 不兼容怎么办？

A: 确保 JWT 签名密钥和 Claims 结构与旧实现一致。

### Q2: 如何处理现有的 Token？

A: 可以保留旧 Token 直到过期，新 Token 使用新模块生成。建议设置过渡期。

### Q3: 迁移会影响现有用户吗？

A: 如果路由和 Token 格式保持兼容，现有用户不会受到影响。

### Q4: 如何验证迁移成功？

A: 运行完整的 OAuth2 流程测试，验证所有功能正常。
