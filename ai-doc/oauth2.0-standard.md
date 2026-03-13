# OAuth 2.0 标准协议文档

## 1. 概述

OAuth 2.0 是一个行业标准的授权协议，允许第三方应用在用户授权下，获取有限的访问权限，而无需暴露用户的凭据（如用户名和密码）。

### 1.1 核心概念

- **Resource Owner（资源所有者）**：用户，能够授权访问受保护资源的实体
- **Client（客户端）**：第三方应用，请求访问受保护资源的应用
- **Authorization Server（授权服务器）**：认证服务器，颁发访问令牌
- **Resource Server（资源服务器）**：存储受保护资源的服务器，使用访问令牌响应请求
- **Access Token（访问令牌）**：用于访问受保护资源的凭证
- **Refresh Token（刷新令牌）**：用于获取新的访问令牌

### 1.2 端点说明

- **Authorization Endpoint（授权端点）**：`/oauth2/authorize`
- **Token Endpoint（令牌端点）**：`/oauth2/token`
- **Refresh Token Endpoint（刷新令牌端点）**：`/oauth2/token`

---

## 2. 授权方式

OAuth 2.0 定义了四种授权方式：

1. **授权码授权模式（Authorization Code Grant）** - 最安全、最常用
2. **隐式授权模式（Implicit Grant）** - 已不推荐使用
3. **密码凭证授权模式（Resource Owner Password Credentials Grant）**
4. **客户端凭证授权模式（Client Credentials Grant）**
5. **刷新令牌模式（Refresh Token Grant）**

---

## 3. 授权码授权模式（Authorization Code Grant）

### 3.1 适用场景

适用于服务器端应用（Web 应用），是最安全、最推荐的授权方式。

### 3.2 流程说明

1. 客户端引导用户到授权服务器
2. 用户授权后，授权服务器返回授权码（Authorization Code）
3. 客户端使用授权码换取访问令牌

### 3.3 请求参数

#### 步骤一：获取授权码

**请求方式**：GET 或 POST  
**端点**：`/oauth2/authorize`

| 参数名 | 必填 | 说明 |
|--------|------|------|
| response_type | 是 | 固定值：`code` |
| client_id | 是 | 客户端ID |
| redirect_uri | 是 | 回调地址，需与注册时一致 |
| scope | 否 | 权限范围，多个用空格分隔 |
| state | 推荐 | 随机字符串，用于防止CSRF攻击 |
| code_challenge | 条件必填 | PKCE扩展，挑战码 |
| code_challenge_method | 条件 | 加密方法：`plain` 或 `S256` |

**请求示例：**

```
GET /oauth2/authorize?
    response_type=code&
    client_id=s6BhdRkqt3&
    redirect_uri=https%3A%2F%2Fclient.example.com%2Fcallback&
    scope=read%20write&
    state=xyzABC123&
    code_challenge=E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM&
    code_challenge_method=S256
Host: authorization-server.com
```

**响应示例：**

```
HTTP/1.1 302 Found
Location: https://client.example.com/callback?
    code=SplxlOBeZQQYbYS6WxSbIA&
    state=xyzABC123
```

#### 步骤二：使用授权码换取令牌

**请求方式**：POST  
**端点**：`/oauth2/token`

| 参数名 | 必填 | 说明 |
|--------|------|------|
| grant_type | 是 | 固定值：`authorization_code` |
| code | 是 | 授权码 |
| redirect_uri | 是 | 必须与获取授权码时一致 |
| client_id | 是 | 客户端ID |
| client_secret | 条件必填 | 客户端密钥 |
| code_verifier | 条件必填 | PKCE扩展，验证码 |

**请求示例：**

```
POST /oauth2/token HTTP/1.1
Host: authorization-server.com
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&
code=SplxlOBeZQQYbYS6WxSbIA&
redirect_uri=https%3A%2F%2Fclient.example.com%2Fcallback&
client_id=s6BhdRkqt3&
client_secret=gX1fQt3QnG&
code_verifier=dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk
```

**响应示例：**

```json
{
  "access_token": "2YotnFZFEjr1zCsicMWpAA",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "tGzv3JOkF0XG5Qx2TlKWIA",
  "scope": "read write"
}
```

