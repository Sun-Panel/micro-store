<script lang="ts" setup>
import type { Theme } from '@/store/modules/app/helper'
import moment from 'moment'
import { NButton, NEllipsis, NImage, NImageGroup, NTooltip, useMessage } from 'naive-ui'
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'
import { getInfo, getInfoByMicroAppId, getVersionList } from '@/api/microApp'
import { SvgIconOnline } from '@/components/common'
import { microAppChargeTypeMap, microAppThirdChargeTypeMap, MicroAppVersionStatus } from '@/enums/panel'
import { isIframe } from '@/utils/cmn'
import { useAppInstallStatus } from '../composables/useAppInstallStatus'
import { getDownloadUrl } from '../composables/useDownload'
import { useIframe } from '../composables/useIframe'
import { useRouterHelper } from '../composables/useRouterHelper'
import 'moment/dist/locale/zh-cn'

// 只显示日期，不显示时间
function dateFormat(timeString?: string) {
  return moment(timeString).format('YYYY-MM-DD')
}

const route = useRoute()
const router = useRouter()
const message = useMessage()
const { sendInstallApp } = useIframe()
const { getHomePath } = useRouterHelper('v1')
const { getAppButtonStatus } = useAppInstallStatus()

// // ==================== 深色模式 ====================
// const themeQuery = computed(() => {
//   // 依赖 route.query 以确保响应式更新
//   void route.query
//   const params = getIframeAllUrlParam()
//   return (params.theme as string) || ''
// })
// const isDarkPage = computed(() => themeQuery.value === 'dark')
// const originalTheme = ref<Theme | null>(null)

// // 根据 iframe URL 参数设置深色模式
// watch(isDarkPage, (dark) => {
//   console.log('dark', dark)
//   if (dark) {
//     // 保存原始主题并切换到深色
//     if (originalTheme.value === null) {
//       originalTheme.value = appStore.theme
//     }
//     appStore.setTheme('dark')
//   }
//   else if (originalTheme.value !== null) {
//     // 恢复原始主题
//     appStore.setTheme(originalTheme.value)
//     originalTheme.value = null
//   }
// }, { immediate: true })

// 路由参数
const routeId = computed(() => route.params.id as string)

// 判断是否为数字 ID
const isNumericId = computed(() => /^\d+$/.test(routeId.value))

// 微应用 ID（数字）
const appRecordId = computed(() => Number(routeId.value))

// 微应用 ID（字符串）
const microAppId = computed(() => routeId.value)

// 数据
const microAppInfo = ref<MicroApp.Info>()
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
  const langList = microAppInfo.value.langList || []
  if (langList.length > 0) {
    return langList.map(l => l.lang)
  }
  return ['zh-CN']
})

