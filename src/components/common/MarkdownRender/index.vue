<script setup lang="ts">
import { watch } from 'vue'
import vditor from 'vditor'
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
  const options: IPreviewOptions = {
    hljs: { style: 'monokai' },
    mode: props.mode,
    after() {
      emits('update:loading', false)
    },
  }

  vditor.preview(document.getElementById(renderId) as HTMLDivElement, md, options)
}

watch(() => props.content, (newValue: string) => {
  previewMarkdown(newValue)
})
</script>

<template>
  <NSpin size="small" :show="loading">
    <div :id="renderId" />
  </NSpin>
</template>
