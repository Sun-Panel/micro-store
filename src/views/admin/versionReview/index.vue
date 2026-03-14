<script lang="ts" setup>
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NInput, NInputGroup, NModal, NSpace, NTag, useMessage } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
import { adminGetPendingVersionList, adminReviewVersion } from '@/api/admin/microAppVersion'
import { MicroAppVersionStatus, microAppVersionStatusMap } from '@/enums/panel'
import { timeFormat } from '@/utils/cmn'

const message = useMessage()

// 数据
const dataList = ref<MicroApp.VersionInfo[]>([])
const loading = ref(false)
const keyWord = ref<string>()

// 审核弹窗
const reviewShow = ref(false)
const reviewLoading = ref(false)
const currentVersion = ref<MicroApp.VersionInfo>()
const reviewForm = ref({
  status: 1,
  reviewNote: '',
})

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    const { data } = await adminGetPendingVersionList<Common.ListResponse<MicroApp.VersionInfo[]>>({
      page: 1,
      limit: 100,
      keyword: keyWord.value,
    })
    dataList.value = data.list || []
  }
  catch (error) {
    message.error('获取列表失败')
  }
  finally {
    loading.value = false
  }
}

// 表格列配置
function createColumns(): DataTableColumns<MicroApp.VersionInfo> {
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
    },
    {
      title: '版本号',
      key: 'version',
      width: 100,
    },
    {
      title: '状态',
      key: 'status',
      width: 100,
      render(row) {
        const type = row.status === MicroAppVersionStatus.APPROVED ? 'success' : row.status === MicroAppVersionStatus.REJECTED ? 'error' : 'warning'
        return h(NTag, { type, size: 'small' }, {
          default: () => microAppVersionStatusMap[row.status] || '未知',
        })
      },
    },
    {
      title: '上传时间',
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
      render(row) {
        return h(NSpace, { size: 'small' }, {
          default: () => [
            row.status === MicroAppVersionStatus.PENDING
              ? h(NButton, { size: 'small', type: 'success', onClick: () => openReview(row, 1) }, {
                  default: () => '通过',
                })
              : null,
            row.status === MicroAppVersionStatus.PENDING
              ? h(NButton, { size: 'small', type: 'error', onClick: () => openReview(row, 2) }, {
                  default: () => '拒绝',
                })
              : null,
          ],
        })
      },
    },
  ]
}

const columns = createColumns()

// 打开审核弹窗
function openReview(row: MicroApp.VersionInfo, status: number) {
  currentVersion.value = row
  reviewForm.value = {
    status,
    reviewNote: '',
  }
  reviewShow.value = true
}

// 提交审核
async function handleReview() {
  if (!currentVersion.value)
    return

  reviewLoading.value = true
  try {
    const { code } = await adminReviewVersion<any>({
      versionId: currentVersion.value.id,
      status: reviewForm.value.status,
      reviewNote: reviewForm.value.reviewNote,
    })

    if (code === 0) {
      message.success(reviewForm.value.status === 1 ? '审核通过' : '已拒绝')
      reviewShow.value = false
      fetchList()
    }
  }
  catch (error) {
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
    <NCard title="待审核版本">
      <NDataTable
        :columns="columns"
        :data="dataList"
        :loading="loading"
        :bordered="false"
      />

      <!-- 无数据提示 -->
      <div v-if="dataList.length === 0 && !loading" class="text-center py-12 text-gray-400">
        暂无待审核的版本
      </div>
    </NCard>

    <!-- 审核弹窗 -->
    <NModal v-model:show="reviewShow" preset="card" style="width: 500px" title="审核版本">
      <div class="space-y-4">
        <div class="flex items-center gap-4">
          <span class="text-gray-500">应用：</span>
          <span class="font-bold">{{ currentVersion?.appName }}</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="text-gray-500">版本：</span>
          <span>{{ currentVersion?.version }}</span>
        </div>
        <div>
          <div class="mb-2">
            审核备注
          </div>
          <NInput v-model:value="reviewForm.reviewNote" type="textarea" placeholder="请输入审核备注（选填）" :rows="3" />
        </div>
      </div>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="reviewShow = false">
            取消
          </NButton>
          <NButton :type="reviewForm.status === 1 ? 'success' : 'error'" :loading="reviewLoading" @click="handleReview">
            {{ reviewForm.status === 1 ? '通过' : '拒绝' }}
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>
