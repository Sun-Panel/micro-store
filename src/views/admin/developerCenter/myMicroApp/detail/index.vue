<script lang="ts" setup>
import { NButton, NCard, NModal, NPopconfirm, NPopover, NSpace, NTag, useMessage } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'
import { cancelReview, deletes, getMicroInfoAndReviewInfoByMicroAppModelId, offline, submitReview } from '@/api/admin/microAppDeveloper'
import { cancelReview as cancelVersionReview, deleteVersion, getVersionList, offlineVersion as offlineVersionApi, submitReview as submitVersionReview } from '@/api/admin/microAppVersion'
import { SvgIcon } from '@/components/common'
import ReviewHistoryModal from '@/components/common/ReviewHistoryModal/index.vue'
import AddVersionModal from '@/components/common/VersionManagement/AddVersionModal.vue'
import VersionDetailModal from '@/components/common/VersionManagement/VersionDetailModal.vue'
import { microAppStatusMap } from '@/enums/panel'
import { apiRespErrMsg } from '@/utils/cmn'
import { getAppDescByLang, getAppNameByLang, getCurrentLang } from '@/utils/functions'
import MicroAppBasicInfo from '../components/MicroAppBasicInfo.vue'
import MicroAppVersionInfo from '../components/MicroAppVersionInfo.vue'
import EditMicroApp from '../EditMicroApp/index.vue'

const route = useRoute()
const router = useRouter()
const message = useMessage()

// 微应用ID
const microAppId = computed(() => Number(route.params.id))

// 数据
const reviewResponse = ref<MicroApp.GetInfoWithReviewResponse | null>(null)
const versionList = ref<MicroApp.VersionInfo[]>([])
const loading = ref(false)
const versionLoading = ref(false)

// 从 langMap 获取当前语言的 appName/appDesc
const currentLangMap = computed(() => reviewResponse.value?.microAppReview?.langMap ?? {})
const currentLangList = computed(() => Object.keys(currentLangMap.value))
const currentLang = computed(() => getCurrentLang(currentLangList.value))
const appName = computed(() => getAppNameByLang(currentLangMap.value, currentLang.value, reviewResponse.value?.microAppReview?.appName))
const appDesc = computed(() => getAppDescByLang(currentLangMap.value, currentLang.value, reviewResponse.value?.microAppReview?.appDesc))

// 编辑弹窗
const editDialogShow = ref(false)
const categoryOptions = ref<Category.Info[]>([])

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

// 拒绝详情弹窗
const refusedPopoverShow = ref(false)

// 获取分类选项
async function fetchCategoryOptions() {
  try {
    const { data } = await getCategoryList<Category.Info[]>()
    categoryOptions.value = data
    // console.log('API 获取的分类数据:', data)
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
    case 1:
      return 'success' // 已通过
    case 0:
      return 'warning' // 审核中
    case 2:
      return 'error' // 已拒绝
    case -1:
      return 'info' // 草稿
    default:
      return 'default'
  }
}

// 获取审核状态文本
function getReviewStatusText(reviewStatus: number) {
  switch (reviewStatus) {
    case 0:
      return '审核中'
    case 1:
      return '已通过'
    case 2:
      return '审核未通过'
    case -1:
      return '草稿'
    default:
      return ''
  }
}

