# OAuth2 模块集成指南

本文档说明如何在当前项目中集成并使用新的 OAuth2 模块。

## 📁 文件修改说明

### 1. 新增文件

#### `service/adapter/oauth2_adapter.go`
新的 OAuth2 客户端适配器，提供与现有代码兼容的接口。

**主要功能：**
- `GetSunStoreClient()` - 获取 SunStore 客户端实例
- `GetAccessToken()` - 使用授权码获取 Token
- `ClientCredentialsAuth()` - 客户端凭证模式
- `PasswordAuth()` - 密码模式
- `RefreshToken()` - 刷新 Token

### 2. 修改文件

#### `service/biz/sunStore.go`

**新增方法:**
```go
// GetOAuth2Client 获取 OAuth2 客户端实例
func (s *SunStoreType) GetOAuth2Client() (*client.OAuth2Client, error)

// GetAccessTokenByCode 使用授权码获取 Access Token（新方法）
func (s *SunStoreType) GetAccessTokenByCode(code string) (*client.TokenResponse, error)

// GetClientCredentialsToken 使用客户端凭证模式获取 Token（新方法）
func (s *SunStoreType) GetClientCredentialsToken() (*client.TokenResponse, error)

// GetAPIClient 获取 API 客户端
func (s *SunStoreType) GetAPIClient() *client.APIClient
```

#### `service/api/api_v1/system/login.go`

**添加注释说明如何迁移:**
- 第 147-161 行：添加了详细的迁移注释
- 展示了新旧方法的对比

## 🚀 使用示例

### 方式一：使用适配器（推荐用于快速迁移）

```go
import "sun-panel/adapter"

// 获取客户端实例
sunStoreClient, err := adapter.GetSunStoreClient()
if err != nil {
    // 处理错误
}

// 使用授权码获取 Token
tokenResp, err := sunStoreClient.GetAccessToken(code)
if err != nil {
    // 处理错误
}

accessToken := tokenResp.AccessToken
refreshToken := tokenResp.RefreshToken
```

### 方式二：直接使用新方法（推荐）

```go
// 使用授权码获取 Token
tokenResp, err := biz.SunStore.GetAccessTokenByCode(code)
if err != nil {
    global.Logger.Errorln("获取access_token失败", err.Error())
    return
}

accessToken := tokenResp.AccessToken
refreshToken := tokenResp.RefreshToken
```

### 方式三：使用客户端凭证模式

```go
// 旧方法
sunapi := sunStore.NewSunStoreApi(apiHost, clientId, clientSecret)
resp, err := sunapi.ClientCredentialsAuth(sunStore.ClientCredentialsParam{}, false)

// 新方法
tokenResp, err := biz.SunStore.GetClientCredentialsToken()
if err != nil {
    return
}

accessToken := tokenResp.AccessToken
```

## 📋 迁移步骤

### 第一步：安装依赖

```bash
cd service
go mod tidy
```

### 第二步：选择迁移方式

#### 选项 A：渐进式迁移（推荐）

1. 保留旧代码
2. 在新功能中使用新模块
3. 逐步替换旧代码

**示例：修改 OAuth2CodeLogin**

```go
// 旧代码（保留）
// sunapi := sunStore.NewSunStoreApi(apiHost, clientId, clientSecret)
// accessToken, err := sunapi.GetAccessToken(sunStore.AccessTokenResquest{
//     Code: req.Code,
// })

// 新代码
tokenResp, err := biz.SunStore.GetAccessTokenByCode(req.Code)
if err != nil {
    global.Logger.Errorln("获取access_token失败", err.Error())
    apiReturn.ErrorByCodeAndMsg(c, -2, err.Error())
    return
}

accessToken := tokenResp.AccessToken
```

#### 选项 B：完全替换

1. 删除旧代码
2. 全面使用新模块

### 第三步：测试验证

1. **测试授权码模式**
   ```bash
   # 测试登录流程
   curl -X POST http://localhost:8080/oAuth2CodeLogin \
     -H "Content-Type: application/json" \
     -d '{"code":"your_code"}'
   ```

2. **测试客户端凭证模式**
   ```bash
   # 测试服务间调用
   curl -X POST http://localhost:8080/api/test
   ```