---

## 4. 隐式授权模式（Implicit Grant）

### 4.1 适用场景

适用于单页应用（SPA）或纯前端应用，但**不推荐使用**，建议改用授权码模式+PKCE。

### 4.2 流程说明

直接在授权端点返回访问令牌，不经过授权码中间步骤。

### 4.3 请求参数

**请求方式**：GET  
**端点**：`/oauth2/authorize`

| 参数名 | 必填 | 说明 |
|--------|------|------|
| response_type | 是 | 固定值：`token` |
| client_id | 是 | 客户端ID |
| redirect_uri | 是 | 回调地址 |
| scope | 否 | 权限范围 |
| state | 推荐 | 随机字符串，防止CSRF攻击 |

**请求示例：**

```
GET /oauth2/authorize?
    response_type=token&
    client_id=s6BhdRkqt3&
    redirect_uri=https%3A%2F%2Fclient.example.com%2Fcallback&
    scope=read%20write&
    state=xyzABC123
Host: authorization-server.com
```

**响应示例：**

```
HTTP/1.1 302 Found
Location: https://client.example.com/callback#
    access_token=2YotnFZFEjr1zCsicMWpAA&
    token_type=Bearer&
    expires_in=3600&
    scope=read%20write&
    state=xyzABC123
```

---

## 5. 密码凭证授权模式（Resource Owner Password Credentials Grant）

### 5.1 适用场景

仅适用于高度信任的应用（如官方应用），**不推荐用于第三方应用**。

### 5.2 流程说明

直接使用用户的用户名和密码获取访问令牌。

### 5.3 请求参数

**请求方式**：POST  
**端点**：`/oauth2/token`

| 参数名 | 必填 | 说明 |
|--------|------|------|
| grant_type | 是 | 固定值：`password` |
| username | 是 | 用户名 |
| password | 是 | 用户密码 |
| scope | 否 | 权限范围 |
| client_id | 条件必填 | 客户端ID（公客户端可不传） |
| client_secret | 条件必填 | 客户端密钥 |

**请求示例：**

```
POST /oauth2/token HTTP/1.1
Host: authorization-server.com
Content-Type: application/x-www-form-urlencoded
Authorization: Basic czZCaGRSa3F0Mzo3RjdzQjNfSzNkX2VfTGl2YV9tSm1B

grant_type=password&
username=johndoe&
password=A3ddj3w&
scope=read%20write
```

**响应示例：**

```json
{
  "access_token": "2YotnFZFEjr1zCsicMWpAA",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "tGzv3JOkF0XG5Qx2TlKWIA",
  "scope": "read write"
}
```

---

## 6. 客户端凭证授权模式（Client Credentials Grant）

### 6.1 适用场景

适用于机器对机器（M2M）通信，不涉及用户授权，客户端以其自身名义请求访问。

### 6.2 流程说明

使用客户端凭据直接获取访问令牌。

### 6.3 请求参数

**请求方式**：POST  
**端点**：`/oauth2/token`

| 参数名 | 必填 | 说明 |
|--------|------|------|
| grant_type | 是 | 固定值：`client_credentials` |
| scope | 否 | 权限范围 |
| client_id | 是 | 客户端ID |
| client_secret | 是 | 客户端密钥 |

**请求示例：**

```
POST /oauth2/token HTTP/1.1
Host: authorization-server.com
Content-Type: application/x-www-form-urlencoded
Authorization: Basic czZCaGRSa3F0Mzo3RjdzQjNfSzNkX2VfTGl2YV9tSm1B

grant_type=client_credentials&
scope=read%20write
```

或使用请求体传递客户端凭据：

```
POST /oauth2/token HTTP/1.1
Host: authorization-server.com
Content-Type: application/x-www-form-urlencoded

grant_type=client_credentials&
client_id=s6BhdRkqt3&
client_secret=7F7wB3_K3d_e_Liva_mJmA&
scope=read%20write
```

**响应示例：**

```json
{
  "access_token": "2YotnFZFEjr1zCsicMWpAA",
  "token_type": "Bearer",
  "expires_in": 3600,
  "scope": "read write"
}
```

