/**
 * 重发机制模块
 * 处理消息超时和重发逻辑
 */

import type { IframeMessage, PendingRequest } from './types'

/**
 * 重发管理器
 */
export class RetryManager {
  private pendingRequests: Map<string, PendingRequest> = new Map()
  private debug: boolean

  constructor(debug: boolean = false) {
    this.debug = debug
  }

  /**
   * 添加待处理请求
   * @param request 待处理请求
   */
  addRequest(request: PendingRequest): void {
    const { message } = request
    this.pendingRequests.set(message.id, request)
    
    if (this.debug) {
      // eslint-disable-next-line no-console
      console.log(`[RetryManager] Added request: ${message.id}, event: ${message.event}`)
    }
  }

  /**
   * 移除待处理请求
   * @param messageId 消息ID
   */
  removeRequest(messageId: string): void {
    const request = this.pendingRequests.get(messageId)
    if (request) {
      // 清除定时器
      if (request.timeoutId) {
        clearTimeout(request.timeoutId)
      }
      if (request.retryId) {
        clearTimeout(request.retryId)
      }
      
      this.pendingRequests.delete(messageId)
      
      if (this.debug) {
        // eslint-disable-next-line no-console
        console.log(`[RetryManager] Removed request: ${messageId}`)
      }
    }
  }

  /**
   * 获取待处理请求
   * @param messageId 消息ID
   * @returns 待处理请求
   */
  getRequest(messageId: string): PendingRequest | undefined {
    return this.pendingRequests.get(messageId)
  }

  /**
   * 检查请求是否存在
   * @param messageId 消息ID
   * @returns 是否存在
   */
  hasRequest(messageId: string): boolean {
    return this.pendingRequests.has(messageId)
  }

  /**
   * 处理响应
   * @param response 响应消息
   * @returns 是否处理成功
   */
  handleResponse(response: IframeMessage): boolean {
    const { requestId } = response
    if (!requestId) {
      return false
    }

    const request = this.pendingRequests.get(requestId)
    if (!request) {
      if (this.debug) {
        // eslint-disable-next-line no-console
        console.log(`[RetryManager] No pending request for response: ${requestId}`)
      }
      return false
    }

    // 清除定时器
    if (request.timeoutId) {
      clearTimeout(request.timeoutId)
    }
    if (request.retryId) {
      clearTimeout(request.retryId)
    }

    // 如果是错误响应
    if (response.type === 'error') {
      request.reject(new Error(response.error || 'Request failed'))
    } else {
      request.resolve(response.data)
    }

    // 移除请求
    this.pendingRequests.delete(requestId)
    
    if (this.debug) {
      // eslint-disable-next-line no-console
      console.log(`[RetryManager] Handled response for request: ${requestId}`)
    }

    return true
  }

  /**
   * 启动超时检测
   * @param messageId 消息ID
   * @param timeout 超时时间（毫秒）
   * @param onTimeout 超时回调
   */
  startTimeout(
    messageId: string,
    timeout: number,
    onTimeout: (message: IframeMessage) => void,
  ): void {
    const request = this.pendingRequests.get(messageId)
    if (!request) {
      return
    }

    // 清除之前的超时定时器
    if (request.timeoutId) {
      clearTimeout(request.timeoutId)
    }

    // 设置新的超时定时器
    request.timeoutId = setTimeout(() => {
      if (this.debug) {
        // eslint-disable-next-line no-console
        console.log(`[RetryManager] Request timeout: ${messageId}, timeout: ${timeout}ms`)
      }
      onTimeout(request.message)
    }, timeout)
  }

  /**
   * 启动重发
   * @param messageId 消息ID
   * @param retryDelay 重发延迟（毫秒）
   * @param onRetry 重发回调
   */
  startRetry(
    messageId: string,
    retryDelay: number,
    onRetry: (message: IframeMessage) => void,
  ): void {
    const request = this.pendingRequests.get(messageId)
    if (!request) {
      return
    }

    // 清除之前的重发定时器
    if (request.retryId) {
      clearTimeout(request.retryId)
    }

    // 计算退避延迟（指数退避）
    const backoffDelay = retryDelay * Math.pow(2, request.retryCount)
    
    // 设置重发定时器
    request.retryId = setTimeout(() => {
      if (this.debug) {
        // eslint-disable-next-line no-console
        console.log(
          `[RetryManager] Retrying request: ${messageId}, `
          + `attempt: ${request.retryCount + 1}, delay: ${backoffDelay}ms`,
        )
      }
      
      // 增加重发次数
      request.retryCount++
      
      // 调用重发回调
      onRetry(request.message)
    }, backoffDelay)
  }

  /**
   * 清理所有待处理请求
   */
  clear(): void {
    this.pendingRequests.forEach((request, messageId) => {
      if (request.timeoutId) {
        clearTimeout(request.timeoutId)
      }
      if (request.retryId) {
        clearTimeout(request.retryId)
      }
    })
    
    this.pendingRequests.clear()
    
    if (this.debug) {
      // eslint-disable-next-line no-console
      console.log('[RetryManager] Cleared all pending requests')
    }
  }

  /**
   * 获取待处理请求数量
   * @returns 待处理请求数量
   */
  getPendingCount(): number {
    return this.pendingRequests.size
  }

  /**
   * 获取所有待处理请求的ID
   * @returns 待处理请求ID列表
   */
  getPendingIds(): string[] {
    return Array.from(this.pendingRequests.keys())
  }
}