/**
 * 版本号比较工具函数
 * 支持格式：1.0.0、1.0.159、2.1.3.4 等
 * 从左向右进行版本号对比大小（语义化版本号标准比较方式）
 */

/**
 * 将版本号字符串解析为数字数组
 * @param version 版本号字符串，如 "1.0.0"
 * @returns 数字数组，如 [1, 0, 0]
 */
export function parseVersion(version: string): number[] {
  if (!version || typeof version !== 'string') {
    return []
  }
  
  return version
    .split('.')
    .map(part => {
      const num = parseInt(part, 10)
      return isNaN(num) ? 0 : num
    })
}

/**
 * 从左向右比较两个版本号（标准语义化版本号比较方式）
 * @param version1 第一个版本号
 * @param version2 第二个版本号
 * @returns -1: version1 < version2, 0: version1 === version2, 1: version1 > version2
 */
export function compareVersions(version1: string, version2: string): -1 | 0 | 1 {
  const parts1 = parseVersion(version1)
  const parts2 = parseVersion(version2)
  
  // 获取最大长度，补零对齐
  const maxLength = Math.max(parts1.length, parts2.length)
  
  // 从左向右逐位比较
  for (let i = 0; i < maxLength; i++) {
    const num1 = parts1[i] || 0
    const num2 = parts2[i] || 0
    
    if (num1 > num2) {
      return 1
    } else if (num1 < num2) {
      return -1
    }
  }
  
  return 0
}

/**
 * 检查版本1是否大于版本2
 * @param version1 第一个版本号
 * @param version2 第二个版本号
 * @returns boolean
 */
export function isVersionGreater(version1: string, version2: string): boolean {
  return compareVersions(version1, version2) === 1
}

/**
 * 检查版本1是否小于版本2
 * @param version1 第一个版本号
 * @param version2 第二个版本号
 * @returns boolean
 */
export function isVersionLess(version1: string, version2: string): boolean {
  return compareVersions(version1, version2) === -1
}

/**
 * 检查版本1是否等于版本2
 * @param version1 第一个版本号
 * @param version2 第二个版本号
 * @returns boolean
 */
export function isVersionEqual(version1: string, version2: string): boolean {
  return compareVersions(version1, version2) === 0
}

/**
 * 检查版本1是否大于等于版本2
 * @param version1 第一个版本号
 * @param version2 第二个版本号
 * @returns boolean
 */
export function isVersionGreaterOrEqual(version1: string, version2: string): boolean {
  const result = compareVersions(version1, version2)
  return result === 1 || result === 0
}

/**
 * 检查版本1是否小于等于版本2
 * @param version1 第一个版本号
 * @param version2 第二个版本号
 * @returns boolean
 */
export function isVersionLessOrEqual(version1: string, version2: string): boolean {
  const result = compareVersions(version1, version2)
  return result === -1 || result === 0
}

/**
 * 从版本字符串中提取基础版本号（主版本.次版本.修订号）
 * 例如：
 * - "2.0.0" → "2.0.0"
 * - "2.0.0-beta260505" → "2.0.0"
 * - "2.0.0-dev-1" → "2.0.0"
 * - "1.0.0-alpha.1" → "1.0.0"
 * @param version 完整版本字符串
 * @returns 基础版本号（x.y.z）
 */
export function extractBaseVersion(version: string): string {
  if (!version || typeof version !== 'string') {
    return '0.0.0'
  }
  // 匹配前三个数字部分：主版本.次版本.修订号
  const match = version.match(/^(\d+)\.(\d+)\.(\d+)/)
  if (match) {
    return `${match[1]}.${match[2]}.${match[3]}`
  }
  return version
}