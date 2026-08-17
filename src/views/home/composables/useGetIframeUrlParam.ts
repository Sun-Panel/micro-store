const SESSION_KEY = 'IframeUrlParam'

/**
 * 判断当前是否处于 iframe 模式（嵌入在主应用中）
 * @returns 是否在 iframe 中运行
 */
export function isIframeMode(): boolean {
  try {
    return window !== window.top
  }
  catch {
    return true
  }
}

/**
 * 从 URL 参数获取所有配置，优先从 URL 读取，若无则回退到 sessionStorage
 * @returns 包含所有配置键值对的对象
 */
export function getIframeAllUrlParam(): Record<string, string> {
  const result: Record<string, string> = {}

  // 1. 优先从 URL GET 参数中读取
  const urlParams = new URLSearchParams(window.location.search)
  urlParams.forEach((value, key) => {
    result[key] = value
  })

  // 2. 如果 URL 中有参数，保存到 sessionStorage
  const urlParamCount = [...urlParams.keys()].length
  if (urlParamCount > 0) {
    window.sessionStorage.setItem(SESSION_KEY, JSON.stringify(result))
  }

  // 3. 如果 URL 中没有参数，从 sessionStorage 中读取
  if (urlParamCount === 0) {
    const stored = window.sessionStorage.getItem(SESSION_KEY)
    if (stored) {
      try {
        return JSON.parse(stored)
      }
      catch {
        return result
      }
    }
  }

  return result
}

import { extractBaseVersion } from '@/utils/functions/version'

/**
 * 获取主应用基础版本号（从 URL 参数 version 中提取）
 * 例如 "2.0.0-beta260505" → "2.0.0"
 * @returns 主应用基础版本号，如 "2.0.0"
 */
export function getHostAppBaseVersion(): string {
  const params = getIframeAllUrlParam()
  const versionParam = params['version'] || ''
  return extractBaseVersion(versionParam)
}

/**
 * 获取主应用完整版本号
 * @returns 完整版本字符串，如 "2.0.0-beta260505"
 */
export function getHostAppVersion(): string {
  const params = getIframeAllUrlParam()
  return params['version'] || ''
}
