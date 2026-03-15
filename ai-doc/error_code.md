# 错误码说明文档

## 文档信息

- **源文件**: `service/api/api_v1/common/apiReturn/ErrorCode.go`
- **前端枚举**: `src/enums/errorCode/index.ts`
- **国际化文件**: `src/locales/zh-CN.json` 和 `src/locales/en-US.json`
- **更新日期**: 2026-03-15

---

## 错误码总览

### 1. 认证与权限类错误 (1000-1009)

| 错误码 | 常量名称 | 英文描述 | 中文说明 |
|--------|----------|----------|----------|
| 1000 | ErrCodeNotLoggedIn | not logged in yet | 还未登录 |
| 1003 | ErrCodeIncorrectUsernameOrPassword | incorrect username or password | 用户名或密码错误 |
| 1004 | ErrCodeAccountDisabledOrNotActivated | account disabled or not activated | 账号已停用或未激活 |
| 1005 | ErrCodeNoCurrentPermission | no current permission for operation | 当前无权限操作 |
| 1006 | ErrCodeAccountDoesNotExist | account does not exist | 账号不存在 |
| 1007 | ErrCodeOldPasswordError | old password error | 旧密码不正确 |
| 1008 | ErrCodeNoPROAuthorization | no PRO authorization | 没有PRO授权 |
| 1009 | ErrCodeCaptchaError | captcha error | 验证码错误 |

### 2. 数据操作类错误 (1200-1203)

| 错误码 | 常量名称 | 英文描述 | 中文说明 |
|--------|----------|----------|----------|
| 1200 | ErrCodeDatabaseError | database error | 数据库错误 |
| 1201 | ErrCodePleaseKeepAtLeastOne | please keep at least one | 请至少保留一个 |
| 1202 | ErrCodeNoDataRecordFound | no data record found | 未找到数据记录 |
| 1203 | ErrCodeDataAlreadyExists | data already exists | 数据已存在 |

### 3. 文件上传类错误 (1300-1301)

| 错误码 | 常量名称 | 英文描述 | 中文说明 |
|--------|----------|----------|----------|
| 1300 | ErrCodeUploadFailed | upload failed | 上传失败 |
| 1301 | ErrCodeUnsupportedFileFormat | unsupported file format | 不被支持的格式文件 |

### 4. 参数与业务类错误 (1400-1402)

| 错误码 | 常量名称 | 英文描述 | 中文说明 |
|--------|----------|----------|----------|
| 1400 | ErrCodeParameterFormatError | parameter format error | 参数格式错误 |
| 1401 | ErrCodeOrderCreateFailed | order create failed | 订单创建失败 |
| 1402 | ErrCodeGoodsNoUsePayPlatform | goods no use pay platform | 商品不支持支付平台 |

### 5. 微应用版本业务错误 (2000-2007)

| 错误码 | 常量名称 | 英文描述 | 中文说明 |
|--------|----------|----------|----------|
| 2000 | ErrCodeAppNotFound | app not found | 应用不存在 |
| 2001 | ErrCodeVersionNotFound | version not found | 版本不存在 |
| 2002 | ErrCodeVersionExists | version already exists | 版本号已存在 |
| 2003 | ErrCodeVersionCodeExists | version code already exists | 版本编号已存在 |
| 2004 | ErrCodeStatusNotAllowed | status not allowed | 状态不允许操作 |
| 2005 | ErrCodeApprovedCannotDelete | approved version cannot be deleted | 已审核版本不能删除 |
| 2006 | ErrCodeNotPendingReview | not pending review | 非待审核状态 |
| 2007 | ErrCodeNoUpdateContent | no update content | 无更新内容 |

---

## 使用说明

### 前端调用方式

```typescript
// 根据错误码获取错误信息
interface ApiResponse {
  code: number;
  message: string;
  data?: any;
}

// 错误处理示例
if (response.code !== 0) {
  // 根据 code 查找对应的中文说明进行友好提示
  const errorMessage = getErrorMessageByCode(response.code);
  showToast(errorMessage);
}
```

### 错误码范围说明

| 范围 | 分类 |
|------|------|
| 1000-1009 | 认证与权限类 |
| 1200-1203 | 数据操作类 |
| 1300-1301 | 文件上传类 |
| 1400-1402 | 参数与业务类 |
| 2000-2099 | 微应用版本业务 |

---

## 备注

- 当遇到未定义的错误码时，后端会返回 `unknown error`
- 前端可根据错误码范围快速定位错误类型
- 所有错误码定义以源文件为准，本文档仅供参考

---

## 前端使用示例

### 错误码枚举导入

```typescript
import { ErrorCode } from '@/enums/errorCode'
```

### 错误处理工具

```typescript
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'
```

### 使用示例

```typescript
// 示例1: 提交审核
async function handleSubmitReview(versionId: number) {
  try {
    const res = await submitReview({ versionId })
    if (res.code === 0) {
      message.success('已提交审核')
    }
    else {
      // 自动根据错误码显示对应的错误提示
      apiRespErrMsg(res)
    }
  }
  catch (error) {
    message.error('提交审核失败')
  }
}

// 示例2: 特定错误码处理
async function handleDeleteVersion(id: number) {
  try {
    const res = await deleteVersion([id])
    if (res.code === 0) {
      message.success('删除成功')
    }
    else if (res.code === ErrorCode.ApprovedCannotDelete) {
      message.error('已审核版本不能删除')
    }
    else {
      apiRespErrMsg(res)
    }
  }
  catch (error) {
    message.error('删除失败')
  }
}
```

### 前端文件修改记录

| 文件路径 | 修改内容 | 日期 |
|---------|---------|------|
| `src/enums/errorCode/index.ts` | 添加所有错误码枚举定义 | 2026-03-15 |
| `src/locales/zh-CN.json` | 添加错误码中文提示 | 2026-03-15 |
| `src/locales/en-US.json` | 添加错误码英文提示 | 2026-03-15 |
| `src/views/admin/myMicroApp/detail/index.vue` | 添加错误码判断和处理 | 2026-03-15 |
