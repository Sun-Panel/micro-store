# OAuth2 模块集成 - 最终步骤

## ✅ 已完成的工作

### 1. 创建了 OAuth2 模块
- ✅ `/pkg/oauth2/` - 完整的 OAuth2 实现
  - `common/` - 公共类型
  - `server/` - 服务端实现
  - `client/` - 客户端实现
  - 所有测试通过 ✅

### 2. 创建了适配器
- ✅ `/service/adapter/oauth2_adapter.go` - 客户端适配器
- ✅ `/service/adapter/INTEGRATION_GUIDE.md` - 集成指南

### 3. 修改了现有代码
- ✅ `/service/biz/sunStore.go` - 添加了新方法
  - `GetOAuth2Client()` - 获取 OAuth2 客户端
  - `GetAccessTokenByCode()` - 授权码模式
  - `GetClientCredentialsToken()` - 客户端凭证模式
  - `GetAPIClient()` - 获取 API 客户端

- ✅ `/service/biz/sunStore_oauth2.go` - OAuth2 客户端封装

- ✅ `/service/api/api_v1/system/login.go` - 添加了迁移注释

## ⚠️ 最后一步：处理模块路径

由于项目结构特殊（根目录和 service 目录各有一个 go.mod），需要将 `pkg/oauth2` 复制到 `service/pkg/oauth2`。

### 方式一：手动复制（推荐）

```bash
# 在项目根目录执行
cd /Users/sunjingliang/my_code/sun-panel_group/sun-panel-micro-store/service
mkdir -p pkg/oauth2/common pkg/oauth2/server pkg/oauth2/client

# 复制文件
cp -r ../pkg/oauth2/common/*.go pkg/oauth2/common/
cp -r ../pkg/oauth2/server/*.go pkg/oauth2/server/
cp -r ../pkg/oauth2/client/*.go pkg/oauth2/client/

# 复制 go.mod
cp ../pkg/oauth2/go.mod pkg/oauth2/
```

### 方式二：修改项目结构（更好）

**建议长期方案：** 将根目录的 `pkg/oauth2` 移动到 `service/pkg/oauth2`，这样所有代码都在同一个模块下。

```bash
# 移动目录
mv /Users/sunjingliang/my_code/sun-panel_group/sun-panel-micro-store/pkg/oauth2 \
   /Users/sunjingliang/my_code/sun-panel_group/sun-panel-micro-store/service/pkg/

# 更新 go.mod 模块名（已在文件中修改）
# module sun-panel/pkg/oauth2
```

## 🚀 完成集成后如何使用

### 1. 使用适配器（最简单）

```go
import "sun-panel/adapter"

// 获取客户端
client, err := adapter.GetSunStoreClient()

// 使用授权码获取 Token
token, err := client.GetAccessToken(code)
```

### 2. 直接使用新方法（推荐）

```go
import "sun-panel/biz"

// 使用授权码获取 Token
tokenResp, err := biz.SunStore.GetAccessTokenByCode(code)

// 使用客户端凭证模式
tokenResp, err := biz.SunStore.GetClientCredentialsToken()
```

### 3. 完整的 OAuth2 流程

```go
// 1. 获取授权 URL
authURL := biz.SunStore.GetOAuth2Client().GetAuthorizationURL(redirectURI, state)

// 2. 用户授权后，使用 code 换取 token
tokenResp, err := biz.SunStore.GetAccessTokenByCode(code)

// 3. 使用 token 调用 API
apiClient := biz.SunStore.GetAPIClient()
var userInfo UserInfo
err := apiClient.Get(ctx, "/api/v1/user/info", tokenResp.AccessToken, &userInfo)

// 4. Token 过期后刷新
newToken, err := biz.SunStore.GetOAuth2Client().RefreshAccessToken(ctx, refreshToken)
```

## 📝 需要修改的文件清单

### 必须修改的文件

1. **`service/api/api_v1/system/login.go`**
   - 第 147-161 行：替换旧的 `sunapi.GetAccessToken()` 为新方法
   - 参考 INTEGRATION_GUIDE.md 中的示例

### 可选修改的文件

1. **`service/biz/sunStore.go`**
   - 可以删除旧方法 `GetClientApiToken()`, `ClientApiAuthLogin()`
   - 保留 `GetMainPlatformUserInfo()` (仍被使用)

2. **`service/lib/sunStore/oAuth2.go`**
   - 迁移完成后可以删除或标记为废弃

## 🎯 下一步操作

### 立即可做

1. **复制模块文件**
   ```bash
   cd service
   mkdir -p pkg/oauth2
   cp -r ../pkg/oauth2/* pkg/oauth2/
   ```

2. **编译测试**
   ```bash
   cd service
   go build ./biz
   go build ./adapter
   ```

3. **运行测试**
   ```bash
   cd pkg/oauth2
   go test ./...
   ```

### 迁移登录代码

1. 打开 `service/api/api_v1/system/login.go`
2. 找到 `OAuth2CodeLogin` 方法（第 140 行）
3. 按照注释替换为新方法
4. 测试登录流程

## 🔍 验证清单

- [ ] 复制 oauth2 模块到 service/pkg
- [ ] 编译通过：`go build ./biz`
- [ ] 编译通过：`go build ./adapter`
- [ ] 测试通过：`cd pkg/oauth2 && go test ./...`
- [ ] 修改登录代码使用新模块
- [ ] 测试 OAuth2 登录流程
- [ ] 测试 Token 刷新
- [ ] 测试 API 调用

## 📚 相关文档

### 已创建的文档
- `/pkg/oauth2/README.md` - 模块说明
- `/pkg/oauth2/USAGE.md` - 使用指南
- `/pkg/oauth2/MIGRATION.md` - 迁移指南
- `/pkg/oauth2/TEST_REPORT.md` - 测试报告
- `/service/adapter/INTEGRATION_GUIDE.md` - 集成指南（本文档）

### 快速链接
```bash
# 查看模块 README
cat /Users/sunjingliang/my_code/sun-panel_group/sun-panel-micro-store/pkg/oauth2/README.md

# 查看集成指南
cat /Users/sunjingliang/my_code/sun-panel_group/sun-panel-micro-store/service/adapter/INTEGRATION_GUIDE.md

# 查看测试报告
cat /Users/sunjingliang/my_code/sun-panel_group/sun-panel-micro-store/pkg/oauth2/TEST_REPORT.md
```

## 💡 提示

1. **暂时保留旧代码** - 迁移期间保留旧的实现作为备份
2. **逐步测试** - 每次修改后都要测试 OAuth2 流程
3. **监控日志** - 注意观察 Token 相关的错误日志
4. **向后兼容** - 新方法可以与旧方法共存

## ⚡ 快速开始

```bash
# 1. 复制模块
cd /Users/sunjingliang/my_code/sun-panel_group/sun-panel-micro-store/service
mkdir -p pkg
cp -r ../pkg/oauth2 pkg/

# 2. 编译检查
go build ./biz
go build ./adapter

# 3. 查看集成指南
cat adapter/INTEGRATION_GUIDE.md
```

完成这些步骤后，你就可以开始使用新的 OAuth2 模块了！🎉
