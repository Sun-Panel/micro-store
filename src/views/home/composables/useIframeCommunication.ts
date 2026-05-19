/**
 * iframe通信 composable
 * 封装 IframeChannel，提供简化的 API，支持自动销毁
 */

import type { ChannelConfig, EventHandler, QuickReplyHandler, RequestHandler } from '../utils/iframeCommunication'
import { onUnmounted } from 'vue'
import { IframeChannel } from '../utils/iframeCommunication'

export interface IframeCommunicationConfig extends ChannelConfig {
  /** 是否自动销毁（组件卸载时），默认 true */
  autoDestroy?: boolean
}

export interface UseIframeCommunicationReturn {
  /** 获取底层 IframeChannel 实例（高级用法） */
  channel: IframeChannel
  /** 发送事件（不需要回复） */
  sendMessage: (targetId: string, event: string, data?: any) => void
  /** 发送请求（需要回复） */
  sendRequest: <T = any>(targetId: string, event: string, data?: any, config?: { timeout?: number, maxRetries?: number, retryDelay?: number }) => Promise<T>
  /** 监听事件 */
  on: (event: string, handler: EventHandler, once?: boolean) => () => void
  /** 监听事件（只触发一次） */
  once: (event: string, handler: EventHandler) => () => void
  /** 注册请求处理器 */
  onRequest: <T = any, R = any>(event: string, handler: RequestHandler<T, R>) => void
  /** 注册快捷回复处理器 */
  onQuickRequest: <T = any, R = any>(event: string, handler: QuickReplyHandler<T, R>) => void
  /** 手动销毁通信通道 */
  destroy: () => void
  /** 重新初始化通信通道（例如更换 postMessageKey） */
  reinit: (newConfig?: IframeCommunicationConfig) => void
}

/**
 * 使用 iframe 通信 composable
 * @param config 通信配置
 * @returns 通信 API
 */
export function useIframeCommunication(config: IframeCommunicationConfig): UseIframeCommunicationReturn {
  // 创建通道实例
  let channel = new IframeChannel(config)

  // 默认自动销毁
  const autoDestroy = config.autoDestroy !== false

  // 组件卸载时自动销毁
  if (autoDestroy) {
    onUnmounted(() => {
      channel.destroy()
    })
  }

  // 包装 API 方法
  const sendMessage: UseIframeCommunicationReturn['sendMessage'] = (targetId, event, data) => {
    channel.send(targetId, event, data)
  }

  const sendRequest: UseIframeCommunicationReturn['sendRequest'] = (targetId, event, data, requestConfig) => {
    return channel.request(targetId, event, data, requestConfig)
  }

  const on: UseIframeCommunicationReturn['on'] = (event, handler, once = false) => {
    return channel.on(event, handler, once)
  }

  const once: UseIframeCommunicationReturn['once'] = (event, handler) => {
    return channel.once(event, handler)
  }

  const onRequest: UseIframeCommunicationReturn['onRequest'] = (event, handler) => {
    channel.onRequest(event, handler)
  }

  const onQuickRequest: UseIframeCommunicationReturn['onQuickRequest'] = (event, handler) => {
    channel.onQuickRequest(event, handler)
  }

  const destroy: UseIframeCommunicationReturn['destroy'] = () => {
    channel.destroy()
  }

  const reinit: UseIframeCommunicationReturn['reinit'] = (newConfig) => {
    // 销毁旧通道
    channel.destroy()
    // 创建新通道
    channel = new IframeChannel(newConfig || config)
  }

  return {
    channel,
    sendMessage,
    sendRequest,
    on,
    once,
    onRequest,
    onQuickRequest,
    destroy,
    reinit,
  }
}
