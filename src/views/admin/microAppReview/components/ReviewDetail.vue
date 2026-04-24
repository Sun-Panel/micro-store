<script lang="ts" setup>
import { NButton, NDescriptions, NDescriptionsItem, NDivider, NImage, NImageGroup, NInput, NModal, NSpace, NTag, useMessage } from 'naive-ui'
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

// 获取待审核信息的多语言列表
const reviewLangList = computed(() => {
  if (!props.reviewInfo?.langMap)
    return ['zh-CN']
  const langMap = props.reviewInfo.langMap
  let parsed: any = langMap

  // 如果是字符串，解析 JSON
  if (typeof langMap === 'string') {
    try {
      parsed = JSON.parse(langMap)
    }
    catch (e) {
      console.error('解析 langMap JSON 失败:', e)
      return ['zh-CN']
    }
  }

  // 处理可能是数组的情况
  if (Array.isArray(parsed))
    return parsed.map((l: any) => l.lang).filter(Boolean)

  return Object.keys(parsed)
})

// 获取当前信息的多语言列表
const currentLangList = computed(() => {
  if (!currentAppInfo.value)
    return ['zh-CN']
  const langList = (currentAppInfo.value as any).langList || []
  if (langList.length > 0)
    return langList.map((l: any) => l.lang)

  return ['zh-CN']
})

// 当前语言
const currentLang = computed(() => {
  const browserLang = getBrowserLang()
  const langs = reviewLangList.value
  return langs.includes(browserLang) ? browserLang : (langs.includes('zh-CN') ? 'zh-CN' : langs[0])
})

// 获取待审核信息的多语言 Map（处理可能是数组或字符串的情况）
const reviewLangMap = computed(() => {
  if (!props.reviewInfo?.langMap)
    return {}
  const langMap = props.reviewInfo.langMap
  let parsed: any = langMap

  // 如果是字符串，解析 JSON
  if (typeof langMap === 'string') {
    try {
      parsed = JSON.parse(langMap)
    }
    catch (e) {
      console.error('解析 langMap JSON 失败:', e)
      return {}
    }
  }

  // 如果是数组，转换为对象
  if (Array.isArray(parsed)) {
    const result: Record<string, any> = {}
    parsed.forEach((l: any) => {
      if (l.lang)
        result[l.lang] = l
    })
    return result
  }

  return parsed
})

// 当前语言下的应用名称（待审核）
const displayReviewAppName = computed(() => {
  if (!props.reviewInfo)
    return ''
  const langMap = reviewLangMap.value
  return langMap[currentLang.value]?.appName
    || langMap['zh-CN']?.appName
    || props.reviewInfo.appName
    || ''
})

// 当前语言下的应用描述（待审核）
const displayReviewAppDesc = computed(() => {
  if (!props.reviewInfo)
    return ''
  const langMap = reviewLangMap.value
  return langMap[currentLang.value]?.appDesc
    || langMap['zh-CN']?.appDesc
    || props.reviewInfo.appDesc
    || ''
})

// 当前语言下的应用名称（当前）
const displayCurrentAppName = computed(() => {
  if (!currentAppInfo.value)
    return ''
  const langList = (currentAppInfo.value as any).langList || []
  const langMap: Record<string, any> = {}
  langList.forEach((l: any) => {
    langMap[l.lang] = l
  })
  return langMap[currentLang.value]?.appName
    || langMap['zh-CN']?.appName
    || currentAppInfo.value.appName
    || ''
})