---

## 7. 刷新令牌模式（Refresh Token Grant）

### 7.1 适用场景

当访问令牌过期时，使用刷新令牌获取新的访问令牌，避免用户重新授权。

### 7.2 流程说明

使用刷新令牌换取新的访问令牌。

### 7.3 请求参数

**请求方式**：POST  
**端点**：`/oauth2/token`

| 参数名 | 必填 | 说明 |
|--------|------|------|
| grant_type | 是 | 固定值：`refresh_token` |
| refresh_token | 是 | 刷新令牌 |
| scope | 否 | 权限范围，不能超出初始授权范围 |
| client_id | 是 | 客户端ID |
| client_secret | 条件必填 | 客户端密钥 |

**请求示例：**

```
POST /oauth2/token HTTP/1.1
Host: authorization-server.com
Content-Type: application/x-www-form-urlencoded
Authorization: Basic czZCaGRSa3F0Mzo3RjdzQjNfSzNkX2VfTGl2YV9tSm1B

grant_type=refresh_token&
refresh_token=tGzv3JOkF0XG5Qx2TlKWIA&
scope=read
```

**响应示例：**

```json
{
  "access_token": "2YotnFZFEjr1zCsicMWpAA",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "tGzv3JOkF0XG5Qx2TlKWIA",
  "scope": "read"
}
```

---

## 8. 使用访问令牌访问资源

### 8.1 请求方式

客户端在请求资源服务器时，需要在请求头中携带访问令牌。

**请求示例：**

```
GET /api/userinfo HTTP/1.1
Host: resource-server.com
Authorization: Bearer 2YotnFZFEjr1zCsicMWpAA
```

### 8.2 响应示例

```json
{
  "sub": "248289761001",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "picture": "https://example.com/photo.jpg"
}
```

---

## 9. 错误响应

### 9.1 授权错误

当授权请求失败时，授权服务器会返回错误信息。

**错误参数：**

| 参数名 | 说明 |
|--------|------|
| error | 错误代码 |
| error_description | 错误描述（可选） |
| error_uri | 错误详细说明页面（可选） |
| state | 原始请求的state值 |

**常见错误代码：**

| 错误代码 | 说明 |
|----------|------|
| invalid_request | 请求缺少必需参数、包含无效参数或格式错误 |
| unauthorized_client | 客户端未授权使用此授权方式 |
| access_denied | 用户或授权服务器拒绝请求 |
| unsupported_response_type | 不支持的响应类型 |
| invalid_scope | 请求的权限范围无效 |
| server_error | 授权服务器遇到意外错误 |
| temporarily_unavailable | 授权服务器暂时不可用 |

**错误响应示例：**

```
HTTP/1.1 302 Found
Location: https://client.example.com/callback?
    error=access_denied&
    error_description=The%20user%20denied%20access%20to%20your%20application&
    state=xyzABC123
```

### 9.2 令牌错误

当令牌请求失败时，返回 JSON 格式的错误信息。

**常见错误代码：**

| 错误代码 | 说明 |
|----------|------|
| invalid_request | 请求缺少必需参数或格式错误 |
| invalid_client | 客户端认证失败 |
| invalid_grant | 授权码或刷新令牌无效或过期 |
| unauthorized_client | 客户端未授权使用此授权方式 |
| unsupported_grant_type | 不支持的授权类型 |
| invalid_scope | 请求的权限范围无效 |

**错误响应示例：**

```json
{
  "error": "invalid_grant",
  "error_description": "The authorization code has expired",
  "error_uri": "https://authorization-server.com/errors/invalid_grant"
}
```

---

## 10. PKCE 扩展（Proof Key for Code Exchange）

### 10.1 概述

PKCE 是授权码模式的扩展，用于防止授权码被拦截攻击，特别适用于移动应用和单页应用。

### 10.2 流程说明

1. 客户端生成一个随机的 `code_verifier`
2. 使用 `code_verifier` 生成 `code_challenge`
3. 在授权请求中携带 `code_challenge`
4. 在令牌请求中携带 `code_verifier` 进行验证

### 10.3 参数说明

