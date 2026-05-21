/**
 * 版本号比较工具函数
 * 支持格式：1.0.0、1.0.159、2.1.3.4 等
 * 从右向左进行版本号对比大小
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
 * 从右向左比较两个版本号
 * @param version1 第一个版本号
 * @param version2 第二个版本号
 * @returns -1: version1 < version2, 0: version1 === version2, 1: version1 > version2
 */
export function compareVersions(version1: string, version2: string): -1 | 0 | 1 {
  const parts1 = parseVersion(version1)
  const parts2 = parseVersion(version2)
  
  // 获取最大长度，确保从右向左比较
  const maxLength = Math.max(parts1.length, parts2.length)
  
  // 从右向左比较
  for (let i = maxLength - 1; i >= 0; i--) {
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