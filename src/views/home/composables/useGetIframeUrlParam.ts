const SESSION_KEY = 'IframeUrlParam'

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
