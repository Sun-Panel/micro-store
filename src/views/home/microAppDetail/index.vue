<script lang="ts" setup>
import { NButton, NCard, NImage, NImageGroup, NTag, useMessage } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'
import { getInfo, getVersionList } from '@/api/microApp'
import { microAppChargeTypeMap, MicroAppVersionStatus } from '@/enums/panel'
import { timeFormat } from '@/utils/cmn'

const route = useRoute()
const router = useRouter()
const message = useMessage()

// 微应用ID
const microAppId = computed(() => Number(route.params.id))

// 数据
const microAppInfo = ref<MicroApp.MicroAppInfo>()
const versionList = ref<MicroApp.VersionInfo[]>([])
const categoryOptions = ref<{ label: string, value: number }[]>([])
const loading = ref(false)

// ==================== 多语言处理 ====================
// 浏览器语言检测
function getBrowserLang(): string {
  const lang = navigator.language || (navigator as any).userLanguage || 'zh-CN'
  if (lang.startsWith('zh'))
    return 'zh-CN'
  if (lang.startsWith('en'))
    return 'en-US'
  if (lang.startsWith('ja'))
    return 'ja-JP'
  if (lang.startsWith('ko'))
    return 'ko-KR'
  return 'zh-CN'
}

// 微应用的多语言列表
const baseInfoLangList = computed(() => {
  if (!microAppInfo.value)
    return ['zh-CN']
  const langList = (microAppInfo.value as any).langList || []
  if (langList.length > 0) {
    return langList.map((l: any) => l.lang)
  }
  return ['zh-CN']
})

// 微应用语言 Map
const baseInfoLangMap = computed(() => {
  const result: Record<string, any> = {}
  if (!microAppInfo.value)
    return result
  const langList = (microAppInfo.value as any).langList || []
  langList.forEach((l: any) => {
    result[l.lang] = l
  })
  return result
})

// 当前语言
const currentLang = computed(() => {
  const browserLang = getBrowserLang()
  const langs = baseInfoLangList.value
  return langs.includes(browserLang) ? browserLang : (langs.includes('zh-CN') ? 'zh-CN' : langs[0])
})

// 当前语言下的应用名称
const displayAppName = computed(() => {
  if (!microAppInfo.value)
    return ''
  const langMap = baseInfoLangMap.value
  return langMap[currentLang.value]?.appName
    || langMap['zh-CN']?.appName
    || microAppInfo.value.appName
    || ''
})

// 当前语言下的应用描述
const displayAppDesc = computed(() => {
  if (!microAppInfo.value)
    return ''
  const langMap = baseInfoLangMap.value
  return langMap[currentLang.value]?.appDesc
    || langMap['zh-CN']?.appDesc
    || microAppInfo.value.appDesc
    || ''
})

// 分类名称
const categoryName = computed(() => {
  if (!microAppInfo.value)
    return ''
  const category = categoryOptions.value.find(c => c.value === microAppInfo.value?.categoryId)
  return category?.label || `ID: ${microAppInfo.value.categoryId}`
})

// ==================== 版本处理 ====================
// 最新审核通过的版本
const latestApprovedVersion = computed(() => {
  const approvedVersions = versionList.value.filter(v => v.status === MicroAppVersionStatus.APPROVED)
  if (approvedVersions.length === 0)
    return null
  // 按创建时间排序，最新的在前面
  return approvedVersions.sort((a, b) => new Date(b.createTime).getTime() - new Date(a.createTime).getTime())[0]
})

// 获取微应用详情
async function fetchMicroAppInfo() {
  loading.value = true
  try {
    const { data } = await getInfo<any>(microAppId.value)
    microAppInfo.value = data
  }
  catch (error) {
    console.error('获取微应用详情失败:', error)
    message.error('获取微应用详情失败')
  }
  finally {
    loading.value = false
  }
}

// 获取分类选项
async function fetchCategoryOptions() {
  try {
    const res = await getCategoryList<any>()
    categoryOptions.value = res.data?.map((item: any) => ({
      label: item.name,
      value: item.id,
    })) || []
  }
  catch (error) {
    console.error(error)
  }
}

// 获取版本列表
async function fetchVersionList() {
  try {
    const { data } = await getVersionList<Common.ListResponse<MicroApp.VersionInfo[]>>({
      appId: microAppId.value,
      page: 1,
      limit: 100,
    })
    versionList.value = data.list || []
  }
  catch (error) {
    console.error('获取版本列表失败:', error)
  }
}

// 返回首页
function handleBack() {
  router.push('/')
}