// 当前语言下的应用描述（当前）
const displayCurrentAppDesc = computed(() => {
  if (!currentAppInfo.value)
    return ''
  const langList = (currentAppInfo.value as any).langList || []
  const langMap: Record<string, any> = {}
  langList.forEach((l: any) => {
    langMap[l.lang] = l
  })
  return langMap[currentLang.value]?.appDesc
    || langMap['zh-CN']?.appDesc
    || currentAppInfo.value.appDesc
    || ''
})

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
      <!-- 对比展示 -->
      <div class="flex gap-6">
        <!-- 原始信息 -->
        <div v-if="currentAppInfo" class="flex-1">
          <div class="text-lg font-semibold mb-4 pb-2 border-b">
            当前发布信息
          </div>
          <NDescriptions bordered :column="1">
            <NDescriptionsItem label="应用名称">
              {{ displayCurrentAppName }}
            </NDescriptionsItem>
            <NDescriptionsItem label="应用图标">
              <img
                v-if="currentAppInfo.appIcon"
                :src="currentAppInfo.appIcon"
                class="w-16 h-16 object-contain rounded"
              >
              <span v-else class="text-gray-400">暂无图标</span>
            </NDescriptionsItem>
            <NDescriptionsItem label="应用描述">
              {{ displayCurrentAppDesc || '-' }}
            </NDescriptionsItem>
            <NDescriptionsItem label="收费方式">
              {{ microAppChargeTypeMap[currentAppInfo.chargeType] || '免费' }}
            </NDescriptionsItem>
            <NDescriptionsItem v-if="currentAppInfo.chargeType === 1" label="价格">
              {{ currentAppInfo.points }} 积分
            </NDescriptionsItem>
            <NDescriptionsItem :label="$t('microApp.thirdCharge')">
              {{ microAppThirdChargeTypeMap[currentAppInfo.thirdCharge || 0] || '不含' }}
            </NDescriptionsItem>
            <NDescriptionsItem label="包含iframe">
              {{ currentAppInfo.haveIframe ? '是' : '否' }}
            </NDescriptionsItem>
            <NDescriptionsItem label="备注">
              {{ currentAppInfo.remark || '暂无备注' }}
            </NDescriptionsItem>
            <NDescriptionsItem label="支持语言">
              <NTag v-for="lang in currentLangList" :key="lang" size="small" class="mr-1">
                {{ lang }}
              </NTag>
            </NDescriptionsItem>
          </NDescriptions>
        </div>

        <!-- 待审核信息 -->
        <div class="flex-1 bg-blue-50 -mx-4 -mt-4 p-4 border-2 border-blue-200 rounded">
          <div class="text-lg font-semibold mb-4 pb-2 border-b text-blue-600">
            待审核信息
          </div>
          <NDescriptions bordered :column="1">
            <NDescriptionsItem label="应用名称">
              {{ displayReviewAppName }}
            </NDescriptionsItem>
            <NDescriptionsItem label="应用图标">
              <img
                v-if="reviewInfo.appIcon"
                :src="reviewInfo.appIcon"
                class="w-16 h-16 object-contain rounded"
              >
              <span v-else class="text-gray-400">暂无图标</span>
            </NDescriptionsItem>
            <NDescriptionsItem label="应用描述">
              {{ displayReviewAppDesc || '-' }}
            </NDescriptionsItem>
            <NDescriptionsItem label="收费方式">
              {{ microAppChargeTypeMap[reviewInfo.chargeType] || '免费' }}
            </NDescriptionsItem>
            <NDescriptionsItem v-if="reviewInfo.chargeType === 1" label="价格">
              {{ reviewInfo.points }} 积分
            </NDescriptionsItem>
            <NDescriptionsItem :label="$t('microApp.thirdCharge')">
              {{ microAppThirdChargeTypeMap[reviewInfo.thirdCharge || 0] || '不含' }}
            </NDescriptionsItem>
            <NDescriptionsItem label="包含iframe">
              {{ reviewInfo.haveIframe ? '是' : '否' }}
            </NDescriptionsItem>
            <NDescriptionsItem label="备注">
              {{ reviewInfo.remark || '暂无备注' }}
            </NDescriptionsItem>
            <NDescriptionsItem label="支持语言">
              <NTag v-for="lang in reviewLangList" :key="lang" type="primary" size="small" class="mr-1">
                {{ lang }}
              </NTag>
            </NDescriptionsItem>
          </NDescriptions>
        </div>
      </div>

      <!-- 多语言详情 -->
      <NDivider title-placement="left">
        多语言详情
      </NDivider>
      <div class="flex gap-6">
        <!-- 当前多语言信息 -->
        <div v-if="currentAppInfo" class="flex-1">
          <div class="text-sm text-gray-500 mb-2">
            当前多语言信息
          </div>
          <div v-if="currentAppInfo?.langList && currentAppInfo.langList.length > 0">
            <div v-for="(langItem, index) in currentAppInfo.langList" :key="`current-lang-${index}`" class="mb-3 p-3 bg-gray-50 rounded">
              <div class="font-semibold text-sm mb-1">
                {{ langItem.lang }}
              </div>
              <div class="text-sm">
                <div class="mb-1">
                  <span class="text-gray-500">名称:</span> {{ langItem.appName }}
                </div>
                <div>
                  <span class="text-gray-500">描述:</span> {{ langItem.appDesc || '-' }}
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-gray-400 text-sm">
            暂无多语言信息
          </div>
        </div>
        <!-- 待审核多语言信息 -->
        <div class="flex-1">
          <div class="text-sm text-blue-600 mb-2">
            待审核多语言信息
          </div>
          <div v-if="Object.keys(reviewLangMap).length > 0">
            <div v-for="(langItem, lang) in reviewLangMap" :key="`review-lang-${lang}`" class="mb-3 p-3 bg-blue-50 rounded border border-blue-200">
              <div class="font-semibold text-sm mb-1 text-blue-700">
                {{ lang }}
              </div>
              <div class="text-sm">
                <div class="mb-1">
                  <span class="text-gray-500">名称:</span> {{ langItem.appName }}
                </div>
                <div>
                  <span class="text-gray-500">描述:</span> {{ langItem.appDesc || '-' }}
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-gray-400 text-sm">
            暂无多语言信息
          </div>
        </div>
      </div>

      <!-- 截图对比 -->
      <NDivider title-placement="left">
        应用截图
      </NDivider>
      <div class="flex gap-6">
        <div class="flex-1">
          <div class="text-sm text-gray-500 mb-2">
            当前截图
          </div>
          <div v-if="currentAppInfo?.screenshots" class="grid grid-cols-4 gap-2">
            <NImageGroup>
              <NImage
                v-for="(screenshot, index) in currentAppInfo?.screenshots.split(',').filter(s => s.trim())"
                :key="`screenshot-${index}`"
                :src="screenshot.trim()"
              />
            </NImageGroup>
          </div>
          <div v-else class="text-gray-400 text-sm">
            暂无截图
          </div>
        </div>
        <div class="flex-1">
          <div class="text-sm text-blue-600 mb-2">
            待审核截图
          </div>
          <div v-if="reviewInfo?.screenshots" class="grid grid-cols-4 gap-2">
            <NImageGroup>
              <NImage
                v-for="(screenshot, index) in reviewInfo.screenshots.split(',').filter(s => s.trim())"
                :key="`review-${index}`"
                :src="screenshot.trim()"
              />
            </NImageGroup>
          </div>
          <div v-else class="text-gray-400 text-sm">
            暂无截图
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
