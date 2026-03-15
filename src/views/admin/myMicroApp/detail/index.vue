<script lang="ts" setup>
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NImage, NImageGroup, NModal, NPopconfirm, NSelect, NSpace, NTag, NUpload, useMessage } from 'naive-ui'
import { computed, h, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { deletes, getInfo as getMicroAppInfo, updateLang, update as updateMicroApp, updateStatus } from '@/api/admin/microApp'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'
import { cancelReview, createVersion, deleteVersion, getVersionList, submitReview, uploadVersionPackage } from '@/api/admin/microAppVersion'
import { ErrorCode } from '@/enums/errorCode'
import { microAppChargeTypeMap, microAppStatusMap, MicroAppVersionStatus, microAppVersionStatusMap } from '@/enums/panel'
import { t } from '@/locales'
import { timeFormat } from '@/utils/cmn'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'
import EditLangModal from '../components/EditLangModal/index.vue'
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
const editLangDialogShow = ref(false)
const categoryOptions = ref<{ label: string, value: number }[]>([])

// 语言编辑信息
const editLangInfo = ref<{ id: number, langMap: Record<string, { appName: string, appDesc: string }> }>()

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
const addVersionLoading = ref(false)

// 添加版本表单
const versionForm = ref({
  version: '',
  versionCode: 0,
  packageUrl: '',
  packageHash: '',
  versionDesc: '',
})
const versionFile = ref<File | null>(null)

// 版本详情弹窗
const versionDetailShow = ref(false)
const versionDetailLoading = ref(false)
const currentVersionDetail = ref<any>(null)

// 打开版本详情
const versionDetailInfo = ref<{
  iconURL: string
  appName: string
  appDesc: string
  microAppId: string
  author: string
  version: string
  versionDesc: string
  permissions: string[]
  dataNodes: Record<string, any>
  networkDomains: string[]
  appInfo: Record<string, { appName: string, appDesc: string }>
} | null>(null)

// 当前选中的语言
const versionDetailLang = ref('zh-CN')

// 可用的语言列表
const versionDetailLangList = computed(() => {
  if (!versionDetailInfo.value?.appInfo)
    return []
  return Object.keys(versionDetailInfo.value.appInfo)
})

// 当前语言下的应用名称
const versionDetailCurrentAppName = computed(() => {
  if (!versionDetailInfo.value?.appInfo)
    return ''
  return versionDetailInfo.value.appInfo[versionDetailLang.value]?.appName
    || versionDetailInfo.value.appInfo['zh-CN']?.appName
    || ''
})

// 当前语言下的应用描述
const versionDetailCurrentAppDesc = computed(() => {
  if (!versionDetailInfo.value?.appInfo)
    return ''
  return versionDetailInfo.value.appInfo[versionDetailLang.value]?.appDesc
    || versionDetailInfo.value.appInfo['zh-CN']?.appDesc
    || ''
})

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
              ? h(NButton, { size: 'small', type: 'warning', onClick: () => handleCancelReview(row.id) }, {
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
async function openVersionDetail(row: MicroApp.VersionInfo) {
  currentVersionDetail.value = row
  versionDetailShow.value = true

  // 从版本配置中获取信息
  const config = row.config

  // 从微应用信息中获取多语言数据作为后备
  const microAppLangList = microAppInfo.value?.langList || []
  const microAppLangMap: Record<string, { appName: string, appDesc: string }> = {}
  microAppLangList.forEach((lang: any) => {
    microAppLangMap[lang.lang] = {
      appName: lang.appName || '',
      appDesc: lang.appDesc || '',
    }
  })

  // 如果微应用没有多语言数据，使用基本信息
  if (Object.keys(microAppLangMap).length === 0 && microAppInfo.value) {
    microAppLangMap['zh-CN'] = {
      appName: microAppInfo.value.appName || '',
      appDesc: microAppInfo.value.appDesc || '',
    }
  }

  // 转换 appInfo 格式
  const versionAppInfo = config?.appInfo || {}
  const formattedAppInfo: Record<string, { appName: string, appDesc: string }> = {}
  for (const [lang, info] of Object.entries(versionAppInfo)) {
    formattedAppInfo[lang] = {
      appName: (info as any).appName || '',
      appDesc: (info as any).description || '',
    }
  }

  versionDetailInfo.value = {
    iconURL: config?.icon || microAppInfo.value?.appIcon || '',
    appName: '',
    appDesc: '',
    microAppId: config?.microAppId || microAppInfo.value?.microAppId || '',
    author: config?.author || '',
    version: row.version || '',
    versionDesc: row.versionDesc || '',
    permissions: config?.permissions || [],
    dataNodes: {},
    networkDomains: config?.networkDomains || [],
    appInfo: Object.keys(formattedAppInfo).length > 0 ? formattedAppInfo : microAppLangMap,
  }

  // 设置默认语言
  const langs = Object.keys(versionDetailInfo.value.appInfo)
  versionDetailLang.value = langs.includes('zh-CN') ? 'zh-CN' : (langs[0] || 'zh-CN')
}

// 设置为主信息 - 同时更新语言和图标
async function handleSetAsMainInfo() {
  if (!versionDetailInfo.value || !microAppInfo.value)
    return

  versionDetailLoading.value = true

  try {
    // 准备语言信息
    const versionLangMap: Record<string, { appName: string, appDesc: string }> = {}

    // 从版本详情中获取多语言信息
    const appInfo = versionDetailInfo.value.appInfo || {}
    for (const [lang, info] of Object.entries(appInfo)) {
      versionLangMap[lang] = {
        appName: (info as any).appName || '',
        appDesc: (info as any).appDesc || '',
      }
    }

    // 如果没有多语言信息，使用默认
    if (Object.keys(versionLangMap).length === 0) {
      versionLangMap['zh-CN'] = {
        appName: versionDetailInfo.value.appName || '',
        appDesc: versionDetailInfo.value.appDesc || '',
      }
    }

    // 1. 使用专门的 updateLang 接口更新语言信息
    const langRes = await updateLang({
      id: microAppInfo.value.id,
      langMap: versionLangMap,
    } as any)

    if (langRes.code !== 0) {
      apiRespErrMsg(langRes)
      return
    }

    // 2. 更新图标（如果版本中有图标）
    const newIcon = versionDetailInfo.value.iconURL
    if (newIcon) {
      const updateRes = await updateMicroApp({
        id: microAppInfo.value.id,
        appName: microAppInfo.value.appName,
        appIcon: newIcon,
        appDesc: microAppInfo.value.appDesc,
        remark: microAppInfo.value.remark,
        categoryId: microAppInfo.value.categoryId,
        chargeType: microAppInfo.value.chargeType,
        price: microAppInfo.value.price || 0,
        screenshots: microAppInfo.value.screenshots || '',
      } as any)

      if (updateRes.code !== 0) {
        apiRespErrMsg(updateRes)
        return
      }
    }

    message.success('已设为主信息')
    versionDetailShow.value = false

    // 刷新微应用信息
    fetchMicroAppInfo()
  }
  catch (error) {
    apiRespErrMsg(error)
  }
  finally {
    versionDetailLoading.value = false
  }
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

// 打开语言编辑弹窗
function handleEditLang() {
  if (!microAppInfo.value)
    return
  const langList = (microAppInfo.value as any).langList || []
  const langMap: Record<string, { appName: string, appDesc: string }> = {}

  if (langList.length > 0) {
    langList.forEach((lang: any) => {
      langMap[lang.lang] = {
        appName: lang.appName || '',
        appDesc: lang.appDesc || '',
      }
    })
  }
  else {
    langMap['zh-CN'] = {
      appName: microAppInfo.value?.appName || '',
      appDesc: microAppInfo.value?.appDesc || '',
    }
  }

  editLangInfo.value = {
    id: microAppInfo.value.id,
    langMap,
  }
  editLangDialogShow.value = true
}

// 语言编辑完成
function handleLangDone() {
  editLangDialogShow.value = false
  message.success('语言保存成功')
  fetchMicroAppInfo()
}

// 处理编辑完成
function handleEditDone() {
  editDialogShow.value = false
  message.success('保存成功')
  fetchMicroAppInfo()
}

// 处理文件选择并上传
const uploadLoading = ref(false)
const uploadedConfig = ref<MicroApp.VersionConfig | null>(null)
async function handleUploadChange(options: { file: any }) {
  const file = options.file.file
  if (!file)
    return

  uploadLoading.value = true
  try {
    const res = await uploadVersionPackage<any>(file)
    if (res.code === 0 && res.data) {
      versionForm.value.packageUrl = res.data.url
      versionForm.value.packageHash = res.data.hash || ''
      // 保存上传的配置信息，并用返回的 IconURL 覆盖 config.icon（完整路径）
      if (res.data.config) {
        res.data.config.icon = res.data.iconURL || res.data.config.icon
      }
      uploadedConfig.value = res.data.config || null
      // 如果配置文件中有版本号，自动填充
      if (res.data.config?.version) {
        versionForm.value.version = res.data.config.version
        handleVersionInput(res.data.config.version)
      }
      message.success('上传成功')
    }
    else {
      apiRespErrMsg(res)
    }
  }
  catch (error) {
    apiRespErrMsg(error)
  }
  finally {
    uploadLoading.value = false
  }
}

// 处理版本号输入
function handleVersionInput(value: string) {
  versionForm.value.version = value
  // 简单版本号转数字：1.0.0 -> 100
  const parts = value.split('.')
  let code = 0
  if (parts.length >= 1)
    code += Number(parts[0]) * 100
  if (parts.length >= 2)
    code += Number(parts[1]) * 10
  if (parts.length >= 3)
    code += Number(parts[2])
  versionForm.value.versionCode = code
}

// 提交添加版本
async function handleAddVersion() {
  if (!versionForm.value.packageUrl || !versionForm.value.version) {
    message.warning('请上传版本包并填写版本号')
    return
  }

  addVersionLoading.value = true
  try {
    // 创建版本，传递完整的配置信息
    const createRes = await createVersion<any>({
      appId: microAppId.value,
      version: versionForm.value.version,
      versionCode: versionForm.value.versionCode,
      packageUrl: versionForm.value.packageUrl,
      packageHash: versionForm.value.packageHash || '',
      versionDesc: versionForm.value.versionDesc,
      config: uploadedConfig.value || undefined,
    })

    if (createRes.code === 0) {
      // 自动提交审核
      if (createRes.data?.id) {
        const reviewRes = await submitReview<any>({ versionId: createRes.data.id })
        if (reviewRes.code !== 0) {
          apiRespErrMsg(reviewRes)
          return
        }
      }

      message.success('版本添加成功，已提交审核')
      addVersionShow.value = false
      versionForm.value = { version: '', versionCode: 0, packageUrl: '', packageHash: '', versionDesc: '' }
      versionFile.value = null
      uploadedConfig.value = null
      fetchVersionList()
    }
    else {
      apiRespErrMsg(createRes)
    }
  }
  catch (error) {
    apiRespErrMsg(error)
  }
  finally {
    addVersionLoading.value = false
  }
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
        </div>
        <NSpace>
          <NButton @click="handlePreview">
            查看公开页面
          </NButton>
          <NButton @click="handleEditLang">
            语言设置
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

    <!-- 语言编辑弹窗 -->
    <EditLangModal
      v-model:visible="editLangDialogShow"
      :micro-app-id="editLangInfo?.id || 0"
      :lang-map="editLangInfo?.langMap || {}"
      @done="handleLangDone"
    />

    <!-- 添加版本弹窗 -->
    <NModal v-model:show="addVersionShow" preset="card" style="width: 500px" title="添加版本">
      <div class="space-y-4">
        <div>
          <div class="mb-2">
            选择版本包 <span class="text-red-500">*</span>
          </div>
          <NUpload
            accept=".zip"
            :max="1"
            :custom-request="(options: any) => handleUploadChange({ file: options.file })"
            :show-file-list="false"
          >
            <NButton :loading="uploadLoading">
              {{ uploadLoading ? '上传中...' : '选择文件' }}
            </NButton>
          </NUpload>
          <div v-if="uploadedConfig !== null" class="text-xs mt-1">
            <div v-if="uploadedConfig.icon" class="mb-2">
              <img :src="uploadedConfig.icon" class="w-12 h-12 object-contain border rounded">
            </div>
            <div v-if="uploadedConfig.appInfo?.['zh-CN']?.appName || uploadedConfig.appInfo?.['en-US']?.appName" class="text-lg font-bold">
              应用名称：{{ uploadedConfig.appInfo?.['zh-CN']?.appName || uploadedConfig.appInfo?.['en-US']?.appName }}
            </div>
            <div v-if="uploadedConfig.microAppId" class="text-gray-500">
              应用ID：{{ uploadedConfig.microAppId }}
            </div>
            <div v-if="uploadedConfig.author" class="text-gray-500">
              作者：{{ uploadedConfig.author }}
            </div>
            <div v-if="uploadedConfig.version" class="text-blue-500">
              版本号：{{ uploadedConfig.version }}
            </div>
            <div v-else class="text-orange-500">
              未检测到版本号，请手动填写
            </div>
          </div>
          <div class="text-xs text-gray-400 mt-1">
            支持 .zip 格式的微应用包
          </div>
        </div>
        <div>
          <div class="mb-2">
            版本号 <span class="text-red-500">*</span>
          </div>
          <NInput v-model:value="versionForm.version" placeholder="如：1.0.0" @update:value="handleVersionInput" />
        </div>
        <div>
          <div class="mb-2">
            版本说明
          </div>
          <NInput v-model:value="versionForm.versionDesc" type="textarea" placeholder="请输入版本说明" :rows="3" />
        </div>
      </div>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="addVersionShow = false">
            取消
          </NButton>
          <NButton type="primary" :loading="addVersionLoading" @click="handleAddVersion">
            添加
          </NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- 版本详情弹窗 -->
    <NModal v-model:show="versionDetailShow" preset="card" style="width: 600px" title="版本详情">
      <div v-if="versionDetailInfo" class="space-y-4">
        <!-- 语言切换 -->
        <div v-if="versionDetailLangList.length > 0" class="flex justify-end">
          <NSelect
            v-model:value="versionDetailLang"
            :options="versionDetailLangList.map(lang => ({ label: lang, value: lang }))"
            style="width: 120px"
            size="small"
          />
        </div>

        <div class="flex items-start gap-4">
          <img v-if="versionDetailInfo.iconURL" :src="versionDetailInfo.iconURL" class="w-20 h-20 object-contain border rounded">
          <div v-else class="w-20 h-20 bg-gray-100 border rounded flex items-center justify-center text-gray-400">
            暂无图标
          </div>
          <div class="flex-1">
            <div class="text-lg font-bold">
              {{ versionDetailCurrentAppName || '未命名' }}
            </div>
            <div v-if="versionDetailCurrentAppDesc" class="text-sm text-gray-500 mt-1">
              {{ versionDetailCurrentAppDesc }}
            </div>
            <div class="text-sm text-gray-500 mt-1">
              ID: {{ versionDetailInfo.microAppId }}
            </div>
            <div class="text-sm text-gray-500">
              作者: {{ versionDetailInfo.author }}
            </div>
            <div class="text-sm text-gray-500">
              版本: {{ versionDetailInfo.version }}
            </div>
          </div>
        </div>

        <!-- 多语言列表 -->
        <div v-if="versionDetailLangList.length > 0" class="text-sm">
          <div class="text-gray-500 mb-2">
            多语言信息：
          </div>
          <NDataTable
            :columns="langTableColumns"
            :data="versionDetailLangList.map(lang => ({
              lang,
              appName: versionDetailInfo?.appInfo?.[lang]?.appName || '-',
              appDesc: versionDetailInfo?.appInfo?.[lang]?.appDesc || '-',
            }))"
            :bordered="true"
            size="small"
          />
        </div>

        <div v-if="versionDetailInfo.versionDesc" class="text-sm">
          <div class="text-gray-500">
            版本说明：
          </div>
          <div>{{ versionDetailInfo.versionDesc }}</div>
        </div>

        <div v-if="versionDetailInfo.permissions?.length" class="text-sm">
          <div class="text-gray-500">
            权限：
          </div>
          <div class="flex flex-wrap gap-1 mt-1">
            <NTag v-for="p in versionDetailInfo.permissions" :key="p" size="small" type="info">
              {{ p }}
            </NTag>
          </div>
        </div>

        <div v-if="versionDetailInfo.networkDomains?.length" class="text-sm">
          <div class="text-gray-500">
            网络域名白名单：
          </div>
          <div class="flex flex-wrap gap-1 mt-1">
            <NTag v-for="d in versionDetailInfo.networkDomains" :key="d" size="small">
              {{ d }}
            </NTag>
          </div>
        </div>
      </div>
      <div v-else class="text-center py-8 text-gray-500">
        暂无详情信息
      </div>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="versionDetailShow = false">
            关闭
          </NButton>
          <NButton type="primary" :loading="versionDetailLoading" @click="handleSetAsMainInfo">
            设为微应用主信息
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>
