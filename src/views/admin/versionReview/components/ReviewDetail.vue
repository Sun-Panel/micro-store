<script lang="ts" setup>
import { NButton, NCard, NCollapse, NCollapseItem, NDescriptions, NDescriptionsItem, NDivider, NInput, NModal, NProgress, NSpace, NSwitch, NTag, useDialog, useMessage } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { getDownloadUrl, getLatestOnlineByAppModelId } from '@/api/admin/microAppVersion'
import { review, triggerSecurityAudit } from '@/api/admin/microAppVersionReview'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'

const props = defineProps<{
  visible: boolean
  versionInfo?: MicroApp.VersionInfo
  microApp?: MicroApp.Info
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'done'): void
}>()

const message = useMessage()
const dialog = useDialog()

// 数据
const reviewLoading = ref(false)
const currentApprovedVersion = ref<MicroApp.VersionInfo>()
const reviewForm = ref({
  status: 1,
  reviewNote: '',
})
const iframeModalVisible = ref(false)
const securityAuditModalVisible = ref(false)
const triggerAuditLoading = ref(false)

// 同步滚动 & 仅显示差异
const syncScrollEnabled = ref(true)
const onlyShowDiff = ref(false)
const leftPanelRef = ref<HTMLDivElement>()
const rightPanelRef = ref<HTMLDivElement>()
let isSyncing = false

function handleLeftScroll() {
  if (!syncScrollEnabled.value || isSyncing)
    return
  isSyncing = true
  requestAnimationFrame(() => {
    if (rightPanelRef.value && leftPanelRef.value)
      rightPanelRef.value.scrollTop = leftPanelRef.value.scrollTop
    isSyncing = false
  })
}

function handleRightScroll() {
  if (!syncScrollEnabled.value || isSyncing)
    return
  isSyncing = true
  requestAnimationFrame(() => {
    if (leftPanelRef.value && rightPanelRef.value)
      leftPanelRef.value.scrollTop = rightPanelRef.value.scrollTop
    isSyncing = false
  })
}

// 差异判断
function isDiff(field: string): boolean {
  if (!currentApprovedVersion.value || !props.versionInfo)
    return true

  const cur = currentApprovedVersion.value
  const rev = props.versionInfo

  switch (field) {
    case 'version': return cur.version !== rev.version
    case 'versionDesc': return getVersionDescContent(cur.versionDesc) !== getVersionDescContent(rev.versionDesc)
    case 'packageUrl': return cur.packageUrl !== rev.packageUrl
    case 'packageHash': return cur.packageHash !== rev.packageHash
    case 'apiVersion': return cur.config?.apiVersion !== rev.config?.apiVersion
    case 'author': return cur.config?.author !== rev.config?.author
    case 'permissions':
      return JSON.stringify(cur.config?.permissions) !== JSON.stringify(rev.config?.permissions)
    default: return true
  }
}

// 获取版本说明（兼容多语言格式）
function getVersionDescContent(versionDesc: Record<string, { content: string }> | string | undefined): string {
  if (!versionDesc)
    return ''
  if (typeof versionDesc === 'string')
    return versionDesc
  return versionDesc['zh-CN']?.content
    || Object.values(versionDesc)[0]?.content
    || ''
}

// 严重程度映射
const severityMap: Record<string, { label: string, color: 'error' | 'warning' | 'default' | 'success', value: number }> = {
  CRITICAL: { label: '高危', color: 'error', value: 4 },
  HIGH: { label: '高', color: 'warning', value: 3 },
  MEDIUM: { label: '中', color: 'default', value: 2 },
  LOW: { label: '低', color: 'success', value: 1 },
  critical: { label: '高危', color: 'error', value: 4 },
  high: { label: '高', color: 'warning', value: 3 },
  medium: { label: '中', color: 'default', value: 2 },
  low: { label: '低', color: 'success', value: 1 },
}

function getSeverityInfo(severity: string) {
  return severityMap[severity] || { label: severity, color: 'default' as const, value: 0 }
}

// 权限分组
interface PermissionGroupItem {
  name: string
  isNew: boolean
  isRemoved: boolean
}

interface PermissionGroup {
  category: string
  items: PermissionGroupItem[]
}

