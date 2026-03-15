<script lang="ts" setup>
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NDescriptions, NDescriptionsItem, NDivider, NImage, NImageGroup, NInput, NInputGroup, NModal, NSpace, NTag, useMessage } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
import { getInfo, getPendingReviewList, reviewApp } from '@/api/admin/microApp'
import { microAppChargeTypeMap, MicroAppReviewStatus, microAppReviewStatusMap } from '@/enums/panel'
import { timeFormat } from '@/utils/cmn'

const message = useMessage()

// 数据
const dataList = ref<MicroApp.MicroAppReviewInfo[]>([])
const loading = ref(false)
const keyWord = ref<string>()

// 审核弹窗
const reviewShow = ref(false)
const reviewLoading = ref(false)
const currentReview = ref<MicroApp.MicroAppReviewInfo>()
const currentAppInfo = ref<MicroApp.MicroAppInfo>()
const reviewForm = ref({
  status: 1,
  reviewNote: '',
})

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    const { data } = await getPendingReviewList<Common.ListResponse<MicroApp.MicroAppReviewInfo[]>>({
      page: 1,
      limit: 100,
    })
    dataList.value = data.list || []
  }
  catch {
    message.error('获取列表失败')
  }
  finally {
    loading.value = false
  }
}

// 表格列配置
function createColumns(): DataTableColumns<MicroApp.MicroAppReviewInfo> {
  return [
    {
      title: 'ID',
      key: 'id',
      width: 60,
    },
    {
      title: '应用名称',
      key: 'appName',
      width: 150,
      ellipsis: {
        tooltip: true,
      },
    },
    {
      title: '审核状态',
      key: 'status',
      width: 100,
      render(row) {
        const type = row.status === MicroAppReviewStatus.APPROVED ? 'success' : row.status === MicroAppReviewStatus.REJECTED ? 'error' : 'warning'
        return h(NTag, { type, size: 'small' }, {
          default: () => microAppReviewStatusMap[row.status] || '未知',
        })
      },
    },
    {
      title: '提交时间',
      key: 'createTime',
      width: 160,
      render(row) {
        return timeFormat(String(row.createTime))
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 120,
      fixed: 'right' as const,
      render(row) {
        return h(NSpace, { size: 'small' }, {
          default: () => [
            h(NButton, { size: 'small', type: 'primary', onClick: () => openReview(row) }, {
              default: () => '审核',
            }),
          ],
        })
      },
    },
  ]
}

const columns = createColumns()

// 打开审核弹窗
async function openReview(row: MicroApp.MicroAppReviewInfo) {
  currentReview.value = row
  reviewForm.value = {
    status: 1,
    reviewNote: '',
  }

  // 获取当前应用信息
  try {
    const { data } = await getInfo<MicroApp.MicroAppInfo>(row.appId)
    currentAppInfo.value = data
  }
  catch {
    message.error('获取应用信息失败')
  }

  reviewShow.value = true
}

// 提交审核
async function handleReview() {
  if (!currentReview.value)
    return

  // 驳回时必须填写原因
  if (reviewForm.value.status === 2 && !reviewForm.value.reviewNote?.trim()) {
    message.error('驳回时必须填写驳回原因')
    return
  }

  reviewLoading.value = true
  try {
    const { code } = await reviewApp<any>({
      reviewId: currentReview.value.id,
      status: reviewForm.value.status,
      reviewNote: reviewForm.value.reviewNote,
    })

    if (code === 0) {
      message.success(reviewForm.value.status === 1 ? '审核通过' : '已拒绝')
      reviewShow.value = false
      fetchList()
    }
  }
  catch {
    message.error('操作失败')
  }
  finally {
    reviewLoading.value = false
  }
}

function handleSearch() {
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>

<template>
  <div>
    <!-- 搜索栏 -->
    <NCard class="mb-[20px]">
      <div class="flex">
        <NInputGroup style="max-width: 500px;">
          <NInput v-model:value="keyWord" placeholder="请输入应用名称搜索" @keyup.enter="handleSearch" />
          <NButton type="primary" @click="handleSearch">
            查询
          </NButton>
        </NInputGroup>
      </div>
    </NCard>

    <!-- 表格 -->
    <NCard title="待审核微应用">
      <NDataTable
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :bordered="false"
      />

      <!-- 无数据提示 -->
      <div v-if="dataList.length === 0 && !loading" class="text-center py-12 text-gray-400">
        暂无待审核的微应用
      </div>
    </NCard>

    <!-- 审核弹窗 -->
    <NModal v-model:show="reviewShow" preset="card" style="width: 1200px;" title="审核微应用">
      <div v-if="currentReview" class="space-y-6">
        <!-- 对比展示 -->
        <div class="flex gap-6">
          <!-- 原始信息 -->
          <div v-if="currentAppInfo" class="flex-1">
            <div class="text-lg font-semibold mb-4 pb-2 border-b">
              当前发布信息
            </div>
            <NDescriptions bordered :column="1">
              <NDescriptionsItem label="应用名称">
                {{ currentAppInfo.appName }}
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
                {{ currentAppInfo.appDesc || '-' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="收费方式">
                {{ microAppChargeTypeMap[currentAppInfo.chargeType] || '免费' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="currentAppInfo.chargeType === 1" label="价格">
                {{ currentAppInfo.price }} 积分
              </NDescriptionsItem>
              <NDescriptionsItem label="备注">
                {{ currentAppInfo.remark || '暂无备注' }}
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
                {{ currentReview.appName }}
              </NDescriptionsItem>
              <NDescriptionsItem label="应用图标">
                <img
                  v-if="currentReview.appIcon"
                  :src="currentReview.appIcon"
                  class="w-16 h-16 object-contain rounded"
                >
                <span v-else class="text-gray-400">暂无图标</span>
              </NDescriptionsItem>
              <NDescriptionsItem label="应用描述">
                {{ currentReview.appDesc || '-' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="收费方式">
                {{ microAppChargeTypeMap[currentReview.chargeType] || '免费' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="currentReview.chargeType === 1" label="价格">
                {{ currentReview.price }} 积分
              </NDescriptionsItem>
              <NDescriptionsItem label="备注">
                {{ currentReview.remark || '暂无备注' }}
              </NDescriptionsItem>
            </NDescriptions>
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
                  v-for="(screenshot, index) in currentReview.screenshots.split(',').filter(s => s.trim())"
                  :key="`review-${index}`"
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
            <div v-if="currentReview?.screenshots" class="grid grid-cols-4 gap-2">
              <NImageGroup>
                <NImage
                  v-for="(screenshot, index) in currentReview.screenshots.split(',').filter(s => s.trim())"
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
          <NButton @click="reviewShow = false">
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
  </div>
</template>
