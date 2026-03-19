<script lang="ts" setup>
import { NButton, NCard, NDropdown, NInput, NInputGroup, NSelect, NSpace, NTag, useDialog, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { cancelReview, deletes, getList, offline } from '@/api/admin/microAppDeveloper'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'
import { checkIsDeveloper, getInfo as getDeveloperInfo } from '@/api/developer'
import ReviewHistoryModal from '@/components/common/ReviewHistoryModal/index.vue'
import { microAppChargeTypeMap, microAppStatusMap } from '@/enums/panel'
import EditMicroApp from './EditMicroApp/index.vue'

const message = useMessage()
const router = useRouter()
const tableIsLoading = ref<boolean>(false)
const editDialogShow = ref<boolean>(false)
const reviewHistoryShow = ref<boolean>(false)
const currentAppId = ref<number>(0) // 当前查看审核历史的应用ID
const keyWord = ref<string>()
const statusFilter = ref<number | null>(null)
const categoryFilter = ref<number | null>(null)
const editInfo = ref<MicroApp.MicroAppInfo>()
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
  // 开发者只能下架自己的应用，不能上架（上架需要审核通过）
  if (newStatus === 1) {
    message.warning('应用需要审核通过后才能上架')
    return
  }

  try {
    const { code } = await offline({ id: row.id, type: 1, reason: '作者主动下架' })
    if (code === 0) {
      message.success('已下架')
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

// 跳转到详情页
function handleViewDetail(item: MicroApp.MicroAppInfo) {
  router.push(`/admin/myMicroApp/detail/${item.id}`)
}

// 打开审核历史弹窗
function handleViewReviewHistory(item: MicroApp.MicroAppInfo) {
  currentAppId.value = item.id
  reviewHistoryShow.value = true
}

// 撤销审核
async function handleCancelReview(item: MicroApp.MicroAppInfo) {
  try {
    const { code } = await cancelAppReview({ id: item.id })
    if (code === 0) {
      message.success('已撤销审核')
      fetchList()
    }
  }
  catch (error) {
    message.error('撤销审核失败')
  }
}

function handleDone() {
  editDialogShow.value = false
  message.success('操作成功')
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
      <NCard v-for="item in dataList" :key="item.id" hoverable @click="handleViewDetail(item)">
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

          <!-- 审核状态和撤销按钮 -->
          <div v-if="item.reviewStatus !== undefined && item.reviewStatus !== 0" class="flex items-center gap-2 mt-2">
            <NTag v-if="item.reviewStatus === 1" type="warning" size="small">
              审核中
            </NTag>
            <NTag v-if="item.reviewStatus === 2" type="success" size="small">
              已通过
            </NTag>
            <NTag v-if="item.reviewStatus === 3" type="error" size="small">
              已拒绝
            </NTag>
            <NButton v-if="item.reviewStatus === 1" text size="small" @click="handleViewReviewHistory(item)">
              查看审核内容
            </NButton>
            <NButton v-if="item.reviewStatus === 1" text size="small" @click="handleCancelReview(item)">
              撤销审核
            </NButton>
          </div>

          <!-- 操作按钮 -->
          <div class="flex justify-end gap-2 pt-2" @click.stop>
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

    <!-- 审核历史弹窗 -->
    <ReviewHistoryModal v-model:visible="reviewHistoryShow" :app-id="currentAppId" />
  </div>
</template>
