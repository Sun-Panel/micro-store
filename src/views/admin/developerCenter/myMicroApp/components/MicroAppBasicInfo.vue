<script lang="ts" setup>
import { NImage, NImageGroup } from 'naive-ui'
import { computed } from 'vue'
import { microAppChargeTypeMap, microAppStatusMap, microAppThirdChargeTypeMap } from '@/enums/panel'
import { timeFormat } from '@/utils/cmn'

const props = defineProps<{
  microAppInfo?: MicroApp.BaseInfo
  createTime?: string
  shelvesStatus?: number
  categoryOptions?: Category.Info[]
  langs?: string[]
  showEditButton?: boolean
  appName: string
  appDesc: string
}>()

// 分类名称
const categoryName = computed(() => {
  if (!props.microAppInfo)
    return ''
  const category = props.categoryOptions?.find(c => c.id === props.microAppInfo?.categoryId)
  return category?.name || `-`
})
</script>

<template>
  <div>
    <div v-if="microAppInfo" class="grid grid-cols-2 gap-4">
      <div class="flex items-center gap-4">
        <img v-if="microAppInfo.appIcon" :src="microAppInfo.appIcon" class="w-16 h-16 object-contain rounded">
        <div v-else class="w-16 h-16 bg-gray-100 rounded flex items-center justify-center text-gray-400">
          暂无图标
        </div>
        <div class="space-y-2">
          <div class="flex items-baseline gap-2">
            <span class="text-sm text-gray-500 whitespace-nowrap">微应用名称:</span>
            <span class="font-bold text-lg">{{ appName }}</span>
          </div>
          <div class="flex items-baseline gap-2">
            <span class="text-sm text-gray-500 whitespace-nowrap">MicroAppID:</span>
            <span class="font-mono text-sm text-gray-700">{{ microAppInfo.microAppId }}</span>
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
        <div>
          <span class="text-gray-500">{{ $t('microApp.thirdCharge') }}：</span>
          <span>{{ microAppThirdChargeTypeMap[microAppInfo.thirdCharge || 0] || '不含' }}</span>
        </div>
        <div>
          <span class="text-gray-500">包含iframe：</span>
          <span>{{ microAppInfo.haveIframe ? '是' : '否' }}</span>
        </div>
      </div>
      <div class="col-span-2">
        <span class="text-gray-500">应用描述：</span>
        <span>{{ appDesc || '-' }}</span>
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
  </div>
</template>
