# iframe跨页面通信模块

基于postMessage的iframe跨页面通信工具，提供类型安全、支持重发机制的通信API。

## 功能特性

✅ **类型安全** - 完整的TypeScript支持  
✅ **对称API** - 主页面和iframe使用相同API  
✅ **消息验证** - 验证消息来源，防止安全风险  
✅ **事件管理** - 支持多个监听器，自动清理  
✅ **请求-响应** - 支持异步请求-响应通信  
✅ **超时重发** - 可配置超时时间和重发策略  
✅ **错误处理** - 完善的错误处理机制  
✅ **调试模式** - 支持详细的调试日志  
✅ **安全密钥** - 支持postMessageKey增强安全性  

## 快速开始

### 1. 安装依赖

```bash
# 无需额外依赖，直接使用即可
```

### 2. 基本用法

#### 主页面端

```typescript
import { createIframeChannel } from '@/utils/iframeCommunication'

// 生成安全密钥（建议使用crypto.randomUUID()）
const postMessageKey = crypto.randomUUID()

// 创建通信通道
const channel = createIframeChannel({
  targetOrigin: '*', // 生产环境建议指定具体origin
  timeout: 5000,
  maxRetries: 3,
  retryDelay: 1000,
  debug: true,
  postMessageKey // 配置安全密钥
})

// 注册iframe窗口
const iframe = document.getElementById('my-iframe') as HTMLIFrameElement
channel.registerTarget('my-iframe', iframe.contentWindow!)

// 发送事件（不需要回复）
channel.send('my-iframe', 'user-login', { userId: 123 })

// 发送请求（需要回复）
try {
  const response = await channel.request('my-iframe', 'get-user-info', { userId: 123 }, {
    timeout: 10000,  // 覆盖默认超时
    maxRetries: 5    // 覆盖默认重发次数
  })
  console.log('用户信息:', response)
} catch (error) {
  console.error('请求失败:', error)
}

// 监听iframe消息
channel.on('iframe-message', (data) => {
  console.log('收到iframe消息:', data)
})

// 一次性监听
channel.once('ready', () => {
  console.log('iframe已就绪')
})

// 清理资源
channel.destroy()
```

#### iframe端

```typescript
import { createIframeChannel } from '@/utils/iframeCommunication'

// 从URL参数中获取安全密钥
const urlParams = new URLSearchParams(window.location.search)
const postMessageKey = urlParams.get('postMessageKey') || ''

// 创建通信通道（iframe模式）
const channel = createIframeChannel({
  isIframe: true,
  targetOrigin: window.location.origin,
  sourceId: 'my-iframe',
  debug: true,
  postMessageKey // 使用相同的安全密钥
})

// 注册主页面窗口
channel.registerTarget('parent', window.parent)

// 监听主页面消息
channel.on('user-login', (data) => {
  console.log('用户登录:', data)
})

// 响应主页面请求
channel.onRequest('get-user-info', async (data) => {
  // 模拟API调用
  const user = await fetchUser(data.userId)
  return user
})

// 向主页面发送事件
channel.send('parent', 'ready', { version: '1.0.0' })

// 向主页面发送请求
const config = await channel.request('parent', 'get-config', { appId: 'clock' })
console.log('配置:', config)

// 清理资源
channel.destroy()
```

### 3. 高级用法

#### 使用单例模式

```typescript
import { getOrCreateDefaultChannel, destroyDefaultChannel } from '@/utils/iframeCommunication'

// 获取或创建默认通道
const channel = getOrCreateDefaultChannel({
  targetOrigin: '*',
  debug: true
})

// 使用通道...

// 应用退出时销毁
destroyDefaultChannel()
```

#### 使用自定义重发策略

```typescript
const channel = createIframeChannel({
  targetOrigin: '*',
  timeout: 5000,
  maxRetries: 3,
  retryDelay: 1000,  // 基础延迟
  debug: true
})

// 发送请求，使用自定义配置
const response = await channel.request('iframe-id', 'critical-operation', data, {
  timeout: 15000,      // 15秒超时
  maxRetries: 10,     // 最多重试10次
  retryDelay: 2000    // 2秒基础延迟
})
```

#### 使用TypeScript泛型

```typescript
// 定义消息类型
interface UserInfo {
  id: number
  name: string
  email: string
}

// 发送请求，指定响应类型
const user = await channel.request<UserInfo>('iframe-id', 'get-user', { id: 123 })
console.log(user.name) // TypeScript知道user是UserInfo类型

// 监听事件，指定数据类型
channel.on<UserInfo>('user-updated', (user) => {
  console.log(user.email) // TypeScript知道user是UserInfo类型
})

// 注册请求处理器
channel.onRequest<{ id: number }, UserInfo>('get-user', async (data) => {
  const user = await fetchUser(data.id)
  return user
})
```

#### 使用快捷回复

