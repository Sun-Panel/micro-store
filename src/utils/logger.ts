/**
 * 日志工具 - 仅在开发模式下输出
 *
 * 特点：
 *   1. 彩色标签 [SunPanel]，方便在控制台过滤
 *   2. 自动捕获调用位置（文件:行号），可点击跳转
 *
 * 用法：
 *   import { log, warn, error } from '@/utils/logger'
 *   log('调试信息', data)
 *   warn('警告', data)
 *   error('错误', data)
 */

const TAG = '[SunPanelMicroAppStore]'
const TAG_STYLE = 'color:#2d8cf0;font-weight:bold'

/** 是否为开发模式 */
const isDev = import.meta.env.DEV

/** 从 Error 栈中解析出实际调用者的文件和行号 */
// function getCallerLocation(): string {
//   const stack = new Error('logger').stack
//   if (!stack)
//     return ''

//   const lines = stack.split('\n')
//   // lines[0] = Error
//   // lines[1] = getCallerLocation
//   // lines[2] = log/warn/error（本文件）
//   // lines[3] = 实际调用者
//   const callerLine = lines[3] || ''
//   const match = callerLine.match(/at\s+([^\s:]+):(\d+):\d+/)
//   if (match) {
//     const file = match[1].split('/').pop() // 取文件名部分
//     const line = match[2]
//     return `${file}:${line}`
//   }
//   return ''
// }

export function log(...args: any[]) {
  if (!isDev)
    return
  // eslint-disable-next-line no-console
  console.log(`%c${TAG}`, TAG_STYLE, ...args)
}

export function warn(...args: any[]) {
  if (!isDev)
    return
  // const loc = getCallerLocation()
  // console.warn(`%c${TAG}`, TAG_STYLE, `@${loc}`, ...args)
  console.warn(`%c${TAG}`, TAG_STYLE, ...args)
}

export function error(...args: any[]) {
  if (!isDev)
    return
  // const loc = getCallerLocation()
  console.error(`%c${TAG}`, TAG_STYLE, ...args)
}
export function debug(eventName: string, ...data: any) {
  log(`[${eventName}]`, ...data)
}
