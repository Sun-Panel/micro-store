<script lang="ts" setup>
import { NButton, NCard, NDropdown, NInput, NInputGroup, NSelect, NSpace, useDialog, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { deletes, getList, updateStatus } from '@/api/admin/microApp'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'
import { checkIsDeveloper, getInfo as getDeveloperInfo } from '@/api/developer'
import { microAppChargeTypeMap, microAppStatusMap } from '@/enums/panel'
import EditLangModal from './components/EditLangModal/index.vue'
import EditMicroApp from './EditMicroApp/index.vue'

const message = useMessage()
const tableIsLoading = ref<boolean>(false)
const editDialogShow = ref<boolean>(false)
const editLangDialogShow = ref<boolean>(false)
const keyWord = ref<string>()
const statusFilter = ref<number | null>(null)
const categoryFilter = ref<number | null>(null)
const editInfo = ref<MicroApp.MicroAppInfo>()
const editLangInfo = ref<{ id: number, langMap: Record<string, { appName: string, appDesc: string }> }>()
const myDeveloperId = ref<number>(0)
const dialog = useDialog()
const categoryOptions = ref<{ label: string, value: number }[]>([])

// 状态选项
const statusOptions = [
  { label: '全部', value: null },
  { label: microAppStatusMap[0], value: 0 },
  { label: microAppStatusMap[1], value: 1 },
  { label: microAppStatusMap[2], value: 2 },
]

// 卡片列表数据
const dataList = ref<MicroApp.MicroAppInfo[]>([])

async function fetchCategoryOptions() {
  try {
    const res = await getCategoryList<any>()
    categoryOptions.value = res.data?.map((item: any) => ({
      label: item.name,
      value: item.id,
    })) || []
  }
  catch (error) {
    console.error(error)
  }
}

async function checkAndGetDeveloper() {
  try {
    const res = await checkIsDeveloper<any>()
    if (!res.data?.isDeveloper) {
      message.warning('您还不是开发者，无法管理微应用')
      return false
    }

    const devInfo = await getDeveloperInfo<any>()
    myDeveloperId.value = devInfo.data?.id
    return true
  }
  catch (error) {
    console.error(error)
    return false
  }
}

async function fetchList() {
  if (!myDeveloperId.value)
    return

  tableIsLoading.value = true
  const req: MicroApp.GetListRequest = {
    page: 1,
    limit: 100, // 卡片形式一次加载更多
  }
  if (keyWord.value)
    req.keyWord = keyWord.value
  if (statusFilter.value !== null)
    req.status = statusFilter.value
  if (categoryFilter.value !== null)
    req.categoryId = categoryFilter.value
  req.authorId = myDeveloperId.value

  try {
    const { data } = await getList<Common.ListResponse<MicroApp.MicroAppInfo[]>>(req)
    dataList.value = data.list || []
  }
  catch (error) {
    message.error('获取列表失败')
  }
  finally {
    tableIsLoading.value = false
  }
}

async function handleDelete(id: number) {
  try {
    const { code } = await deletes([id])
    if (code === 0) {
      message.success('删除成功')
      fetchList()
    }
  }
  catch (error) {
    message.error('删除失败')
  }
}

async function handleChangeStatus(row: MicroApp.MicroAppInfo, status?: number) {
  const newStatus = status ?? (row.status === 1 ? 0 : 1)
  try {
    const { code } = await updateStatus({ id: row.id, status: newStatus })
    if (code === 0) {
      message.success(newStatus === 1 ? '已上架' : '已下架')
      fetchList()
    }
  }
  catch (error) {
    message.error('操作失败')
  }
}

function handleSelect() {
  fetchList()
}

function handleAdd() {
  editInfo.value = undefined
  editDialogShow.value = true
}

// 处理下拉菜单选择
function handleDropdownSelect(key: string, item: MicroApp.MicroAppInfo) {
  if (key === 'offline')
    handleChangeStatus(item, 0)
  if (key === 'online')
    handleChangeStatus(item, 1)
  if (key === 'delete') {
    dialog.warning({
      title: '警告',
      content: `确定删除微应用"${item.appName}"吗？`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => handleDelete(item.id),
    })
  }
}

// 打开编辑语言弹窗
function handleEditLang(row: MicroApp.MicroAppInfo) {
  const langList = (row as any).langList || []
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
    // 没有多语言数据时，使用默认的 appName/appDesc
    langMap['zh-CN'] = {
      appName: row.appName || '',
      appDesc: row.appDesc || '',
    }
  }

  editLangInfo.value = {
    id: row.id,
    langMap,
  }
  editLangDialogShow.value = true
}

