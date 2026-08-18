<script lang="ts" setup>
import { NButton, NDescriptions, NDescriptionsItem, NDivider, NImage, NImageGroup, NInput, NModal, NSpace, NSwitch, NTag, useMessage } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { approve, getMicroAppInfo } from '@/api/admin/microAppReview'
import { microAppChargeTypeMap, microAppThirdChargeTypeMap } from '@/enums/panel'

const props = defineProps<{
  visible: boolean
  reviewInfo?: MicroApp.MicroAppReviewInfo
  microAppModelId: number
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'done'): void
}>()

const message = useMessage()

// 数据
const reviewLoading = ref(false)
const currentAppInfo = ref<MicroApp.Info>()
const reviewForm = ref({
  status: 1,
  reviewNote: '',
})
const iframeModalVisible = ref(false)

// ==================== 同步滚动 ====================
const syncScrollEnabled = ref(true)
const onlyShowDiff = ref(false)
const leftPanelRef = ref<HTMLDivElement>()
const rightPanelRef = ref<HTMLDivElement>()
let isSyncing = false

function handleLeftScroll() {
  if (!syncScrollEnabled.value || isSyncing) return
  isSyncing = true
  if (leftPanelRef.value && rightPanelRef.value) {
    rightPanelRef.value.scrollTop = leftPanelRef.value.scrollTop
  }
  requestAnimationFrame(() => {
    isSyncing = false
  })
}

function handleRightScroll() {
  if (!syncScrollEnabled.value || isSyncing) return
  isSyncing = true
  if (leftPanelRef.value && rightPanelRef.value) {
    leftPanelRef.value.scrollTop = rightPanelRef.value.scrollTop
  }
  requestAnimationFrame(() => {
    isSyncing = false
  })
}

// ==================== 多语言处理工具函数 ====================
interface LangItem {
  lang: string
  appName?: string
  appDesc?: string
}

// 解析 langMap（可能为字符串、数组或对象）为 Record<string, LangItem>
function parseLangMap(langMap: MicroApp.MicroAppReviewInfo['langMap']): Record<string, LangItem> {
  if (!langMap)
    return {}
  let parsed: any = langMap
  if (typeof langMap === 'string') {
    try {
      parsed = JSON.parse(langMap)
    }
    catch {
      return {}
    }
  }
  if (Array.isArray(parsed)) {
    const result: Record<string, LangItem> = {}
    parsed.forEach((l: any) => {
      if (l.lang)
        result[l.lang] = l
    })
    return result
  }
  return parsed
}

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

// ==================== 计算属性 ====================
// 待审核语言 Map（解析一次，其他计算属性复用）
const reviewLangMap = computed(() => parseLangMap(props.reviewInfo?.langMap))

// 待审核语言列表（从已解析的 map 提取）
const reviewLangList = computed(() => {
  const keys = Object.keys(reviewLangMap.value)
  return keys.length > 0 ? keys : ['zh-CN']
})

// 当前语言 Map（从 langList 数组转换）
const currentLangMap = computed(() => {
  if (!currentAppInfo.value)
    return {}
  const langList = currentAppInfo.value.langList || []
  const map: Record<string, LangItem> = {}
  langList.forEach((l) => {
    if (l.lang)
      map[l.lang] = l
  })
  return map
})

// 当前语言列表（从已解析的 map 提取）
const currentLangList = computed(() => {
  const keys = Object.keys(currentLangMap.value)
  return keys.length > 0 ? keys : ['zh-CN']
})

// 当前语言（优先匹配浏览器语言）
const currentLang = computed(() => {
  const browserLang = getBrowserLang()
  const langs = reviewLangList.value
  return langs.includes(browserLang) ? browserLang : (langs.includes('zh-CN') ? 'zh-CN' : langs[0])
})

// 显示名称/描述（带语言回退）
function getLocalizedField(
  langMap: Record<string, LangItem>,
  fallbackInfo: { appName?: string; appDesc?: string },
  field: 'appName' | 'appDesc',
): string {
  const lang = currentLang.value
  return langMap[lang]?.[field]
    || langMap['zh-CN']?.[field]
    || fallbackInfo[field]
    || ''
}

const displayReviewAppName = computed(() => {
  if (!props.reviewInfo)
    return ''
  return getLocalizedField(reviewLangMap.value, props.reviewInfo, 'appName')
})

const displayReviewAppDesc = computed(() => {
  if (!props.reviewInfo)
    return ''
  return getLocalizedField(reviewLangMap.value, props.reviewInfo, 'appDesc')
})

