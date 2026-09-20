/**
 * 代码复制加工
 * - JS：复制时检测是否已有闭包匿名函数（IIFE），没有则补上
 * - JS / CSS / 页脚HTML：复制时加上来源注释（标题、作者、更新时间、地址）
 * - 其他语言：原样复制
 */

import type { SnipType } from '@/utils/customCodeSnippet'
import { buildSnippetMarkers } from '@/utils/customCodeSnippet'

export type CodeKind = 'js' | 'css' | 'html' | 'other'

export interface CopyCodeMeta {
  title: string
  author: string
  updateTime: string
  url: string
  /** 自定义代码帖子 id（用于复制标识） */
  id?: number
  /** 块唯一标识（用于复制标识） */
  onlyId?: string
  /** 片段更新时间（Unix 秒，用于复制标识） */
  updateTimeUnix?: number
}

export interface CopyCodeResult {
  text: string
  /** JS 是否补了闭包 */
  wrapped: boolean
  /** 需要提示用户的说明 */
  note: string
}

const LANGUAGE_KIND_MAP: Record<string, CodeKind> = {
  javascript: 'js',
  js: 'js',
  typescript: 'js',
  ts: 'js',
  jsx: 'js',
  css: 'css',
  scss: 'css',
  less: 'css',
  markup: 'html',
  html: 'html',
  xml: 'html',
  svg: 'html',
}

export function detectCodeKind(language?: string): CodeKind {
  if (!language)
    return 'other'
  return LANGUAGE_KIND_MAP[language.toLowerCase()] ?? 'other'
}

// 块的代码类型 -> 复制加工使用的语言类型（结构化存储后不再需要猜测）
const CODE_TYPE_KIND_MAP: Record<number, CodeKind> = {
  1: 'js', // 自定义JS
  2: 'css', // 自定义CSS
  3: 'html', // 自定义页脚
}

// 语言类型 -> 标记文本类型（复制标识用；页脚在标记中写作 FOOTER）
// 注：'other' 已在 buildCopyText 提前返回，不需映射
const KIND_MARKER_TYPE: Record<Exclude<CodeKind, 'other'>, SnipType> = { js: 'JS', css: 'CSS', html: 'FOOTER' }

export function kindFromCodeType(codeType: number): CodeKind {
  return CODE_TYPE_KIND_MAP[codeType] ?? 'other'
}

/** 是否已经是闭包匿名函数(IIFE) */
export function hasIIFE(code: string): boolean {
  const trimmed = code.trim()
  if (!trimmed)
    return false

  // 结尾形式：})() / })() ; / }()) / )()
  const callTail = /(?:\}|\)|=>\s*\})\s*\(\s*\)\s*(?:\)\s*)?;?\s*$/.test(trimmed)
  if (!callTail)
    return false

  // 开头形式：(function / (async function / (() => / (async () => / !function / ~+- function
  const heads: RegExp[] = [
    /^\s*(?:;\s*)?\(\s*(?:async\s+)?function\b/,
    /^\s*(?:;\s*)?\(\s*(?:async\s*)?\(\s*\)\s*=>/,
    /^\s*(?:;\s*)?\(\s*(?:async\s*)?\w+\s*=>/,
    /^\s*[!~+-]\s*function\b/,
  ]

  return heads.some(head => head.test(trimmed))
}

/** 是否为 ES Module（含顶层 import / export，不能包闭包） */
export function hasESModule(code: string): boolean {
  return /^\s*(?:import|export)\s/m.test(code)
}

/** 包一层闭包匿名函数 */
export function wrapWithIIFE(code: string): string {
  return `;(function () {\n  'use strict';\n\n${code}\n})();`
}

function buildCommentHeader(kind: CodeKind, meta: CopyCodeMeta): string {
  const lines = [
    `Title: ${meta.title}`,
    `Author: ${meta.author}`,
    `UpdateTime: ${meta.updateTime}`,
    `URL: ${meta.url}`,
  ].join('\n')

  // 注释符号按语言区分：HTML 用 <!-- -->，JS/CSS 用 /* */
  if (kind === 'html')
    return `<!--\n${lines}\n-->`

  return `/*\n${lines}\n*/`
}

/** 生成最终复制到剪贴板的文本 */
export function buildCopyText(code: string, kind: CodeKind, meta: CopyCodeMeta): CopyCodeResult {
  if (kind === 'other')
    return { text: code, wrapped: false, note: '' }

  let body = code
  let wrapped = false
  let note = ''

  if (kind === 'js') {
    if (hasESModule(code)) {
      note = '该代码为 ES Module，未添加闭包'
    }
    else if (!hasIIFE(code)) {
      body = wrapWithIIFE(code)
      wrapped = true
    }
  }

  let text = `${buildCommentHeader(kind, meta)}\n${body}`

  // 包裹复制标识头尾（提供 id / onlyId 时生效），用于跨项目粘贴去重
  if (meta.id && meta.onlyId) {
    const type = KIND_MARKER_TYPE[kind as Exclude<CodeKind, 'other'>]
    const { start, end } = buildSnippetMarkers(type, meta.id, meta.onlyId, meta.updateTimeUnix ?? 0)
    text = `${start}\n${text}\n${end}`
  }

  return {
    text,
    wrapped,
    note,
  }
}

/** 写入剪贴板（优先 Clipboard API，失败降级 execCommand） */
export async function copyToClipboard(text: string): Promise<boolean> {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text)
      return true
    }
  }
  catch {
    // 继续走降级方案
  }

  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.setAttribute('readonly', 'readonly')
  textarea.style.position = 'fixed'
  textarea.style.top = '-9999px'
  document.body.appendChild(textarea)
  textarea.select()
  const ok = document.execCommand('copy')
  document.body.removeChild(textarea)
  return ok
}
