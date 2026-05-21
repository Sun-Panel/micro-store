/**
 * 版本比较函数测试用例
 * 由于项目可能没有测试框架，这里提供手动测试示例
 */

import { 
  parseVersion, 
  compareVersions, 
  isVersionEqual, 
  isVersionGreater, 
  isVersionLess,
  isVersionGreaterOrEqual,
  isVersionLessOrEqual
} from './version'

// 测试 parseVersion
console.log('=== 测试 parseVersion ===')
console.log('parseVersion("1.0.0"):', parseVersion('1.0.0')) // [1, 0, 0]
console.log('parseVersion("1.0.159"):', parseVersion('1.0.159')) // [1, 0, 159]
console.log('parseVersion("2.1.3.4"):', parseVersion('2.1.3.4')) // [2, 1, 3, 4]
console.log('parseVersion("1.0"):', parseVersion('1.0')) // [1, 0]
console.log('parseVersion(""):', parseVersion('')) // []
console.log('parseVersion(undefined as any):', parseVersion(undefined as any)) // []

// 测试 compareVersions
console.log('\n=== 测试 compareVersions ===')
console.log('compareVersions("1.0.0", "1.0.0"):', compareVersions('1.0.0', '1.0.0')) // 0
console.log('compareVersions("1.0.0", "1.0.1"):', compareVersions('1.0.0', '1.0.1')) // -1
console.log('compareVersions("1.0.1", "1.0.0"):', compareVersions('1.0.1', '1.0.0')) // 1
console.log('compareVersions("1.0.0", "1.1.0"):', compareVersions('1.0.0', '1.1.0')) // -1
console.log('compareVersions("1.1.0", "1.0.0"):', compareVersions('1.1.0', '1.0.0')) // 1
console.log('compareVersions("1.0.0", "2.0.0"):', compareVersions('1.0.0', '2.0.0')) // -1
console.log('compareVersions("2.0.0", "1.0.0"):', compareVersions('2.0.0', '1.0.0')) // 1
console.log('compareVersions("1.0.0", "1.0.0.0"):', compareVersions('1.0.0', '1.0.0.0')) // 0
console.log('compareVersions("1.0.0.0", "1.0.0"):', compareVersions('1.0.0.0', '1.0.0')) // 0
console.log('compareVersions("1.0.0", "1.0"):', compareVersions('1.0.0', '1.0')) // 0
console.log('compareVersions("1.0", "1.0.0"):', compareVersions('1.0', '1.0.0')) // 0
console.log('compareVersions("1.0.159", "1.0.2"):', compareVersions('1.0.159', '1.0.2')) // 1
console.log('compareVersions("1.0.2", "1.0.159"):', compareVersions('1.0.2', '1.0.159')) // -1

// 测试辅助函数
console.log('\n=== 测试辅助函数 ===')
console.log('isVersionEqual("1.0.0", "1.0.0"):', isVersionEqual('1.0.0', '1.0.0')) // true
console.log('isVersionEqual("1.0.0", "1.0.1"):', isVersionEqual('1.0.0', '1.0.1')) // false
console.log('isVersionGreater("1.0.1", "1.0.0"):', isVersionGreater('1.0.1', '1.0.0')) // true
console.log('isVersionGreater("1.0.0", "1.0.1"):', isVersionGreater('1.0.0', '1.0.1')) // false
console.log('isVersionLess("1.0.0", "1.0.1"):', isVersionLess('1.0.0', '1.0.1')) // true
console.log('isVersionLess("1.0.1", "1.0.0"):', isVersionLess('1.0.1', '1.0.0')) // false
console.log('isVersionGreaterOrEqual("1.0.0", "1.0.0"):', isVersionGreaterOrEqual('1.0.0', '1.0.0')) // true
console.log('isVersionGreaterOrEqual("1.0.1", "1.0.0"):', isVersionGreaterOrEqual('1.0.1', '1.0.0')) // true
console.log('isVersionGreaterOrEqual("1.0.0", "1.0.1"):', isVersionGreaterOrEqual('1.0.0', '1.0.1')) // false
console.log('isVersionLessOrEqual("1.0.0", "1.0.0"):', isVersionLessOrEqual('1.0.0', '1.0.0')) // true
console.log('isVersionLessOrEqual("1.0.0", "1.0.1"):', isVersionLessOrEqual('1.0.0', '1.0.1')) // true
console.log('isVersionLessOrEqual("1.0.1", "1.0.0"):', isVersionLessOrEqual('1.0.1', '1.0.0')) // false

console.log('\n=== 测试完成 ===')