const displayCurrentAppName = computed(() => {
  if (!currentAppInfo.value)
    return ''
  return getLocalizedField(currentLangMap.value, currentAppInfo.value, 'appName')
})

const displayCurrentAppDesc = computed(() => {
  if (!currentAppInfo.value)
    return ''
  return getLocalizedField(currentLangMap.value, currentAppInfo.value, 'appDesc')
})

// 截图列表（预处理为数组，避免模板中重复计算）
const currentScreenshots = computed(() => {
  const screenshots = currentAppInfo.value?.screenshots || ''
  return screenshots.split(',').filter(s => s.trim())
})

const reviewScreenshots = computed(() => {
  const screenshots = props.reviewInfo?.screenshots || ''
  return screenshots.split(',').filter(s => s.trim())
})

// ==================== 字段差异对比 ====================
const fieldDiff = computed(() => {
  const curr = currentAppInfo.value
  const rev = props.reviewInfo
  if (!curr || !rev) return {} as Record<string, boolean>

  // 复用已计算的语言映射
  const currLang = currentLangMap.value[currentLang.value] || currentLangMap.value['zh-CN'] || {}
  const revLang = reviewLangMap.value[currentLang.value] || reviewLangMap.value['zh-CN'] || {}

  return {
    appName: (currLang.appName || curr.appName || '') !== (revLang.appName || rev.appName || ''),
    appIcon: (curr.appIcon || '') !== (rev.appIcon || ''),
    appDesc: (currLang.appDesc || curr.appDesc || '') !== (revLang.appDesc || rev.appDesc || ''),
    chargeType: curr.chargeType !== rev.chargeType,
    points: curr.points !== rev.points,
    thirdCharge: (curr.thirdCharge || 0) !== (rev.thirdCharge || 0),
    haveIframe: curr.haveIframe !== rev.haveIframe,
    openSourceUrl: (curr.openSourceUrl || '') !== (rev.openSourceUrl || ''),
    feedbackChannel: (curr.feedbackChannel || '') !== (rev.feedbackChannel || ''),
    remark: (curr.remark || '') !== (rev.remark || ''),
    screenshots: (curr.screenshots || '') !== (rev.screenshots || ''),
    langList: JSON.stringify(currentLangList.value.slice().sort())
      !== JSON.stringify(reviewLangList.value.slice().sort()),
  }
})

function isDiff(field: string): boolean {
  return !!fieldDiff.value[field]
}

// 逐语言差异对比
const langDiffMap = computed<Record<string, Record<string, boolean>>>(() => {
  const currLang = currentLangMap.value
  const revLang = reviewLangMap.value
  if (!currLang && !revLang) return {}

  const allLangs = new Set([...Object.keys(currLang), ...Object.keys(revLang)])
  const result: Record<string, Record<string, boolean>> = {}

  allLangs.forEach((lang) => {
    const c = currLang[lang]
    const r = revLang[lang]
    result[lang] = {
      appName: (c?.appName || '') !== (r?.appName || ''),
      appDesc: (c?.appDesc || '') !== (r?.appDesc || ''),
      missingInCurrent: !c && !!r,
      missingInReview: !!c && !r,
    }
  })

  return result
})

function isLangDiff(lang: string, field: string): boolean {
  return !!langDiffMap.value[lang]?.[field]
}

// 监听弹窗打开，获取应用信息
watch(() => props.visible, async (visible) => {
  if (visible && props.reviewInfo) {
    reviewForm.value = {
      status: 1,
      reviewNote: '',
    }

    // 获取当前应用信息
    try {
      const { data } = await getMicroAppInfo<MicroApp.Info>(props.reviewInfo.appRecordId)
      currentAppInfo.value = data
    }
    catch {
      message.error('获取应用信息失败')
    }
  }
})

// 打开应用公开页面
function handleOpenMicroAppPublic() {
  iframeModalVisible.value = true
}

// 提交审核
async function handleReview() {
  if (!props.reviewInfo)
    return

  // 驳回时必须填写原因
  if (reviewForm.value.status === 2 && !reviewForm.value.reviewNote?.trim()) {
    message.error('驳回时必须填写驳回原因')
    return
  }

  reviewLoading.value = true
  try {
    const { code } = await approve<any>({
      reviewId: props.reviewInfo.id,
      status: reviewForm.value.status,
      reviewNote: reviewForm.value.reviewNote,
    })

    if (code === 0) {
      message.success(reviewForm.value.status === 1 ? '审核通过' : '已拒绝')
      emit('update:visible', false)
      emit('done')
    }
  }
  catch {
    message.error('操作失败')
  }
  finally {
    reviewLoading.value = false
  }
}
</script>

