<script lang="ts" setup>
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NDescriptions, NDescriptionsItem, NDivider, NInput, NInputGroup, NModal, NSpace, NTag, useMessage } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
import { adminGetPendingVersionList, adminReviewVersion, getVersionList } from '@/api/admin/microAppVersion'
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
const currentApprovedVersion = ref<MicroApp.VersionInfo>()
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
async function openReview(row: MicroApp.VersionInfo) {
  currentVersion.value = row
  reviewForm.value = {
    status: 1,
    reviewNote: '',
  }

  // 获取当前最新已通过版本
  try {
    const { data } = await getVersionList<Common.ListResponse<MicroApp.VersionInfo[]>>({
      appId: row.appId,
      page: 1,
      limit: 1,
      status: 1, // 已通过
    })
    currentApprovedVersion.value = data.list && data.list.length > 0 ? data.list[0] : undefined
  }
  catch (error) {
    message.error('获取版本信息失败')
  }

  reviewShow.value = true
}

// 提交审核
async function handleReview() {
  if (!currentVersion.value)
    return

  // 驳回时必须填写原因
  if (reviewForm.value.status === 2 && !reviewForm.value.reviewNote?.trim()) {
    message.error('驳回时必须填写驳回原因')
    return
  }

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

// 下载版本包
function handleDownload(url: string) {
  window.open(url, '_blank')
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
    <NModal v-model:show="reviewShow" preset="card" style="width: 1200px;" title="审核版本">
      <div v-if="currentVersion" class="space-y-6">
        <!-- 对比展示 -->
        <div class="flex gap-6">
          <!-- 当前已发布版本 -->
          <div class="flex-1">
            <div class="text-lg font-semibold mb-4 pb-2 border-b">
              当前已发布版本
              <span v-if="!currentApprovedVersion" class="text-sm font-normal text-gray-400">（暂无）</span>
            </div>
            <NDescriptions v-if="currentApprovedVersion" bordered :column="1">
              <NDescriptionsItem label="应用名称">
                {{ currentVersion.appName }}
              </NDescriptionsItem>
              <NDescriptionsItem label="版本号">
                {{ currentApprovedVersion.version }}
              </NDescriptionsItem>
              <NDescriptionsItem label="版本代码">
                {{ currentApprovedVersion.versionCode }}
              </NDescriptionsItem>
              <NDescriptionsItem label="版本说明">
                {{ currentApprovedVersion.versionDesc || '暂无说明' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="包地址">
                <a :href="currentApprovedVersion.packageUrl" target="_blank" class="text-blue-600 hover:underline">
                  {{ currentApprovedVersion.packageUrl }}
                </a>
              </NDescriptionsItem>
              <NDescriptionsItem label="包校验值">
                {{ currentApprovedVersion.packageHash }}
              </NDescriptionsItem>
              <NDescriptionsItem label="API 版本" v-if="currentApprovedVersion.config?.apiVersion">
                {{ currentApprovedVersion.config.apiVersion }}
              </NDescriptionsItem>
              <NDescriptionsItem label="作者" v-if="currentApprovedVersion.config?.author">
                {{ currentApprovedVersion.config.author }}
              </NDescriptionsItem>
            </NDescriptions>
            <div v-else class="text-center py-8 text-gray-400">
              暂无已发布的版本
            </div>
          </div>

          <!-- 待审核版本 -->
          <div class="flex-1 bg-blue-50 -mx-4 -mt-4 p-4 border-2 border-blue-200 rounded">
            <div class="text-lg font-semibold mb-4 pb-2 border-b text-blue-600">
              待审核版本
            </div>
            <NDescriptions bordered :column="1">
              <NDescriptionsItem label="应用名称">
                {{ currentVersion.appName }}
              </NDescriptionsItem>
              <NDescriptionsItem label="版本号" :class="{ 'font-bold text-red-600': !currentApprovedVersion || currentVersion.version !== currentApprovedVersion.version }">
                {{ currentVersion.version }}
                <span v-if="!currentApprovedVersion || currentVersion.version !== currentApprovedVersion.version" class="ml-2 text-xs bg-red-100 text-red-600 px-2 py-1 rounded">新版本</span>
              </NDescriptionsItem>
              <NDescriptionsItem label="版本代码">
                {{ currentVersion.versionCode }}
              </NDescriptionsItem>
              <NDescriptionsItem label="版本说明">
                {{ currentVersion.versionDesc || '暂无说明' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="包地址">
                <a :href="currentVersion.packageUrl" target="_blank" class="text-blue-600 hover:underline">
                  {{ currentVersion.packageUrl }}
                </a>
                <NButton
                  size="tiny"
                  type="primary"
                  class="ml-2"
                  @click="handleDownload(currentVersion.packageUrl)"
                >
                  下载
                </NButton>
              </NDescriptionsItem>
              <NDescriptionsItem label="包校验值">
                {{ currentVersion.packageHash }}
              </NDescriptionsItem>
              <NDescriptionsItem label="API 版本" v-if="currentVersion.config?.apiVersion">
                {{ currentVersion.config.apiVersion }}
              </NDescriptionsItem>
              <NDescriptionsItem label="作者" v-if="currentVersion.config?.author">
                {{ currentVersion.config.author }}
              </NDescriptionsItem>
            </NDescriptions>
          </div>
        </div>

        <!-- 权限对比 -->
        <NDivider v-if="currentVersion.config?.permissions?.length || (currentApprovedVersion?.config?.permissions?.length)" title-placement="left">
          权限对比
        </NDivider>
        <div v-if="currentVersion.config?.permissions?.length || (currentApprovedVersion?.config?.permissions?.length)" class="flex gap-6">
          <div class="flex-1">
            <div class="text-sm text-gray-500 mb-2">
              当前版本权限
            </div>
            <div class="space-y-1">
              <div v-if="currentApprovedVersion?.config?.permissions?.length" v-for="(perm, index) in currentApprovedVersion.config.permissions" :key="index" class="px-3 py-1 bg-gray-100 rounded text-sm">
                {{ perm }}
              </div>
              <div v-else class="text-gray-400 text-sm">
                无权限要求
              </div>
            </div>
          </div>
          <div class="flex-1">
            <div class="text-sm text-blue-600 mb-2">
              待审核版本权限
            </div>
            <div class="space-y-1">
              <div v-for="(perm, index) in (currentVersion.config?.permissions || [])" :key="index" class="px-3 py-1 bg-blue-100 text-blue-700 rounded text-sm">
                {{ perm }}
                <span v-if="!currentApprovedVersion?.config?.permissions?.includes(perm)" class="ml-2 text-xs bg-red-100 text-red-600 px-2 py-1 rounded">新增</span>
              </div>
              <div v-if="!currentVersion.config?.permissions?.length" class="text-gray-400 text-sm">
                无权限要求
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
