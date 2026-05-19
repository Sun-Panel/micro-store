/**
 * iframe跨页面通信类型定义
 */

// =======================
// 消息类型
// =======================

/**
 * 消息类型
 */
export type MessageType = 'event' | 'request' | 'response' | 'error'

/**
 * 消息状态
 */
export type MessageStatus = 'pending' | 'success' | 'error' | 'timeout'

// =======================
// 消息接口
// =======================

/**
 * iframe通信消息
 */
export interface IframeMessage {
  /** 消息唯一ID */
  id: string
  /** 消息类型 */
  type: MessageType
  /** 事件名称 */
  event: string
  /** 消息数据 */
  data?: any
  /** 来源标识 */
  source: string
  /** 目标标识 */
  target: string
  /** 时间戳 */
  timestamp: number
  /** 请求ID（用于请求-响应） */
  requestId?: string
  /** 错误信息 */
  error?: string
  /** 是否需要回复 */
  needResponse?: boolean
  /** 超时时间（毫秒） */
  timeout?: number
  /** 重发次数 */
  retryCount?: number
  /** 最大重发次数 */
  maxRetries?: number
  /** 重发延迟（毫秒） */
  retryDelay?: number
  /** 关联ID（用于重发时匹配） */
  correlationId?: string
  /** postMessage安全密钥 */
  postMessageKey?: string
}

/**
 * 消息响应
 */
export interface MessageResponse {
  /** 消息ID */
  id: string
  /** 请求ID */
  requestId: string
  /** 事件名称 */
  event: string
  /** 响应数据 */
  data?: any
  /** 错误信息 */
  error?: string
  /** 状态 */
  status: MessageStatus
  /** 时间戳 */
  timestamp: number
}

// =======================
// 配置类型
// =======================

/**
 * 通道配置
 */
export interface ChannelConfig {
  /** 目标origin（'*'表示接受所有origin，生产环境建议指定具体origin） */
  targetOrigin: string
  /** 默认超时时间（毫秒） */
  timeout?: number
  /** 默认最大重发次数 */
  maxRetries?: number
  /** 重发延迟（毫秒） */
  retryDelay?: number
  /** 是否是iframe端 */
  isIframe?: boolean
  /** 来源标识 */
  sourceId?: string
  /** 是否启用调试模式 */
  debug?: boolean
  /** postMessage安全密钥，用于验证消息来源 */
  postMessageKey?: string
}

/**
 * 请求配置
 */
export interface RequestConfig {
  /** 超时时间（毫秒） */
  timeout?: number
  /** 最大重发次数 */
  maxRetries?: number
  /** 重发延迟（毫秒） */
  retryDelay?: number
  /** 是否需要回复 */
  needResponse?: boolean
}

// =======================
// 事件处理器类型
// =======================

/**
 * 事件处理器
 */
export type EventHandler<T = any> = (data: T) => void | Promise<void>

/**
 * 请求处理器
 */
export type RequestHandler<T = any, R = any> = (data: T) => R | Promise<R>

/**
 * 快捷回复上下文
 */
export interface ReplyContext<T = any> {
  /** 回复数据 */
  reply: (data: T) => void
  /** 回复错误 */
  replyError: (error: string) => void
  /** 原始请求消息 */
  message: IframeMessage
}

/**
 * 快捷回复处理器
 */
export type QuickReplyHandler<T = any, R = any> = (data: T, ctx: ReplyContext<R>) => void | Promise<void>

// =======================
// 内部类型
// =======================

/**
 * 待处理请求
 */
export interface PendingRequest {
  /** 请求消息 */
  message: IframeMessage
  /** Promise resolve函数 */
  resolve: (value: any) => void
  /** Promise reject函数 */
  reject: (reason: any) => void
  /** 超时定时器 */
  timeoutId?: ReturnType<typeof setTimeout>
  /** 重发定时器 */
  retryId?: ReturnType<typeof setTimeout>
  /** 当前重发次数 */
  retryCount: number
}

/**
 * 事件监听器
 */
export interface EventListener {
  /** 事件名称 */
  event: string
  /** 处理函数 */
  handler: EventHandler
  /** 是否只触发一次 */
  once?: boolean
}

// =======================
// 常量
// =======================

/**
 * 默认配置
 */
export const DEFAULT_CONFIG: Required<ChannelConfig> = {
  targetOrigin: '*',
  timeout: 5000,
  maxRetries: 3,
  retryDelay: 1000,
  isIframe: false,
  sourceId: 'main',
  debug: false,
  postMessageKey: '',
}

/**
 * 消息类型常量
 */
export const MESSAGE_TYPES = {
  EVENT: 'event',
  REQUEST: 'request',
  RESPONSE: 'response',
  ERROR: 'error',
} as const

/**
 * 消息状态常量
 */
export const MESSAGE_STATUS = {
  PENDING: 'pending',
  SUCCESS: 'success',
  ERROR: 'error',
  TIMEOUT: 'timeout',
} as const