<template>
  <NModal :show="visible" preset="card" style="width: 1200px;" @update:show="emit('update:visible', $event)">
    <template #header>
      <div class="flex gap-2 items-center">
        <div>
          审核微应用
        </div>
        <div>
          <NButton size="small" @click="handleOpenMicroAppPublic">
            查看应用公开页面
          </NButton>
          <span class="text-gray-400 text-base ml-2">作者：{{ currentAppInfo?.developer?.name }} ({{ currentAppInfo?.developer?.developerName }})</span>
        </div>
      </div>
    </template>
    <div v-if="reviewInfo" class="space-y-6">
      <!-- 工具栏 -->
      <div class="flex items-center gap-6 px-2 py-1 bg-gray-50 rounded">
        <div class="flex items-center gap-2">
          <NSwitch v-model:value="syncScrollEnabled" size="small" />
          <span class="text-sm text-gray-600">同步滚动</span>
        </div>
        <div class="flex items-center gap-2">
          <NSwitch v-model:value="onlyShowDiff" size="small" />
          <span class="text-sm text-gray-600">仅显示有变动的项</span>
        </div>
      </div>
      <!-- 对比展示 -->
      <div class="flex gap-6">
        <!-- 原始信息 -->
        <div v-if="currentAppInfo" class="flex-1 min-w-0">
          <div class="text-lg font-semibold mb-4 pb-2 border-b">
            当前发布信息
          </div>
          <div ref="leftPanelRef" class="scroll-panel" @scroll="handleLeftScroll">
            <NDescriptions bordered :column="1">
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('appName')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('appName') }">应用名称</span>
                </template>
                {{ displayCurrentAppName }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('appIcon')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('appIcon') }">应用图标</span>
                </template>
                <img
                  v-if="currentAppInfo.appIcon"
                  :src="currentAppInfo.appIcon"
                  class="w-16 h-16 object-contain rounded"
                >
                <span v-else class="text-gray-400">暂无图标</span>
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('appDesc')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('appDesc') }">应用描述</span>
                </template>
                {{ displayCurrentAppDesc || '-' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('chargeType')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('chargeType') }">收费方式</span>
                </template>
                {{ microAppChargeTypeMap[currentAppInfo.chargeType] || '免费' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="(!onlyShowDiff || isDiff('points')) && currentAppInfo.chargeType === 1">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('points') }">价格</span>
                </template>
                {{ currentAppInfo.points }} 积分
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('thirdCharge')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('thirdCharge') }">{{ $t('microApp.thirdCharge') }}</span>
                </template>
                {{ microAppThirdChargeTypeMap[currentAppInfo.thirdCharge || 0] || '不含' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('haveIframe')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('haveIframe') }">包含iframe</span>
                </template>
                {{ currentAppInfo.haveIframe ? '是' : '否' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('openSourceUrl')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('openSourceUrl') }">开源地址</span>
                </template>
                <a
                  v-if="currentAppInfo.openSourceUrl"
                  :href="currentAppInfo.openSourceUrl"
                  target="_blank"
                  class="text-blue-600 hover:text-blue-800"
                >
                  {{ currentAppInfo.openSourceUrl }}
                </a>
                <span v-else class="text-gray-400">未提供</span>
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('feedbackChannel')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('feedbackChannel') }">反馈信息</span>
                </template>
                {{ currentAppInfo.feedbackChannel || '未提供' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('remark')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('remark') }">备注</span>
                </template>
                {{ currentAppInfo.remark || '暂无备注' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('langList')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('langList') }">支持语言</span>
                </template>
                <NTag v-for="lang in currentLangList" :key="lang" size="small" class="mr-1">
                  {{ lang }}
                </NTag>
              </NDescriptionsItem>
            </NDescriptions>

            <!-- 当前多语言详情 -->
            <NDivider v-if="!onlyShowDiff || isDiff('langList')" title-placement="left">
              <span :class="{ 'diff-label': isDiff('langList') }">多语言详情</span>
            </NDivider>
            <div v-if="currentAppInfo?.langList && currentAppInfo.langList.length > 0">
              <div
                v-for="(langItem, index) in currentAppInfo.langList"
                v-show="!onlyShowDiff || isLangDiff(langItem.lang, 'appName') || isLangDiff(langItem.lang, 'appDesc') || isLangDiff(langItem.lang, 'missingInReview')"
                :key="`current-lang-${index}`"
                class="mb-3 p-3 bg-gray-50 rounded"
              >
                <div class="font-semibold text-sm mb-1">
                  {{ langItem.lang }}
                  <span v-if="isLangDiff(langItem.lang, 'missingInReview')" class="diff-tag added">审核中无此语言</span>
                </div>
                <div class="text-sm">
                  <div class="mb-1">
                    <span :class="{ 'diff-label': isLangDiff(langItem.lang, 'appName') }">名称:</span> {{ langItem.appName }}
                  </div>
                  <div>
                    <span :class="{ 'diff-label': isLangDiff(langItem.lang, 'appDesc') }">描述:</span> {{ langItem.appDesc || '-' }}
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="text-gray-400 text-sm">
              暂无多语言信息
            </div>

            <!-- 当前截图 -->
            <NDivider v-if="!onlyShowDiff || isDiff('screenshots')" title-placement="left">
              <span :class="{ 'diff-label': isDiff('screenshots') }">应用截图</span>
            </NDivider>
            <div v-if="currentScreenshots.length > 0" class="grid grid-cols-2 gap-2">
              <NImageGroup>
                <NImage
                  v-for="(screenshot, index) in currentScreenshots"
                  :key="`screenshot-${index}`"
                  :src="screenshot"
                />
              </NImageGroup>
            </div>
            <div v-else-if="!onlyShowDiff || isDiff('screenshots')" class="text-gray-400 text-sm">
              暂无截图
            </div>
          </div>
        </div>

        <!-- 待审核信息 -->
        <div class="flex-1 bg-blue-50 -mx-4 -mt-4 p-4 border-2 border-blue-200 rounded min-w-0">
          <div class="text-lg font-semibold mb-4 pb-2 border-b text-blue-600">
            待审核信息
          </div>
          <div ref="rightPanelRef" class="scroll-panel" @scroll="handleRightScroll">
            <NDescriptions bordered :column="1">
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('appName')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('appName') }">应用名称</span>
                </template>
                {{ displayReviewAppName }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('appIcon')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('appIcon') }">应用图标</span>
                </template>
                <img
                  v-if="reviewInfo.appIcon"
                  :src="reviewInfo.appIcon"
                  class="w-16 h-16 object-contain rounded"
                >
                <span v-else class="text-gray-400">暂无图标</span>
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('appDesc')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('appDesc') }">应用描述</span>
                </template>
                {{ displayReviewAppDesc || '-' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('chargeType')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('chargeType') }">收费方式</span>
                </template>
                {{ microAppChargeTypeMap[reviewInfo.chargeType] || '免费' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="(!onlyShowDiff || isDiff('points')) && reviewInfo.chargeType === 1">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('points') }">价格</span>
                </template>
                {{ reviewInfo.points }} 积分
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('thirdCharge')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('thirdCharge') }">{{ $t('microApp.thirdCharge') }}</span>
                </template>
                {{ microAppThirdChargeTypeMap[reviewInfo.thirdCharge || 0] || '不含' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('haveIframe')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('haveIframe') }">包含iframe</span>
                </template>
                {{ reviewInfo.haveIframe ? '是' : '否' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('openSourceUrl')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('openSourceUrl') }">开源地址</span>
                </template>
                <a
                  v-if="reviewInfo.openSourceUrl"
                  :href="reviewInfo.openSourceUrl"
                  target="_blank"
                  class="text-blue-600 hover:text-blue-800"
                >
                  {{ reviewInfo.openSourceUrl }}
                </a>
                <span v-else class="text-gray-400">未提供</span>
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('feedbackChannel')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('feedbackChannel') }">反馈信息</span>
                </template>
                {{ reviewInfo.feedbackChannel || '未提供' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('remark')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('remark') }">备注</span>
                </template>
                {{ reviewInfo.remark || '暂无备注' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('langList')">
                <template #label>
                  <span :class="{ 'diff-label': isDiff('langList') }">支持语言</span>
                </template>
                <NTag v-for="lang in reviewLangList" :key="lang" type="primary" size="small" class="mr-1">
                  {{ lang }}
                </NTag>
              </NDescriptionsItem>
            </NDescriptions>

            <!-- 待审核多语言详情 -->
            <NDivider v-if="!onlyShowDiff || isDiff('langList')" title-placement="left">
              <span :class="{ 'diff-label': isDiff('langList') }">多语言详情</span>
            </NDivider>
            <div v-if="reviewLangList.length > 0">
              <div
                v-for="(langItem, lang) in reviewLangMap"
                v-show="!onlyShowDiff || isLangDiff(lang as string, 'appName') || isLangDiff(lang as string, 'appDesc') || isLangDiff(lang as string, 'missingInCurrent')"
                :key="`review-lang-${lang}`"
                class="mb-3 p-3 bg-blue-50 rounded border border-blue-200"
              >
                <div class="font-semibold text-sm mb-1 text-blue-700">
                  {{ lang }}
                  <span v-if="isLangDiff(lang as string, 'missingInCurrent')" class="diff-tag added">新语言</span>
                </div>
                <div class="text-sm">
                  <div class="mb-1">
                    <span :class="{ 'diff-label': isLangDiff(lang as string, 'appName') }">名称:</span> {{ langItem.appName }}
                  </div>
                  <div>
                    <span :class="{ 'diff-label': isLangDiff(lang as string, 'appDesc') }">描述:</span> {{ langItem.appDesc || '-' }}
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="text-gray-400 text-sm">
              暂无多语言信息
            </div>

            <!-- 待审核截图 -->
            <NDivider v-if="!onlyShowDiff || isDiff('screenshots')" title-placement="left">
              <span :class="{ 'diff-label': isDiff('screenshots') }">应用截图</span>
            </NDivider>
            <div v-if="reviewScreenshots.length > 0" class="grid grid-cols-2 gap-2">
              <NImageGroup>
                <NImage
                  v-for="(screenshot, index) in reviewScreenshots"
                  :key="`review-${index}`"
                  :src="screenshot"
                />
              </NImageGroup>
            </div>
            <div v-else-if="!onlyShowDiff || isDiff('screenshots')" class="text-gray-400 text-sm">
              暂无截图
            </div>
          </div>
        </div>
      </div>

      <!-- 审核表单 -->
      <NDivider title-placement="left">
        审核操作
      </NDivider>
      <div class="space-y-4">
        <div>
          <div class="mb-2 font-semibold">
            审核决定
          </div>
          <NSpace>
            <NButton
              :type="reviewForm.status === 1 ? 'success' : 'default'"
              @click="reviewForm.status = 1"
            >
              通过
            </NButton>
            <NButton
              :type="reviewForm.status === 2 ? 'error' : 'default'"
              @click="reviewForm.status = 2"
            >
              拒绝
            </NButton>
          </NSpace>
        </div>
        <div v-if="reviewForm.status === 2">
          <div class="mb-2">
            <span class="text-red-500">*</span> 驳回原因
          </div>
          <NInput
            v-model:value="reviewForm.reviewNote"
            type="textarea"
            placeholder="请输入驳回原因（必填）"
            :rows="4"
          />
        </div>
        <div v-else>
          <div class="mb-2">
            审核备注（选填）
          </div>
          <NInput
            v-model:value="reviewForm.reviewNote"
            type="textarea"
            placeholder="请输入审核备注"
            :rows="3"
          />
        </div>
      </div>
    </div>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="emit('update:visible', false)">
          取消
        </NButton>
        <NButton
          :type="reviewForm.status === 1 ? 'success' : 'error'"
          :loading="reviewLoading"
          @click="handleReview"
        >
          {{ reviewForm.status === 1 ? '通过审核' : '拒绝申请' }}
        </NButton>
      </NSpace>
    </template>
  </NModal>

 <!-- 应用公开页面 Modal -->
  <NModal
    :show="iframeModalVisible"
    preset="card"
    style="width: 1200px; height: 800px;"
    title="应用公开页面"
    @update:show="iframeModalVisible = $event"
  >
    <div class="h-full">
      <iframe
        v-if="reviewInfo?.microAppId"
        :src="`/microApp/${reviewInfo.appRecordId}`"
        frameborder="0"
        style="width: 100%; height: 700px; border: none;"
      />
    </div>
  </NModal>
</template>

<style scoped>
.diff-label {
  color: #e53e3e;
  font-weight: 700;
  position: relative;
}
.diff-label::after {
  content: '✦';
  margin-left: 4px;
  font-size: 0.7em;
}
.diff-tag {
  display: inline-block;
  margin-left: 6px;
  padding: 0 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  line-height: 20px;
  vertical-align: middle;
}
.diff-tag.added {
  background-color: #e53e3e;
  color: #fff;
}
.scroll-panel {
  max-height: calc(100vh - 320px);
  overflow-y: auto;
  scrollbar-width: thin;
}
.scroll-panel::-webkit-scrollbar {
  width: 6px;
}
.scroll-panel::-webkit-scrollbar-thumb {
  background-color: #c1c1c1;
  border-radius: 3px;
}
.scroll-panel::-webkit-scrollbar-track {
  background-color: transparent;
}
</style>