function buildPermissionGroups(perms: string[], otherPerms: string[], mode: 'current' | 'review'): PermissionGroup[] {
  if (perms.length === 0 && otherPerms.length === 0)
    return []

  const groups: Record<string, PermissionGroup> = {}

  for (const perm of perms) {
    const prefix = perm.includes(':') ? perm.split(':')[0] : '其他'
    if (!groups[prefix]) {
      groups[prefix] = { category: prefix, items: [] }
    }
    groups[prefix].items.push({
      name: perm,
      isNew: mode === 'review' && !otherPerms.includes(perm),
      isRemoved: mode === 'current' && !otherPerms.includes(perm),
    })
  }

  return Object.values(groups).sort((a, b) => a.category.localeCompare(b.category))
}

const currentPermissionGroups = computed<PermissionGroup[]>(() => {
  const curPerms = currentApprovedVersion.value?.config?.permissions || []
  const revPerms = props.versionInfo?.config?.permissions || []
  return buildPermissionGroups(curPerms, revPerms, 'current')
})

const reviewPermissionGroups = computed<PermissionGroup[]>(() => {
  const curPerms = currentApprovedVersion.value?.config?.permissions || []
  const revPerms = props.versionInfo?.config?.permissions || []
  return buildPermissionGroups(revPerms, curPerms, 'review')
})

const hasPermissionDiff = computed(() => {
  const curPerms = currentApprovedVersion.value?.config?.permissions || []
  const revPerms = props.versionInfo?.config?.permissions || []
  return JSON.stringify(curPerms) !== JSON.stringify(revPerms)
})

// 始终显示所有权限分组，不做过滤
const filteredCurrentPermissionGroups = computed(() => {
  return currentPermissionGroups.value
})

const filteredReviewPermissionGroups = computed(() => {
  return reviewPermissionGroups.value
})

// iframe 权限检查（基于微应用信息的 haveIframe 字段）
const hasIframePermission = computed(() => {
  return props.microApp?.haveIframe ?? false
})

const currentHasIframe = computed(() => {
  return currentApprovedVersion.value?.config?.permissions?.some(p => p.includes('iframe') || p.includes('IFRAME')) ?? false
})

const reviewHasIframe = computed(() => {
  return props.versionInfo?.config?.permissions?.some(p => p.includes('iframe') || p.includes('IFRAME')) ?? false
})

// 监听弹窗打开，获取已发布版本信息
watch(() => props.visible, async (visible) => {
  if (visible && props.versionInfo) {
    reviewForm.value = {
      status: 1,
      reviewNote: '',
    }

    // 获取当前最新已通过版本
    try {
      const { data } = await getLatestOnlineByAppModelId<MicroApp.VersionInfo>(props.versionInfo.appRecordId)
      currentApprovedVersion.value = data
    }
    catch (error: any) {
      if (error.code !== 1200) {
        apiRespErrMsg(error)
      }
    }
  }
})

// 提交审核
async function handleReview() {
  if (!props.versionInfo)
    return

  if (reviewForm.value.status === 2 && !reviewForm.value.reviewNote?.trim()) {
    message.error('驳回时必须填写驳回原因')
    return
  }

  if (reviewForm.value.status === 1 && !props.versionInfo.codeSecurityAudit) {
    dialog.warning({
      title: '确认通过',
      content: '当前版本没有安全审核报告，如果执意通过后，该版本出现问题将扣除相应的积分，确定要继续通过吗？',
      positiveText: '确定通过',
      negativeText: '取消',
      onPositiveClick: () => {
        doReview()
      },
    })
    return
  }

  await doReview()
}

// 执行审核操作
async function doReview() {
  if (!props.versionInfo)
    return

  reviewLoading.value = true
  try {
    const { code } = await review<any>({
      versionId: props.versionInfo.id,
      status: reviewForm.value.status,
      reviewNote: reviewForm.value.reviewNote,
    })

    if (code === 0) {
      message.success(reviewForm.value.status === 1 ? '审核通过' : '已拒绝')
      emit('update:visible', false)
      emit('done')
    }
  }
  catch (error: any) {
    apiRespErrMsg(error)
  }
  finally {
    reviewLoading.value = false
  }
}

// 下载版本包
async function handleDownloadByVersionId(versionId: number) {
  await getDownloadUrl<string>(versionId).then(({ data }) => {
    window.open(data, '_blank')
  }).catch(() => {
    message.error('下载失败，请重试')
  })
}

// 打开应用公开页面
function handleOpenMicroAppPublic() {
  iframeModalVisible.value = true
}

