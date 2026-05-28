# Sun Panel Micro Store - 客户端 API 协议文档

## 1. 认证机制

### 请求头

| Header | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `Authorization` | string | 是 | `Bearer <JWT Token>` |
| `Client-App-Version` | string | 是 | 客户端版本号，如 `1.0.0` |

### 认证流程

1. 客户端在每次请求中携带 JWT Token 和版本号
2. 服务端根据版本号获取对应的密钥
3. 使用 HS256 算法验证 JWT Token
4. 验证通过后，将 `userID` 和 `jwtClaims` 注入请求上下文

---

## 2. 响应格式

### 成功响应

```json
{
  "code": 0,
  "msg": "OK",
  "data": { ... },
  "responseId": "2026052620010012345"
}
```

### 列表响应

```json
{
  "code": 0,
  "msg": "OK",
  "data": {
    "list": [ ... ],
    "count": 100
  },
  "responseId": "2026052620010012345"
}
```

### 错误响应

```json
{
  "code": -1,
  "msg": "错误描述",
  "responseId": "2026052620010012345"
}
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | int | 业务码，`0` 成功，非 `0` 失败 |
| `msg` | string | 消息 |
| `data` | object/array/null | 数据 |
| `responseId` | string | 请求追踪ID |

---

## 3. 错误码

| 错误码 | HTTP状态码 | 含义 |
|--------|-----------|------|
| `0` | 200 | 成功 |
| `-1` | 200 | 通用错误 |
| `1001` | 401 | Token 过期 |
| `1002` | 401 | Token 无效 |
| `1005` | 200 | 无权限 |
| `1200` | 200 | 数据库错误 |

---

## 4. 请求示例

```bash
curl -X GET "https://domain/api/v1/xxx" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Client-App-Version: 1.0.0" \
  -H "Content-Type: application/json"
```

---

## 5. 接口列表

### 接口权限说明
jwt是否强制需要包含登录用户
| 接口 | 是否需要登录 |
|------|-------------|
| `POST developer/batchGetInfo` | 否 |
| `POST microApp/batchGetInfo` | 否 |
| `POST developer/checkIsDeveloper` | **都可以** |
| `POST microApp/installRecord` | **是** |

---

### 5.1 批量获取作者昵称

**请求**

```
POST /api/clientApp/developer/batchGetInfo
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `developerNames` | string[] | 是 | 开发者标识数组，最多100个 |

```json
{
  "developerNames": ["devname1", "devname2", "devname3"]
}
```

**响应 data**

| 参数 | 类型 | 说明 |
|------|------|------|
| `developerInfos` | map\<string, DeveloperInfo\> | key 为 developerName |

**DeveloperInfo 结构**

| 参数 | 类型 | 说明 |
|------|------|------|
| `developerName` | string | 开发者标识 |
| `name` | string | 开发者昵称 |

**响应示例**

```json
{
  "code": 0,
  "msg": "OK",
  "data": {
    "developerInfos": {
      "devname1": {
        "developerName": "devname1",
        "name": "开发者一"
      },
      "devname2": {
        "developerName": "devname2",
        "name": "开发者二"
      }
    }
  },
  "responseId": "2026052620010012345"
}
```

**说明**
- 请求中的 `developerNames` 会自动去重
- 不存在的开发者标识不会出现在响应中（不会报错）
- 最多支持查询 100 个开发者

---

### 5.2 微应用安装触发记录（批量）

**请求**

```
POST /api/clientApp/microApp/installRecord
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `installList` | InstallRecord[] | 是 | 安装记录数组，最多100条 |

**InstallRecord 结构**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `microAppId` | string | 是 | 微应用唯一标识 |
| `microAppVersion` | string | 是 | 微应用版本号 |
| `installTimestamp` | int64 | 是 | 安装时间戳 |

```json
{
  "installList": [
    {
      "microAppId": "app-xxx",
      "microAppVersion": "1.0.0",
      "installTimestamp": 1685200000
    },
    {
      "microAppId": "app-yyy",
      "microAppVersion": "2.0.0",
      "installTimestamp": 1685200100
    }
  ]
}
```

**响应 data**

无（成功时返回空对象）

**说明**
- **需要登录**：通过 `LoginInterceptor` 中间件验证
- 用于记录用户安装微应用的行为
- 目前为预留接口，具体实现待完善

---

### 5.3 批量查询微应用信息

**请求**

```
POST /api/clientApp/microApp/batchGetInfo
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `microAppIds` | string[] | 是 | 微应用唯一标识数组，最多100个 |

```json
{
  "microAppIds": ["app-xxx", "app-yyy"]
}
```

**响应 data**

| 参数 | 类型 | 说明 |
|------|------|------|
| `microAppInfos` | map\<string, MicroAppInfo\> | key 为 microAppId，不存在的应用不返回 |

**MicroAppInfo 结构**

| 参数 | 类型 | 说明 |
|------|------|------|
| `microAppId` | string | 微应用唯一标识 |
| `appIcon` | string | 应用图标URL |
| `chargeType` | int | 收费方式：0-免费 1-积分 2-订阅PRO免费 |
| `points` | int | 价格（积分数值） |
| `developerName` | string | 开发者名称 |
| `developer` | string | 开发者标识 |

**响应示例**

```json
{
  "code": 0,
  "msg": "OK",
  "data": {
    "microAppInfos": {
      "app-xxx": {
        "microAppId": "app-xxx",
        "appIcon": "/uploads/icon.png",
        "chargeType": 0,
        "points": 0,
        "developerName": "开发者名称",
        "developer": "devname"
      }
    }
  },
  "responseId": "2026052620010012345"
}
```

**说明**
- 请求中的 `microAppIds` 会自动去重
- 不存在的应用不会出现在响应中（不会报错）
- 最多支持查询 100 个应用
- 单个查询传数组即可：`{ "microAppIds": ["app-xxx"] }`

---

### 5.4 查询当前账号是否为开发者

**请求**

```
POST /api/clientApp/developer/checkIsDeveloper
```

**请求参数**

无

**响应 data**

| 参数 | 类型 | 说明 |
|------|------|------|
| `isDeveloper` | bool | 是否为开发者 |

**响应示例**

```json
{
  "code": 0,
  "msg": "OK",
  "data": {
    "isDeveloper": true
  },
  "responseId": "2026052620010012345"
}
```

**说明**
- **需要登录**：通过 `LoginInterceptor` 中间件验证
- 未登录将返回 401 错误

---

## 6. 注意事项

1. **缺少 `Client-App-Version` 将返回 401**
2. **Token 过期 (1001) 需重新登录获取新 Token**
3. **`responseId` 用于问题排查，请记录**
4. **POST/PUT 请求使用 `Content-Type: application/json`**
5. **所有接口基础路径: `/api/clientApp/`**
