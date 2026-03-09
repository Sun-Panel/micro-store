<script setup lang="ts">
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { onMounted, onUnmounted, ref, toRaw, unref, watch } from 'vue'

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

onMounted(() => {
  contentEditor.value = new Vditor(editorRef.value as HTMLElement, {
    ...props.options,
    value: props.modelValue,
    after() {
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
  <div ref="editorRef" />
</template>
