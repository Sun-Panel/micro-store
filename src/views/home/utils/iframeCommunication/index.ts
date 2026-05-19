/**
 * iframe跨页面通信模块
 * 提供基于postMessage的跨页面通信功能
 */

import { IframeChannel } from './channel'

// 重新导出类型
export type {
  ChannelConfig,
  EventHandler,
  EventListener,
  IframeMessage,
  MessageResponse,
  MessageStatus,
  MessageType,
  QuickReplyHandler,
  ReplyContext,
  RequestConfig,
  RequestHandler,
} from './types'

// 重新导出常量
export { DEFAULT_CONFIG, MESSAGE_STATUS, MESSAGE_TYPES } from './types'

// 重新导出工具函数
export {
  generateCorrelationId,
  generateMessageId,
  sanitizeObject,
  sanitizeString,
  validateMessageFormat,
  validateMessageOrigin,
  validateMessageSecurity,
} from './security'

// 重新导出重发管理器
export { RetryManager } from './retry'

// 重新导出通道类
export { IframeChannel } from './channel'

/**
 * 创建iframe通信通道
 * @param config 通道配置
 * @returns IframeChannel实例
 */
export function createIframeChannel(config: import('./types').ChannelConfig): IframeChannel {
  return new IframeChannel(config)
}

/**
 * 创建iframe通信通道（单例模式）
 * @param config 通道配置
 * @returns IframeChannel实例
 */
let defaultChannel: IframeChannel | null = null

export function getOrCreateDefaultChannel(config?: import('./types').ChannelConfig): IframeChannel {
  if (!defaultChannel) {
    defaultChannel = new IframeChannel(config || {
      targetOrigin: '*',
      sourceId: 'main',
    })
  }
  return defaultChannel
}

/**
 * 销毁默认通道
 */
export function destroyDefaultChannel(): void {
  if (defaultChannel) {
    defaultChannel.destroy()
    defaultChannel = null
  }
}