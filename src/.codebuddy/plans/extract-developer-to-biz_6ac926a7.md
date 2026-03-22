---
name: extract-developer-to-biz
overview: 将 micro_app_developer.go 中的业务逻辑提取到 biz 层，创建 MicroAppDeveloperService，使 API 层只负责参数绑定和响应返回。
todos:
  - id: create-biz-errors
    content: 在 biz/errors.go 新增 5 个微应用开发者错误码常量并注册到 validCodes
    status: completed
  - id: create-biz-service
    content: 创建 biz/micro_app_developer.go，实现 MicroAppDeveloperService 全部 7 个业务方法
    status: completed
    dependencies:
      - create-biz-errors
  - id: refactor-handler
    content: 重构 micro_app_developer.go Handler 为薄壳模式，调用 biz 层
    status: completed
    dependencies:
      - create-biz-service
  - id: update-error-mapping
    content: 更新 developer_version.go 的 bizCodeToInt 映射表，追加 3000-3003 错误码
    status: completed
    dependencies:
      - create-biz-service
---

## 用户需求

将 `micro_app_developer.go` 中直接写在 API Handler 里的业务逻辑（事务处理、状态检查、权限验证、数据组装、多语言合并）提取到 biz 层，遵循项目既有的分层架构规范。

## 产品概述

对后端代码进行分层重构，不改变 API 接口行为，仅将业务逻辑从 API 层迁移到 biz 层。

## 核心功能

- 创建 `biz/micro_app_developer.go`，定义 `MicroAppDeveloperService` 结构体及便捷函数
- 新增微应用开发者相关的 biz 错误码（APP_ID_EXISTS, NO_PERMISSION, PENDING_REVIEW_EXISTS, NO_PENDING_REVIEW）
- 在 `biz/errors.go` 的 `IsBizError` validCodes 中注册新错误码
- 将 Create/Update/UpdateLang/CancelReview/GetReviewHistory 的事务和业务逻辑下沉到 biz 层
- 重构 `micro_app_developer.go` API Handler 为"参数绑定 + 调用 biz + 返回响应"的薄壳模式
- 更新 `developer_version.go` 中的 `bizCodeToInt` 映射表和 `ErrorCode.go` 常量

## 技术栈

- Go (Gin + GORM)，现有后端架构
- 分层模式：API 层 → Biz 层 → Models 层

## 实现方案

参照 `biz/micro_app_version.go` 已有的 Service 结构体 + 便捷函数模式，为开发者微应用操作创建对应的 biz 服务。biz 层接收基础类型参数（uint、string 等），不依赖 Gin Context 或 HTTP DTO，确保纯业务逻辑可测试。

### biz 层方法设计

| Handler 方法 | biz 方法 | 职责 |
| --- | --- | --- |
| GetMyList | `GetDeveloperAppList(db, page, limit, status, categoryId, developerId, keyword)` | 查询开发者应用列表 |
| GetMyInfo | `GetDeveloperAppInfo(db, appId, developerId)` | 权限验证 + 查询详情 + 作者名组装 |
| Create | `CreateApp(db, microAppId, appName, appIcon, appDesc, remark, categoryId, chargeType, price, developerId, screenshots, langMap)` | ID 生成、唯一性检查、事务创建 |
| Update | `SubmitAppUpdate(db, id, developerId, appName, appIcon, appDesc, remark, categoryId, chargeType, price, screenshots, langMap)` | 权限验证、待审核检查、多语言合并、创建审核快照 |
| UpdateLang | `SubmitLangUpdate(db, id, developerId, langMap)` | 权限验证、待审核检查、语言合并、创建审核快照 |
| CancelReview | `CancelAppReview(db, id, developerId)` | 权限验证、事务删除审核记录 |
| GetReviewHistory | `GetAppReviewHistory(db, appId, developerId, page, limit)` | 权限验证、分页查询 |


### biz 输入参数设计（脱离 HTTP DTO）

- `langMap` 使用 `map[string]struct { AppName, AppDesc string }` 独立定义在 biz 包中
- biz 层的 `generateMicroAppId` 函数移入 biz 包
- biz 层只返回 `(data, error)` 或 `(data, total, error)`，不依赖 gin.H

### 错误码设计

新增 5 个微应用开发者相关错误码，数字编号从 3000 起始（避免与版本模块 2000-2007 冲突）：

- `APP_ID_EXISTS` (3000) - 应用ID已存在
- `NO_PERMISSION` (3001) - 无权操作此应用
- `PENDING_REVIEW_EXISTS` (3002) - 已有待审核记录
- `NO_PENDING_REVIEW` (3003) - 没有待审核记录
- `APP_NOT_FOUND` 复用已有的 ErrCodeAppNotFound (2000)

## 实现注意事项

- `handleBizError` 和 `bizCodeToInt` 已在 `developer_version.go` 中定义，新的错误码映射需追加到该 map 中
- `ErrorCode.go`（如果存在）也需追加对应的数字常量
- API Handler 重构后 import 需添加 `sun-panel/biz`，移除不再需要的 `crypto/rand`、`encoding/hex`、`encoding/json`、`time` 等依赖
- `MicroAppLangInfo` DTO 当前在 `admin` 包中定义，biz 层需要独立定义自己的 LangInput 结构体，避免循环依赖
- `isBizError` 的 validCodes 切片需要添加新的错误码常量

## 目录结构

```
service/
├── biz/
│   ├── errors.go                    # [MODIFY] 新增 5 个错误码常量 + 注册到 validCodes
│   └── micro_app_developer.go       # [NEW] MicroAppDeveloperService + 便捷函数
├── api/api_v1/admin/
│   ├── micro_app_developer.go       # [MODIFY] Handler 瘦身为薄壳模式
│   └── developer_version.go         # [MODIFY] bizCodeToInt 追加新错误码映射
```