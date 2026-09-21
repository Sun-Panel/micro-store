/**
 * 自定义代码片段「复制标识」的封装与解析（跨项目复用）
 *
 * 标记格式（以 JS / CSS 为例，页脚类型用 <!-- -->）：
 *   /* ====START===={type}===={key}-{onlyId}===={updateTime}==== *\/
 *   <源码（含来源注释头，JS 可能被包 IIFE）>
 *   /* ====END===={type}===={key}-{onlyId}===={updateTime}==== *\/
 *
 * 字段说明：
 *   type      - 代码类型：JS / CSS / FOOTER(页脚)
 *   key       - 自定义代码唯一标识（开发者标识-后缀，可能含 -，故按最后一个 - 切分）
 *   onlyId    - 块唯一标识（仅含字母/数字/下划线，不含 -，避免与连接符冲突）
 *   updateTime- 片段更新时间（Unix 秒），用于判断是否过期，不参与去重
 *
 * 去重依据：{key}-{onlyId}（同一自定义代码内唯一）。复制时包裹头尾，
 * 粘贴时解析头尾提取 onlyId，已存在则更新、不存在则新增。
 */

export const SNIPPET_SEP = '===='

/** 代码类型（标记文本用，非数字）：JS / CSS / FOOTER(页脚) */
export type SnipType = 'JS' | 'CSS' | 'FOOTER'

export interface SnippetMarkerMeta {
  /** JS / CSS / FOOTER */
  type: SnipType
  /** 自定义代码唯一标识（开发者标识-后缀，可能含 -） */
  key: string
  /** 块唯一标识（不含 -） */
  onlyId: string
  /** 片段更新时间（Unix 秒） */
  updateTime: number
}

export interface ParsedSnippet extends SnippetMarkerMeta {
  /** START 与 END 之间的原始内容（含来源注释头，未剥离 IIFE） */
  code: string
}

/** 唯一标识合法字符集：字母/数字/下划线，最长 40，不含 - */
export function isValidOnlyId(onlyId: string): boolean {
  return /^[A-Za-z0-9_]{1,40}$/.test(onlyId)
}

/** 生成块唯一标识（12 位十六进制） */
export function genOnlyId(): string {
  const buf = new Uint8Array(6)
  if (typeof crypto !== 'undefined' && crypto.getRandomValues) {
    crypto.getRandomValues(buf)
  }
  else {
    for (let i = 0; i < buf.length; i++)
      buf[i] = Math.floor(Math.random() * 256)
  }
  return Array.from(buf, b => b.toString(16).padStart(2, '0')).join('')
}

function commentWrap(inner: string, type: SnipType): string {
  return type === 'FOOTER' ? `<!-- ${inner} -->` : `/* ${inner} */`
}

/** 生成头/尾标记行 */
export function buildSnippetMarkers(
  type: SnipType,
  key: string,
  onlyId: string,
  updateTime: number,
): { start: string; end: string } {
  const body = `${key}-${onlyId}`
  const inner = (token: 'START' | 'END') =>
    `${SNIPPET_SEP}${token}${SNIPPET_SEP}${type}${SNIPPET_SEP}${body}${SNIPPET_SEP}${updateTime}${SNIPPET_SEP}`
  return {
    start: commentWrap(inner('START'), type),
    end: commentWrap(inner('END'), type),
  }
}

function parseMarkerLine(line: string): { token: 'START' | 'END'; meta: SnippetMarkerMeta } | null {
  // 去掉注释定界符 /* */ 或 <!-- -->（注意页脚用的是 <!-- -->，不是 ---）
  const cm = line.match(/^\s*(?:\/\*|<!--)\s*([\s\S]*?)\s*(?:\*\/|-->\s*)$/)
  if (!cm)
    return null
  const inner = cm[1]
  const sm = inner.match(/^====(START|END)====([\s\S]*)====$/)
  if (!sm)
    return null

  const parts = sm[2].split(SNIPPET_SEP)
  if (parts.length !== 3)
    return null

  const type = parts[0] as SnipType
  if (type !== 'JS' && type !== 'CSS' && type !== 'FOOTER')
    return null
  // 唯一标识可能含 -，按最后一个 - 切分（onlyId 不含 -）
  const idOnlyId = parts[1]
  const dash = idOnlyId.lastIndexOf('-')
  if (dash <= 0)
    return null
  const key = idOnlyId.slice(0, dash)
  const onlyId = idOnlyId.slice(dash + 1)
  const updateTime = Number(parts[2])
  if (!Number.isFinite(updateTime) || !isValidOnlyId(onlyId))
    return null

  return {
    token: sm[1] as 'START' | 'END',
    meta: { type, key, onlyId, updateTime },
  }
}

/** 解析被复制包裹过的代码片段；无法识别时返回 null */
export function parseSnippet(text: string): ParsedSnippet | null {
  const lines = text.split(/\r?\n/)
  let startMeta: SnippetMarkerMeta | null = null
  let startIdx = -1

  for (let i = 0; i < lines.length; i++) {
    const r = parseMarkerLine(lines[i])
    if (!r)
      continue

    if (r.token === 'START') {
      startMeta = r.meta
      startIdx = i
    }
    else if (r.token === 'END' && startMeta && startIdx >= 0) {
      // 校验 END 与 START 元信息一致，避免代码体内误含标记行被误判为结尾
      if (r.meta.key === startMeta.key && r.meta.onlyId === startMeta.onlyId
        && r.meta.type === startMeta.type && r.meta.updateTime === startMeta.updateTime) {
        const code = lines.slice(startIdx + 1, i).join('\n')
        return { ...startMeta, code }
      }
    }
  }

  return null
}

/** 去除开头的来源注释头（/* ... *\/ 或 <!-- ... -->） */
export function stripSourceComment(code: string): string {
  const m = code.match(/^\s*(\/\*[\s\S]*?\*\/|<!--[\s\S]*?-->)/)
  if (!m)
    return code
  return code.slice(m[0].length).replace(/^\r?\n/, '')
}

/** 剥离复制时为 JS 补的 IIFE 闭包（还原原始代码） */
export function unwrapIIFE(code: string): string {
  const t = code.trim()
  const m = t.match(/^;\s*\(function\s*\(\)\s*\{\s*'use strict';\s*([\s\S]*?)\s*\}\)\(\);\s*$/)
  if (m)
    return m[1]
  return code
}

/** 还原为干净代码：去来源注释头 + 去 IIFE 包裹 */
export function cleanSnippetCode(code: string): string {
  return unwrapIIFE(stripSourceComment(code))
}