// 获取微应用详情
async function fetchreviewResponse() {
  loading.value = true
  try {
    const res = await getMicroInfoAndReviewInfoByMicroAppModelId<MicroApp.GetInfoWithReviewResponse>(microAppId.value)
    reviewResponse.value = res.data
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
    const res = await getVersionList<Common.ListResponse<MicroApp.VersionInfo[]>>({
      appRecordId: microAppId.value,
      page: 1,
      limit: 100,
    })
    versionList.value = res.data.list || []
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
    const res = await submitVersionReview<Common.Response<null>>({ versionId })
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
    const res = await cancelVersionReview<Common.Response<null>>({ versionId })
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
    const res = await deleteVersion<Common.Response<null>>(ids)
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
    const res = await offlineVersionApi<Common.Response<null>>({
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
  if (!reviewResponse.value?.microApp?.id)
    return
  // 开发者只能下架自己的应用，不能上架（上架需要审核通过）
  if (status === 1) {
    message.warning('应用需要审核通过后才能上架')
    return
  }

  try {
    const res = await offline<Common.Response<null>>({ id: reviewResponse.value.microApp.id, offlineType: 1, reason: '作者主动下架' })
    if (res.code === 0) {
      message.success('已下架')
      fetchreviewResponse()
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
  if (!reviewResponse.value?.microApp?.id)
    return
  try {
    const res = await deletes<Common.Response<null>>([reviewResponse.value.microApp.id])
    if (res.code === 0) {
      message.success('删除成功')
      router.push({ name: 'AdminMyMicroApp' })
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
  if (!reviewResponse.value?.microApp?.id)
    return
  try {
    const { code } = await cancelReview<Common.Response<null>>({ reviewId: reviewResponse.value.microAppReview?.id || 0 })
    if (code === 0) {
      message.success('已撤销审核')
      fetchreviewResponse()
    }
  }
  catch {
    message.error('撤销审核失败')
  }
}

// 提交审核
async function handleSubmitReview() {
  if (!reviewResponse.value?.microApp?.id)
    return
  try {
    await submitReview<Common.Response<null>>({ reviewId: reviewResponse.value.microAppReview?.id || 0 })

    message.success('已提交审核')
    fetchreviewResponse()
  }
  catch (error: any) {
    if (error.code === 3004) {
      message.error('信息不完整，请点击「编辑信息」补全信息后再提交审核', { duration: 50000, closable: true })
      return
    }
    apiRespErrMsg(error)
  }
}

// 处理编辑完成
function handleEditDone() {
  editDialogShow.value = false
  message.success('保存成功')
  fetchreviewResponse()
}

// // 查看审核历史
// function handleViewReviewHistory() {
//   reviewHistoryShow.value = true
// }

// 返回列表
function handleBack() {
  router.push({ name: 'AdminMyMicroApp' })
}

// 预览应用（跳转到前台公开页面）
function handlePreview() {
  const url = `/microApp/${microAppId.value}`
  window.open(url, '_blank')
}

// // 显示拒绝详情
// function handleRefusedInfo() {
//   refusedPopoverShow.value = true
// }

onMounted(async () => {
  await fetchCategoryOptions()
  fetchreviewResponse()
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
          <span class="text-lg font-bold">{{ reviewResponse?.microAppReview?.adminName || '微应用详情' }}</span>
          <!-- 应用状态和审核状态 -->
          <div class="flex items-center gap-2">
            <!-- 应用状态 -->
            <NPopover v-if="reviewResponse?.microApp?.status === 0 && reviewResponse.microApp.offlineReason" trigger="hover">
              <template #trigger>
                <NTag :type="getStatusTagType(reviewResponse.microApp.status)" size="small" style="cursor: pointer;">
                  {{ microAppStatusMap[reviewResponse.microApp.status] }}
                  <template v-if="reviewResponse.microApp.offlineReason !== ''" #icon>
                    <SvgIcon icon="lucide:info" />
                  </template>
                </NTag>
              </template>
              {{ reviewResponse.microApp.offlineReason }}
            </NPopover>
            <NTag v-else-if="reviewResponse?.microApp?.status !== undefined" :type="getStatusTagType(reviewResponse.microApp.status)" size="small">
              {{ microAppStatusMap[reviewResponse.microApp.status] }}
            </NTag>
            <!-- 审核状态 -->
            <NTag
              v-if="reviewResponse?.microAppReview?.status !== undefined && reviewResponse?.microAppReview.status !== 1"
              :type="getReviewStatusTagType(reviewResponse.microAppReview.status)" size="small"
            >
              {{ getReviewStatusText(reviewResponse.microAppReview.status) }}
            </NTag>
            <NPopover
              v-if="reviewResponse?.microAppReview?.status === 2"
              v-model:show="refusedPopoverShow"
              trigger="click"
            >
              <template #trigger>
                <NButton size="tiny" type="error">
                  查看详情
                </NButton>
              </template>

              <div class="font-bold">
                审核未通过原因
              </div>
              <div class="text-sm">
                {{ reviewResponse?.microAppReview?.reviewNote || '暂无原因' }}
              </div>
            </NPopover>

            <!-- 草稿显示提交按钮 -->
            <NButton
              v-if="reviewResponse?.microAppReview?.status === -1"
              size="tiny"
              type="success"
              @click="handleSubmitReview"
            >
              提交审核（基本信息）
            </NButton>
            <!-- 草稿和已拒绝状态显示提交审核按钮 -->
            <NButton v-else-if="reviewResponse?.microAppReview?.status === 2" size="tiny" type="primary" @click="handleSubmitReview">
              重新提交审核
            </NButton>

            <!-- 审核中显示撤销按钮 -->
            <NButton
              v-if="reviewResponse?.microAppReview?.status === 0"
              size="tiny"
              type="error"
              text
              @click="handleCancelAppReview"
            >
              撤销
            </NButton>
          </div>
        </div>
        <NSpace>
          <!-- <NButton type="primary" @click="editDialogShow = true">
            编辑信息
          </NButton> -->
          <NButton :disabled="reviewResponse?.microApp?.status === 0" @click="handlePreview">
            查看公开页面
          </NButton>
          <NButton type="primary" @click="addVersionShow = true">
            添加版本
          </NButton>
          <!-- 数据加载中 -->
          <template v-if="loading">
            <span class="text-gray-400">加载中...</span>
          </template>
          <!-- 审核通过状态：显示编辑、上架/下架、删除 -->
          <template v-else>
            <!-- 编辑: 非审核中状态显示 -->
            <NButton type="primary" :disabled="reviewResponse?.microAppReview?.status === 0" @click="editDialogShow = true">
              编辑信息
            </NButton>
            <NPopconfirm v-if="reviewResponse?.microApp?.status === 1" @positive-click="handleChangeStatus(0)">
              <template #trigger>
                <NButton>
                  下架
                </NButton>
              </template>
              下架后重新上架需要再次审核，确定要下架吗？
            </NPopconfirm>
            <!-- 上架按钮，仅在非草稿状态，非审核状态，已下架状态显示 -->
            <NButton
              v-else-if="reviewResponse?.microApp?.status === 0 && reviewResponse?.microAppReview?.status !== 0 && reviewResponse?.microAppReview?.status !== -1"
              type="success" @click="handleSubmitReview"
            >
              上架（提交审核）
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
          <!-- <template v-else-if="!loading && reviewResponse">
            <span class="text-red-400">状态异常 (reviewStatus: {{ reviewResponse.microAppReview?.status }})</span>
          </template> -->
        </NSpace>
      </div>
    </NCard>

    <!-- 基本信息组件 -->
    <NCard title="基本信息" size="small" class="mb-2">
      <MicroAppBasicInfo
        class="mb-[20px]"
        :micro-app-info="reviewResponse?.microAppReview ?? undefined"
        :shelves-status="reviewResponse?.microApp.status"
        :create-time="reviewResponse?.microApp.createTime"
        :app-name="appName"
        :app-desc="appDesc"
        :category-options="categoryOptions"
      />
    </NCard>

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
      :micro-app-info="reviewResponse?.microAppReview ?? undefined"
      :author-id="reviewResponse?.microApp?.developerId || 0"
      :category-options="categoryOptions"
      @done="handleEditDone"
    />

    <!-- 添加版本弹窗 -->
    <AddVersionModal
      v-model:visible="addVersionShow"
      :app-record-id="microAppId"
      @done="fetchVersionList"
    />

    <!-- 版本详情弹窗 -->
    <VersionDetailModal
      v-model:visible="versionDetailShow"
      :version-info="currentVersionDetail"
      :micro-app-info="reviewResponse?.microApp ?? undefined"
      @done="fetchreviewResponse"
    />

    <!-- 审核历史弹窗 -->
    <ReviewHistoryModal v-model:visible="reviewHistoryShow" :app-record-id="microAppId" />

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
          下架版本：{{ offlineVersion?.version }}，下架后再次上架需要重新提交审核
        </div>
        <!-- <div class="mb-4">
          <span class="text-gray-600">下架原因（选填）：</span>
        </div> -->
        <!-- <NInput
          v-model:value="offlineReason"
          type="textarea"
          placeholder="请输入下架原因"
          :rows="3"
        /> -->
      </div>
    </NModal>
  </div>
</template>