```typescript
// iframe端使用快捷回复
channel.onQuickRequest<{ userId: number }, UserInfo>('get-user-quick', (data, ctx) => {
  // 直接使用ctx.reply回复
  ctx.reply({ id: data.userId, name: 'John', email: 'john@example.com' })
})

// 异步快捷回复
channel.onQuickRequest<{ userId: number }, UserInfo>('get-user-async', async (data, ctx) => {
  const user = await fetchUser(data.userId)
  ctx.reply(user)
})

// 错误回复
channel.onQuickRequest<{ userId: number }, UserInfo>('get-user-error', (data, ctx) => {
  if (!data.userId) {
    ctx.replyError('User ID is required')
    return
  }
  // 正常回复...
})
```

### 4. 安全配置

#### 使用postMessageKey增强安全性

```typescript
// 主页面生成随机密钥
const postMessageKey = crypto.randomUUID()

// 创建通信通道时配置密钥
const channel = createIframeChannel({
  targetOrigin: 'https://your-domain.com',
  postMessageKey // 配置安全密钥
})

// iframe需要从URL参数获取相同的密钥
const urlParams = new URLSearchParams(window.location.search)
const iframePostMessageKey = urlParams.get('postMessageKey') || ''

const iframeChannel = createIframeChannel({
  isIframe: true,
  targetOrigin: window.location.origin,
  postMessageKey: iframePostMessageKey // 使用相同的密钥
})
```

#### 安全配置最佳实践

```typescript
// 生产环境建议配置具体的origin和安全密钥
const channel = createIframeChannel({
  targetOrigin: 'https://your-domain.com', // 指定具体origin
  // 或者允许多个origin
  // targetOrigin: ['https://domain1.com', 'https://domain2.com']
  postMessageKey: crypto.randomUUID(), // 使用随机密钥
  debug: false // 生产环境关闭调试模式
})
```

## API 参考

### createIframeChannel(config: ChannelConfig): IframeChannel

创建iframe通信通道。

**参数：**
- `config`: 通道配置

**返回值：** IframeChannel实例

### IframeChannel

#### 方法

| 方法 | 说明 | 参数 |
|------|------|------|
| `registerTarget(id, window)` | 注册目标窗口 | `id: string`, `window: Window` |
| `unregisterTarget(id)` | 移除目标窗口 | `id: string` |
| `send(targetId, event, data)` | 发送事件（不需要回复） | `targetId: string`, `event: string`, `data?: any` |
| `request(targetId, event, data, config?)` | 发送请求（需要回复） | `targetId: string`, `event: string`, `data?: any`, `config?: RequestConfig` |
| `on(event, handler, once?)` | 监听事件 | `event: string`, `handler: EventHandler`, `once?: boolean` |
| `once(event, handler)` | 监听事件（只触发一次） | `event: string`, `handler: EventHandler` |
| `off(event, handler)` | 移除事件监听 | `event: string`, `handler: EventHandler` |
| `onRequest(event, handler)` | 注册请求处理器 | `event: string`, `handler: RequestHandler` |
| `onQuickRequest(event, handler)` | 注册快捷回复处理器 | `event: string`, `handler: QuickReplyHandler` |
| `offRequest(event)` | 移除请求处理器 | `event: string` |
| `destroy()` | 销毁通道 | 无 |

#### 属性

| 属性 | 说明 | 类型 |
|------|------|------|
| `config` | 通道配置 | `Required<ChannelConfig>` |

### ChannelConfig

| 字段 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| `targetOrigin` | 目标origin | `string` | `'*'` |
| `timeout` | 默认超时时间（毫秒） | `number` | `5000` |
| `maxRetries` | 默认最大重发次数 | `number` | `3` |
| `retryDelay` | 重发延迟（毫秒） | `number` | `1000` |
| `isIframe` | 是否是iframe端 | `boolean` | `false` |
| `sourceId` | 来源标识 | `string` | `'main'` |
| `debug` | 是否启用调试模式 | `boolean` | `false` |
| `postMessageKey` | postMessage安全密钥 | `string` | `''` |

### RequestConfig

| 字段 | 说明 | 类型 |
|------|------|------|
| `timeout` | 超时时间（毫秒） | `number` |
| `maxRetries` | 最大重发次数 | `number` |
| `retryDelay` | 重发延迟（毫秒） | `number` |
| `needResponse` | 是否需要回复 | `boolean` |

### QuickReplyHandler

```typescript
type QuickReplyHandler<T = any, R = any> = (data: T, ctx: ReplyContext<R>) => void | Promise<void>
```

### ReplyContext

```typescript
interface ReplyContext<T = any> {
  reply: (data: T) => void
  replyError: (error: string) => void
  message: IframeMessage
}
```

## 消息格式

```typescript
interface IframeMessage {
  id: string           // 消息唯一ID
  type: 'event' | 'request' | 'response' | 'error'
  event: string        // 事件名称
  data?: any           // 消息数据
  source: string       // 来源标识
  target: string       // 目标标识
  timestamp: number    // 时间戳
  requestId?: string   // 请求ID（用于请求-响应）
  error?: string       // 错误信息
  needResponse?: boolean
  timeout?: number
  retryCount?: number
  maxRetries?: number
  correlationId?: string
  postMessageKey?: string // postMessage安全密钥
}
```

