/**
 * 通信核心逻辑模块
 * 实现iframe跨页面通信的核心功能
 */

import type {
  ChannelConfig,
  EventHandler,
  EventListener,
  IframeMessage,
  QuickReplyHandler,
  ReplyContext,
  RequestConfig,
  RequestHandler,
} from './types'
import { RetryManager } from './retry'
import { generateCorrelationId, generateMessageId, validateMessageSecurity } from './security'
import { DEFAULT_CONFIG, MESSAGE_TYPES } from './types'

/**
 * iframe通信通道
 */
export class IframeChannel {
  private config: Required<ChannelConfig>
  private retryManager: RetryManager
  private eventListeners: Map<string, EventListener[]> = new Map()
  private requestHandlers: Map<string, RequestHandler> = new Map()
  private quickRequestHandlers: Map<string, QuickReplyHandler> = new Map()
  private targetWindows: Map<string, Window> = new Map()
  private isDestroyed: boolean = false

  constructor(config: ChannelConfig) {
    this.config = {
      ...DEFAULT_CONFIG,
      ...config,
    }
    this.retryManager = new RetryManager(this.config.debug)

    // 绑定消息处理函数
    this.handleMessage = this.handleMessage.bind(this)

    // 添加消息监听
    window.addEventListener('message', this.handleMessage)

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Initialized with config:`, this.config)
    }
  }

  /**
   * 注册目标窗口
   * @param id 窗口标识
   * @param window 目标窗口
   */
  registerTarget(id: string, window: Window): void {
    this.targetWindows.set(id, window)

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Registered target: ${id}`)
    }
  }

  /**
   * 移除目标窗口
   * @param id 窗口标识
   */
  unregisterTarget(id: string): void {
    this.targetWindows.delete(id)

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Unregistered target: ${id}`)
    }
  }

  /**
   * 获取目标窗口
   * @param id 窗口标识
   * @returns 目标窗口
   */
  getTarget(id: string): Window | undefined {
    return this.targetWindows.get(id)
  }

  /**
   * 发送事件（不需要回复）
   * @param targetId 目标标识
   * @param event 事件名称
   * @param data 数据
   */
  send(targetId: string, event: string, data?: any): void {
    if (this.isDestroyed) {
      throw new Error('Channel is destroyed')
    }

    const target = this.targetWindows.get(targetId)
    if (!target) {
      throw new Error(`Target not found: ${targetId}`)
    }

    const message: IframeMessage = {
      id: generateMessageId(),
      type: MESSAGE_TYPES.EVENT,
      event,
      data,
      source: this.config.sourceId,
      target: targetId,
      timestamp: Date.now(),
      needResponse: false,
      postMessageKey: this.config.postMessageKey,
    }

    this.postMessage(target, message)

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Sent event to ${targetId}:`, message)
    }
  }

  /**
   * 发送请求（需要回复）
   * @param targetId 目标标识
   * @param event 事件名称
   * @param data 数据
   * @param config 请求配置
   * @returns Promise，解析为响应数据
   */
  request<T = any>(
    targetId: string,
    event: string,
    data?: any,
    config?: RequestConfig,
  ): Promise<T> {
    if (this.isDestroyed) {
      return Promise.reject(new Error('Channel is destroyed'))
    }

    const target = this.targetWindows.get(targetId)
    if (!target) {
      return Promise.reject(new Error(`Target not found: ${targetId}`))
    }

    const requestId = generateMessageId()
    const correlationId = generateCorrelationId()
    const timeout = config?.timeout ?? this.config.timeout
    const maxRetries = config?.maxRetries ?? this.config.maxRetries
    const retryDelay = config?.retryDelay ?? this.config.retryDelay

    const message: IframeMessage = {
      id: requestId,
      type: MESSAGE_TYPES.REQUEST,
      event,
      data,
      source: this.config.sourceId,
      target: targetId,
      timestamp: Date.now(),
      requestId,
      needResponse: true,
      timeout,
      retryCount: 0,
      maxRetries,
      retryDelay,
      correlationId,
      postMessageKey: this.config.postMessageKey,
    }

    return new Promise<T>((resolve, reject) => {
      // 创建待处理请求
      const pendingRequest = {
        message,
        resolve: resolve as (value: any) => void,
        reject,
        retryCount: 0,
      }

      // 添加到重发管理器
      this.retryManager.addRequest(pendingRequest)

      // 启动超时检测
      this.retryManager.startTimeout(requestId, timeout, (msg) => {
        this.handleTimeout(msg)
      })

      // 发送消息
      this.postMessage(target, message)

      if (this.config.debug) {
        // eslint-disable-next-line no-console
        console.log(`[IframeChannel] Sent request to ${targetId}:`, message)
      }
    })
  }

  /**
   * 监听事件
   * @param event 事件名称
   * @param handler 处理函数
   * @param once 是否只触发一次
   * @returns 取消监听函数
   */
  on(event: string, handler: EventHandler, once: boolean = false): () => void {
    if (this.isDestroyed) {
      throw new Error('Channel is destroyed')
    }

    const listeners = this.eventListeners.get(event) || []
    const listener: EventListener = { event, handler, once }
    listeners.push(listener)
    this.eventListeners.set(event, listeners)

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Added listener for event: ${event}`)
    }

    // 返回取消监听函数
    return () => {
      this.off(event, handler)
    }
  }

  /**
   * 监听事件（只触发一次）
   * @param event 事件名称
   * @param handler 处理函数
   * @returns 取消监听函数
   */
  once(event: string, handler: EventHandler): () => void {
    return this.on(event, handler, true)
  }

  /**
   * 移除事件监听
   * @param event 事件名称
   * @param handler 处理函数
   */
  off(event: string, handler: EventHandler): void {
    const listeners = this.eventListeners.get(event)
    if (listeners) {
      const index = listeners.findIndex(l => l.handler === handler)
      if (index !== -1) {
        listeners.splice(index, 1)
        if (listeners.length === 0) {
          this.eventListeners.delete(event)
        }
      }
    }

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Removed listener for event: ${event}`)
    }
  }

  /**
   * 注册请求处理器
   * @param event 事件名称
   * @param handler 处理函数
   */
  onRequest<T = any, R = any>(event: string, handler: RequestHandler<T, R>): void {
    if (this.isDestroyed) {
      throw new Error('Channel is destroyed')
    }

    this.requestHandlers.set(event, handler as RequestHandler)

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Registered request handler for event: ${event}`)
    }
  }

  /**
   * 注册快捷回复处理器
   * @param event 事件名称
   * @param handler 处理函数
   */
  onQuickRequest<T = any, R = any>(event: string, handler: QuickReplyHandler<T, R>): void {
    if (this.isDestroyed) {
      throw new Error('Channel is destroyed')
    }

    this.quickRequestHandlers.set(event, handler as QuickReplyHandler)

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Registered quick reply handler for event: ${event}`)
    }
  }

  /**
   * 移除请求处理器
   * @param event 事件名称
   */
  offRequest(event: string): void {
    this.requestHandlers.delete(event)
    this.quickRequestHandlers.delete(event)

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Removed request handler for event: ${event}`)
    }
  }

  /**
   * 处理接收到的消息
   * @param event MessageEvent
   */
  private handleMessage(event: MessageEvent): void {
    if (this.isDestroyed) {
      return
    }

    // 验证消息安全性
    const validation = validateMessageSecurity(event, this.config.targetOrigin)
    if (!validation.valid) {
      if (this.config.debug) {
        console.warn(`[IframeChannel] Invalid message:`, validation.error)
      }
      return
    }

    const message = validation.message!

    // 验证postMessageKey（如果配置了）
    if (this.config.postMessageKey) {
      if (message.postMessageKey !== this.config.postMessageKey) {
        if (this.config.debug) {
          console.warn(`[IframeChannel] Invalid postMessageKey:`, message.postMessageKey)
        }
        return
      }
    }

    // 处理响应消息
    if (message.type === MESSAGE_TYPES.RESPONSE || message.type === MESSAGE_TYPES.ERROR) {
      this.retryManager.handleResponse(message)
      return
    }

    // 处理请求消息
    if (message.type === MESSAGE_TYPES.REQUEST) {
      if (event.source) {
        this.handleRequestMessage(message, event.source as Window)
      }
      return
    }

    // 处理事件消息
    if (message.type === MESSAGE_TYPES.EVENT) {
      this.handleEventMessage(message)
    }
  }

  /**
   * 处理事件消息
   * @param message 消息
   */
  private handleEventMessage(message: IframeMessage): void {
    const listeners = this.eventListeners.get(message.event)
    if (listeners) {
      listeners.forEach((listener) => {
        try {
          listener.handler(message.data)
        }
        catch (error) {
          console.error(`[IframeChannel] Error in event handler:`, error)
        }
      })

      // 移除一次性监听器
      const onceListeners = listeners.filter(l => l.once)
      onceListeners.forEach((listener) => {
        this.off(message.event, listener.handler)
      })
    }

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Handled event: ${message.event}`)
    }
  }

  /**
   * 处理请求消息
   * @param message 消息
   * @param sourceWindow 源窗口
   */
  private async handleRequestMessage(message: IframeMessage, sourceWindow: Window): Promise<void> {
    // 检查快捷回复处理器
    const quickHandler = this.quickRequestHandlers.get(message.event)
    if (quickHandler) {
      try {
        // 创建回复上下文
        const ctx: ReplyContext = {
          reply: (data: any) => {
            this.sendResponse(sourceWindow, message, data)
          },
          replyError: (error: string) => {
            this.sendErrorResponse(sourceWindow, message, error)
          },
          message,
        }

        // 调用快捷回复处理器
        await quickHandler(message.data, ctx)
      }
      catch (error) {
        // 发送错误响应
        const errorMessage = error instanceof Error ? error.message : String(error)
        this.sendErrorResponse(sourceWindow, message, errorMessage)
      }

      if (this.config.debug) {
        // eslint-disable-next-line no-console
        console.log(`[IframeChannel] Handled quick request: ${message.event}`)
      }
      return
    }

    // 检查标准请求处理器
    const handler = this.requestHandlers.get(message.event)
    if (!handler) {
      // 发送错误响应
      this.sendErrorResponse(
        sourceWindow,
        message,
        `No handler registered for event: ${message.event}`,
      )
      return
    }

    try {
      // 调用处理器
      const result = await handler(message.data)

      // 发送成功响应
      this.sendResponse(sourceWindow, message, result)
    }
    catch (error) {
      // 发送错误响应
      const errorMessage = error instanceof Error ? error.message : String(error)
      this.sendErrorResponse(sourceWindow, message, errorMessage)
    }

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Handled request: ${message.event}`)
    }
  }

  /**
   * 发送响应
   * @param targetWindow 目标窗口
   * @param request 原始请求
   * @param data 响应数据
   */
  private sendResponse(targetWindow: Window, request: IframeMessage, data: any): void {
    const response: IframeMessage = {
      id: generateMessageId(),
      type: MESSAGE_TYPES.RESPONSE,
      event: request.event,
      data,
      source: this.config.sourceId,
      target: request.source,
      timestamp: Date.now(),
      requestId: request.id,
      postMessageKey: request.postMessageKey,
    }

    this.postMessage(targetWindow, response)
  }

  /**
   * 发送错误响应
   * @param targetWindow 目标窗口
   * @param request 原始请求
   * @param error 错误信息
   */
  private sendErrorResponse(targetWindow: Window, request: IframeMessage, error: string): void {
    const response: IframeMessage = {
      id: generateMessageId(),
      type: MESSAGE_TYPES.ERROR,
      event: request.event,
      error,
      source: this.config.sourceId,
      target: request.source,
      timestamp: Date.now(),
      requestId: request.id,
      postMessageKey: request.postMessageKey,
    }

    this.postMessage(targetWindow, response)
  }

  /**
   * 处理超时
   * @param message 消息
   */
  private handleTimeout(message: IframeMessage): void {
    const request = this.retryManager.getRequest(message.id)
    if (!request) {
      return
    }

    // 检查是否还可以重发
    if (request.retryCount < (message.maxRetries || this.config.maxRetries)) {
      // 启动重发
      this.retryManager.startRetry(
        message.id,
        message.retryDelay || this.config.retryDelay,
        (msg) => {
          this.retryMessage(msg)
        },
      )
    }
    else {
      // 超过最大重发次数，reject Promise
      request.reject(new Error(`Request timeout after ${message.maxRetries || this.config.maxRetries} retries`))
      this.retryManager.removeRequest(message.id)
    }
  }

  /**
   * 重发消息
   * @param message 消息
   */
  private retryMessage(message: IframeMessage): void {
    const target = this.targetWindows.get(message.target)
    if (!target) {
      const request = this.retryManager.getRequest(message.id)
      if (request) {
        request.reject(new Error(`Target not found: ${message.target}`))
        this.retryManager.removeRequest(message.id)
      }
      return
    }

    // 更新消息
    const updatedMessage: IframeMessage = {
      ...message,
      timestamp: Date.now(),
      retryCount: message.retryCount || 0,
    }

    // 重新启动超时检测
    this.retryManager.startTimeout(
      message.id,
      message.timeout || this.config.timeout,
      (msg) => {
        this.handleTimeout(msg)
      },
    )

    // 发送消息
    this.postMessage(target, updatedMessage)

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log(`[IframeChannel] Retried message: ${message.id}, attempt: ${updatedMessage.retryCount}`)
    }
  }

  /**
   * 发送postMessage
   * @param targetWindow 目标窗口
   * @param message 消息
   */
  private postMessage(targetWindow: Window, message: IframeMessage): void {
    targetWindow.postMessage(message, this.config.targetOrigin)
  }

  /**
   * 销毁通道
   */
  destroy(): void {
    if (this.isDestroyed) {
      return
    }

    this.isDestroyed = true

    // 移除消息监听
    window.removeEventListener('message', this.handleMessage)

    // 清理重发管理器
    this.retryManager.clear()

    // 清理监听器
    this.eventListeners.clear()
    this.requestHandlers.clear()
    this.quickRequestHandlers.clear()

    // 清理目标窗口
    this.targetWindows.clear()

    if (this.config.debug) {
      // eslint-disable-next-line no-console
      console.log('[IframeChannel] Channel destroyed')
    }
  }
}
