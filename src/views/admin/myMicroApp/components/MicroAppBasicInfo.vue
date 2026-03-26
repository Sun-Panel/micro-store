<script lang="ts" setup>
import { NCard, NImage, NImageGroup, NSelect, NTag } from 'naive-ui'
import { computed, ref } from 'vue'
import { microAppChargeTypeMap, microAppStatusMap } from '@/enums/panel'
import { timeFormat } from '@/utils/cmn'
import { getAppDescByLang, getAppNameByLang, getCurrentLang, getLangListFromAppInfo, getLangMapFromAppInfo } from '@/utils/functions'

const props = defineProps<{
  microAppInfo?: MicroApp.BaseInfo
  createTime?: string
  shelvesStatus?: number
  categoryOptions?: { label: string, value: number }[]
  showEditButton?: boolean
}>()

// 分类名称
const categoryName = computed(() => {
  if (!props.microAppInfo)
    return ''
  const category = props.categoryOptions?.find(c => c.value === props.microAppInfo?.categoryId)
  return category?.label || `ID: ${props.microAppInfo.categoryId}`
})

// ==================== 多语言处理 ====================
const baseInfoLang = ref('zh-CN')

// 微应用的多语言列表
const baseInfoLangList = computed(() => getLangListFromAppInfo(props.microAppInfo))

// 微应用语言 Map
const baseInfoLangMap = computed(() => getLangMapFromAppInfo(props.microAppInfo))

// 初始化语言
function initLang() {
  baseInfoLang.value = getCurrentLang(baseInfoLangList.value)
}

// 初始化
initLang()

// 当前语言下的应用名称
const baseInfoAppName = computed(() => getAppNameByLang(baseInfoLangMap.value, baseInfoLang.value, props.microAppInfo?.appName))

// 当前语言下的应用描述
const baseInfoAppDesc = computed(() => getAppDescByLang(baseInfoLangMap.value, baseInfoLang.value, props.microAppInfo?.appDesc))
</script>

<template>
  <NCard title="基本信息">
    <template #header>
      <div class="flex items-center gap-2">
        基本信息
        <NSelect
          v-model:value="baseInfoLang"
          :options="baseInfoLangList.map((lang: string) => ({ label: lang, value: lang }))"
          style="width: 140px"
          size="small"
        />
      </div>
    </template>
    <div v-if="microAppInfo" class="grid grid-cols-2 gap-4">
      <div class="flex items-center gap-4">
        <img v-if="microAppInfo.appIcon" :src="microAppInfo.appIcon" class="w-16 h-16 object-contain rounded">
        <div v-else class="w-16 h-16 bg-gray-100 rounded flex items-center justify-center text-gray-400">
          暂无图标
        </div>
        <div class="space-y-2">
          <div class="flex items-baseline gap-2">
            <span class="text-sm text-gray-500 whitespace-nowrap">MicroAppID:</span>
            <span class="font-mono text-sm text-gray-700">{{ microAppInfo.microAppId }}</span>
          </div>
          <div class="flex items-baseline gap-2">
            <span class="text-sm text-gray-500 whitespace-nowrap">微应用名称:</span>
            <span class="font-bold text-lg">{{ baseInfoAppName || microAppInfo.appName }}</span>
          </div>
        </div>
      </div>
      <div v-if="shelvesStatus !== undefined" class="grid grid-cols-2 gap-2 text-sm">
        <div>
          <span class="text-gray-500">状态：</span>
          <span
            :class="{
              'text-blue-500': shelvesStatus === -1,
              'text-green-500': shelvesStatus === 1,
              'text-yellow-500': shelvesStatus === 2,
              'text-gray-500': shelvesStatus === 0,
            }"
          >{{ microAppStatusMap[shelvesStatus] }}</span>
        </div>
        <div>
          <span class="text-gray-500">收费方式：</span>
          <span>{{ microAppChargeTypeMap[microAppInfo.chargeType] || '免费' }}</span>
        </div>
        <div>
          <span class="text-gray-500">分类：</span>
          <span>{{ categoryName }}</span>
        </div>
        <div>
          <span class="text-gray-500">创建时间：</span>
          <span>{{ timeFormat(String(createTime)) }}</span>
        </div>
      </div>
      <div class="col-span-2">
        <span class="text-gray-500">应用描述：</span>
        <span>{{ baseInfoAppDesc || microAppInfo.appDesc || '暂无描述' }}</span>
      </div>
      <!-- 图集 -->
      <div v-if="microAppInfo.screenshots" class="col-span-2">
        <span class="text-gray-500">图集：</span>
        <NImageGroup>
          <div class="flex flex-wrap gap-2 mt-2">
            <NImage
              v-for="(screenshot, index) in microAppInfo.screenshots.split(',')"
              :key="index"
              :src="screenshot"
              width="96"
              height="96"
              style="object-fit: cover; border-radius: 6px;"
            />
          </div>
        </NImageGroup>
      </div>
    </div>
  </NCard>
</template>