// 下载版本
function handleDownload() {
  if (latestApprovedVersion.value?.packageUrl) {
    window.location.href = latestApprovedVersion.value.packageUrl
  }
}

// 安装版本
function handleInstall() {
  message.info('安装功能开发中，敬请期待')
}

onMounted(async () => {
  await fetchCategoryOptions()
  await Promise.all([
    fetchMicroAppInfo(),
    fetchVersionList(),
  ])
})
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 返回按钮 -->
    <div class="max-w-[1200px] mx-auto px-5 pt-5">
      <NButton quaternary @click="handleBack">
        ← 返回首页
      </NButton>
    </div>

    <div v-if="loading" class="max-w-[1200px] mx-auto px-5 py-12 text-center text-gray-400">
      加载中...
    </div>

    <div v-else-if="!microAppInfo" class="max-w-[1200px] mx-auto px-5 py-12 text-center text-gray-400">
      微应用不存在
    </div>

    <div v-else class="max-w-[1200px] mx-auto px-5 pb-12">
      <!-- 应用主信息卡片 -->
      <NCard class="mb-6" :bordered="false" shadow="hover">
        <div class="flex flex-col md:flex-row gap-6">
          <!-- 应用图标 -->
          <div class="flex-shrink-0">
            <img
              v-if="microAppInfo.appIcon"
              :src="microAppInfo.appIcon"
              class="w-24 h-24 object-contain rounded-lg shadow-sm"
            >
            <div v-else class="w-24 h-24 bg-gray-100 rounded-lg flex items-center justify-center text-gray-400">
              暂无图标
            </div>
          </div>

          <!-- 应用信息 -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-3 mb-2">
              <h1 class="text-2xl font-bold text-gray-800">
                {{ displayAppName }}
              </h1>
              <NTag v-if="microAppInfo.status !== 1" size="small">
                已下架
              </NTag>
            </div>

            <div class="flex flex-wrap gap-x-6 gap-y-2 text-sm text-gray-500 mb-3">
              <span>AppID: {{ microAppInfo.microAppId }}</span>
              <span>作者: {{ microAppInfo.authorName || '未知' }}</span>
              <span>分类: {{ categoryName }}</span>
              <span>收费: {{ microAppChargeTypeMap[microAppInfo.chargeType] || '免费' }}</span>
              <!-- <span>创建时间: {{ timeFormat(String(microAppInfo.createTime)) }}</span> -->
            </div>

            <p class="text-gray-600">
              权限:
            </p>
            <p class="text-gray-600 leading-relaxed">
              -
            </p>

            <p class="text-gray-600">
              支持的语言:
            </p>
            <p class="text-gray-600 leading-relaxed">
              -
            </p>

            <p class="text-gray-600">
              介绍:
            </p>
            <p v-if="displayAppDesc" class="text-gray-600 leading-relaxed">
              {{ displayAppDesc }}
            </p>
            <p v-else class="text-gray-400">
              -
            </p>
          </div>
        </div>
      </NCard>

      <!-- 截图展示 -->
      <NCard v-if="microAppInfo.screenshots" title="应用截图" :bordered="false" shadow="hover">
        <NImageGroup>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <NImage
              v-for="(screenshot, index) in microAppInfo.screenshots.split(',')"
              :key="index"
              :src="screenshot"
              class="rounded-lg"
            />
          </div>
        </NImageGroup>
      </NCard>

      <!-- 版本信息卡片 -->
      <NCard v-if="latestApprovedVersion" class="mb-6" title="最新版本" :bordered="false" shadow="hover">
        <div class="flex flex-col md:flex-row md:items-center gap-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
              <span class="text-blue-600 font-bold">V</span>
            </div>
            <div>
              <div class="text-lg font-semibold text-gray-800">
                {{ latestApprovedVersion.version }}
              </div>
              <div v-if="latestApprovedVersion.versionDesc" class="text-sm text-gray-500 mt-1">
                {{ latestApprovedVersion.versionDesc }}
              </div>
            </div>
          </div>
          <div class="md:ml-auto flex items-center gap-3">
            <div class="text-sm text-gray-400 mr-2">
              发布时间: {{ timeFormat(String(latestApprovedVersion.createTime)) }}
            </div>
            <NButton type="primary" @click="handleDownload">
              下载
            </NButton>
            <NButton @click="handleInstall">
              安装
            </NButton>
          </div>
        </div>
      </NCard>

      <!-- 无版本提示 -->
      <NCard v-else class="mb-6" :bordered="false" shadow="hover">
        <div class="text-center py-8 text-gray-400">
          暂无审核通过的版本
        </div>
      </NCard>
    </div>
  </div>
</template>
