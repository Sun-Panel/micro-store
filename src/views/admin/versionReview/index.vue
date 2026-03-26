<script lang="ts" setup>
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NInput, NInputGroup, NSpace, NTag, useMessage } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
import { getPendingList } from '@/api/admin/microAppVersionReview'
import { MicroAppVersionStatus, microAppVersionStatusMap } from '@/enums/panel'
import { timeFormat } from '@/utils/cmn'
import ReviewDetail from './components/ReviewDetail.vue'

const message = useMessage()

// 数据
const dataList = ref<MicroApp.VersionInfo[]>([])
const loading = ref(false)
const keyWord = ref<string>()

// 审核弹窗
const reviewShow = ref(false)
const currentVersion = ref<MicroApp.VersionInfo>()

// 获取列表
async function fetchList() {
  loading.value = true
  try {
    const { data } = await getPendingList<Common.ListResponse<MicroApp.VersionInfo[]>>({
      page: 1,
      limit: 100,
      keyword: keyWord.value,
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
      key: 'createdAt',
      width: 160,
      render(row) {
        return timeFormat(String(row.createdAt))
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
function openReview(row: MicroApp.VersionInfo) {
  currentVersion.value = row
  reviewShow.value = true
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

    <!-- 审核详情弹窗 -->
    <ReviewDetail
      v-model:visible="reviewShow"
      :version-info="currentVersion"
      @done="fetchList"
    />
  </div>
</template>
