# OAuth2 "need to Retry" 错误排查指南

## 错误说明

**错误信息**：`need to Retry`

**错误来源**：`service/api/api_v1/system/login.go:397`

**触发场景**：使用 Access Token 调用主平台 API 获取用户信息失败

## 错误流程

```
前端发送授权码 (code)
    ↓
服务端用 code 换取 access_token  ← 新模块处理
    ↓
服务端用 access_token 调用主平台 API
    ↓
调用失败 → 返回 "need to Retry"
```

## 排查步骤

### 1. 检查主平台连通性

**检查配置**：
```ini
[sun_store]
client_id=xs1ubmt25p
client_secret=pKcHMOKBUOvwsOmokvn82JMTC2zR4mIy
api_host=http://192.168.3.101:3088
auth_endpoint=http://192.168.3.101:3088/oauth2/v1/authorize
```

**测试连通性**：
```bash
# 测试主平台是否可达
curl -I http://192.168.3.101:3088

# 测试授权端点
curl -I http://192.168.3.101:3088/oauth2/v1/authorize

# 测试 API 端点
curl -X POST http://192.168.3.101:3088/openApi/v1/u/user/getCurrentUserInfo \
  -H "Content-Type: application/json"
```

### 2. 查看详细错误日志

**启用调试日志**：

修改 `service/conf/conf.ini`：
```ini
[base]
# 设置为 debug 模式
RUNCODE=debug
```

**查看日志**：
```bash
# 查看运行日志
tail -f /Users/sunjingliang/my_code/sun-panel_group/sun-panel-micro-store/service/runtime/logs/*.log

# 或查看标准输出
# 如果服务运行在终端，直接查看输出
```

### 3. 检查 Token 获取过程

**在新模块中添加调试日志**：

编辑 `service/biz/sunStore.go`：

```go
func (s *SunStoreType) GetAccessTokenByCode(code string) (*client.TokenResponse, error) {
    oauthClient, err := s.GetOAuth2Client()
    if err != nil {
        global.Logger.Errorln("创建 OAuth2 客户端失败:", err)
        return nil, err
    }

    ctx := context.Background()
    tokenResp, err := oauthClient.GetAccessTokenByCode(ctx, code)
    
    // 添加调试日志
    if err != nil {
        global.Logger.Errorln("获取 Access Token 失败:", err)
        return nil, err
    }
    
    global.Logger.Debugln("成功获取 Access Token:", tokenResp.AccessToken)
    global.Logger.Debugln("Token 类型:", tokenResp.TokenType)
    global.Logger.Debugln("过期时间:", tokenResp.ExpiresIn)
    
    return tokenResp, nil
}
```

### 4. 检查 API 调用

**测试 API 调用**：

编辑 `service/api/api_v1/system/login.go`，在 `oAuth2CodeNoLoggedAuthProcess` 中添加日志：

```go
func (l *LoginApi) oAuth2CodeNoLoggedAuthProcess(apiHost string, accessToken sunStore.AccessTokenResponse) (models.User, error) {
    userInfo := models.User{}
    
    // 添加调试日志
    global.Logger.Debugln("调用主平台 API 获取用户信息")
    global.Logger.Debugln("API Host:", apiHost)
    global.Logger.Debugln("Access Token:", accessToken.AccessToken)
    
    openUser, err := biz.SunStore.GetMainPlatformUserInfo(apiHost, accessToken.AccessToken)
    if err != nil {
        global.Logger.Errorln("GetMainPlatformUserInfo 失败:", err)
        global.Logger.Errorln("详细错误信息:", err.Error())
        return userInfo, ErrOAuth2CodeRetry
    }
    
    // ... 后续逻辑
}
```

### 5. 验证新模块是否正常工作

**测试 Token 获取**：

```go
// 在控制器中添加测试代码
func (l *LoginApi) TestOAuth2Token(c *gin.Context) {
    // 测试客户端凭证模式
    tokenResp, err := biz.SunStore.GetClientCredentialsToken()
    if err != nil {
        global.Logger.Errorln("客户端凭证模式失败:", err)
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{
        "access_token": tokenResp.AccessToken,
        "token_type":   tokenResp.TokenType,
        "expires_in":   tokenResp.ExpiresIn,
    })
}
```

## 常见问题

### 问题 1：Access Token 立即过期

**症状**：
- Token 获取成功
- 但立即调用 API 失败（401 Unauthorized）

**原因**：
- 主平台的 Token 过期时间配置过短
- 主平台 Token 验证逻辑有问题

