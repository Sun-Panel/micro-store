<script lang="ts" setup>
import { NButton, NCard, NInput, NModal, NPopconfirm, NSpace, NTag, useMessage } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'
import { cancelReview, deletes, getInfo as getMicroAppInfo, offline, submitReview } from '@/api/admin/microAppDeveloper'
import { cancelReview as cancelVersionReview, deleteVersion, getVersionList, offlineVersion as offlineVersionApi, submitReview as submitVersionReview } from '@/api/admin/microAppVersion'
import ReviewHistoryModal from '@/components/common/ReviewHistoryModal/index.vue'
import AddVersionModal from '@/components/common/VersionManagement/AddVersionModal.vue'
import VersionDetailModal from '@/components/common/VersionManagement/VersionDetailModal.vue'
import { microAppStatusMap } from '@/enums/panel'
import { apiRespErrMsg } from '@/utils/cmn'
import MicroAppBasicInfo from '../components/MicroAppBasicInfo.vue'
import MicroAppVersionInfo from '../components/MicroAppVersionInfo.vue'
import EditMicroApp from '../EditMicroApp/index.vue'

const route = useRoute()
const router = useRouter()
const message = useMessage()

// 微应用ID
const microAppId = computed(() => Number(route.params.id))

// 数据
const microAppInfo = ref<MicroApp.MicroAppInfo>()
const versionList = ref<MicroApp.VersionInfo[]>([])
const loading = ref(false)
const versionLoading = ref(false)

// 编辑弹窗
const editDialogShow = ref(false)
const categoryOptions = ref<{ label: string, value: number }[]>([])

// 审核历史
const reviewHistoryShow = ref(false)

// 添加版本弹窗
const addVersionShow = ref(false)

// 版本详情弹窗
const versionDetailShow = ref(false)
const currentVersionDetail = ref<MicroApp.VersionInfo | null>(null)

// 版本下架弹窗
const offlineDialogShow = ref(false)
const offlineVersion = ref<MicroApp.VersionInfo | null>(null)
const offlineReason = ref('')

// 获取分类选项
async function fetchCategoryOptions() {
  try {
    const res = await getCategoryList<any>()
    categoryOptions.value = res.data?.map((item: any) => ({
      label: item.name,
      value: item.ID,
    })) || []
  }
  catch (error) {
    apiRespErrMsg(error)
  }
}

// 获取应用状态标签类型
function getStatusTagType(status: number) {
  switch (status) {
    case -1:
      return 'info' // 草稿
    case 0:
      return 'default' // 下架
    case 1:
      return 'success' // 上架
    case 2:
      return 'warning' // 审核中
    default:
      return 'default'
  }
}

// 获取审核状态标签类型
function getReviewStatusTagType(reviewStatus: number) {
  switch (reviewStatus) {
    case 0:
      return 'success' // 已通过
    case 1:
      return 'warning' // 审核中
    case 2:
      return 'error' // 已拒绝
    case 3:
      return 'info' // 草稿
    default:
      return 'default'
  }
}

// 获取审核状态文本
function getReviewStatusText(reviewStatus: number) {
  switch (reviewStatus) {
    case 0:
      return '已通过'
    case 1:
      return '审核中'
    case 2:
      return '已拒绝'
    case 3:
      return '草稿'
    default:
      return ''
  }
}

// 获取微应用详情
async function fetchMicroAppInfo() {
  loading.value = true
  try {
    const { data } = await getMicroAppInfo<any>(microAppId.value)
    microAppInfo.value = data
  }
  catch (error) {
    apiRespErrMsg(error)
  }
  finally {
    loading.value = false
  }
}

// 获取版本列表
async function fetchVersionList() {
  versionLoading.value = true
  try {
    const { data } = await getVersionList<Common.ListResponse<MicroApp.VersionInfo[]>>({
      appId: microAppId.value,
      page: 1,
      limit: 100,
    })
    versionList.value = data.list || []
  }
  catch (error) {
    apiRespErrMsg(error)
  }
  finally {
    versionLoading.value = false
  }
}

// 提交版本审核
async function handleVersionSubmitReview(versionId: number) {
  try {
    const res = await submitVersionReview<any>({ versionId })
    if (res.code === 0) {
      message.success('已提交审核')
      fetchVersionList()
    }
    else {
      apiRespErrMsg(res)
    }
  }
  catch (error) {
    apiRespErrMsg(error)
  }
}

// 撤销版本审核
async function handleVersionCancelReview(versionId: number) {
  try {
    const res = await cancelVersionReview<any>({ versionId })
    if (res.code === 0) {
      message.success('已撤销审核')
      fetchVersionList()
    }
    else {
      apiRespErrMsg(res)
    }
  }
  catch (error) {
    apiRespErrMsg(error)
  }
}

