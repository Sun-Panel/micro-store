<script lang="ts" setup>
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NImage, NImageGroup, NModal, NPopconfirm, NSelect, NSpace, NTag, NUpload, useMessage } from 'naive-ui'
import { computed, h, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { cancelAppReview, deletes, getInfo as getMicroAppInfo, updateStatus } from '@/api/admin/microApp'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'
import { cancelReview, deleteVersion, getVersionList, submitReview } from '@/api/admin/microAppVersion'
import ReviewHistoryModal from '@/components/common/ReviewHistoryModal/index.vue'
import AddVersionModal from '@/components/common/VersionManagement/AddVersionModal.vue'
import VersionDetailModal from '@/components/common/VersionManagement/VersionDetailModal.vue'
import { microAppChargeTypeMap, microAppStatusMap, MicroAppVersionStatus, microAppVersionStatusMap } from '@/enums/panel'
import { timeFormat } from '@/utils/cmn'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'
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

// 获取分类选项
async function fetchCategoryOptions() {
  try {
    const res = await getCategoryList<any>()
    categoryOptions.value = res.data?.map((item: any) => ({
      label: item.name,
      value: item.id,
    })) || []
  }
  catch (error) {
    apiRespErrMsg(error)
  }
}

// 分类名称
const categoryName = computed(() => {
  if (!microAppInfo.value)
    return ''
  const category = categoryOptions.value.find(c => c.value === microAppInfo.value?.categoryId)
  return category?.label || `ID: ${microAppInfo.value.categoryId}`
})

// 添加版本弹窗
const addVersionShow = ref(false)

// 版本详情弹窗
const versionDetailShow = ref(false)
const currentVersionDetail = ref<any>(null)

// 获取微应用详情
async function fetchMicroAppInfo() {
  loading.value = true
  try {
    const { data } = await getMicroAppInfo<any>(microAppId.value)
    microAppInfo.value = data
    // 设置默认语言
    initBaseInfoLang()
  }
  catch (error) {
    apiRespErrMsg(error)
  }
  finally {
    loading.value = false
  }
}

// ==================== 基础信息卡片多语言 ====================
const baseInfoLang = ref('zh-CN')

// 微应用的多语言列表
const baseInfoLangList = computed(() => {
  if (!microAppInfo.value)
    return ['zh-CN']
  const langList = (microAppInfo.value as any).langList || []
  if (langList.length > 0) {
    return langList.map((l: any) => l.lang)
  }
  return ['zh-CN']
})

// 微应用语言 Map（避免重复计算）
const baseInfoLangMap = computed(() => {
  const result: Record<string, any> = {}
  if (!microAppInfo.value)
    return result
  const langList = (microAppInfo.value as any).langList || []
  langList.forEach((l: any) => {
    result[l.lang] = l
  })
  return result
})

// 浏览器语言检测
function getBrowserLang(): string {
  const lang = navigator.language || (navigator as any).userLanguage || 'zh-CN'
  if (lang.startsWith('zh'))
    return 'zh-CN'
  if (lang.startsWith('en'))
    return 'en-US'
  if (lang.startsWith('ja'))
    return 'ja-JP'
  if (lang.startsWith('ko'))
    return 'ko-KR'
  return 'zh-CN'
}

// 初始化基础信息的语言
function initBaseInfoLang() {
  const browserLang = getBrowserLang()
  const langs = baseInfoLangList.value
  baseInfoLang.value = langs.includes(browserLang) ? browserLang : (langs.includes('zh-CN') ? 'zh-CN' : langs[0])
}

// 当前语言下的应用名称
const baseInfoAppName = computed(() => {
  if (!microAppInfo.value)
    return ''
  const langMap = baseInfoLangMap.value
  return langMap[baseInfoLang.value]?.appName
    || langMap['zh-CN']?.appName
    || microAppInfo.value.appName
    || ''
})

// 当前语言下的应用描述
const baseInfoAppDesc = computed(() => {
  if (!microAppInfo.value)
    return ''
  const langMap = baseInfoLangMap.value
  return langMap[baseInfoLang.value]?.appDesc
    || langMap['zh-CN']?.appDesc
    || microAppInfo.value.appDesc
    || ''
})
// ====================

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

// 版本状态颜色
function getVersionStatusType(status: number): 'default' | 'success' | 'error' {
  if (status === MicroAppVersionStatus.APPROVED)
    return 'success'
  if (status === MicroAppVersionStatus.REJECTED)
    return 'error'
  return 'default'
}

// 表格列配置
function createColumns(): DataTableColumns<MicroApp.VersionInfo> {
  return [
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
        return h(NTag, { type: getVersionStatusType(row.status), size: 'small' }, {
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
      title: '审核时间',
      key: 'reviewTime',
      width: 160,
      render(row) {
        return row.reviewTime ? timeFormat(String(row.reviewTime)) : '-'
      },
    },
    {
      title: '审核备注',
      key: 'reviewNote',
      ellipsis: { tooltip: true },
    },
    {
      title: '操作',
      key: 'actions',
      width: 280,
      render(row) {
        return h(NSpace, { size: 'small' }, {
          default: () => [
            // 查看详情
            h(NButton, { size: 'small', onClick: () => openVersionDetail(row) }, {
              default: () => '查看',
            }),
            // 草稿状态，可以提交审核
            row.status === MicroAppVersionStatus.DRAFT
              ? h(NButton, { size: 'small', type: 'primary', onClick: () => handleSubmitReview(row.id) }, {
                  default: () => '提交审核',
                })
              : null,
            // 待审核状态，可以撤销
            row.status === MicroAppVersionStatus.PENDING
              ? h(NButton, { size: 'small', type: 'warning', onClick: () => handleVersionCancelReview(row.id) }, {
                  default: () => '撤销',
                })
              : null,
            // 非通过状态，可以删除
            row.status !== MicroAppVersionStatus.APPROVED
              ? h(NPopconfirm, { onPositiveClick: () => handleDeleteVersion([row.id]) }, {
                  trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
                  default: () => `确定删除版本 ${row.version} 吗？`,
                })
              : null,
          ],
        })
      },
    },
  ]
}

const columns = createColumns()

// 多语言表格列配置
const langTableColumns: DataTableColumns<{ lang: string, appName: string, appDesc: string }> = [
  { title: '语言', key: 'lang', width: 100 },
  { title: '应用名称', key: 'appName' },
  { title: '应用描述', key: 'appDesc', ellipsis: { tooltip: true } },
]

// 提交审核
async function handleSubmitReview(versionId: number) {
  try {
    const res = await submitReview<any>({ versionId })
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

// 撤销审核
async function handleCancelReview(versionId: number) {
  try {
    const res = await cancelReview<any>({ versionId })
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

// 打开版本详情弹窗
function openVersionDetail(row: MicroApp.VersionInfo) {
  currentVersionDetail.value = row
  versionDetailShow.value = true
}

// 上架/下架
async function handleChangeStatus(status: number) {
  if (!microAppInfo.value)
    return
  try {
    const res = await updateStatus({ id: microAppInfo.value.id, status })
    if (res.code === 0) {
      message.success(status === 1 ? '已上架' : '已下架')
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
    const { code } = await cancelAppReview({ id: microAppInfo.value.id })
    if (code === 0) {
      message.success('已撤销审核')
      fetchMicroAppInfo()
    }
  }
  catch (error) {
    message.error('撤销审核失败')
  }
}

// 版本撤销审核
async function handleVersionCancelReview(versionId: number) {
  if (!versionId)
    return
  try {
    const { code } = await cancelReview({ versionId })
    if (code === 0) {
      message.success('已撤销审核')
      fetchVersionList()
    }
  }
  catch (error) {
    message.error('撤销审核失败')
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

// 预览应用
function handlePreview() {
  // TODO: 跳转到微应用公开首页
  message.info('预览功能待开发')
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
          <!-- 审核状态 -->
          <div v-if="microAppInfo?.reviewStatus && microAppInfo.reviewStatus !== 0" class="flex items-center gap-2">
            <NTag v-if="microAppInfo.reviewStatus === 1" type="warning" size="small">
              审核中
            </NTag>
            <NTag v-if="microAppInfo.reviewStatus === 2" type="success" size="small">
              已通过
            </NTag>
            <NTag v-if="microAppInfo.reviewStatus === 3" type="error" size="small">
              已拒绝
            </NTag>
            <NButton v-if="microAppInfo.reviewStatus === 1" size="tiny" @click="handleViewReviewHistory">
              查看审核内容
            </NButton>
          </div>
        </div>
        <NSpace>
          <NButton @click="handlePreview">
            查看公开页面
          </NButton>
          <NButton v-if="microAppInfo?.reviewStatus === 1" @click="handleCancelAppReview">
            撤销审核
          </NButton>
          <NButton type="primary" @click="editDialogShow = true">
            编辑信息
          </NButton>
          <NButton v-if="microAppInfo?.status === 0" type="success" @click="handleChangeStatus(1)">
            上架
          </NButton>
          <NButton v-else-if="microAppInfo?.status === 1" @click="handleChangeStatus(0)">
            下架
          </NButton>
          <NPopconfirm @positive-click="handleDelete">
            <template #trigger>
              <NButton type="error">
                删除
              </NButton>
            </template>
            确定删除该微应用吗？删除后无法恢复。
          </NPopconfirm>
        </NSpace>
      </div>
    </NCard>

    <!-- 基本信息 -->
    <NCard class="mb-[20px]" title="">
      <template #header>
        <div class="flex items-center gap-2">
          基本信息
          <!-- 语言切换 -->
          <!-- <div class="col-span-2 flex justify-end"> -->
          <NSelect
            v-model:value="baseInfoLang"
            :options="baseInfoLangList.map((lang: string) => ({ label: lang, value: lang }))"
            style="width: 140px"
            size="small"
          />
        </div>
      </template>
      <div v-if="microAppInfo" class="grid grid-cols-2 gap-4">
        <div class="flex items-center gap-4">
          <img v-if="microAppInfo.appIcon" :src="microAppInfo.appIcon" class="w-16 h-16 object-contain rounded">
          <div v-else class="w-16 h-16 bg-gray-100 rounded flex items-center justify-center text-gray-400">
            暂无图标
          </div>
          <div class="space-y-2">
            <div class="flex items-baseline gap-2">
              <span class="text-sm text-gray-500 whitespace-nowrap">MicroAppID:</span>
              <span class="font-mono text-sm text-gray-700">{{ microAppInfo.microAppId }}</span>
            </div>
            <div class="flex items-baseline gap-2">
              <span class="text-sm text-gray-500 whitespace-nowrap">微应用名称:</span>
              <span class="font-bold text-lg">{{ baseInfoAppName || microAppInfo.appName }}</span>
            </div>
          </div>
        </div>
        <div class="grid grid-cols-2 gap-2 text-sm">
          <div>
            <span class="text-gray-500">状态：</span>
            <span
              :class="{
                'text-green-500': microAppInfo.status === 1,
                'text-yellow-500': microAppInfo.status === 2,
                'text-gray-500': microAppInfo.status === 0,
              }"
            >{{ microAppStatusMap[microAppInfo.status] }}</span>
          </div>
          <div>
            <span class="text-gray-500">收费方式：</span>
            <span>{{ microAppChargeTypeMap[microAppInfo.chargeType] || '免费' }}</span>
          </div>
          <div>
            <span class="text-gray-500">分类：</span>
            <span>{{ categoryName }}</span>
          </div>
          <div>
            <span class="text-gray-500">创建时间：</span>
            <span>{{ timeFormat(String(microAppInfo.createTime)) }}</span>
          </div>
        </div>
        <div class="col-span-2">
          <span class="text-gray-500">应用描述：</span>
          <span>{{ baseInfoAppDesc || microAppInfo.appDesc || '暂无描述' }}</span>
        </div>
        <!-- 图集 -->
        <div v-if="microAppInfo.screenshots" class="col-span-2">
          <span class="text-gray-500">图集：</span>
          <NImageGroup>
            <div class="flex flex-wrap gap-2 mt-2">
              <NImage
                v-for="(screenshot, index) in microAppInfo.screenshots.split(',')"
                :key="index"
                :src="screenshot"
                width="96"
                height="96"
                style="object-fit: cover; border-radius: 6px;"
              />
            </div>
          </NImageGroup>
        </div>
      </div>
    </NCard>

    <!-- 版本列表 -->
    <NCard title="版本管理">
      <template #header-extra>
        <NButton type="primary" @click="addVersionShow = true">
          添加版本
        </NButton>
      </template>

      <NDataTable
        :columns="columns"
        :data="versionList"
        :loading="versionLoading"
        :bordered="false"
      >
        <template #empty>
          <div class="text-center py-12 text-gray-400">
            暂无版本，点击"添加版本"上传第一个版本
          </div>
        </template>
      </NDataTable>

      <!-- 无数据提示 -->
    </NCard>

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
  </div>
</template>
