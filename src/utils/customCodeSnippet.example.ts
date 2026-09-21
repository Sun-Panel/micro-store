/**
 * 其他项目「粘贴导入自定义代码片段」示例
 * --------------------------------------------------------------------------
 * 本文件仅作示例，不参与本项目业务（本项目不实现粘贴触发）。
 * 其他项目引入 `customCodeSnippet.ts` 后，可照此实现：用户从本商店复制带
 * 标识的代码，粘贴到「自定义板块」时自动识别，按 onlyId 去重：
 *   - 已存在相同 onlyId 的块 -> 更新（删旧代码块，写入新代码块）
 *   - 不存在 -> 新增一个块
 * 从而避免不知情地添加多个相同代码块。
 *
 * 复用方式：
 *   import { parseSnippet, cleanSnippetCode } from '@/utils/customCodeSnippet'
 *   import { applyPastedSnippet } from '@/utils/customCodeSnippet.example'
 */

import type { ParsedSnippet, SnipType } from './customCodeSnippet'
import { cleanSnippetCode, parseSnippet } from './customCodeSnippet'

/** 其他项目里「自定义板块」的代码块结构（按需替换字段名） */
export interface OtherProjectBlock {
  /** 唯一标识（对应本商店 {key}-{onlyId} 中的 onlyId） */
  onlyId: string
  /** 代码类型：JS / CSS / FOOTER（与来源标记保持一致；FOOTER 即页脚） */
  codeType: SnipType
  /** 干净代码（已剥离来源注释头与 IIFE 包裹） */
  code: string
  /** 可选：保留来源信息，方便展示 */
  source?: {
    /** 来源自定义代码唯一标识（开发者标识-后缀） */
    key: string
    title?: string
    updateTime?: number
  }
}

export interface PasteResult {
  /** 是否识别为本商店复制格式 */
  matched: boolean
  /** 操作类型 */
  action: 'added' | 'updated' | 'ignored'
  /** 提示文案 */
  message: string
  /** 处理后的块列表（已更新） */
  blocks: OtherProjectBlock[]
}

/**
 * 将剪贴板文本应用到现有块列表：解析 -> 去重 -> 新增/更新
 * @param text       剪贴板文本
 * @param blocks     当前「自定义板块」的代码块列表（会被原地合并后返回新数组）
 */
export function applyPastedSnippet(
  text: string,
  blocks: OtherProjectBlock[],
): PasteResult {
  const parsed: ParsedSnippet | null = parseSnippet(text)
  if (!parsed) {
    // 非本商店复制格式：交还给上层按普通粘贴处理
    return { matched: false, action: 'ignored', message: '', blocks }
  }

  const cleanCode = cleanSnippetCode(parsed.code)

  const idx = blocks.findIndex(b => b.onlyId === parsed.onlyId)
  if (idx >= 0) {
    // 已存在 -> 更新
    const next = blocks.slice()
    next[idx] = {
      ...next[idx],
      codeType: parsed.type,
      code: cleanCode,
      source: { key: parsed.key, updateTime: parsed.updateTime },
    }
    return {
      matched: true,
      action: 'updated',
      message: `已更新片段（${parsed.onlyId}）`,
      blocks: next,
    }
  }

  // 不存在 -> 新增
  const next = blocks.slice()
  next.push({
    onlyId: parsed.onlyId,
    codeType: parsed.type,
    code: cleanCode,
    source: { key: parsed.key, updateTime: parsed.updateTime },
  })
  return {
    matched: true,
    action: 'added',
    message: `已新增片段（${parsed.onlyId}）`,
    blocks: next,
  }
}

/*
 * ⚠️ 关于 IIFE 还原：本商店复制 JS 时会补 `;(function () { 'use strict'; ... })();`
 *    `cleanSnippetCode` 已帮你剥掉该闭包，得到开发者原始代码。若你的项目直接
 *    原样注入（不二次包闭包），使用 cleanSnippetCode 返回的 code 即可；若你的
 *    项目有自己独立的沙箱包裹逻辑，也可直接用 parsed.code（含 IIFE）自行处理。
 *
 * 使用示意（伪代码）：
 *   onPaste(event) {
 *     const text = event.clipboardData.getData('text')
 *     const res = applyPastedSnippet(text, myBlocks.value)
 *     if (res.matched) {
 *       event.preventDefault()        // 拦截默认粘贴，避免把标识头尾也贴进去
 *       myBlocks.value = res.blocks
 *       message.success(res.message)
 *     }
 *     // 未匹配则走默认粘贴（普通代码）
 *   }
 */