// 打开安全审核报告弹窗
function handleOpenSecurityAudit() {
  securityAuditModalVisible.value = true
}

// 主动触发安全审核
async function handleTriggerSecurityAudit() {
  if (!props.versionInfo)
    return

  triggerAuditLoading.value = true
  try {
    const { code } = await triggerSecurityAudit<any>(props.versionInfo.id)
    if (code === 0) {
      message.success('已触发安全审核，请稍后刷新查看结果')
    }
  }
  catch (error: any) {
    apiRespErrMsg(error)
  }
  finally {
    triggerAuditLoading.value = false
  }
}

// 打开外部链接
function openExternalUrl(url: string) {
  window.open(url, '_blank')
}
</script>

<template>
  <NModal :show="visible" preset="card" style="width: 1200px;" title="审核版本" @update:show="emit('update:visible', $event)">
    <template #header>
      <div class="flex gap-2 items-center">
        <div class="flex justify-between">
          {{ versionInfo?.microApp?.appName }} - 版本审核
        </div>
        <div>
          <NButton size="small" @click="handleOpenMicroAppPublic">
            查看应用公开页面
          </NButton>
        </div>
      </div>
    </template>
    <div v-if="versionInfo" class="space-y-6">
      <!-- 工具栏 -->
      <div class="flex items-center gap-6 px-2">
        <div class="flex items-center gap-2">
          <NSwitch v-model:value="syncScrollEnabled" />
          <span class="text-sm">同步滚动</span>
        </div>
        <div class="flex items-center gap-2">
          <NSwitch v-model:value="onlyShowDiff" />
          <span class="text-sm">仅显示修改</span>
        </div>
      </div>

      <!-- 对比展示 -->
      <div class="flex gap-6">
        <!-- 当前已发布版本 -->
        <div class="flex-1 min-w-0">
          <div class="text-lg font-semibold mb-4 pb-2 border-b">
            当前已发布版本
            <span v-if="!currentApprovedVersion" class="text-sm font-normal text-gray-400">（暂无）</span>
          </div>
          <div v-if="currentApprovedVersion" ref="leftPanelRef" class="scroll-panel" @scroll="handleLeftScroll">
            <NDescriptions bordered :column="1">
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('version')">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('version') }">版本号</span>
                </template>
                {{ currentApprovedVersion.version }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('versionDesc')">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('versionDesc') }">版本说明</span>
                </template>
                {{ getVersionDescContent(currentApprovedVersion.versionDesc) || '暂无说明' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('packageUrl')">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('packageUrl') }">包地址</span>
                </template>
                <a :href="currentApprovedVersion.packageUrl" target="_blank" class="text-blue-600 hover:underline">
                  {{ currentApprovedVersion.packageUrl }}
                </a>
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('packageHash')">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('packageHash') }">包校验值</span>
                </template>
                {{ currentApprovedVersion.packageHash }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="(!onlyShowDiff || isDiff('apiVersion')) && currentApprovedVersion.config?.apiVersion">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('apiVersion') }">API 版本</span>
                </template>
                {{ currentApprovedVersion.config.apiVersion }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="(!onlyShowDiff || isDiff('author')) && currentApprovedVersion.config?.author">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('author') }">作者</span>
                </template>
                {{ currentApprovedVersion.config.author }}
              </NDescriptionsItem>
              <NDescriptionsItem label="安全审核">
                <NButton size="small" @click="handleOpenSecurityAudit">
                  {{ currentApprovedVersion.codeSecurityAudit ? '查看报告' : '无报告' }}
                </NButton>
              </NDescriptionsItem>
            </NDescriptions>

          </div>
          <div v-else class="text-center py-8 text-gray-400">
            暂无已发布的版本
          </div>
          <!-- 当前版本权限（始终显示） -->
          <div>
            <NDivider title-placement="left" class="!my-3">
              权限（{{ currentApprovedVersion?.config?.permissions?.length || 0 }}项）
              <NTag v-if="currentHasIframe" size="small" type="warning" class="ml-2">
                含 iframe
              </NTag>
            </NDivider>
            <div v-for="group in filteredCurrentPermissionGroups" :key="group.category" class="mb-3">
              <div class="text-xs text-gray-400 mb-1 uppercase tracking-wide">
                {{ group.category }}
              </div>
              <div class="space-y-1">
                <div
                  v-for="item in group.items"
                  :key="item.name"
                  class="px-3 py-1 rounded text-sm flex items-center justify-between"
                  :class="item.isRemoved ? 'bg-red-50 text-red-500 line-through' : 'bg-gray-100'"
                >
                  <span>{{ item.name }}</span>
                  <div class="flex items-center gap-1">
                    <NTag v-if="item.isRemoved" size="tiny" type="error">
                      已移除
                    </NTag>
                    <NTag v-if="item.name.includes('iframe') || item.name.includes('IFRAME')" size="tiny" type="warning">
                      iframe
                    </NTag>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="!currentApprovedVersion?.config?.permissions?.length" class="text-gray-400 text-sm">
              无权限要求
            </div>
          </div>
        </div>

        <!-- 待审核版本 -->
        <div class="flex-1 min-w-0 bg-blue-50 -mx-4 -mt-4 p-4 border-2 border-blue-200 rounded">
          <div class="text-lg font-semibold mb-4 pb-2 border-b text-blue-600">
            待审核版本
          </div>
          <div ref="rightPanelRef" class="scroll-panel" @scroll="handleRightScroll">
            <NDescriptions bordered :column="1">
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('version')">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('version') }">版本号</span>
                </template>
                {{ versionInfo.version }}
                <span v-if="!currentApprovedVersion || versionInfo.version !== currentApprovedVersion.version" class="ml-2 text-xs bg-red-100 text-red-600 px-2 py-1 rounded">新版本</span>
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('versionDesc')">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('versionDesc') }">版本说明</span>
                </template>
                {{ getVersionDescContent(versionInfo.versionDesc) || '暂无说明' }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('packageUrl')">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('packageUrl') }">包地址</span>
                </template>
                <div class="flex items-center">
                  <a :href="versionInfo.packageUrl" target="_blank" class="text-blue-600 hover:underline">
                    {{ versionInfo.packageUrl }}
                  </a>
                  <NButton
                    size="tiny"
                    type="primary"
                    class="ml-2"
                    @click="handleDownloadByVersionId(versionInfo!.id)"
                  >
                    下载
                  </NButton>
                </div>
              </NDescriptionsItem>
              <NDescriptionsItem v-if="!onlyShowDiff || isDiff('packageHash')">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('packageHash') }">包校验值</span>
                </template>
                {{ versionInfo.packageHash }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="(!onlyShowDiff || isDiff('apiVersion')) && versionInfo.config?.apiVersion">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('apiVersion') }">API 版本</span>
                </template>
                {{ versionInfo.config.apiVersion }}
              </NDescriptionsItem>
              <NDescriptionsItem v-if="(!onlyShowDiff || isDiff('author')) && versionInfo.config?.author">
                <template #label>
                  <span :class="{ 'font-bold text-red-600': isDiff('author') }">作者</span>
                </template>
                {{ versionInfo.config.author }}
              </NDescriptionsItem>
              <NDescriptionsItem label="安全审核">
                <NButton size="small" @click="handleOpenSecurityAudit">
                  {{ versionInfo.codeSecurityAudit ? '查看报告' : '无报告' }}
                </NButton>
                <NButton
                  size="small"
                  type="warning"
                  style="margin-left: 10px;"
                  :loading="triggerAuditLoading"
                  @click="handleTriggerSecurityAudit"
                >
                  重新审核
                </NButton>
              </NDescriptionsItem>
            </NDescriptions>
          </div>
          <!-- 待审核版本权限（始终显示） -->
          <div>
            <NDivider title-placement="left" class="!my-3">
              权限（{{ versionInfo.config?.permissions?.length || 0 }}项）
              <NTag v-if="reviewHasIframe" size="small" type="warning" class="ml-2">
                含 iframe
              </NTag>
            </NDivider>
            <div v-for="group in filteredReviewPermissionGroups" :key="group.category" class="mb-3">
              <div class="text-xs text-gray-400 mb-1 uppercase tracking-wide">
                {{ group.category }}
              </div>
              <div class="space-y-1">
                <div
                  v-for="item in group.items"
                  :key="item.name"
                  class="px-3 py-1 rounded text-sm flex items-center justify-between"
                  :class="item.isNew ? 'bg-blue-100 text-blue-700' : 'bg-gray-100'"
                >
                  <span>{{ item.name }}</span>
                  <div class="flex items-center gap-1">
                    <NTag v-if="item.isNew" size="tiny" type="success">
                      新增
                    </NTag>
                    <NTag v-if="item.name.includes('iframe') || item.name.includes('IFRAME')" size="tiny" type="warning">
                      iframe
                    </NTag>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="!versionInfo.config?.permissions?.length" class="text-gray-400 text-sm">
              无权限要求
            </div>
          </div>
        </div>
      </div>

      <!-- 微应用 iframe 权限信息 -->
      <template v-if="microApp">
        <div :class="microApp.haveIframe ? 'p-3 bg-yellow-50 border border-yellow-200 rounded text-sm text-yellow-700' : 'p-3 bg-gray-50 border border-gray-200 rounded text-sm text-gray-700'">
          <strong>iframe 权限信息：</strong>
          <template v-if="microApp.haveIframe">
            当前微应用信息中标记为包含 iframe，请注意安全审核。
          </template>
          <template v-else>
            当前微应用信息中未包含 iframe。
          </template>
          <NTag size="small" :type="microApp.haveIframe ? 'warning' : 'default'" class="ml-2">
            have_iframe: {{ microApp.haveIframe ? '是' : '否' }}
          </NTag>
        </div>
      </template>

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
        <NButton @click="emit('update:visible', false)">
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

  <!-- 应用公开页面 Modal -->
  <NModal
    :show="iframeModalVisible"
    preset="card"
    style="width: 1200px; height: 800px;"
    title="应用公开页面"
    @update:show="iframeModalVisible = $event"
  >
    <div class="h-full">
      <iframe
        v-if="microApp?.id"
        :src="`/microApp/${microApp?.id}`"
        frameborder="0"
        style="width: 100%; height: 700px; border: none;"
      />
    </div>
  </NModal>

  <!-- 安全审核报告对比弹窗 -->
  <NModal
    :show="securityAuditModalVisible"
    preset="card"
    style="width: 1400px;"
    title="安全审核报告对比"
    @update:show="securityAuditModalVisible = $event"
  >
    <div class="flex gap-6">
      <!-- 当前版本安全审核 -->
      <div class="flex-1">
        <div class="text-lg font-semibold mb-4 pb-2 border-b">
          当前已发布版本
        </div>
        <template v-if="currentApprovedVersion?.codeSecurityAudit">
          <NCard class="mb-4">
            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <span>审核状态：</span>
                <NTag :type="currentApprovedVersion.codeSecurityAudit.isPassed ? 'success' : 'error'">
                  {{ currentApprovedVersion.codeSecurityAudit.isPassed ? '通过' : '未通过' }}
                </NTag>
              </div>
              <div>
                <div class="mb-2">
                  安全评分：
                </div>
                <NProgress
                  type="line"
                  :percentage="currentApprovedVersion.codeSecurityAudit.score"
                  :color="currentApprovedVersion.codeSecurityAudit.score >= 80 ? '#18a058' : currentApprovedVersion.codeSecurityAudit.score >= 60 ? '#f0a020' : '#d03050'"
                />
              </div>
              <div class="flex gap-4">
                <NTag type="error">
                  高危: {{ currentApprovedVersion.codeSecurityAudit.highRiskCount }}
                </NTag>
                <NTag type="warning">
                  中危: {{ currentApprovedVersion.codeSecurityAudit.mediumRiskCount }}
                </NTag>
                <NTag type="default">
                  低危: {{ currentApprovedVersion.codeSecurityAudit.lowRiskCount }}
                </NTag>
              </div>
              <div class="text-sm text-gray-500">
                扫描时间：{{ new Date(currentApprovedVersion.codeSecurityAudit.scanTime).toLocaleString() }}
              </div>
              <div v-if="currentApprovedVersion.codeSecurityAudit.reportUrl" class="mt-2">
                <NButton size="small" type="primary" @click="openExternalUrl(currentApprovedVersion!.codeSecurityAudit?.reportUrl!)">
                  查看完整报告
                </NButton>
              </div>
            </div>
          </NCard>

          <NCard v-if="currentApprovedVersion.codeSecurityAudit.vulnerabilities.length > 0">
            <div class="mb-3 font-semibold">
              漏洞列表（{{ currentApprovedVersion.codeSecurityAudit.vulnerabilities.length }}）
            </div>
            <NCollapse>
              <NCollapseItem
                v-for="(vuln, index) in currentApprovedVersion.codeSecurityAudit.vulnerabilities"
                :key="index"
                :name="index"
              >
                <template #header>
                  <div class="flex items-center gap-2">
                    <NTag :type="getSeverityInfo(vuln.severity).color as any">
                      {{ getSeverityInfo(vuln.severity).label }}
                    </NTag>
                    <span>{{ vuln.title }}</span>
                  </div>
                </template>
                <div class="space-y-2 text-sm">
                  <div><strong>描述：</strong>{{ vuln.description }}</div>
                  <div><strong>位置：</strong>{{ vuln.location }}:{{ vuln.lineNumber }}</div>
                  <div><strong>修复建议：</strong>{{ vuln.remediation }}</div>
                </div>
              </NCollapseItem>
            </NCollapse>
          </NCard>
        </template>
        <div v-else class="text-center py-8 text-gray-400">
          暂无安全审核报告
        </div>
      </div>

      <!-- 待审核版本安全审核 -->
      <div class="flex-1 bg-blue-50 -mx-4 -mt-4 p-4 border-2 border-blue-200 rounded">
        <div class="text-lg font-semibold mb-4 pb-2 border-b text-blue-600">
          待审核版本
        </div>
        <template v-if="versionInfo?.codeSecurityAudit">
          <NCard class="mb-4">
            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <span>审核状态：</span>
                <NTag :type="versionInfo.codeSecurityAudit.isPassed ? 'success' : 'error'">
                  {{ versionInfo.codeSecurityAudit.isPassed ? '通过' : '未通过' }}
                </NTag>
              </div>
              <div>
                <div class="mb-2">
                  安全评分：
                </div>
                <NProgress
                  type="line"
                  :percentage="versionInfo.codeSecurityAudit.score"
                  :color="versionInfo.codeSecurityAudit.score >= 80 ? '#18a058' : versionInfo.codeSecurityAudit.score >= 60 ? '#f0a020' : '#d03050'"
                />
              </div>
              <div class="flex gap-4">
                <NTag type="error">
                  高危: {{ versionInfo.codeSecurityAudit.highRiskCount }}
                </NTag>
                <NTag type="warning">
                  中危: {{ versionInfo.codeSecurityAudit.mediumRiskCount }}
                </NTag>
                <NTag type="default">
                  低危: {{ versionInfo.codeSecurityAudit.lowRiskCount }}
                </NTag>
              </div>
              <div class="text-sm text-gray-500">
                扫描时间：{{ new Date(versionInfo.codeSecurityAudit.scanTime).toLocaleString() }}
              </div>
              <div v-if="versionInfo.codeSecurityAudit.reportUrl" class="mt-2">
                <NButton size="small" type="primary" @click="openExternalUrl(versionInfo!.codeSecurityAudit?.reportUrl!)">
                  查看完整报告
                </NButton>
              </div>
            </div>
          </NCard>

          <NCard v-if="versionInfo.codeSecurityAudit.vulnerabilities.length > 0">
            <div class="mb-3 font-semibold">
              漏洞列表（{{ versionInfo.codeSecurityAudit.vulnerabilities.length }}）
            </div>
            <NCollapse>
              <NCollapseItem
                v-for="(vuln, index) in versionInfo.codeSecurityAudit.vulnerabilities"
                :key="index"
                :name="index"
              >
                <template #header>
                  <div class="flex items-center gap-2">
                    <NTag :type="getSeverityInfo(vuln.severity).color as any">
                      {{ getSeverityInfo(vuln.severity).label }}
                    </NTag>
                    <span>{{ vuln.title }}</span>
                  </div>
                </template>
                <div class="space-y-2 text-sm">
                  <div><strong>描述：</strong>{{ vuln.description }}</div>
                  <div><strong>位置：</strong>{{ vuln.location }}:{{ vuln.lineNumber }}</div>
                  <div><strong>修复建议：</strong>{{ vuln.remediation }}</div>
                </div>
              </NCollapseItem>
            </NCollapse>
          </NCard>
        </template>
        <div v-else class="text-center py-8 text-gray-400">
          暂无安全审核报告
        </div>
      </div>
    </div>
  </NModal>
</template>

<style scoped>
.scroll-panel {
  max-height: calc(100vh - 320px);
  overflow-y: auto;
  scrollbar-width: thin;
}

.scroll-panel::-webkit-scrollbar {
  width: 6px;
}

.scroll-panel::-webkit-scrollbar-thumb {
  background-color: #d1d5db;
  border-radius: 3px;
}

.scroll-panel::-webkit-scrollbar-track {
  background-color: transparent;
}
</style>
