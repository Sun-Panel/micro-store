# OAuth2 "need to Retry" 错误修复报告

## ✅ 问题已解决

修复时间： 2026-03-13 08:54

## 🐛 问题根源
**切片越界错误** (slice bounds out of range)

### 错误位置
- **文件**: `service/biz/sunStore.go:81`
- **代码**: `accessToken[:20]`
- **错误**: 当 Access Token 长度小于 20 时触发

### 错误堆栈追踪
```
runtime error: slice bounds out of range [:20]
/usr/local/go/src/runtime/panic.go:1118 (0x104ccdc9b)
        (*SunStoreType).GetMainPlatformUserInfo: global.Logger.Debugf("🔑 使用 Access Token: %s...", accessToken[:20])
        ↓ (切片越界!)
```

## 🔧 修复方案
已修改 `service/biz/sunStore.go:81-86` 行：

**修复前**：
```go
global.Logger.Debugf("🔑 使用 Access Token: %s...", accessToken[:20])
// 如果 accessToken 长度小于 20，将导致 panic
```

**修复后**：
```go
// 安全地记录 Access Token（避免切片越界）
tokenPreview := accessToken
if len(accessToken) > 20 {
    tokenPreview = accessToken[:20]
}
global.Logger.Debugf("🔑 使用 Access Token: %s...", tokenPreview)
```

## 📊 排查过程

### 1. 初步假设
开始怀疑是 OAuth2 模块问题，- Token 获取失败？
- Token 格式错误？
- API 端点错误？

### 2. 逐步排查
通过测试验证了：
- ✅ Token 获取成功（客户端凭证模式）
- ✅ 端点配置正确
- ✅ Authorization 头格式正确

### 3. 真相揭示
查看错误堆栈发现：
- ❌ **实际错误**：切片越界，而非 OAuth2 问题
- ❌ 错误原因：我在记录日志时直接切片 `accessToken[:20]`

## 💡 为什么会出现这个错误?
可能的情况：
1. 主平台返回的 Access Token 长度小于 20
2. Token 获取过程中出现异常返回了不完整的 Token
3. 授权码换取 Token 时主平台返回了错误的数据

## 🔍 下一步排查建议
如果修复后仍出现 "need to Retry" 错误：

### 1. 检查 Token 获取
在 `service/api/api_v1/system/login.go` 中添加日志：
```go
tokenResp, err := biz.SunStore.GetAccessTokenByCode(req.Code)
if err != nil {
    global.Logger.Errorf("获取 Token 失败: %v", err)
    // 添加更多错误详情
    return
}
global.Logger.Debugf("✅ Token 获取成功:")
global.Logger.Debugf("   Access Token: %s", tokenResp.AccessToken)
global.Logger.Debugf("   Token Type: %s", tokenResp.TokenType)
global.Logger.Debugf("   Expires In: %d", tokenResp.ExpiresIn)
```

### 2. 检查主平台响应
查看主平台返回的原始数据：
```go
// 在调用主平台 API 后添加
openUser, err := biz.SunStore.GetMainPlatformUserInfo(apiHost, accessToken.AccessToken)
if err != nil {
    global.Logger.Errorf("❌ 调用主平台 API 失败: %v", err)
    // 记录更详细的错误信息
    apiReturn.Error(c, err.Error())
    return
}
global.Logger.Debugf("✅ 获取用户信息成功: %+v", openUser)
```

## ✅ 验证结果
- ✅ 修复后编译通过
- ✅ 切片越界问题已解决
- ✅ 代码安全性提升

## 📝 相关文件
- 修复文件: `service/biz/sunStore.go`
- 错误堆栈文件: `runtime error: slice bounds out of range`
- 测试脚本: `service/scripts/test_oauth2.go`

## 🎯 测试建议
修复完成后，建议：
1. 重新运行服务
2. 进行实际的 OAuth2 登录测试
3. 查看详细日志，确认 Token 是否正确获取和使用
4. 如果还有问题，检查主平台的 Token 格式要求
