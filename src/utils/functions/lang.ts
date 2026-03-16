/**
 * 获取浏览器语言
 */
export function getBrowserLang(): string {
  const lang = navigator.language || (navigator as any).userLanguage || 'zh-CN'
  if (lang.startsWith('zh'))
    return 'zh-CN'
  if (lang.startsWith('en'))
    return 'en-US'
  if (lang.startsWith('ja'))
    return 'ja-JP'
  if (lang.startsWith('ko'))
    return 'ko-KR'
  return 'zh-CN'
}

/**
 * 根据支持的语言列表获取当前语言
 * @param langList 支持的语言列表
 */
export function getCurrentLang(langList: string[]): string {
  const browserLang = getBrowserLang()
  return langList.includes(browserLang) ? browserLang : (langList.includes('zh-CN') ? 'zh-CN' : langList[0])
}

/**
 * 从应用信息中提取语言列表
 * @param appInfo 包含 langList 的应用信息
 */
export function getLangListFromAppInfo(appInfo: any): string[] {
  if (!appInfo)
    return ['zh-CN']
  const langList = appInfo.langList || []
  if (langList.length > 0) {
    return langList.map((l: any) => l.lang)
  }
  return ['zh-CN']
}

/**
 * 从应用信息中提取语言Map
 * @param appInfo 包含 langList 的应用信息
 */
export function getLangMapFromAppInfo(appInfo: any): Record<string, any> {
  const result: Record<string, any> = {}
  if (!appInfo)
    return result
  const langList = appInfo.langList || []
  langList.forEach((l: any) => {
    result[l.lang] = l
  })
  return result
}

/**
 * 获取指定语言下的应用名称（支持多语言回退）
 */
export function getAppNameByLang(langMap: Record<string, any>, currentLang: string, fallbackName?: string): string {
  return langMap[currentLang]?.appName
    || langMap['zh-CN']?.appName
    || fallbackName
    || ''
}

/**
 * 获取指定语言下的应用描述（支持多语言回退）
 */
export function getAppDescByLang(langMap: Record<string, any>, currentLang: string, fallbackDesc?: string): string {
  return langMap[currentLang]?.appDesc
    || langMap['zh-CN']?.appDesc
    || fallbackDesc
    || ''
}