| 参数名 | 说明 |
|--------|------|
| code_verifier | 43-128位的随机字符串，使用 URL 安全字符 |
| code_challenge | 对 `code_verifier` 进行转换后的值 |
| code_challenge_method | 转换方法：`plain`（明文）或 `S256`（SHA-256，推荐） |

### 10.4 生成示例

**code_verifier 示例：**
```
dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk
```

**使用 S256 生成 code_challenge：**
```
code_challenge = BASE64URL(SHA256(ASCII(code_verifier)))
```

**结果：**
```
E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM
```

---

## 11. 安全最佳实践

### 11.1 客户端安全

- **机密客户端**：能够安全存储 `client_secret` 的服务器端应用
  - 使用授权码模式
  - 使用 HTTP Basic 认证或请求体传递客户端凭据

- **公开客户端**：无法安全存储密钥的应用（移动应用、SPA）
  - 使用授权码模式 + PKCE
  - 不要使用隐式授权模式

### 11.2 State 参数

- 每次授权请求都应包含 `state` 参数
- `state` 应为随机生成的不可预测字符串
- 验证回调中的 `state` 是否与请求时一致，防止 CSRF 攻击

### 11.3 令牌安全

- 访问令牌应设置合理的过期时间（建议 1-2 小时）
- 刷新令牌应长期有效或设置较长过期时间
- 刷新令牌应进行安全存储，建议绑定客户端和用户
- 使用 HTTPS 传输所有请求

### 11.4 Redirect URI 验证

- 必须精确匹配注册的 `redirect_uri`
- 不允许使用通配符
- 不允许使用 HTTP（生产环境）

### 11.5 Scope 管理

- 定义清晰的权限范围
- 遵循最小权限原则
- 记录用户授权的权限范围

---

## 12. 客户端认证方式

### 12.1 HTTP Basic 认证

将 `client_id` 和 `client_secret` 进行 Base64 编码后放在 Authorization 请求头中。

**格式：**
```
Authorization: Basic Base64(client_id:client_secret)
```

**示例：**
```
Authorization: Basic czZCaGRSa3F0Mzo3RjdzQjNfSzNkX2VfTGl2YV9tSm1B
```

### 12.2 请求体认证

将 `client_id` 和 `client_secret` 放在请求体中（不推荐，安全性较低）。

**示例：**
```
POST /oauth2/token HTTP/1.1
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&
code=SplxlOBeZQQYbYS6WxSbIA&
client_id=s6BhdRkqt3&
client_secret=gX1fQt3QnG
```

---

## 13. 常见问题

### Q1: 授权码和访问令牌有什么区别？

- **授权码**：临时的、一次性的中间凭证，用于换取访问令牌，生命周期短（通常几分钟）
- **访问令牌**：用于访问受保护资源的凭证，有效期较长（通常几小时）

### Q2: 什么时候使用 PKCE？

- 移动应用
- 单页应用（SPA）
- 无法安全存储 `client_secret` 的公开客户端
- 建议所有授权码流程都使用 PKCE

### Q3: 隐式授权模式为什么被弃用？

- 访问令牌暴露在 URL fragment 中，容易被窃取
- 无法使用刷新令牌
- 授权码模式 + PKCE 提供了更好的安全性

### Q4: 如何处理令牌过期？

1. 在请求资源时收到 401 错误
2. 使用刷新令牌获取新的访问令牌
3. 重试原请求
4. 如果刷新令牌也过期，引导用户重新授权

---

## 14. 参考资料

- [RFC 6749 - The OAuth 2.0 Authorization Framework](https://tools.ietf.org/html/rfc6749)
- [RFC 7636 - Proof Key for Code Exchange (PKCE)](https://tools.ietf.org/html/rfc7636)
- [RFC 6750 - The OAuth 2.0 Authorization Framework: Bearer Token Usage](https://tools.ietf.org/html/rfc6750)
- [RFC 6819 - OAuth 2.0 Threat Model and Security Considerations](https://tools.ietf.org/html/rfc6819)
- [OAuth 2.0 Security Best Current Practice](https://tools.ietf.org/html/draft-ietf-oauth-security-topics)

---

**文档版本**: 1.0  
**最后更新**: 2026-03-13