3. **检查 Token 格式**
   - 确认 AccessToken 格式正确
   - 验证 RefreshToken 可用

## 🔄 新旧方法对照表

| 功能 | 旧方法 | 新方法 | 说明 |
|------|--------|--------|------|
| 创建客户端 | `sunStore.NewSunStoreApi(apiHost, clientId, clientSecret)` | `biz.SunStore.GetOAuth2Client()` | 更简洁 |
| 授权码换Token | `sunapi.GetAccessToken(sunStore.AccessTokenResquest{Code: code})` | `biz.SunStore.GetAccessTokenByCode(code)` | 接口更清晰 |
| 客户端凭证 | `sunapi.ClientCredentialsAuth(sunStore.ClientCredentialsParam{}, false)` | `biz.SunStore.GetClientCredentialsToken()` | 无需参数 |
| 密码模式 | `sunapi.PasswordAuth(sunStore.PasswordAuthRequest{...})` | `biz.SunStore.GetOAuth2Client().GetAccessTokenByPassword(ctx, username, password)` | 需要上下文 |
| 刷新Token | 未实现 | `oauthClient.RefreshAccessToken(ctx, refreshToken)` | 新增功能 |
| 授权URL | 手动构造 | `oauthClient.GetAuthorizationURL(redirectURI, state)` | 新增功能 |

## ✨ 新模块优势

### 1. 代码更简洁
```go
// 旧代码
sunapi := sunStore.NewSunStoreApi(apiHost, clientId, clientSecret)
accessToken, err := sunapi.GetAccessToken(sunStore.AccessTokenResquest{
    Code: req.Code,
})

// 新代码
tokenResp, err := biz.SunStore.GetAccessTokenByCode(req.Code)
```

### 2. 类型更明确
```go
// 旧代码 - AccessTokenResponse 类型不明确
type AccessTokenResponse struct {
    AccessToken string `json:"access_token"`
    Scope       string `json:"scope"`
}

// 新代码 - 标准的 Token 响应
type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token,omitempty"`
    Scope        string `json:"scope,omitempty"`
}
```

### 3. 错误处理更完善
- 标准的 OAuth2 错误响应
- 更详细的错误信息
- 统一的错误处理机制

### 4. 功能更完整
- ✅ 支持刷新 Token
- ✅ 自动生成授权 URL
- ✅ 完整的 API 客户端
- ✅ Context 支持

## 📌 注意事项

### 1. 向后兼容
- 旧代码可以继续使用
- 新旧代码可以共存
- 建议渐进式迁移

### 2. 配置不变
- 配置文件路径：`service/conf/conf.ini`
- 配置项：
  ```ini
  [sun_store]
  client_id = your_client_id
  client_secret = your_client_secret
  api_host = http://auth-server
  auth_endpoint = http://auth-server/oauth2/v1/authorize
  ```

### 3. Token 缓存
- 现有的 Token 缓存机制保持不变
- 可以继续使用 `global.SystemSetting` 存储 Token

### 4. 依赖关系
```go
import (
    "sun-panel-micro-store/pkg/oauth2/client"  // 新模块
    "sun-panel/biz"                            // 现有业务逻辑
    "sun-panel/lib/sunStore"                   // 旧模块（可选）
)
```

## 🧪 测试清单

- [ ] 测试授权码模式登录
- [ ] 测试客户端凭证模式
- [ ] 测试密码模式
- [ ] 测试 Token 刷新
- [ ] 验证 Token 格式
- [ ] 检查错误处理
- [ ] 确认向后兼容性

## 🔗 相关文档

- [OAuth2 模块 README](../pkg/oauth2/README.md)
- [使用指南](../pkg/oauth2/USAGE.md)
- [迁移指南](../pkg/oauth2/MIGRATION.md)
- [测试报告](../pkg/oauth2/TEST_REPORT.md)

## 💡 建议

1. **优先使用新方法** - 新方法更简洁、更标准
2. **保留旧代码** - 迁移期间保留旧代码作为备份
3. **逐步迁移** - 不要一次性替换所有代码
4. **充分测试** - 每次修改后都要测试所有 OAuth2 流程
5. **监控日志** - 注意观察 Token 相关的错误日志

---

**迁移完成后，建议删除旧的 `service/lib/sunStore/oAuth2.go` 文件。**