## 重发机制

1. **超时检测**：发送请求后启动超时计时器
2. **自动重发**：超时后自动重发，使用指数退避策略
3. **最大次数**：达到最大重发次数后，reject Promise
4. **关联匹配**：使用correlationId匹配请求和响应

## 注意事项

1. **安全配置**：生产环境建议指定具体的`targetOrigin`，避免使用`'*'`
2. **资源清理**：组件卸载时调用`destroy()`方法清理资源
3. **错误处理**：始终使用try-catch处理`request`方法的Promise
4. **调试模式**：开发时开启`debug: true`查看详细日志
5. **性能考虑**：避免频繁发送大量消息，合理设置超时时间
6. **安全密钥**：使用`postMessageKey`增强安全性，确保密钥安全传递和存储
7. **密钥管理**：每次通信会话建议生成新的密钥，避免密钥泄露

## 与现有系统集成

### 微应用商店

```typescript
// 主页面
const postMessageKey = crypto.randomUUID()

const channel = createIframeChannel({
  targetOrigin: '*',
  sourceId: 'main',
  postMessageKey // 使用安全密钥
})

// 注册微应用商店iframe
const storeIframe = document.getElementById('store-iframe') as HTMLIFrameElement
channel.registerTarget('micro-app-store', storeIframe.contentWindow!)

// 监听微应用商店事件
channel.on('app-installed', (app) => {
  console.log('应用已安装:', app)
})

// 向微应用商店发送请求
const apps = await channel.request('micro-app-store', 'get-installed-apps')
```

### 微应用通信

```typescript
// 微应用端
const urlParams = new URLSearchParams(window.location.search)
const postMessageKey = urlParams.get('postMessageKey') || ''

const channel = createIframeChannel({
  isIframe: true,
  targetOrigin: window.location.origin,
  sourceId: 'micro-app-xxx',
  postMessageKey // 使用从URL获取的安全密钥
})

// 注册主页面
channel.registerTarget('parent', window.parent)

// 响应主页面请求
channel.onRequest('get-app-info', async () => {
  return { id: 'xxx', name: 'My App', version: '1.0.0' }
})

// 向主页面发送事件
channel.send('parent', 'app-ready', { appId: 'xxx' })
```

## 安全机制详解

### postMessageKey 安全密钥

`postMessageKey` 是一个安全增强机制，通过为每个通信会话生成唯一的随机密钥，防止消息劫持和伪造。

#### 工作原理

1. **密钥生成**：主页面生成随机密钥（推荐使用`crypto.randomUUID()`）
2. **密钥传递**：通过URL参数传递给iframe：`https://example.com?postMessageKey=xxx`
3. **密钥验证**：所有消息都携带密钥，接收端验证密钥是否匹配
4. **会话隔离**：每次通信会话使用不同的密钥，确保会话隔离

#### 安全优势

1. **防止消息劫持**：即使攻击者能注入iframe，没有密钥也无法通信
2. **多实例隔离**：每个iframe实例有唯一密钥，避免消息混淆
3. **纵深防御**：即使origin验证被绕过，还有密钥验证
4. **会话安全**：每次打开弹窗生成新密钥，防止重放攻击

#### 实现示例

```typescript
// 主页面
const postMessageKey = crypto.randomUUID()

// 将密钥添加到iframe URL
const iframe = document.getElementById('my-iframe') as HTMLIFrameElement
const url = new URL(iframe.src)
url.searchParams.set('postMessageKey', postMessageKey)
iframe.src = url.toString()

// 创建通信通道时配置密钥
const channel = createIframeChannel({
  targetOrigin: '*',
  postMessageKey
})

// iframe端
const urlParams = new URLSearchParams(window.location.search)
const iframePostMessageKey = urlParams.get('postMessageKey') || ''

const iframeChannel = createIframeChannel({
  isIframe: true,
  targetOrigin: window.location.origin,
  postMessageKey: iframePostMessageKey
})
```

#### 注意事项

1. **密钥安全**：确保密钥在传递过程中不被泄露
2. **密钥存储**：不要将密钥存储在localStorage等不安全的地方
3. **密钥轮换**：每次通信会话建议生成新的密钥
4. **兼容性**：确保两端都支持postMessageKey功能

## 故障排除

### 消息收不到

1. 检查`targetOrigin`配置是否正确
2. 检查目标窗口是否已注册
3. 开启`debug`模式查看日志

### 请求超时

1. 检查目标窗口是否在线
2. 检查请求处理器是否正确注册
3. 调整`timeout`和`maxRetries`配置

### 安全错误

1. 检查origin是否匹配
2. 生产环境避免使用`'*'`
3. 使用HTTPS协议