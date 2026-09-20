<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { NSpin } from 'naive-ui'
import DOMPurify from 'dompurify'
import Prism from '@/utils/cmn/prism'
// 代码块统一使用深色主题，浅色/深色界面下均有稳定对比度
import 'prismjs/themes/prism-tomorrow.min.css'
import {
  buildCopyText,
  copyToClipboard,
  detectCodeKind,
} from '@/utils/cmn/copyCode'

interface Props {
  // 富文本内容(HTML)
  content?: string
  // 内容加载中时显示占位动画
  loading?: boolean
  // 是否关闭自带的内边距
  noPadding?: boolean
  // 是否显示代码块复制按钮
  copyable?: boolean
  // 复制时写入注释的来源信息
  copyMeta?: {
    title: string
    author: string
    updateTime: string
    url: string
  }
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  noPadding: false,
  copyable: false,
})

// 编辑器内置 Prism 语言有限，用户常只能选默认的 markup，
// 因此当声明语言解析不出任何 token 时，用候选语言重新识别
const CANDIDATE_LANGUAGES = [
  'markup',
  'css',
  'clike',
  'javascript',
  'typescript',
  'json',
  'yaml',
  'bash',
  'go',
  'python',
  'java',
  'c',
  'cpp',
  'csharp',
  'php',
  'ruby',
  'sql',
  'markdown',
]

// 语言特征指纹，比单纯比较 token 数量更能凑准语言（例如 Go 与 Java 的 token 数接近）
const LANGUAGE_FINGERPRINTS: Array<[string, RegExp]> = [
  ['bash', /^\s*#!.*\b(bash|sh|zsh)\b/m],
  ['php', /^\s*<\?php/],
  ['go', /^\s*package\s+\w+[\s\S]*?\bimport\s*\(/m],
  ['python', /^\s*(?:from\s+[\w.]+\s+import|def\s+\w+\s*\()/m],
  ['java', /^\s*(?:public|protected|private)\s+(?:static\s+)?(?:final\s+)?(?:class|void|int|String)\b/m],
  // 对象必须以带引号的键开头，避免把 JS/Go 的对象字面量误判为 JSON
  ['json', /^\s*\{\s*"[^"]*"\s*:[\s\S]*\}\s*$|^\s*\[[\s\S]*\]\s*$/],
]

// 逐语言试解析的开销随内容长度增长，超长文本仅用少量候选语言
const AUTO_DETECT_MAX_LENGTH = 20000
const AUTO_DETECT_CANDIDATES_FOR_LARGE = ['clike', 'javascript', 'json', 'python', 'go']

const contentRef = ref<HTMLElement>()

// UGC 内容渲染前再做一次消毒（后端已消毒一次，这里是第二层）
const safeContent = computed(() => (props.content ? DOMPurify.sanitize(props.content) : ''))

function countTokens(html: string) {
  return html.match(/class="token/g)?.length ?? 0
}

// 声明语言与实际内容不符时（例如 Go 源码被标注为 markup），重新判断语言
function detectLanguage(code: string) {
  for (const [language, pattern] of LANGUAGE_FINGERPRINTS) {
    if (Prism.languages[language] && pattern.test(code))
      return language
  }

  // 指纹未命中时退化为按候选语言的解析效果择优
  const languages = code.length > AUTO_DETECT_MAX_LENGTH ? AUTO_DETECT_CANDIDATES_FOR_LARGE : CANDIDATE_LANGUAGES
  let bestLanguage: string | undefined
  let bestScore = 0
  for (const language of languages) {
    const grammar = Prism.languages[language]
    if (!grammar)
      continue
    const score = countTokens(Prism.highlight(code, grammar, language))
    if (score > bestScore) {
      bestScore = score
      bestLanguage = language
    }
  }
  return bestLanguage
}

function highlightCodeElement(el: HTMLElement) {
  const className = `${el.className} ${el.parentElement?.className ?? ''}`
  const declared = /language-([\w+#-]+)/.exec(className)?.[1]
  const code = el.textContent ?? ''

  let language = declared && Prism.languages[declared] ? declared : undefined
  let html = language ? Prism.highlight(code, Prism.languages[language], language) : ''

  if (!language || countTokens(html) === 0) {
    language = detectLanguage(code)
    html = language ? Prism.highlight(code, Prism.languages[language], language) : ''
  }

  if (html) {
    el.innerHTML = html
    el.dataset.highlighted = 'yes'
    // 复制时以实际使用的语言为准
    if (language)
      el.dataset.language = language
  }
}

function addCopyButton(pre: HTMLElement) {
  if (!props.copyable)
    return

  const wrapper = document.createElement('div')
  wrapper.className = 'html-render-code-block'
  pre.parentNode?.insertBefore(wrapper, pre)
  wrapper.appendChild(pre)

  const button = document.createElement('button')
  button.type = 'button'
  button.className = 'html-render-code-copy'
  button.textContent = '复制'
  button.addEventListener('click', async () => {
    const codeEl = pre.querySelector('code')
    const language = (codeEl as HTMLElement | null)?.dataset.language
      || /language-([\w+#-]+)/.exec(`${pre.className} ${codeEl?.className ?? ''}`)?.[1]
    const kind = detectCodeKind(language)
    const result = buildCopyText(codeEl?.textContent ?? '', kind, {
      title: props.copyMeta?.title ?? '',
      author: props.copyMeta?.author ?? '',
      updateTime: props.copyMeta?.updateTime ?? '',
      url: props.copyMeta?.url ?? '',
    })

    const ok = await copyToClipboard(result.text)
    button.textContent = ok ? '已复制' : '复制失败'
    window.setTimeout(() => {
      button.textContent = '复制'
    }, 1500)
  })

  wrapper.appendChild(button)
}

function highlightCodeBlocks() {
  const root = contentRef.value
  if (!root)
    return

  // 兼容：编辑器「预格式化」等写法生成的 <pre> 没有内层 <code>，补一层后再交给 Prism
  root.querySelectorAll('pre').forEach((pre) => {
    if (pre.querySelector('code'))
      return
    const code = document.createElement('code')
    const language = /language-([\w+#-]+)/.exec(pre.className)?.[1]
    if (language)
      code.className = `language-${language}`
    while (pre.firstChild)
      code.appendChild(pre.firstChild)
    pre.appendChild(code)
  })

  root.querySelectorAll('code[class*="language-"], [class*="language-"] code').forEach((el) => {
    if (el instanceof HTMLElement)
      highlightCodeElement(el)
  })

  // 复制按钮（放在高亮之后，避免影响代码内容）
  if (props.copyable)
    root.querySelectorAll('pre').forEach(pre => addCopyButton(pre as HTMLElement))
}

// flush: 'post' 确保在 v-html 渲染完成后再处理
watch(() => props.content, () => {
  highlightCodeBlocks()
}, { flush: 'post' })

onMounted(() => {
  highlightCodeBlocks()
})
</script>

<template>
  <div>
    <div v-if="loading" class="flex justify-center">
      <NSpin size="large" />
    </div>
    <div v-show="!loading" ref="contentRef" class="markdown-body" :class="{ 'p-[20px]': !noPadding }" v-html="safeContent" />
  </div>
</template>

<style>
.html-render-code-block {
  position: relative;
}

.html-render-code-copy {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 2px 8px;
  font-size: 12px;
  line-height: 20px;
  color: #c5c8c6;
  background-color: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 4px;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.15s ease-in-out;
}

.html-render-code-block:hover .html-render-code-copy {
  opacity: 1;
}

.html-render-code-copy:hover {
  background-color: rgba(255, 255, 255, 0.24);
}
</style>
