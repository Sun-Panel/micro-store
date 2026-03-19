<script lang="ts" setup>
import type { DataTableColumns } from 'naive-ui'
import { NButton, NDataTable, NModal, NSelect, NSpace, NTag, useMessage } from 'naive-ui'
import { computed, h, ref, watch } from 'vue'
import { updateLang, update as updateMicroApp } from '@/api/admin/microAppDeveloper'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'

interface Props {
  visible: boolean
  versionInfo: any
  microAppInfo: MicroApp.MicroAppInfo | undefined
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:visible': [value: boolean]
  'done': []
}>()

const message = useMessage()
const loading = ref(false)
const currentLang = ref('zh-CN')

// 版本详情信息
const detailInfo = ref<{
  iconURL: string
  appName: string
  appDesc: string
  microAppId: string
  author: string
  version: string
  versionDesc: string
  permissions: string[]
  dataNodes: Record<string, any>
  networkDomains: string[]
  appInfo: Record<string, { appName: string, appDesc: string }>
} | null>(null)

// 可用的语言列表
const langList = computed(() => {
  if (!detailInfo.value?.appInfo)
    return []
  return Object.keys(detailInfo.value.appInfo)
})

// 当前语言下的应用名称
const currentAppName = computed(() => {
  if (!detailInfo.value?.appInfo)
    return ''
  return detailInfo.value.appInfo[currentLang.value]?.appName
    || detailInfo.value.appInfo['zh-CN']?.appName
    || ''
})

// 当前语言下的应用描述
const currentAppDesc = computed(() => {
  if (!detailInfo.value?.appInfo)
    return ''
  return detailInfo.value.appInfo[currentLang.value]?.appDesc
    || detailInfo.value.appInfo['zh-CN']?.appDesc
    || ''
})

// 多语言表格列配置
const langTableColumns: DataTableColumns<{ lang: string, appName: string, appDesc: string }> = [
  { title: '语言', key: 'lang', width: 100 },
  { title: '应用名称', key: 'appName' },
  { title: '应用描述', key: 'appDesc', ellipsis: { tooltip: true } },
]

// 双向绑定
const show = ref(props.visible)
watch(() => props.visible, (val) => {
  show.value = val
  if (val) {
    initDetailInfo()
  }
})
watch(show, (val) => {
  emit('update:visible', val)
})

// 初始化详情信息
function initDetailInfo() {
  if (!props.versionInfo || !props.microAppInfo)
    return

  const config = props.versionInfo.config

  // 从微应用信息中获取多语言数据作为后备
  const microAppLangList = props.microAppInfo.langList || []
  const microAppLangMap: Record<string, { appName: string, appDesc: string }> = {}
  microAppLangList.forEach((lang: any) => {
    microAppLangMap[lang.lang] = {
      appName: lang.appName || '',
      appDesc: lang.appDesc || '',
    }
  })

  // 如果微应用没有多语言数据，使用基本信息
  if (Object.keys(microAppLangMap).length === 0 && props.microAppInfo) {
    microAppLangMap['zh-CN'] = {
      appName: props.microAppInfo.appName || '',
      appDesc: props.microAppInfo.appDesc || '',
    }
  }

  // 转换 appInfo 格式
  const versionAppInfo = config?.appInfo || {}
  const formattedAppInfo: Record<string, { appName: string, appDesc: string }> = {}
  for (const [lang, info] of Object.entries(versionAppInfo)) {
    formattedAppInfo[lang] = {
      appName: (info as any).appName || '',
      appDesc: (info as any).description || '',
    }
  }

  detailInfo.value = {
    iconURL: config?.icon || props.microAppInfo.appIcon || '',
    appName: '',
    appDesc: '',
    microAppId: config?.microAppId || props.microAppInfo.microAppId || '',
    author: config?.author || '',
    version: props.versionInfo.version || '',
    versionDesc: props.versionInfo.versionDesc || '',
    permissions: config?.permissions || [],
    dataNodes: {},
    networkDomains: config?.networkDomains || [],
    appInfo: Object.keys(formattedAppInfo).length > 0 ? formattedAppInfo : microAppLangMap,
  }

  // 设置默认语言
  const langs = Object.keys(detailInfo.value.appInfo)
  currentLang.value = langs.includes('zh-CN') ? 'zh-CN' : (langs[0] || 'zh-CN')
}

