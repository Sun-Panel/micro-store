import { isIframe } from '@/utils/cmn'
import { extractBaseVersion } from '@/utils/functions/version'

const SESSION_KEY = 'IframeUrlParam'

/** 缓存解析后的 URL 参数，避免重复解析 */
let cachedUrlParams: Record<string, string> | null = null

/**
 * 判断当前是否处于 iframe 模式（嵌入在主应用中）
 * 委托给共享的 isIframe 函数，保持单一实现
 * @returns 是否在 iframe 中运行
 */
export const isIframeMode = isIframe

/**
 * 从 URL 参数获取所有配置，优先从 URL 读取，若无则回退到 sessionStorage
 * 结果会被缓存，整个页面生命周期内只解析一次
 * @returns 包含所有配置键值对的对象
 */
export function getIframeAllUrlParam(): Record<string, string> {
  if (cachedUrlParams) {
    return cachedUrlParams
  }

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
        Object.assign(result, JSON.parse(stored))
      }
      catch {
        // 解析失败使用空结果
      }
    }
  }

  cachedUrlParams = result
  return result
}

const DEFAULT_HOST_APP_VERSION = '2.0.0'

/** 缓存提取后的基础版本号 */
let cachedHostBaseVersion: string | null = null

/**
 * 获取主应用基础版本号（从 URL 参数 version 中提取）
 * 例如 "2.0.0-beta260505" → "2.0.0"
 * @returns 主应用基础版本号，如 "2.0.0"，未传时默认 "2.0.0"
 */
export function getHostAppBaseVersion(): string {
  if (cachedHostBaseVersion) {
    return cachedHostBaseVersion
  }
  const params = getIframeAllUrlParam()
  const versionParam = params.version || ''
  const baseVersion = extractBaseVersion(versionParam)
  cachedHostBaseVersion = baseVersion === '0.0.0' ? DEFAULT_HOST_APP_VERSION : baseVersion
  return cachedHostBaseVersion
}

/**
 * 获取主应用完整版本号
 * @returns 完整版本字符串，如 "2.0.0-beta260505"
 */
export function getHostAppVersion(): string {
  const params = getIframeAllUrlParam()
  return params.version || ''
}
