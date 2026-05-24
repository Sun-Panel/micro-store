/**
 * 安全验证模块
 * 用于验证消息来源和内容安全性
 */

import type { IframeMessage } from './types'

/**
 * 验证消息来源
 * @param event MessageEvent
 * @param allowedOrigins 允许的origin列表
 * @returns 是否验证通过
 */
export function validateMessageOrigin(
  event: MessageEvent,
  allowedOrigins: string[] | string,
): boolean {
  const { origin } = event

  // 如果允许所有origin
  if (allowedOrigins === '*') {
    return true
  }

  // 如果是数组，检查是否在允许列表中
  if (Array.isArray(allowedOrigins)) {
    return allowedOrigins.includes(origin)
  }

  // 如果是字符串，直接比较
  return origin === allowedOrigins
}

/**
 * 验证消息格式
 * @param data 消息数据
 * @returns 是否是有效的IframeMessage
 */
export function validateMessageFormat(data: any): data is IframeMessage {
  if (!data || typeof data !== 'object') {
    return false
  }

  // 检查必需字段
  const requiredFields = ['id', 'type', 'event', 'source', 'timestamp']
  for (const field of requiredFields) {
    if (!(field in data)) {
      return false
    }
  }

  // 检查字段类型
  if (typeof data.id !== 'string' || data.id.length === 0) {
    return false
  }

  if (!['event', 'request', 'response', 'error'].includes(data.type)) {
    return false
  }

  if (typeof data.event !== 'string' || data.event.length === 0) {
    return false
  }

  if (typeof data.source !== 'string') {
    return false
  }

  if (typeof data.timestamp !== 'number') {
    return false
  }

  // 检查时间戳合理性（不超过1小时）
  const now = Date.now()
  const oneHour = 60 * 60 * 1000
  if (Math.abs(now - data.timestamp) > oneHour) {
    return false
  }

  return true
}

/**
 * 生成唯一消息ID
 * @returns 唯一ID
 */
export function generateMessageId(): string {
  const timestamp = Date.now().toString(36)
  const random = Math.random().toString(36).substring(2, 8)
  return `${timestamp}-${random}`
}

/**
 * 生成关联ID（用于重发时匹配）
 * @returns 关联ID
 */
export function generateCorrelationId(): string {
  return `corr-${Date.now()}-${Math.random().toString(36).substring(2, 8)}`
}

/**
 * 防止XSS攻击的字符串清理
 * @param input 输入字符串
 * @returns 清理后的字符串
 */
export function sanitizeString(input: string): string {
  if (typeof input !== 'string') {
    return ''
  }

  // 只转义 HTML 特殊字符，不转义斜杠（避免破坏 URL）
  return input
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#x27;')
}

/**
 * 深度清理对象中的字符串
 * @param obj 输入对象
 * @returns 清理后的对象
 */
export function sanitizeObject<T>(obj: T): T {
  if (typeof obj !== 'object' || obj === null) {
    return obj
  }

  if (Array.isArray(obj)) {
    return obj.map(item => sanitizeObject(item)) as T
  }

  const sanitized: any = {}
  for (const [key, value] of Object.entries(obj)) {
    if (typeof value === 'string') {
      sanitized[key] = sanitizeString(value)
    }
    else if (typeof value === 'object' && value !== null) {
      sanitized[key] = sanitizeObject(value)
    }
    else {
      sanitized[key] = value
    }
  }

  return sanitized as T
}

/**
 * 验证消息安全性
 * @param event MessageEvent
 * @param allowedOrigins 允许的origin列表
 * @returns 验证结果
 */
export function validateMessageSecurity(
  event: MessageEvent,
  allowedOrigins: string[] | string,
): {
  valid: boolean
  error?: string
  message?: IframeMessage
} {
  // 验证origin
  if (!validateMessageOrigin(event, allowedOrigins)) {
    return {
      valid: false,
      error: `Invalid origin: ${event.origin}`,
    }
  }

  // 验证消息格式
  if (!validateMessageFormat(event.data)) {
    return {
      valid: false,
      error: 'Invalid message format',
    }
  }

  // 清理消息数据（防止XSS）
  const message = sanitizeObject(event.data) as IframeMessage

  return {
    valid: true,
    message,
  }
}