// 设置为主信息
async function handleSetAsMainInfo() {
  if (!detailInfo.value || !props.microAppInfo)
    return

  loading.value = true

  try {
    // 准备语言信息
    const versionLangMap: Record<string, { appName: string, appDesc: string }> = {}

    // 从版本详情中获取多语言信息
    const appInfo = detailInfo.value.appInfo || {}
    for (const [lang, info] of Object.entries(appInfo)) {
      versionLangMap[lang] = {
        appName: (info as any).appName || '',
        appDesc: (info as any).appDesc || '',
      }
    }

    // 如果没有多语言信息，使用默认
    if (Object.keys(versionLangMap).length === 0) {
      versionLangMap['zh-CN'] = {
        appName: detailInfo.value.appName || '',
        appDesc: detailInfo.value.appDesc || '',
      }
    }

    // 1. 更新语言信息
    const langRes = await updateLang({
      id: props.microAppInfo.id,
      langMap: versionLangMap,
    } as any)

    if (langRes.code !== 0) {
      apiRespErrMsg(langRes)
      return
    }

    // 2. 更新图标（如果版本中有图标）
    const newIcon = detailInfo.value.iconURL
    if (newIcon) {
      const updateRes = await updateMicroApp({
        id: props.microAppInfo.id,
        appName: props.microAppInfo.appName,
        appIcon: newIcon,
        appDesc: props.microAppInfo.appDesc,
        remark: props.microAppInfo.remark,
        categoryId: props.microAppInfo.categoryId,
        chargeType: props.microAppInfo.chargeType,
        price: props.microAppInfo.price || 0,
        screenshots: props.microAppInfo.screenshots || '',
      } as any)

      if (updateRes.code !== 0) {
        apiRespErrMsg(updateRes)
        return
      }
    }

    message.success('已设为主信息')
    show.value = false
    emit('done')
  }
  catch (error) {
    apiRespErrMsg(error)
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px" title="版本详情">
    <div v-if="detailInfo" class="space-y-4">
      <!-- 语言切换 -->
      <div v-if="langList.length > 0" class="flex justify-end">
        <NSelect
          v-model:value="currentLang"
          :options="langList.map(lang => ({ label: lang, value: lang }))"
          style="width: 120px"
          size="small"
        />
      </div>

      <div class="flex items-start gap-4">
        <img v-if="detailInfo.iconURL" :src="detailInfo.iconURL" class="w-20 h-20 object-contain border rounded">
        <div v-else class="w-20 h-20 bg-gray-100 border rounded flex items-center justify-center text-gray-400">
          暂无图标
        </div>
        <div class="flex-1">
          <div class="text-lg font-bold">
            {{ currentAppName || '未命名' }}
          </div>
          <div v-if="currentAppDesc" class="text-sm text-gray-500 mt-1">
            {{ currentAppDesc }}
          </div>
          <div class="text-sm text-gray-500 mt-1">
            ID: {{ detailInfo.microAppId }}
          </div>
          <div class="text-sm text-gray-500">
            作者: {{ detailInfo.author }}
          </div>
          <div class="text-sm text-gray-500">
            版本: {{ detailInfo.version }}
          </div>
        </div>
      </div>

      <!-- 多语言列表 -->
      <div v-if="langList.length > 0" class="text-sm">
        <div class="text-gray-500 mb-2">
          多语言信息：
        </div>
        <NDataTable
          :columns="langTableColumns"
          :data="langList.map(lang => ({
            lang,
            appName: detailInfo?.appInfo?.[lang]?.appName || '-',
            appDesc: detailInfo?.appInfo?.[lang]?.appDesc || '-',
          }))"
          :bordered="true"
          size="small"
        />
      </div>

      <div v-if="detailInfo.versionDesc" class="text-sm">
        <div class="text-gray-500">
          版本说明：
        </div>
        <div>{{ detailInfo.versionDesc }}</div>
      </div>

      <div v-if="detailInfo.permissions?.length" class="text-sm">
        <div class="text-gray-500">
          权限：
        </div>
        <div class="flex flex-wrap gap-1 mt-1">
          <NTag v-for="p in detailInfo.permissions" :key="p" size="small" type="info">
            {{ p }}
          </NTag>
        </div>
      </div>

      <div v-if="detailInfo.networkDomains?.length" class="text-sm">
        <div class="text-gray-500">
          网络域名白名单：
        </div>
        <div class="flex flex-wrap gap-1 mt-1">
          <NTag v-for="d in detailInfo.networkDomains" :key="d" size="small">
            {{ d }}
          </NTag>
        </div>
      </div>
    </div>
    <div v-else class="text-center py-8 text-gray-500">
      暂无详情信息
    </div>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="show = false">
          关闭
        </NButton>
        <NButton type="primary" :loading="loading" @click="handleSetAsMainInfo">
          设为微应用主信息
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
