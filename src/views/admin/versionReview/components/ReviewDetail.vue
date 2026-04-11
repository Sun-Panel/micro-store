<script lang="ts" setup>
import { NButton, NCard, NCollapse, NCollapseItem, NDescriptions, NDescriptionsItem, NDivider, NInput, NModal, NProgress, NSpace, NTag, useMessage } from 'naive-ui'
import { ref, watch } from 'vue'
import { getDownloadUrl, getLatestOnlineByAppModelId } from '@/api/admin/microAppVersion'
import { review } from '@/api/admin/microAppVersionReview'
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

// 数据
const reviewLoading = ref(false)
const currentApprovedVersion = ref<MicroApp.VersionInfo>()
const reviewForm = ref({
  status: 1,
  reviewNote: '',
})
const iframeModalVisible = ref(false)
const securityAuditModalVisible = ref(false)

// 严重程度映射
const severityMap: Record<string, { label: string, color: 'error' | 'warning' | 'default' | 'success', value: number }> = {
  CRITICAL: { label: '高危', color: 'error', value: 4 },
  HIGH: { label: '高', color: 'warning', value: 3 },
  MEDIUM: { label: '中', color: 'default', value: 2 },
  LOW: { label: '低', color: 'success', value: 1 },
}

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

  // 驳回时必须填写原因
  if (reviewForm.value.status === 2 && !reviewForm.value.reviewNote?.trim()) {
    message.error('驳回时必须填写驳回原因')
    return
  }

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
// function handleDownload(url: string) {
//   window.open(url, '_blank')
// }

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
      <!-- 对比展示 -->
      <div class="flex gap-6">
        <!-- 当前已发布版本 -->
        <div class="flex-1">
          <div class="text-lg font-semibold mb-4 pb-2 border-b">
            当前已发布版本
            <span v-if="!currentApprovedVersion" class="text-sm font-normal text-gray-400">（暂无）</span>
          </div>
          <NDescriptions v-if="currentApprovedVersion" bordered :column="1">
            <NDescriptionsItem label="版本号">
              {{ currentApprovedVersion.version }}
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
            <NDescriptionsItem v-if="currentApprovedVersion.config?.apiVersion" label="API 版本">
              {{ currentApprovedVersion.config.apiVersion }}
            </NDescriptionsItem>
            <NDescriptionsItem v-if="currentApprovedVersion.config?.author" label="作者">
              {{ currentApprovedVersion.config.author }}
            </NDescriptionsItem>
            <NDescriptionsItem label="安全审核">
              <NButton size="small" @click="handleOpenSecurityAudit">
                {{ currentApprovedVersion.codeSecurityAudit ? '查看报告' : '无报告' }}
              </NButton>
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
            <NDescriptionsItem label="版本号" :class="{ 'font-bold text-red-600': !currentApprovedVersion || versionInfo.version !== currentApprovedVersion.version }">
              {{ versionInfo.version }}
              <span v-if="!currentApprovedVersion || versionInfo.version !== currentApprovedVersion.version" class="ml-2 text-xs bg-red-100 text-red-600 px-2 py-1 rounded">新版本</span>
            </NDescriptionsItem>
            <NDescriptionsItem label="版本说明">
              {{ versionInfo.versionDesc || '暂无说明' }}
            </NDescriptionsItem>
            <NDescriptionsItem label="包地址">
              <a :href="versionInfo.packageUrl" target="_blank" class="text-blue-600 hover:underline">
                {{ versionInfo.packageUrl }}
              </a>
              <NButton
                size="tiny"
                type="primary"
                class="ml-2"
                @click="handleDownloadByVersionId(versionInfo.id)"
              >
                下载
              </NButton>
            </NDescriptionsItem>
            <NDescriptionsItem label="包校验值">
              {{ versionInfo.packageHash }}
            </NDescriptionsItem>
            <NDescriptionsItem v-if="versionInfo.config?.apiVersion" label="API 版本">
              {{ versionInfo.config.apiVersion }}
            </NDescriptionsItem>
            <NDescriptionsItem v-if="versionInfo.config?.author" label="作者">
              {{ versionInfo.config.author }}
            </NDescriptionsItem>
            <NDescriptionsItem label="安全审核">
              <NButton size="small" @click="handleOpenSecurityAudit">
                {{ versionInfo.codeSecurityAudit ? '查看报告' : '无报告' }}
              </NButton>
            </NDescriptionsItem>
          </NDescriptions>
        </div>
      </div>

      <!-- 权限对比 -->
      <NDivider v-if="versionInfo.config?.permissions?.length || (currentApprovedVersion?.config?.permissions?.length)" title-placement="left">
        权限对比
      </NDivider>
      <div v-if="versionInfo.config?.permissions?.length || (currentApprovedVersion?.config?.permissions?.length)" class="flex gap-6">
        <div class="flex-1">
          <div class="text-sm text-gray-500 mb-2">
            当前版本权限
          </div>
          <template v-if="currentApprovedVersion?.config?.permissions?.length">
            <div class="space-y-1">
              <div v-for="(perm, index) in currentApprovedVersion.config.permissions" :key="index" class="px-3 py-1 bg-gray-100 rounded text-sm">
                {{ perm }}
              </div>
            </div>
          </template>
          <div v-else class="text-gray-400 text-sm">
            无权限要求
          </div>
        </div>
        <div class="flex-1">
          <div class="text-sm text-blue-600 mb-2">
            待审核版本权限
          </div>
          <div class="space-y-1">
            <div v-for="(perm, index) in (versionInfo.config?.permissions || [])" :key="index" class="px-3 py-1 bg-blue-100 text-blue-700 rounded text-sm">
              {{ perm }}
              <span v-if="!currentApprovedVersion?.config?.permissions?.includes(perm)" class="ml-2 text-xs bg-red-100 text-red-600 px-2 py-1 rounded">新增</span>
            </div>
            <div v-if="!versionInfo.config?.permissions?.length" class="text-gray-400 text-sm">
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
          <!-- 概览卡片 -->
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
                <div class="flex items-center gap-2">
                  <NTag type="error">
                    高危: {{ currentApprovedVersion.codeSecurityAudit.highRiskCount }}
                  </NTag>
                </div>
                <div class="flex items-center gap-2">
                  <NTag type="warning">
                    中危: {{ currentApprovedVersion.codeSecurityAudit.mediumRiskCount }}
                  </NTag>
                </div>
                <div class="flex items-center gap-2">
                  <NTag type="default">
                    低危: {{ currentApprovedVersion.codeSecurityAudit.lowRiskCount }}
                  </NTag>
                </div>
              </div>
              <div class="text-sm text-gray-500">
                扫描时间：{{ new Date(currentApprovedVersion.codeSecurityAudit.scanTime).toLocaleString() }}
              </div>
              <div v-if="currentApprovedVersion.codeSecurityAudit.reportUrl" class="mt-2">
                <NButton size="small" type="primary" @click="openExternalUrl(currentApprovedVersion.codeSecurityAudit.reportUrl!)">
                  查看完整报告
                </NButton>
              </div>
            </div>
          </NCard>

          <!-- 漏洞列表 -->
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
                    <NTag :type="severityMap[vuln.severity].color as any">
                      {{ severityMap[vuln.severity].label }}
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
          <!-- 概览卡片 -->
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
                <div class="flex items-center gap-2">
                  <NTag type="error">
                    高危: {{ versionInfo.codeSecurityAudit.highRiskCount }}
                  </NTag>
                </div>
                <div class="flex items-center gap-2">
                  <NTag type="warning">
                    中危: {{ versionInfo.codeSecurityAudit.mediumRiskCount }}
                  </NTag>
                </div>
                <div class="flex items-center gap-2">
                  <NTag type="default">
                    低危: {{ versionInfo.codeSecurityAudit.lowRiskCount }}
                  </NTag>
                </div>
              </div>
              <div class="text-sm text-gray-500">
                扫描时间：{{ new Date(versionInfo.codeSecurityAudit.scanTime).toLocaleString() }}
              </div>
              <div v-if="versionInfo.codeSecurityAudit.reportUrl" class="mt-2">
                <NButton size="small" type="primary" @click="openExternalUrl(versionInfo.codeSecurityAudit.reportUrl!)">
                  查看完整报告
                </NButton>
              </div>
            </div>
          </NCard>

          <!-- 漏洞列表 -->
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
                    <NTag :type="severityMap[vuln.severity].color as any">
                      {{ severityMap[vuln.severity].label }}
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