**解决**：
```go
// 检查 Token 过期时间
global.Logger.Debugln("Token 过期时间(秒):", tokenResp.ExpiresIn)
```

### 问题 2：主平台 API 不可达

**症状**：
- 网络超时
- 连接被拒绝

**原因**：
- `api_host` 配置错误
- 主平台服务未启动
- 网络防火墙阻止

**解决**：
```bash
# 检查主平台服务状态
curl http://192.168.3.101:3088/health

# 检查网络连通性
ping 192.168.3.101
telnet 192.168.3.101 3088
```

### 问题 3：Token 格式不正确

**症状**：
- API 返回 401
- 提示 Token 无效

**原因**：
- 主平台期望的 Token 格式不同
- 需要添加前缀（如 "Bearer "）

**检查**：
```go
// 查看 lib/sunStore/openApi/openApi.go
// 确认 Token 是否正确添加到请求头
func (o *OpenApi) Post(url string, data interface{}, resp interface{}) (int, []byte, error) {
    // 检查是否正确设置了 Authorization 头
    // 例如：Bearer YOUR_ACCESS_TOKEN
}
```

### 问题 4：新模块兼容性问题

**症状**：
- 新模块返回的 Token 格式与旧代码不兼容

**检查**：
```go
// service/api/api_v1/system/login.go:159-172
// 确认转换逻辑是否正确
accessToken := sunStore.AccessTokenResponse{
    AccessToken: tokenResp.AccessToken,
    Scope:       tokenResp.Scope,
}
```

## 快速诊断脚本

创建 `service/scripts/test_oauth2.go`：

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "cnb.cool/hslr-s/go-pkg/oauth2-go/client"
)

func main() {
    // 配置
    config := &client.Config{
        AuthServerURL: "http://192.168.3.101:3088",
        APIServerURL:  "http://192.168.3.101:3088",
        ClientID:      "xs1ubmt25p",
        ClientSecret:  "pKcHMOKBUOvwsOmokvn82JMTC2zR4mIy",
        Timeout:       30,
    }
    
    // 创建客户端
    oauthClient, err := client.NewOAuth2Client(config)
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }
    
    // 测试客户端凭证模式
    ctx := context.Background()
    tokenResp, err := oauthClient.GetClientCredentialsToken(ctx)
    if err != nil {
        log.Fatal("获取 Token 失败:", err)
    }
    
    fmt.Printf("✅ Token 获取成功:\n")
    fmt.Printf("   Access Token: %s\n", tokenResp.AccessToken)
    fmt.Printf("   Token Type: %s\n", tokenResp.TokenType)
    fmt.Printf("   Expires In: %d\n", tokenResp.ExpiresIn)
    
    // 测试 API 调用
    apiClient := client.NewAPIClient("http://192.168.3.101:3088", 30)
    
    type UserInfo struct {
        Username string `json:"username"`
        Mail     string `json:"mail"`
    }
    
    var userInfo UserInfo
    err = apiClient.Call(ctx, "POST", "/openApi/v1/u/user/getCurrentUserInfo", 
        tokenResp.AccessToken, nil, &userInfo)
    
    if err != nil {
        log.Fatal("❌ API 调用失败:", err)
    }
    
    fmt.Printf("✅ API 调用成功:\n")
    fmt.Printf("   Username: %s\n", userInfo.Username)
    fmt.Printf("   Mail: %s\n", userInfo.Mail)
}
```

运行测试：
```bash
cd service/scripts
go run test_oauth2.go
```

## 解决方案总结

### 立即可做的检查

1. ✅ **检查主平台连通性**
   ```bash
   curl http://192.168.3.101:3088
   ```

2. ✅ **启用调试日志**
   ```ini
   [base]
   RUNCODE=debug
   ```

3. ✅ **查看错误详情**
   ```bash
   tail -f runtime/logs/*.log
   ```

4. ✅ **测试 Token 获取**
   - 使用上述测试脚本
   - 或在代码中添加日志

### 如果问题持续

1. **回退到旧实现**（临时）
   ```go
   // service/api/api_v1/system/login.go
   // 注释掉新代码，恢复旧代码测试
   ```

2. **联系主平台开发者**
   - 确认 API 端点是否正确
   - 确认 Token 格式要求
   - 确认认证流程

3. **抓包分析**
   ```bash
   # 使用 tcpdump 或 Wireshark 抓包
   tcpdump -i any host 192.168.3.101 and port 3088 -w oauth2.pcap
   ```

---

**创建时间**：2026-03-13  
**状态**：待排查  
**优先级**：高
