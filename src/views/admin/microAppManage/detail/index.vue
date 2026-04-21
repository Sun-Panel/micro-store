<script lang="ts" setup>
import { NButton, NCard, NInput, NModal, NPopover, NSpace, NTag, useMessage } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { deletes, getInfo as getMicroAppInfo, offline, updateStatus } from '@/api/admin/microApp'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'
import { offlineVersion as adminOfflineVersion, cancelReview, getVersionList } from '@/api/admin/microAppVersion'
import { SvgIcon } from '@/components/common'
import VersionDetailModal from '@/components/common/VersionManagement/VersionDetailModal.vue'
import { microAppStatusMap } from '@/enums/panel'
import { apiRespErrMsg } from '@/utils/cmn'
import { getAppDescByLang, getAppNameByLang, getCurrentLang, getLangMapFromAppInfo } from '@/utils/functions'
import MicroAppBasicInfo from '../../developerCenter/myMicroApp/components/MicroAppBasicInfo.vue'
import MicroAppVersionInfo from '../../developerCenter/myMicroApp/components/MicroAppVersionInfo.vue'

const route = useRoute()
const router = useRouter()
const message = useMessage()

// 应用状态标签类型
function getStatusTagType(status: number) {
  switch (status) {
    case -1: return 'info' // 草稿
    case 0: return 'default' // 下架
    case 1: return 'success' // 上架
    case 2: return 'warning' // 审核中
    default: return 'default'
  }
}

// 审核状态标签类型
function getReviewStatusTagType(reviewStatus: number) {
  switch (reviewStatus) {
    case 1: return 'success' // 已通过
    case 0: return 'warning' // 审核中
    case 2: return 'error' // 已拒绝
    case -1: return 'info' // 草稿
    default: return 'default'
  }
}

// 审核未通过弹窗
const refusedPopoverShow = ref(false)

// 微应用ID
const microAppId = computed(() => Number(route.params.id))

// 数据
const microAppInfo = ref<MicroApp.Info>()
const versionList = ref<MicroApp.VersionInfo[]>([])
const loading = ref(false)
const versionLoading = ref(false)

// 从 langList 获取当前语言的 appName/appDesc
const currentLangMap = computed(() => getLangMapFromAppInfo(microAppInfo.value))
const currentLangList = computed(() => Object.keys(currentLangMap.value))
const currentLang = computed(() => getCurrentLang(currentLangList.value))
const appName = computed(() => getAppNameByLang(currentLangMap.value, currentLang.value, microAppInfo.value?.appName))
const appDesc = computed(() => getAppDescByLang(currentLangMap.value, currentLang.value, microAppInfo.value?.appDesc))

const categoryOptions = ref<Category.Info[]>([])

// 版本详情弹窗
const versionDetailShow = ref(false)
const currentVersionDetail = ref<MicroApp.VersionInfo | null>(null)

// 删除确认弹窗
const deleteDialogShow = ref(false)

// 下架弹窗
const offlineDialogShow = ref(false)
const offlineReason = ref('')

// 版本下架弹窗
const versionOfflineDialogShow = ref(false)
const versionOfflineVersion = ref<MicroApp.VersionInfo | null>(null)
const versionOfflineReason = ref('')