// 删除版本
async function handleDeleteVersion(ids: number[]) {
  try {
    const res = await deleteVersion<any>(ids)
    if (res.code === 0) {
      message.success('删除成功')
      fetchVersionList()
    }
    else {
      apiRespErrMsg(res)
    }
  }
  catch (error) {
    apiRespErrMsg(error)
  }
}

// 打开版本下架弹窗
function openOfflineDialog(version: MicroApp.VersionInfo) {
  offlineVersion.value = version
  offlineReason.value = ''
  offlineDialogShow.value = true
}

// 确认下架版本
async function handleOfflineVersion() {
  if (!offlineVersion.value)
    return
  try {
    const res = await offlineVersionApi<any>({
      id: offlineVersion.value.id,
      type: 1, // 作者下架
      reason: offlineReason.value || undefined,
    })
    if (res.code === 0) {
      message.success('已下架')
      offlineDialogShow.value = false
      fetchVersionList()
    }
    else {
      apiRespErrMsg(res)
    }
  }
  catch (error) {
    apiRespErrMsg(error)
  }
}

// 打开版本详情弹窗
function openVersionDetail(version: MicroApp.VersionInfo) {
  currentVersionDetail.value = version
  versionDetailShow.value = true
}

// 上架/下架
async function handleChangeStatus(status: number) {
  if (!microAppInfo.value)
    return
  // 开发者只能下架自己的应用，不能上架（上架需要审核通过）
  if (status === 1) {
    message.warning('应用需要审核通过后才能上架')
    return
  }

  try {
    const res = await offline({ id: microAppInfo.value.id, type: 1, reason: '作者主动下架' })
    if (res.code === 0) {
      message.success('已下架')
      fetchMicroAppInfo()
    }
    else {
      apiRespErrMsg(res)
    }
  }
  catch (error) {
    apiRespErrMsg(error)
  }
}

// 删除微应用
async function handleDelete() {
  if (!microAppInfo.value)
    return
  try {
    const res = await deletes([microAppInfo.value.id])
    if (res.code === 0) {
      message.success('删除成功')
      router.push('/admin/myMicroApp')
    }
    else {
      apiRespErrMsg(res)
    }
  }
  catch (error) {
    apiRespErrMsg(error)
  }
}

// 撤销微应用主信息审核
async function handleCancelAppReview() {
  if (!microAppInfo.value)
    return
  try {
    const { code } = await cancelReview({ id: microAppInfo.value.id })
    if (code === 0) {
      message.success('已撤销审核')
      fetchMicroAppInfo()
    }
  }
  catch (error) {
    message.error('撤销审核失败')
  }
}

// 提交审核
async function handleSubmitReview() {
  if (!microAppInfo.value)
    return
  try {
    const { code, msg } = await submitReview({ id: microAppInfo.value.id })
    if (code === 0) {
      message.success('已提交审核')
      fetchMicroAppInfo()
    }
    else {
      message.error(msg || '提交审核失败')
    }
  }
  catch {
    message.error('提交审核失败')
  }
}

// 处理编辑完成
function handleEditDone() {
  editDialogShow.value = false
  message.success('保存成功')
  fetchMicroAppInfo()
}

// 查看审核历史
function handleViewReviewHistory() {
  reviewHistoryShow.value = true
}

// 返回列表
function handleBack() {
  router.push('/admin/myMicroApp')
}

// 预览应用（跳转到前台公开页面）
function handlePreview() {
  const url = `/microApp/${microAppId.value}`
  window.open(url, '_blank')
}

onMounted(async () => {
  await fetchCategoryOptions()
  fetchMicroAppInfo()
  fetchVersionList()
})
</script>

