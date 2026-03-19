<script lang="ts" setup>
import { NButton, NCard, NModal, NTag } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { getReviewHistory } from '@/api/admin/microAppDeveloper'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'

interface Props {
  visible: boolean
  appId: number
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

const reviewHistoryList = ref<MicroApp.MicroAppReviewInfo[]>([])
const categoryList = ref<MicroAppCategory.CategoryInfo[]>([])

// 分类ID到名称的映射
const categoryMap = computed(() => {
  const map: Record<number, string> = {}
  categoryList.value.forEach((cat) => {
    map[cat.id] = cat.name
  })
  return map
})

// 双向绑定
const show = ref(props.visible)
watch(() => props.visible, (val) => {
  show.value = val
  if (val && props.appId) {
    fetchCategoryList()
    fetchReviewHistory()
  }
})
watch(show, (val) => {
  emit('update:visible', val)
})

// 获取分类列表
async function fetchCategoryList() {
  try {
    const { data } = await getCategoryList<any>()
    categoryList.value = data || []
  }
  catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 获取审核历史
async function fetchReviewHistory() {
  try {
    const { data } = await getReviewHistory<Common.ListResponse<MicroApp.MicroAppReviewInfo[]>>({
      appId: props.appId,
    })
    // 解析 langMap，确保是对象格式
    const processedList = (data?.list || []).map((review: any) => {
      if (typeof review.langMap === 'string' && review.langMap) {
        try {
          review.langMap = JSON.parse(review.langMap)
        }
        catch (e) {
          review.langMap = {}
        }
      }
      return review
    })
    reviewHistoryList.value = processedList
  }
  catch (error) {
    console.error('获取审核历史失败:', error)
  }
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 900px" title="审核历史">
    <div v-if="reviewHistoryList.length === 0" class="text-center py-8 text-gray-400">
      暂无审核历史
    </div>
    <div v-else class="space-y-4">
      <NCard v-for="review in reviewHistoryList" :key="review.id" size="small" :class="{ 'border-l-4 border-blue-500': review.status === 0 }">
        <template #header>
          <div class="flex justify-between items-center">
            <div class="flex items-center gap-2">
              <NTag v-if="review.status === 0" type="info" size="small">
                待审核
              </NTag>
              <NTag v-else-if="review.status === 1" type="success" size="small">
                已通过
              </NTag>
              <NTag v-else-if="review.status === 2" type="error" size="small">
                已拒绝
              </NTag>
            </div>
            <span class="text-sm text-gray-400">{{ new Date(review.reviewTime || review.createTime).toLocaleString() }}</span>
          </div>
        </template>

        <div class="space-y-3">
          <!-- 审核快照数据 -->
          <div class="grid grid-cols-2 gap-4">
            <!-- 左侧：基本信息 -->
            <div class="space-y-2">
              <div class="text-sm text-gray-600">
                应用名称
              </div>
              <div class="font-medium">
                {{ review.appName }}
              </div>

              <div class="text-sm text-gray-600 mt-2">
                应用图标
              </div>
              <div class="flex items-center gap-2">
                <img :src="review.appIcon" width="40" height="40" class="rounded">
              </div>

              <div class="text-sm text-gray-600 mt-2">
                所属分类
              </div>
              <div class="font-medium">
                {{ categoryMap[review.categoryId] || `分类ID: ${review.categoryId}` }}
              </div>
            </div>

            <!-- 右侧：描述、备注、截图 -->
            <div class="space-y-2">
              <div>
                <div class="text-sm text-gray-600">
                  应用描述
                </div>
                <div class="text-sm">
                  {{ review.appDesc || '-' }}
                </div>
              </div>

              <div>
                <div class="text-sm text-gray-600">
                  应用备注
                </div>
                <div class="text-sm">
                  {{ review.remark || '-' }}
                </div>
              </div>

              <div v-if="review.screenshots">
                <div class="text-sm text-gray-600">
                  应用截图
                </div>
                <div class="flex flex-wrap gap-2">
                  <img
                    v-for="(url, index) in review.screenshots.split(',')"
                    :key="index"
                    :src="url"
                    width="80"
                    height="60"
                    class="rounded border"
                  >
                </div>
              </div>
            </div>
          </div>

          <!-- 多语言信息 -->
          <div v-if="review.langMap">
            <div class="text-sm text-gray-600 mb-2">
              多语言信息：
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div v-for="(langInfo, lang) in review.langMap" :key="lang" class="space-y-1">
                <div class="text-sm text-gray-500">
                  {{ lang }} 应用
                </div>
                <div class="text-sm font-medium">
                  {{ langInfo.appName || '-' }}
                </div>
                <div class="text-sm text-gray-500">
                  {{ lang }} 描述
                </div>
                <div class="text-sm">
                  {{ langInfo.appDesc || '-' }}
                </div>
              </div>
            </div>
          </div>

          <!-- 审核信息 -->
          <div v-if="review.status !== 0" class="mt-4 pt-4 border-t border-gray-200">
            <div class="text-sm text-gray-600 mb-2">
              审核信息
            </div>
            <div class="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span class="text-gray-500">审核人ID：</span>
                <span class="ml-2">{{ review.reviewerId || '-' }}</span>
              </div>
              <div>
                <span class="text-gray-500">审核备注：</span>
                <span class="ml-2">{{ review.reviewNote || '-' }}</span>
              </div>
            </div>
          </div>
        </div>
      </NCard>
    </div>
  </NModal>
</template>
