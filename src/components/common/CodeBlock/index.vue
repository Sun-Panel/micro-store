<script setup lang="ts">
import { NButton, NSpace, NTag } from 'naive-ui'
import { computed, ref } from 'vue'
import { buildCopyText, copyToClipboard, kindFromCodeType } from '@/utils/cmn/copyCode'
import Prism from '@/utils/cmn/prism'
import 'prismjs/themes/prism-tomorrow.min.css'

interface Props {
  /** 代码片段（纯文本） */
  code: string
  /** 代码类型：1-JS 2-CSS 3-页脚 */
  codeType: number
  /** 是否默认折叠 */
  defaultCollapsed?: boolean
  /** 复制到注释里的来源信息 */
  copyMeta?: {
    title: string
    author: string
    updateTime: string
    url: string
    /** 自定义代码唯一标识（开发者标识-后缀） */
    uniqueKey?: string
    /** 块唯一标识 */
    onlyId?: string
    /** 片段更新时间（Unix 秒） */
    updateTimeUnix?: number
  }
}

const props = withDefaults(defineProps<Props>(), {
  defaultCollapsed: true,
})

const collapsed = ref(props.defaultCollapsed)
const copying = ref(false)
const copied = ref(false)

const codeTypeLabelMap: Record<number, string> = {
  1: 'JS',
  2: 'CSS',
  3: '页脚',
}

const languageMap: Record<number, string> = {
  1: 'javascript',
  2: 'css',
  3: 'markup',
}

const codeTypeLabel = computed(() => codeTypeLabelMap[props.codeType] ?? '代码')

// Prism 会对内容做转义，因此可以安全地用于 v-html
const highlightedHtml = computed(() => {
  const language = languageMap[props.codeType]
  const grammar = language ? Prism.languages[language] : undefined
  if (!grammar)
    return ''
  return Prism.highlight(props.code, grammar, language)
})

async function handleCopy() {
  const result = buildCopyText(props.code, kindFromCodeType(props.codeType), {
    title: props.copyMeta?.title ?? '',
    author: props.copyMeta?.author ?? '',
    updateTime: props.copyMeta?.updateTime ?? '',
    url: props.copyMeta?.url ?? '',
    uniqueKey: props.copyMeta?.uniqueKey,
    onlyId: props.copyMeta?.onlyId,
    updateTimeUnix: props.copyMeta?.updateTimeUnix,
  })

  copying.value = true
  const ok = await copyToClipboard(result.text)
  copying.value = false
  copied.value = ok

  if (result.note)
    window.alert(result.note)

  if (ok) {
    window.setTimeout(() => {
      copied.value = false
    }, 1500)
  }
}
</script>

<template>
  <div class="border border-slate-200 dark:border-slate-700 rounded-[8px] overflow-hidden">
    <div class="flex items-center gap-[8px] px-[12px] py-[8px] bg-slate-50 dark:bg-slate-800">
      <NTag size="small">
        {{ codeTypeLabel }}
      </NTag>
      <NSpace :size="8" class="ml-auto">
        <!-- 复制按钮常驻，折叠状态下同样可用 -->
        <NButton size="tiny" :loading="copying" @click="handleCopy">
          {{ copied ? '已复制' : '复制' }}
        </NButton>
        <NButton size="tiny" @click="collapsed = !collapsed">
          {{ collapsed ? '展开代码' : '收起代码' }}
        </NButton>
      </NSpace>
    </div>

    <div v-show="!collapsed" class="code-block-content">
      <pre><code v-if="!highlightedHtml">{{ code }}</code><code v-else v-html="highlightedHtml" /></pre>
    </div>
  </div>
</template>

<style>
.code-block-content pre {
  margin: 0;
  padding: 12px;
  overflow-x: auto;
  background-color: #2d2d2d;
  color: #ccc;
}
</style>