<template>
  <div>
    <!-- 头部操作栏 -->
    <NCard class="mb-[20px]">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-4">
          <NButton @click="handleBack">
            返回列表
          </NButton>
          <span class="text-lg font-bold">{{ microAppInfo?.appName || '微应用详情' }}</span>
          <!-- 应用状态和审核状态 -->
          <div class="flex items-center gap-2">
            <!-- 应用状态 -->
            <NTag v-if="microAppInfo?.status !== undefined" :type="getStatusTagType(microAppInfo.status)" size="small">
              {{ microAppStatusMap[microAppInfo.status] }}
            </NTag>
            <!-- 审核状态 -->
            <NTag v-if="microAppInfo?.reviewStatus && microAppInfo.reviewStatus !== 0" :type="getReviewStatusTagType(microAppInfo.reviewStatus)" size="small">
              {{ getReviewStatusText(microAppInfo.reviewStatus) }}
            </NTag>
          </div>
        </div>
        <NSpace>
          <NButton @click="handlePreview">
            查看公开页面
          </NButton>
          <!-- 数据加载中 -->
          <template v-if="loading">
            <span class="text-gray-400">加载中...</span>
          </template>
          <!-- 审核中状态：显示撤销审核、查看审核内容 -->
          <template v-else-if="microAppInfo?.reviewStatus === 1">
            <NButton @click="handleCancelAppReview">
              撤销审核
            </NButton>
            <NButton @click="handleViewReviewHistory">
              查看审核内容
            </NButton>
          </template>
          <!-- 草稿状态：显示提交审核、编辑、删除 -->
          <template v-else-if="microAppInfo?.reviewStatus === 3">
            <NButton type="primary" @click="handleSubmitReview">
              提交审核
            </NButton>
            <NButton type="primary" @click="editDialogShow = true">
              编辑信息
            </NButton>
            <NPopconfirm @positive-click="handleDelete">
              <template #trigger>
                <NButton type="error">
                  删除应用
                </NButton>
              </template>
              确定要删除此应用吗？删除后将无法恢复。
            </NPopconfirm>
          </template>
          <!-- 审核拒绝状态：显示重新提交审核、编辑、删除 -->
          <template v-else-if="microAppInfo?.reviewStatus === 2">
            <NButton type="primary" @click="handleSubmitReview">
              重新提交审核
            </NButton>
            <NButton type="primary" @click="editDialogShow = true">
              编辑信息
            </NButton>
            <NPopconfirm @positive-click="handleDelete">
              <template #trigger>
                <NButton type="error">
                  删除应用
                </NButton>
              </template>
              确定要删除此应用吗？删除后将无法恢复。
            </NPopconfirm>
          </template>
          <!-- 审核通过状态：显示编辑、上架/下架、删除 -->
          <template v-else-if="microAppInfo?.reviewStatus === 0">
            <NButton type="primary" @click="editDialogShow = true">
              编辑信息
            </NButton>
            <NButton v-if="microAppInfo?.status === 1" @click="handleChangeStatus(0)">
              下架
            </NButton>
            <NButton v-else-if="microAppInfo?.status === 0" type="success" @click="handleChangeStatus(1)">
              上架
            </NButton>
            <NPopconfirm @positive-click="handleDelete">
              <template #trigger>
                <NButton type="error">
                  删除应用
                </NButton>
              </template>
              确定要删除此应用吗？删除后将无法恢复。
            </NPopconfirm>
          </template>
          <!-- 数据加载完成但无有效状态 -->
          <template v-else-if="!loading && microAppInfo">
            <span class="text-red-400">状态异常 (reviewStatus: {{ microAppInfo.reviewStatus }})</span>
          </template>
        </NSpace>
      </div>
    </NCard>

    <!-- 基本信息组件 -->
    <MicroAppBasicInfo
      class="mb-[20px]"
      :micro-app-info="microAppInfo"
      :category-options="categoryOptions"
    />

    <!-- 版本管理组件 -->
    <MicroAppVersionInfo
      :version-list="versionList"
      :loading="versionLoading"
      :can-add-version="true"
      :can-delete-version="true"
      :can-submit-review="true"
      :can-offline-version="true"
      @add-version="addVersionShow = true"
      @view-detail="openVersionDetail"
      @submit-review="handleVersionSubmitReview"
      @cancel-review="handleVersionCancelReview"
      @delete-version="handleDeleteVersion"
      @offline-version="openOfflineDialog"
    />

    <!-- 编辑弹窗 -->
    <EditMicroApp
      v-model:visible="editDialogShow"
      :micro-app-info="microAppInfo"
      :author-id="microAppInfo?.authorId || 0"
      :category-options="categoryOptions"
      @done="handleEditDone"
    />

    <!-- 添加版本弹窗 -->
    <AddVersionModal
      v-model:visible="addVersionShow"
      :app-id="microAppId"
      @done="fetchVersionList"
    />

    <!-- 版本详情弹窗 -->
    <VersionDetailModal
      v-model:visible="versionDetailShow"
      :version-info="currentVersionDetail"
      :micro-app-info="microAppInfo"
      @done="fetchMicroAppInfo"
    />

    <!-- 审核历史弹窗 -->
    <ReviewHistoryModal v-model:visible="reviewHistoryShow" :app-id="microAppId" />

    <!-- 版本下架弹窗 -->
    <NModal
      v-model:show="offlineDialogShow"
      preset="dialog"
      title="下架版本"
      positive-text="确认下架"
      negative-text="取消"
      @positive-click="handleOfflineVersion"
    >
      <div class="py-4">
        <div class="mb-2 text-gray-600">
          下架版本：{{ offlineVersion?.version }}
        </div>
        <div class="mb-4">
          <span class="text-gray-600">下架原因（选填）：</span>
        </div>
        <NInput
          v-model:value="offlineReason"
          type="textarea"
          placeholder="请输入下架原因"
          :rows="3"
        />
      </div>
    </NModal>
  </div>
</template>
