<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import vditor from 'vditor'
import DOMPurify from 'dompurify'
import { NSpin } from 'naive-ui'

interface Props {
  content: string
  mode: 'dark' | 'light'
  loading: boolean
}

const props = withDefaults(defineProps<Props>(), { loading: false })
// 定义事件
const emits = defineEmits([
  'update:loading',
])

// vditor.preview 是异步且较慢的渲染，渲染期间展示 loading，避免弹窗内长时间空白
const isRendering = ref(false)
const showSpin = computed(() => props.loading || isRendering.value)

const renderId = `vditor-render-${generateRandomString(5)}-${generateRandomString(10)}`

function generateRandomString(length: number) {
  const characters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let result = ''
  for (let i = 0; i < length; i++) {
    const randomIndex = Math.floor(Math.random() * characters.length)
    result += characters[randomIndex]
  }
  return result
}

function previewMarkdown(md: string) {
  isRendering.value = true
  const options: IPreviewOptions = {
    hljs: { style: 'monokai' },
    mode: props.mode,
    after() {
      isRendering.value = false
      emits('update:loading', false)
    },
  }
  // 净化渲染内容，防止用户提交的 Markdown 注入 XSS
  // （vditor 的类型未声明 sanitize 选项，使用断言注入）
  ;(options as any).sanitize = (html: string) => DOMPurify.sanitize(html)

  vditor.preview(document.getElementById(renderId) as HTMLDivElement, md, options)
}

watch(() => props.content, (newValue: string) => {
  previewMarkdown(newValue)
})

// 初始内容需在挂载时渲染：NModal 弹窗懒加载，组件挂载时 content 已就绪，
// 而 watch 默认不在首次触发，若不在此渲染弹窗内将始终空白。
onMounted(() => {
  if (props.content) {
    previewMarkdown(props.content)
  }
})
</script>

<template>
  <NSpin size="small" :show="showSpin">
    <div :id="renderId" />
  </NSpin>
</template>
