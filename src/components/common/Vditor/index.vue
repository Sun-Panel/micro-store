<script setup lang="ts">
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { NSpin } from 'naive-ui'
import { computed, onMounted, onUnmounted, ref, toRaw, unref, watch } from 'vue'

const props = defineProps({
  options: {
    type: Object,
  },
  modelValue: {
    type: String,
    default: '',
  },
})

const emit = defineEmits([
  'update:modelValue',
  'after',
  'focus',
  'blur',
  'esc',
  'ctrlEnter',
  'select',
])

const contentEditor = ref<Vditor | null>()
const editorRef = ref<string | HTMLElement>()
// 编辑器异步初始化完成标志：初始化完成前写入的值会丢失，需等就绪后再写入
const isReady = ref(false)
// 初始化期间预留高度，避免 loading 区域塌陷（取 options.height，缺省 460）
const placeholderHeight = computed(() => (props.options as any)?.height || 460)

onMounted(() => {
  contentEditor.value = new Vditor(editorRef.value as HTMLElement, {
    ...props.options,
    value: props.modelValue,
    after() {
      isReady.value = true
      // 就绪后再写入最新值，避免初始化完成前 setValue 被覆盖导致内容丢失
      contentEditor.value?.setValue(props.modelValue)
      emit('after', toRaw(contentEditor.value))
    },
    input(value: string) {
      emit('update:modelValue', value)
    },
    focus(value: string) {
      emit('focus', value)
    },
    blur(value: string) {
      emit('blur', value)
    },
    esc(value: string) {
      emit('esc', value)
    },
    ctrlEnter(value: string) {
      emit('ctrlEnter', value)
    },
    select(value: string) {
      emit('select', value)
    },
  })
})

watch(
  () => props.modelValue,
  (newVal) => {
    // 编辑器未就绪时不写入，等待 after() 统一处理，避免竞争丢失内容
    if (!isReady.value)
      return
    if (newVal !== contentEditor.value?.getValue())
      contentEditor.value?.setValue(newVal)
  },
)

onUnmounted(() => {
  const editorInstance = unref(contentEditor)
  if (!editorInstance)
    return
  try {
    editorInstance?.destroy?.()
  }
  catch (error) {
    console.log(error)
  }
})
</script>

<template>
  <NSpin :show="!isReady" size="small">
    <!-- 未就绪时预留高度，避免 spinner 区域塌陷 -->
    <div ref="editorRef" :style="{ minHeight: isReady ? 'auto' : placeholderHeight + 'px' }" />
  </NSpin>
</template>