function handleDone() {
  editDialogShow.value = false
  message.success('操作成功')
  fetchList()
}

function handleLangDone() {
  editLangDialogShow.value = false
  message.success('语言保存成功')
  fetchList()
}

onMounted(async () => {
  await fetchCategoryOptions()
  const isDev = await checkAndGetDeveloper()
  if (isDev) {
    fetchList()
  }
})
</script>

<template>
  <div>
    <!-- 搜索栏 -->
    <NCard class="mb-[20px]">
      <div class="flex">
        <NInputGroup style="max-width: 700px;">
          <NInput v-model:value="keyWord" :style="{ width: '30%' }" placeholder="请输入应用名称搜索" />
          <NSelect v-model:value="statusFilter" :options="statusOptions" :style="{ width: '100px' }" placeholder="状态" />
          <NSelect v-model:value="categoryFilter" :options="categoryOptions" :style="{ width: '120px' }" placeholder="分类" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup>
        <span class="flex ml-auto">
          <NButton type="primary" ghost @click="handleAdd">创建微应用</NButton>
        </span>
      </div>
    </NCard>

    <!-- 卡片列表 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <NCard v-for="item in dataList" :key="item.id" hoverable>
        <template #cover>
          <div class="h-40 overflow-hidden flex items-center justify-center bg-gray-50">
            <img v-if="item.appIcon" :src="item.appIcon" class="w-20 h-20 object-contain">
            <span v-else class="text-gray-400">暂无图标</span>
          </div>
        </template>

        <div class="space-y-2">
          <!-- 应用名称 -->
          <div class="font-bold text-lg truncate">
            {{ item.appName || '未命名' }}
          </div>

          <!-- 应用ID -->
          <div class="text-xs text-gray-500 truncate">
            ID: {{ item.microAppId }}
          </div>

          <!-- 备注 -->
          <div v-if="item.remark" class="text-sm text-gray-400 truncate">
            {{ item.remark }}
          </div>

          <!-- 状态和收费 -->
          <div class="flex items-center justify-between">
            <NSpace>
              <span
                :class="{
                  'text-green-500': item.status === 1,
                  'text-yellow-500': item.status === 2,
                  'text-gray-500': item.status === 0,
                }"
              >{{ microAppStatusMap[item.status] || '未知' }}</span>
              <span class="text-gray-400">{{ microAppChargeTypeMap[item.chargeType] || '免费' }}</span>
            </NSpace>
          </div>

          <!-- 操作按钮 -->
          <div class="flex justify-end gap-2 pt-2">
            <NButton size="small" @click="handleEditLang(item)">
              语言
            </NButton>
            <NButton size="small" type="primary" quaternary @click="editInfo = item; editDialogShow = true">
              编辑
            </NButton>
            <NDropdown
              trigger="click"
              :options="[
                { label: '下架', key: 'offline' },
                { label: '上架', key: 'online' },
                { label: '删除', key: 'delete' },
              ]"
              @select="(key: string) => handleDropdownSelect(key, item)"
            >
              <NButton size="small">
                更多
              </NButton>
            </NDropdown>
          </div>
        </div>
      </NCard>
    </div>

    <!-- 无数据提示 -->
    <div v-if="dataList.length === 0 && !tableIsLoading" class="text-center py-12 text-gray-400">
      暂无微应用，点击"创建微应用"开始添加
    </div>

    <!-- 编辑弹窗 -->
    <EditMicroApp
      v-model:visible="editDialogShow"
      :micro-app-info="editInfo"
      :author-id="myDeveloperId"
      :category-options="categoryOptions"
      @done="handleDone"
    />

    <!-- 编辑语言弹窗 -->
    <EditLangModal
      v-model:visible="editLangDialogShow"
      :micro-app-id="editLangInfo?.id || 0"
      :lang-map="editLangInfo?.langMap || {}"
      @done="handleLangDone"
    />
  </div>
</template>
