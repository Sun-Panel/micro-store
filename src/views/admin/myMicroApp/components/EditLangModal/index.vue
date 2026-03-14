<script setup lang="ts">
import { NButton, NCard, NForm, NFormItem, NInput, NModal, NSelect, NSpace } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { updateLang } from '@/api/admin/microApp'

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}>()

// 语言标签映射
const langLabelMap: Record<string, string> = {
  'zh-CN': '中文',
  'en-US': 'English',
  'ja-JP': '日本語',
  'ko-KR': '한국어',
}

interface Props {
  visible: boolean
  microAppId: number
  langMap: Record<string, { appName: string, appDesc: string }>
}

// 语言选项
const langOptions = [
  { label: '中文', value: 'zh-CN' },
  { label: 'English', value: 'en-US' },
  { label: '日本語', value: 'ja-JP' },
  { label: '한국어', value: 'ko-KR' },
]

// 已添加的语言列表
const addedLangs = ref<string[]>([])
const newLang = ref<string | null>(null)

// 本地语言数据
const localLangMap = ref<Record<string, { appName: string, appDesc: string }>>({})

// 弹窗显示状态
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => emit('update:visible', visible),
})

// 初始化数据
function initData() {
  localLangMap.value = JSON.parse(JSON.stringify(props.langMap || {}))
  addedLangs.value = Object.keys(localLangMap.value)

  // 如果没有语言，默认添加中文
  if (addedLangs.value.length === 0) {
    addedLangs.value = ['zh-CN']
    localLangMap.value = { 'zh-CN': { appName: '', appDesc: '' } }
  }
}

// 监听弹窗打开
watch(() => props.visible, (val) => {
  if (val) {
    initData()
  }
})

// 添加语言
function addLang(lang: string) {
  if (lang && !addedLangs.value.includes(lang)) {
    addedLangs.value.push(lang)
    localLangMap.value[lang] = { appName: '', appDesc: '' }
  }
  newLang.value = null
}

// 移除语言
function removeLang(lang: string) {
  if (addedLangs.value.length <= 1)
    return
  const index = addedLangs.value.indexOf(lang)
  if (index > -1) {
    addedLangs.value.splice(index, 1)
    delete localLangMap.value[lang]
  }
}

// 保存语言
async function submit() {
  try {
    await updateLang({
      id: props.microAppId,
      langMap: localLangMap.value,
    })
    emit('done')
    show.value = false
  }
  catch (error) {
    console.error(error)
  }
}

// 获取可用语言选项
const availableLangOptions = computed(() => {
  return langOptions.filter(l => !addedLangs.value.includes(l.value))
})
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 800px" title="编辑多语言">
    <!-- 语言列表：长卡片形式 -->
    <div class="space-y-4">
      <NCard v-for="lang in addedLangs" :key="lang" size="small">
        <template #header>
          <div class="flex justify-between items-center">
            <span class="font-bold">{{ langLabelMap[lang] || lang }}</span>
            <NButton v-if="addedLangs.length > 1" text size="small" type="error" @click="removeLang(lang)">
              删除
            </NButton>
          </div>
        </template>

        <NForm label-placement="top">
          <NFormItem :label="`${langLabelMap[lang] || lang} 应用名称`">
            <NInput
              v-model:value="localLangMap[lang].appName"
              :placeholder="`请输入应用名称 (${lang})`"
            />
          </NFormItem>
          <NFormItem :label="`${langLabelMap[lang] || lang} 应用简介`">
            <NInput
              v-model:value="localLangMap[lang].appDesc"
              type="textarea"
              :placeholder="`请输入应用简介 (${lang})`"
              :rows="2"
            />
          </NFormItem>
        </NForm>
      </NCard>
    </div>

    <!-- 添加语言 -->
    <div class="mt-4 flex items-center gap-2">
      <span>添加语言：</span>
      <NSelect
        v-model:value="newLang"
        :options="availableLangOptions"
        placeholder="选择语言"
        style="width: 150px;"
        @update:value="addLang"
      />
    </div>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="show = false">
          取消
        </NButton>
        <NButton type="success" @click="submit">
          保存
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