// 获取分类选项
async function fetchCategoryOptions() {
  try {
    const res = await getCategoryList<Category.Info[]>()
    categoryOptions.value = res.data || []
  }
  catch (error) {
    apiRespErrMsg(error)
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
      appRecordId: microAppId.value,
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

// 撤销审核
async function handleVersionCancelReview(versionId: number) {
  try {
    const res = await cancelReview({ versionId })
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

// 打开版本详情弹窗
function openVersionDetail(version: MicroApp.VersionInfo) {
  currentVersionDetail.value = version
  versionDetailShow.value = true
}

// 打开下架弹窗
function openOfflineDialog() {
  offlineReason.value = ''
  offlineDialogShow.value = true
}

// 打开版本下架弹窗
function openVersionOfflineDialog(version: MicroApp.VersionInfo) {
  versionOfflineVersion.value = version
  versionOfflineReason.value = ''
  versionOfflineDialogShow.value = true
}

// 确认下架版本
async function handleVersionOffline() {
  if (!versionOfflineVersion.value)
    return
  try {
    const res = await adminOfflineVersion<any>({
      id: versionOfflineVersion.value.id,
      type: 2, // 平台下架
      reason: versionOfflineReason.value || undefined,
    })
    if (res.code === 0) {
      message.success('已下架')
      versionOfflineDialogShow.value = false
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

// 确认下架
async function handleOffline() {
  if (!microAppInfo.value || !microAppInfo.value.id)
    return
  try {
    const res = await offline<any>({
      id: microAppInfo.value.id,
      offlineType: 2, // 平台下架
      reason: offlineReason.value || undefined,
    })
    if (res.code === 0) {
      message.success('已下架')
      offlineDialogShow.value = false
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

// 上架
async function handleChangeStatus(status: number) {
  if (!microAppInfo.value || !microAppInfo.value.id)
    return
  try {
    const res = await updateStatus({ id: microAppInfo.value.id, status })
    if (res.code === 0) {
      message.success('已上架')
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
  if (!microAppInfo.value || !microAppInfo.value.id)
    return
  deleteDialogShow.value = true
}

// 执行删除
async function handleConfirmDelete() {
  if (!microAppInfo.value || !microAppInfo.value.id)
    return

  try {
    const res = await deletes([microAppInfo.value.id])
    if (res.code === 0) {
      message.success('删除成功')
      router.push('/admin/microAppManage')
    }
    else {
      apiRespErrMsg(res)
    }
  }
  catch (error) {
    apiRespErrMsg(error)
  }
  finally {
    deleteDialogShow.value = false
  }
}

// 返回列表
function handleBack() {
  router.push('/admin/microAppManage')
}

// 预览应用
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
          <span class="text-lg font-bold">{{ appName || '微应用详情' }}</span>
          <!-- 应用状态 -->
          <div class="flex items-center gap-2">
            <NPopover v-if="microAppInfo?.status === 0 && microAppInfo.offlineReason" trigger="hover">
              <template #trigger>
                <NTag :type="getStatusTagType(microAppInfo.status)" size="small" style="cursor: pointer;">
                  {{ microAppStatusMap[microAppInfo.status] }}
                  <template v-if="microAppInfo.offlineReason !== ''" #icon>
                    <SvgIcon icon="lucide:info" />
                  </template>
                </NTag>
              </template>
              {{ microAppInfo.offlineReason }}
            </NPopover>
            <NTag v-else-if="microAppInfo?.status !== undefined" :type="getStatusTagType(microAppInfo.status)" size="small">
              {{ microAppStatusMap[microAppInfo.status] }}
            </NTag>
            <!-- 审核状态 -->
            <NTag
              v-if="microAppInfo?.reviewStatus !== undefined && microAppInfo.reviewStatus !== 0 && microAppInfo.reviewStatus !== 1"
              :type="getReviewStatusTagType(microAppInfo.reviewStatus === 3 ? 2 : microAppInfo.reviewStatus)"
              size="small"
            >
              {{ microAppInfo.reviewStatus === 2 ? '已通过' : microAppInfo.reviewStatus === 3 ? '审核未通过' : '' }}
            </NTag>
            <NPopover
              v-if="microAppInfo?.reviewStatus === 3"
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
                {{ microAppInfo?.offlineReason || '暂无原因' }}
              </div>
            </NPopover>
          </div>
        </div>
        <NSpace>
          <NButton @click="handlePreview">
            查看公开页面
          </NButton>
          <!-- 草稿状态 -->
          <template v-if="microAppInfo?.status === -1">
            <NButton type="primary" disabled>
              等待开发者提交审核
            </NButton>
          </template>
          <!-- 审核中状态 -->
          <!-- <template v-else-if="microAppInfo?.reviewStatus === 1">
            <NButton type="warning" @click="handleReject">
              拒绝
            </NButton>
            <NButton type="success" @click="handleApprove">
              通过
            </NButton>
          </template> -->
          <!-- 已上架状态 -->
          <NButton v-else-if="microAppInfo?.status === 1" @click="openOfflineDialog">
            下架
          </NButton>
          <!-- 已下架状态 -->
          <NButton v-else-if="microAppInfo?.status === 0" type="success" @click="handleChangeStatus(1)">
            上架
          </NButton>
          <NButton type="error" @click="handleDelete">
            删除
          </NButton>
        </NSpace>
      </div>
    </NCard>

    <!-- 基本信息组件（不显示编辑按钮） -->
    <NCard title="基本信息" size="small" class="mb-[20px]">
      <MicroAppBasicInfo
        :micro-app-info="microAppInfo"
        :shelves-status="microAppInfo?.status"
        :create-time="microAppInfo?.createTime"
        :category-options="categoryOptions"
        :show-edit-button="false"
        :app-name="appName"
        :app-desc="appDesc"
      />
    </NCard>

    <!-- 版本管理组件（只读，不能添加版本，不能删除，只能查看详情和下架） -->
    <MicroAppVersionInfo
      :version-list="versionList"
      :loading="versionLoading"
      :can-add-version="false"
      :can-delete-version="false"
      :can-submit-review="false"
      :can-offline-version="true"
      @view-detail="openVersionDetail"
      @cancel-review="handleVersionCancelReview"
      @offline-version="openVersionOfflineDialog"
    />

    <!-- 版本详情弹窗 -->
    <VersionDetailModal
      v-model:visible="versionDetailShow"
      :version-info="currentVersionDetail"
      :micro-app-info="microAppInfo"
      @done="fetchMicroAppInfo"
    />

    <!-- 下架弹窗 -->
    <NModal
      v-model:show="offlineDialogShow"
      preset="dialog"
      title="下架微应用"
      positive-text="确认下架"
      negative-text="取消"
      @positive-click="handleOffline"
    >
      <div class="py-4">
        <div class="mb-4">
          <span class="text-gray-600">请输入下架原因（必填）：</span>
        </div>
        <NInput
          v-model:value="offlineReason"
          type="textarea"
          placeholder="请输入下架原因"
          :rows="3"
        />
      </div>
    </NModal>

    <!-- 版本下架弹窗 -->
    <NModal
      v-model:show="versionOfflineDialogShow"
      preset="dialog"
      title="下架版本"
      positive-text="确认下架"
      negative-text="取消"
      @positive-click="handleVersionOffline"
    >
      <div class="py-4">
        <div class="mb-2 text-gray-600">
          下架版本：{{ versionOfflineVersion?.version }}
        </div>
        <div class="mb-4">
          <span class="text-gray-600">请输入下架原因（必填）：</span>
        </div>
        <NInput
          v-model:value="versionOfflineReason"
          type="textarea"
          placeholder="请输入下架原因"
          :rows="3"
        />
      </div>
    </NModal>

    <!-- 删除确认弹窗 -->
    <NModal
      v-model:show="deleteDialogShow"
      preset="dialog"
      title="删除微应用"
      type="warning"
      positive-text="确认删除"
      negative-text="取消"
      @positive-click="handleConfirmDelete"
    >
      <div class="py-4">
        <div class="text-gray-600">
          确定要删除该微应用吗？删除后将无法恢复。
        </div>
      </div>
    </NModal>
  </div>
</template>
