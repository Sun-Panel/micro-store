<script lang="ts" setup>
import { NButton, NCard, NSpace, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'
import { getList } from '@/api/admin/microAppDeveloper'
import ReviewHistoryModal from '@/components/common/ReviewHistoryModal/index.vue'
import { microAppChargeTypeMap, microAppStatusMap } from '@/enums/panel'
import SimpleCreateMicroApp from './SimpleCreateMicroApp/index.vue'

const message = useMessage()
const router = useRouter()
const tableIsLoading = ref<boolean>(false)
const createDialogShow = ref<boolean>(false)
const reviewHistoryShow = ref<boolean>(false)
const currentAppId = ref<number>(0) // 当前查看审核历史的应用ID
const keyWord = ref<string>()
const statusFilter = ref<number | null>(null)
const categoryFilter = ref<number | null>(null)
const sortBy = ref<string>('id') // 排序字段
const sortOrder = ref<string>('desc') // 排序方式
// const dialog = useDialog()
const categoryOptions = ref<{ label: string, value: number }[]>([])

// 状态选项
// const statusOptions = [
//   { label: '全部', value: null },
//   { label: microAppStatusMap[0], value: 0 },
//   { label: microAppStatusMap[1], value: 1 },
//   { label: microAppStatusMap[2], value: 2 },
// ]

// // 排序选项
// const sortOptions = [
//   { label: '默认排序', value: 'id' },
//   { label: '下载量', value: 'download_count' },
//   { label: '安装量', value: 'install_count' },
// ]

// // 卡片列表数据
const dataList = ref<MicroApp.Info[]>([])

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

async function fetchList() {
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
  if (sortBy.value)
    req.sortBy = sortBy.value
  if (sortOrder.value)
    req.sortOrder = sortOrder.value

  try {
    const { data } = await getList<Common.ListResponse<MicroApp.Info[]>>(req)
    dataList.value = data.list || []
  }
  catch (error) {
    console.error(`获取列表失败${error}`)
  }
  finally {
    tableIsLoading.value = false
  }
}

// function handleSelect() {
//   fetchList()
// }

// // 获取审核状态标签类型
// function getReviewStatusTagType(reviewStatus: number) {
//   switch (reviewStatus) {
//     case 0:
//       return 'success' // 已通过
//     case 1:
//       return 'warning' // 审核中
//     case 2:
//       return 'error' // 已拒绝
//     case 3:
//       return 'info' // 草稿
//     default:
//       return 'default'
//   }
// }

function handleAdd() {
  createDialogShow.value = true
}

// 跳转到详情页
function handleViewDetail(item: MicroApp.Info) {
  router.push(`/admin/developerCenter/myMicroApp/detail/${item.id}`)
}

// // 打开审核历史弹窗
// function handleViewReviewHistory(item: MicroApp.Info) {
//   if (item.id) {
//     currentAppId.value = item.id
//     reviewHistoryShow.value = true
//   }
// }

function handleDone() {
  createDialogShow.value = false
  message.success('操作成功')
  fetchList()
}

onMounted(async () => {
  await fetchCategoryOptions()
  fetchList()
})
</script>

<template>
  <div>
    <!-- 搜索栏 -->
    <NCard class="mb-[20px]">
      <div class="flex">
        <!-- <NInputGroup style="max-width: 850px;">
          <NInput v-model:value="keyWord" :style="{ width: '30%' }" placeholder="请输入应用名称搜索" />
          <NSelect v-model:value="statusFilter" :options="statusOptions" :style="{ width: '100px' }" placeholder="状态" />
          <NSelect v-model:value="categoryFilter" :options="categoryOptions" :style="{ width: '120px' }" placeholder="分类" />
          <NSelect v-model:value="sortBy" :options="sortOptions" :style="{ width: '100px' }" placeholder="排序" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup> -->
        <span class="flex">
          <NButton type="primary" ghost @click="handleAdd">创建微应用</NButton>
        </span>
      </div>
    </NCard>

    <!-- 卡片列表 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 cursor-pointer">
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
            {{ item.adminName || '未命名' }}
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

          <!-- 统计数据 -->
          <div class="flex items-center justify-between text-xs text-gray-500">
            <span>下载: {{ item.downloadCount || 0 }}</span>
            <span>安装: {{ item.installCount || 0 }}</span>
          </div>

          <!-- 操作按钮 -->
          <!-- <div class="flex justify-end gap-2 pt-2" @click.stop>
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
          </div> -->
        </div>
      </NCard>
    </div>

    <!-- 无数据提示 -->
    <div v-if="dataList.length === 0 && !tableIsLoading" class="text-center py-12 text-gray-400">
      暂无微应用，点击"创建微应用"开始添加
    </div>

    <!-- 简化创建弹窗 -->
    <SimpleCreateMicroApp
      v-model:visible="createDialogShow"
      @done="handleDone"
    />

    <!-- 审核历史弹窗 -->
    <ReviewHistoryModal v-model:visible="reviewHistoryShow" :app-record-id="currentAppId" />
  </div>
</template>