// 微应用语言 Map
const baseInfoLangMap = computed(() => {
  const result: Record<string, MicroApp.LangInfo> = {}
  if (!microAppInfo.value)
    return result
  const langList = microAppInfo.value.langList || []
  langList.forEach((l) => {
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

// 截图列表
const screenshotsList = computed(() => {
  if (!microAppInfo.value?.screenshots)
    return []
  return microAppInfo.value.screenshots.split(',').map(s => s.trim()).filter(s => s)
})

// 截图轮播控制
const screenshotScrollRef = ref<HTMLElement | null>(null)
const screenshotScrollStep = 260 // 每次滚动的像素距离
const canScrollLeft = ref(false)
const canScrollRight = ref(false)

function updateScrollButtons() {
  const el = screenshotScrollRef.value
  if (!el)
    return
  canScrollLeft.value = el.scrollLeft > 2
  canScrollRight.value = el.scrollLeft + el.clientWidth < el.scrollWidth - 2
}

function scrollScreenshots(direction: 'left' | 'right') {
  const el = screenshotScrollRef.value
  if (!el)
    return
  const scrollAmount = direction === 'left' ? -screenshotScrollStep : screenshotScrollStep
  el.scrollBy({ left: scrollAmount, behavior: 'smooth' })
}

// 监听截图列表变化，更新按钮状态
watch(screenshotsList, async () => {
  await nextTick()
  updateScrollButtons()
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
  return approvedVersions.sort((a, b) => {
    const timeA = new Date(a.createTime ?? 0).getTime()
    const timeB = new Date(b.createTime ?? 0).getTime()
    return timeB - timeA
  })[0]
})

// 获取应用状态信息
const appStatusInfo = computed(() => {
  return getAppButtonStatus(microAppInfo.value?.microAppId, latestApprovedVersion.value?.version)
})

// 获取版本说明（兼容多语言格式）
function getVersionDescContent(versionDesc: Record<string, { content: string }> | string | undefined): string {
  if (!versionDesc)
    return ''
  // 如果是字符串（旧格式），直接返回
  if (typeof versionDesc === 'string')
    return versionDesc
  // 新格式：优先使用当前语言，否则用 zh-CN，最后用第一个可用语言
  return versionDesc[currentLang.value]?.content
    || versionDesc['zh-CN']?.content
    || Object.values(versionDesc)[0]?.content
    || ''
}

// 获取微应用详情
async function fetchMicroAppInfo() {
  loading.value = true
  try {
    if (isNumericId.value) {
      // 数字 ID，使用原有接口
      const { data } = await getInfo<MicroApp.Info>(appRecordId.value)
      microAppInfo.value = data
    }
    else {
      // 字符串 ID，使用新接口
      const { data } = await getInfoByMicroAppId<MicroApp.Info>(microAppId.value)
      microAppInfo.value = data
    }
  }
  catch (error) {
    console.error('获取微应用详情失败:', error)
    message.error('获取微应用详情失败')
  }
  finally {
    loading.value = false
  }
}

// 获取分类选项（该接口虽在 admin 路径下，但未挂载认证中间件，公开可访问）
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
  if (!microAppInfo.value?.id)
    return

  try {
    const { data } = await getVersionList<Common.ListResponse<MicroApp.VersionInfo[]>>({
      appRecordId: microAppInfo.value.id,
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
  router.push(getHomePath())
}

// // 下载版本
// function handleDownload() {
//   if (latestApprovedVersion.value?.packageUrl) {
//     window.location.href = latestApprovedVersion.value.packageUrl
//   }
// }

// 下载版本包
async function handleDownloadByVersionId(version?: string) {
  if (microAppInfo.value?.microAppId) {
    await getDownloadUrl(microAppInfo.value?.microAppId, version).then((data) => {
      window.open(data, '_blank')
    })
  }
}

// 安装版本
async function handleInstall() {
  if (microAppInfo.value?.microAppId) {
    if (isIframe()) {
      const url = await getDownloadUrl(microAppInfo.value?.microAppId)
      // 在 iframe 中，发送安装消息给父窗口
      sendInstallApp({ microAppId: microAppInfo.value.microAppId, url })
    }
    else {
      message.info('请在私有部署项目中打开微应用商店安装')
    }
  }
}

onMounted(async () => {
  await fetchCategoryOptions()
  await fetchMicroAppInfo()
  await fetchVersionList()
  await nextTick()
  // 监听截图容器的滚动事件，更新按钮状态
  const el = screenshotScrollRef.value
  if (el) {
    el.addEventListener('scroll', updateScrollButtons, { passive: true })
    window.addEventListener('resize', updateScrollButtons)
    updateScrollButtons()
  }
})

onUnmounted(() => {
  const el = screenshotScrollRef.value
  if (el) {
    el.removeEventListener('scroll', updateScrollButtons)
  }
  window.removeEventListener('resize', updateScrollButtons)
  // // 恢复原始主题
  // if (originalTheme.value !== null) {
  //   appStore.setTheme(originalTheme.value)
  //   originalTheme.value = null
  // }
})
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 via-white to-slate-100 dark:from-slate-900 dark:via-slate-800 dark:to-slate-900">
    <!-- 顶部导航栏 -->
    <div class="sticky top-0 z-10 bg-white/80 dark:bg-slate-800/80 backdrop-blur-md border-b border-slate-200/60 dark:border-slate-700/60">
      <div class="mx-auto ">
        <div class="flex items-center justify-between h-14">
          <NButton quaternary class="flex items-center gap-1 text-slate-600 dark:text-slate-300 hover:text-slate-900 dark:hover:text-white" @click="handleBack">
            <template #icon>
              <SvgIconOnline icon="ph:arrow-left" />
            </template>
            <span class="hidden sm:inline">返回</span>
          </NButton>
          <h1 class="text-sm font-medium text-slate-500 dark:text-slate-400 truncate max-w-[200px] sm:max-w-none">
            {{ displayAppName || '应用详情' }}
          </h1>
          <div class="w-10" />
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-24 text-center">
      <div class="inline-flex items-center gap-3 text-slate-400 dark:text-slate-500">
        <SvgIconOnline icon="eos-icons:loading" class="text-xl animate-spin" />
        <span class="text-sm">加载中...</span>
      </div>
    </div>

    <!-- 不存在状态 -->
    <div v-else-if="!microAppInfo" class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-24 text-center">
      <div class="w-16 h-16 mx-auto mb-4 bg-slate-100 dark:bg-slate-700 rounded-2xl flex items-center justify-center">
        <SvgIconOnline icon="ph:magnifying-glass" class="text-2xl text-slate-400 dark:text-slate-500" />
      </div>
      <p class="text-slate-400 dark:text-slate-500">
        微应用不存在
      </p>
    </div>

    <!-- 主内容 -->
    <div v-else class="mx-auto py-6 sm:py-8">
      <!-- 应用头部卡片 -->
      <div class="bg-white dark:bg-slate-800 rounded-3xl shadow-sm border border-slate-200/60 dark:border-slate-700/60 overflow-hidden">
        <div class="p-5 sm:p-8">
          <div class="flex flex-col sm:flex-row gap-5 sm:gap-8">
            <!-- 应用图标 -->
            <div class="flex-shrink-0 flex flex-col items-center sm:items-start gap-3">
              <div class="relative group">
                <img
                  v-if="microAppInfo.appIcon"
                  :src="microAppInfo.appIcon"
                  class="w-28 h-28 sm:w-32 sm:h-32 object-contain rounded-2xl  ring-slate-100 dark:ring-slate-700 group-hover:shadow-xl transition-shadow duration-300"
                >
                <div v-else class="w-28 h-28 sm:w-32 sm:h-32 bg-gradient-to-br from-slate-100 to-slate-50 dark:from-slate-700 dark:to-slate-600 rounded-2xl flex items-center justify-center border border-slate-200/60 dark:border-slate-600/60">
                  <SvgIconOnline icon="ph:app-window" class="text-3xl text-slate-300 dark:text-slate-500" />
                </div>
                <!-- 状态角标 -->
                <div
                  v-if="microAppInfo.status !== 1"
                  class="absolute -top-1 -right-1 px-2 py-0.5 bg-red-500 text-white text-[10px] font-medium rounded-full shadow-sm"
                >
                  已下架
                </div>
              </div>
              <!-- 数据统计 -->
              <div class="flex items-center gap-4 text-xs text-slate-500 dark:text-slate-400">
                <div class="flex items-center gap-1.5" :title="$t('common.downloadCount')">
                  <SvgIconOnline icon="ph:arrow-down" class="text-sm" />
                  <span>{{ microAppInfo.downloadCount || 0 }}</span>
                </div>
                <div class="w-px h-3 bg-slate-200 dark:bg-slate-600" />
                <div class="flex items-center gap-1.5" :title="$t('common.installCount')">
                  <SvgIconOnline icon="ph:package" class="text-sm" />
                  <span>{{ microAppInfo.installCount || 0 }}</span>
                </div>
              </div>
            </div>

            <!-- 应用信息 -->
            <div class="flex-1 min-w-0 text-center sm:text-left">
              <!-- 标题行 -->
              <div class="mb-3">
                <h1 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white leading-tight mb-2">
                  {{ displayAppName }}
                </h1>
                <div class="flex flex-wrap items-center justify-center sm:justify-start gap-2 text-sm text-slate-500 dark:text-slate-400">
                  <span v-if="microAppInfo.developer?.name" class="inline-flex items-center gap-1">
                    <SvgIconOnline icon="ph:user" class="text-sm" />
                    {{ microAppInfo.developer.name }}
                  </span>
                  <span v-if="categoryName" class="inline-flex items-center gap-1">
                    <SvgIconOnline icon="ph:tag" class="text-sm" />
                    {{ categoryName }}
                  </span>
                </div>
              </div>

              <!-- 版本与收费信息 -->
              <div class="flex flex-wrap items-center justify-center sm:justify-start gap-2 mb-4">
                <span v-if="latestApprovedVersion" class="inline-flex items-center px-2.5 py-1 bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 text-xs font-medium rounded-lg">
                  v{{ latestApprovedVersion.version }}
                </span>
                <span class="inline-flex items-center px-2.5 py-1 bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-300 text-xs font-medium rounded-lg">
                  {{ microAppChargeTypeMap[microAppInfo.chargeType] || '免费' }}
                </span>
                <span v-if="microAppInfo.thirdCharge" class="inline-flex items-center px-2.5 py-1 bg-amber-50 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400 text-xs font-medium rounded-lg">
                  {{ microAppThirdChargeTypeMap[microAppInfo.thirdCharge] || '第三方收费' }}
                </span>
                <span v-if="microAppInfo.haveIframe" class="inline-flex items-center px-2.5 py-1 bg-purple-50 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 text-xs font-medium rounded-lg">
                  Iframe
                </span>
              </div>

              <!-- 操作按钮 -->
              <div v-if="latestApprovedVersion" class="flex flex-wrap items-center justify-center sm:justify-start gap-3">
                <NButton
                  v-if="!isIframe()"
                  size="large"
                  class="download-btn"
                  @click="handleDownloadByVersionId()"
                >
                  <template #icon>
                    <SvgIconOnline icon="ph:arrow-down" />
                  </template>
                  下载
                </NButton>
                <NButton
                  size="large"
                  :type="appStatusInfo.type"
                  :disabled="appStatusInfo.disabled"
                  class="install-btn"
                  @click="handleInstall"
                >
                  <template #icon>
                    <SvgIconOnline icon="ph:download-simple" />
                  </template>
                  {{ appStatusInfo.text }}
                </NButton>
                <NButton
                  v-if="appStatusInfo.isInstalled && !appStatusInfo.hasUpdate"
                  size="large"
                  @click="handleInstall"
                >
                  覆盖安装
                </NButton>
              </div>
              <div v-else class="text-sm text-slate-400 dark:text-slate-500 italic">
                暂无可用版本
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 两栏布局：左侧截图+介绍，右侧应用信息 -->
      <div class="mt-6 grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 左侧：截图与介绍 -->
        <div class="lg:col-span-2 space-y-6">
          <!-- 截图预览 -->
          <div v-if="screenshotsList.length > 0">
            <h2 class="text-base font-semibold text-slate-700 dark:text-slate-200 mb-3 flex items-center gap-2">
              <SvgIconOnline icon="ph:images-square" class="text-lg text-slate-400 dark:text-slate-500" />
              应用截图
            </h2>
            <div class="relative rounded-2xl overflow-hidden group/screenshot">
              <NImageGroup>
                <div
                  ref="screenshotScrollRef"
                  class="flex gap-3 overflow-x-auto scrollbar-hide snap-x snap-mandatory"
                >
                  <NImage
                    v-for="(screenshot, index) in screenshotsList"
                    :key="index"
                    :src="screenshot"
                    class="flex-shrink-0 w-64 h-40 sm:w-80 sm:h-52 object-cover rounded-xl shadow-md ring-1 ring-slate-200/60 dark:ring-slate-600/60 snap-center hover:shadow-lg transition-shadow duration-300 cursor-pointer"
                  />
                </div>
              </NImageGroup>
              <!-- 左右滑动按钮 -->
              <template v-if="screenshotsList.length > 1">
                <button
                  v-show="canScrollLeft"
                  class="absolute left-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-black/30 hover:bg-black/50 backdrop-blur-sm shadow-md rounded-full flex items-center justify-center text-white hover:text-white transition-all"
                  @click="scrollScreenshots('left')"
                >
                  <SvgIconOnline icon="ph:caret-left" class="text-lg" />
                </button>
                <button
                  v-show="canScrollRight"
                  class="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-black/30 hover:bg-black/50 backdrop-blur-sm shadow-md rounded-full flex items-center justify-center text-white hover:text-white transition-all"
                  @click="scrollScreenshots('right')"
                >
                  <SvgIconOnline icon="ph:caret-right" class="text-lg" />
                </button>
              </template>
            </div>
          </div>

          <!-- 应用介绍 -->
          <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200/60 dark:border-slate-700/60 p-5 sm:p-6">
            <h2 class="text-base font-semibold text-slate-700 dark:text-slate-200 mb-4 flex items-center gap-2">
              <div class="w-1 h-5 bg-blue-500 rounded-full" />
              应用介绍
            </h2>
            <div v-if="displayAppDesc" class="text-sm text-slate-600 dark:text-slate-300 leading-relaxed whitespace-pre-line">
              {{ displayAppDesc }}
            </div>
            <div v-else class="text-sm text-slate-400 dark:text-slate-500 italic py-4 text-center bg-slate-50 dark:bg-slate-700/50 rounded-xl">
              暂无介绍
            </div>
          </div>

          <!-- 新版本特性 -->
          <div v-if="latestApprovedVersion" class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200/60 dark:border-slate-700/60 p-5 sm:p-6">
            <h2 class="text-base font-semibold text-slate-700 dark:text-slate-200 mb-4 flex items-center gap-2">
              <div class="w-1 h-5 bg-green-500 rounded-full" />
              更新日志
              <span class="text-xs font-normal text-slate-400 dark:text-slate-500 ml-auto">
                v{{ latestApprovedVersion.version }}
              </span>
            </h2>
            <div v-if="getVersionDescContent(latestApprovedVersion.versionDesc)" class="text-sm text-slate-600 dark:text-slate-300 leading-relaxed whitespace-pre-line">
              {{ getVersionDescContent(latestApprovedVersion.versionDesc) }}
            </div>
            <div v-else class="text-sm text-slate-400 dark:text-slate-500 italic py-4 text-center bg-slate-50 dark:bg-slate-700/50 rounded-xl">
              暂无版本说明
            </div>
          </div>
        </div>

        <!-- 右侧：信息面板 -->
        <div class="space-y-6">
          <!-- 基本信息 -->
          <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200/60 dark:border-slate-700/60 p-5 sm:p-6">
            <h2 class="text-base font-semibold text-slate-700 dark:text-slate-200 mb-4 flex items-center gap-2">
              <div class="w-1 h-5 bg-indigo-500 rounded-full" />
              应用信息
            </h2>
            <div class="space-y-3">
              <div class="flex items-center justify-between py-2 border-b border-slate-100 dark:border-slate-700 last:border-0">
                <span class="text-sm text-slate-500 dark:text-slate-400 mr-5">AppID</span>
                <NEllipsis :title="microAppInfo.microAppId" tooltip class="text-sm text-slate-700 dark:text-slate-200 font-medium font-mono max-w-[180px]">
                  {{ microAppInfo.microAppId }}
                </NEllipsis>
              </div>
              <div class="flex items-center justify-between py-2 border-b border-slate-100 dark:border-slate-700 last:border-0">
                <span class="text-sm text-slate-500 dark:text-slate-400">版本</span>
                <span class="text-sm text-slate-700 dark:text-slate-200 font-medium">
                  {{ latestApprovedVersion ? `v${latestApprovedVersion.version}` : '未发布' }}
                </span>
              </div>
              <div class="flex items-center justify-between py-2 border-b border-slate-100 dark:border-slate-700 last:border-0">
                <span class="text-sm text-slate-500 dark:text-slate-400">更新时间</span>
                <span class="text-sm text-slate-700 dark:text-slate-200">
                  {{ latestApprovedVersion ? dateFormat(String(latestApprovedVersion.createTime)) : '-' }}
                </span>
              </div>
              <div class="flex items-center justify-between py-2 border-b border-slate-100 dark:border-slate-700 last:border-0">
                <span class="text-sm text-slate-500 dark:text-slate-400">收费类型</span>
                <span class="text-sm text-slate-700 dark:text-slate-200">{{ microAppChargeTypeMap[microAppInfo.chargeType] || '免费' }}</span>
              </div>
              <div class="flex items-center justify-between py-2 border-b border-slate-100 dark:border-slate-700 last:border-0">
                <span class="text-sm text-slate-500 dark:text-slate-400">第三方收费</span>
                <span class="text-sm text-slate-700 dark:text-slate-200">{{ microAppThirdChargeTypeMap[microAppInfo.thirdCharge || 0] || '不含' }}</span>
              </div>
              <div class="flex items-center justify-between py-2 border-b border-slate-100 dark:border-slate-700 last:border-0">
                <span class="text-sm text-slate-500 dark:text-slate-400">Iframe</span>
                <span class="text-sm text-slate-700 dark:text-slate-200">{{ microAppInfo.haveIframe ? '支持' : '不支持' }}</span>
              </div>
              <!-- 开源信息 -->
              <div v-if="microAppInfo.openSourceUrl" class="flex items-center justify-between py-2">
                <span class="text-sm text-slate-500 dark:text-slate-400">开源</span>
                <a
                  :href="microAppInfo.openSourceUrl"
                  target="_blank"
                  class="inline-flex items-center gap-1 text-sm text-blue-500 dark:text-blue-400 hover:text-blue-600 dark:hover:text-blue-300 hover:underline font-medium"
                >
                  <SvgIconOnline icon="ph:git-branch" class="text-sm" />
                  查看仓库
                  <SvgIconOnline icon="ph:arrow-up-right" class="text-xs" />
                </a>
              </div>
            </div>
          </div>

          <!-- 反馈信息 -->
          <div v-if="microAppInfo.feedbackChannel" class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200/60 dark:border-slate-700/60 p-5 sm:p-6">
            <h2 class="text-base font-semibold text-slate-700 dark:text-slate-200 mb-3 flex items-center gap-2">
              <div class="w-1 h-5 bg-pink-500 rounded-full" />
              反馈渠道
              <NTooltip trigger="hover">
                <template #trigger>
                  <SvgIconOnline icon="ph:info" class="text-sm text-slate-400 dark:text-slate-500 hover:text-slate-600 dark:hover:text-slate-300 cursor-help" />
                </template>
                反馈渠道由开发者自行提供，请自行甄别安全性。如内部产生任何金钱交易，与 Sun-Panel 官方无关。
              </NTooltip>
            </h2>
            <p class="text-sm text-slate-600 dark:text-slate-300 leading-relaxed whitespace-pre-line break-words">
              {{ microAppInfo.feedbackChannel }}
            </p>
          </div>

          <!-- 开发者信息 -->
          <div v-if="microAppInfo.developer" class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200/60 dark:border-slate-700/60 p-5 sm:p-6">
            <h2 class="text-base font-semibold text-slate-700 dark:text-slate-200 mb-4 flex items-center gap-2">
              <div class="w-1 h-5 bg-orange-500 rounded-full" />
              开发者
            </h2>
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-gradient-to-br from-orange-400 to-amber-400 rounded-xl flex items-center justify-center text-white font-medium text-sm">
                {{ (microAppInfo.developer.name || 'U')[0] }}
              </div>
              <div>
                <div class="text-sm font-medium text-slate-700 dark:text-slate-200">
                  {{ microAppInfo.developer.name || '未知开发者' }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 版本历史 -->
      <div v-if="versionList.length > 0" class="mt-6">
        <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200/60 dark:border-slate-700/60 p-5 sm:p-6">
          <h2 class="text-base font-semibold text-slate-700 dark:text-slate-200 mb-4 flex items-center gap-2">
            <div class="w-1 h-5 bg-violet-500 rounded-full" />
            版本历史
          </h2>
          <div class="space-y-4">
            <div
              v-for="(version, index) in versionList.filter(v => v.status === MicroAppVersionStatus.APPROVED)"
              :key="version.id"
              class="relative pl-6 pb-4 last:pb-0"
              :class="{ 'border-l-2 border-slate-100 dark:border-slate-700': index < versionList.filter(v => v.status === MicroAppVersionStatus.APPROVED).length - 1 }"
            >
              <!-- 时间线圆点 -->
              <div class="absolute left-0 top-0 w-3 h-3 bg-white dark:bg-slate-800 border-2 border-violet-400 rounded-full -translate-x-[7px] translate-y-1" />
              <div class="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-3 mb-2">
                <span class="text-sm font-semibold text-slate-800 dark:text-slate-200">v{{ version.version }}</span>
                <span class="text-xs text-slate-400 dark:text-slate-500">{{ dateFormat(String(version.createTime)) }}</span>
              </div>
              <div v-if="getVersionDescContent(version.versionDesc)" class="text-sm text-slate-600 dark:text-slate-300 leading-relaxed whitespace-pre-line">
                {{ getVersionDescContent(version.versionDesc) }}
              </div>
              <div v-else class="text-sm text-slate-400 dark:text-slate-500 italic">
                暂无版本说明
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.install-btn :deep(.n-button__icon) {
  margin-right: 4px;
}
</style>
