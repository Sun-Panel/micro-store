<script lang="ts" setup>
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import { NButton, NCard, NDataTable, NDropdown, NInput, NInputGroup, NSelect, useDialog, useMessage } from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { deletes, getList, updateStatus } from '@/api/admin/microApp'
import { getEnabledList as getCategoryList } from '@/api/admin/microAppCategory'
import { SvgIcon } from '@/components/common'
import { microAppChargeTypeMap, microAppStatusMap } from '@/enums/panel'
import { timeFormat } from '@/utils/cmn'

const message = useMessage()
const router = useRouter()
const tableIsLoading = ref<boolean>(false)
const keyWord = ref<string>()
const statusFilter = ref<number | null>(null)
const categoryFilter = ref<number | null>(null)
const dialog = useDialog()
const categoryOptions = ref<{ label: string, value: number }[]>([])

// 状态选项
const statusOptions = [
  { label: '全部', value: null },
  { label: microAppStatusMap[0], value: 0 },
  { label: microAppStatusMap[1], value: 1 },
  { label: microAppStatusMap[2], value: 2 },
]

// 表格列配置
function createColumns({
  handleStatus,
  handleView,
}: {
  handleStatus: (row: MicroApp.MicroAppInfo, status: number) => void
  handleView: (row: MicroApp.MicroAppInfo) => void
}): DataTableColumns<MicroApp.MicroAppInfo> {
  return [
    {
      title: 'ID',
      key: 'id',
      width: 60,
    },
    {
      title: '图标',
      key: 'appIcon',
      width: 60,
      render(row) {
        return h('img', {
          src: row.appIcon,
          style: { width: '32px', height: '32px', objectFit: 'cover' },
        })
      },
    },
    {
      title: '应用名称',
      key: 'appName',
    },
    {
      title: '应用ID',
      key: 'microAppId',
      width: 150,
    },
    {
      title: '收费',
      key: 'chargeType',
      width: 80,
      render(row) {
        return microAppChargeTypeMap[row.chargeType] || '免费'
      },
    },
    {
      title: '状态',
      key: 'status',
      width: 90,
      render(row) {
        return h('span', {
          class: row.status === 1 ? 'text-green-500' : row.status === 2 ? 'text-yellow-500' : 'text-gray-500',
        }, microAppStatusMap[row.status] || '未知')
      },
    },
    {
      title: '创建时间',
      key: 'createTime',
      width: 160,
      render(row) {
        return timeFormat(String(row.createTime))
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 120,
      render(row) {
        const btn = h(NButton, { strong: true, tertiary: true, size: 'small' }, {
          default: () => h(SvgIcon, { icon: 'mingcute:more-1-fill' }),
        })

        return h(NDropdown, {
          trigger: 'click',
          onSelect(key: string | number) {
            switch (key) {
              case 'view':
                handleView(row)
                break
              case 'pass':
                handleStatus(row, 1)
                break
              case 'reject':
                handleStatus(row, 0)
                break
              case 'delete':
                dialog.warning({
                  title: '警告',
                  content: `确定删除微应用"${row.appName}"吗？`,
                  positiveText: '确定',
                  negativeText: '取消',
                  onPositiveClick: () => handleDelete([row.id]),
                })
                break
            }
          },
          options: [
            { label: '查看详情', key: 'view' },
            { label: '审核通过', key: 'pass' },
            { label: '审核拒绝', key: 'reject' },
            { label: '下架', key: 'offline' },
            { label: '删除', key: 'delete' },
          ],
        }, { default: () => btn })
      },
    },
  ]
}

const dataList = ref<MicroApp.MicroAppInfo[]>([])

const columns = createColumns({
  handleStatus(row: MicroApp.MicroAppInfo, status: number) {
    handleChangeStatus(row, status)
  },
  handleView(row: MicroApp.MicroAppInfo) {
    router.push(`/admin/microAppManage/detail/${row.id}`)
  },
})

// 分页配置
const pagination = reactive({
  page: 1,
  showSizePicker: true,
  pageSizes: [10, 30, 50, 100],
  pageSize: 10,
  itemCount: 0,
  onChange: (page: number) => {
    pagination.page = page
    fetchList()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    fetchList()
  },
  prefix(item: PaginationProps) {
    return `共 ${item.itemCount} 个微应用`
  },
})

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
    console.error(error)
  }
}

// 获取列表
async function fetchList() {
  tableIsLoading.value = true
  const req: MicroApp.GetListRequest = {
    page: pagination.page,
    limit: pagination.pageSize,
  }
  if (keyWord.value)
    req.keyWord = keyWord.value
  if (statusFilter.value !== null)
    req.status = statusFilter.value
  if (categoryFilter.value !== null)
    req.categoryId = categoryFilter.value

  try {
    const { data } = await getList<Common.ListResponse<MicroApp.MicroAppInfo[]>>(req)
    pagination.itemCount = data.count
    dataList.value = data.list || []
  }
  catch (error) {
    message.error('获取列表失败')
  }
  finally {
    tableIsLoading.value = false
  }
}

// 删除
async function handleDelete(ids: number[]) {
  try {
    const { code } = await deletes(ids)
    if (code === 0) {
      message.success('删除成功')
      fetchList()
    }
  }
  catch (error) {
    message.error('删除失败')
  }
}

// 修改状态
async function handleChangeStatus(row: MicroApp.MicroAppInfo, status: number) {
  try {
    const { code } = await updateStatus({ id: row.id, status })
    if (code === 0) {
      message.success(status === 1 ? '已通过审核' : '已拒绝/下架')
      fetchList()
    }
  }
  catch (error) {
    message.error('操作失败')
  }
}

function handleSelect() {
  pagination.page = 1
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
        <NInputGroup style="max-width: 700px;">
          <NInput v-model:value="keyWord" :style="{ width: '30%' }" placeholder="请输入应用名称搜索" />
          <NSelect v-model:value="statusFilter" :options="statusOptions" :style="{ width: '100px' }" placeholder="状态" />
          <NSelect v-model:value="categoryFilter" :options="categoryOptions" :style="{ width: '120px' }" placeholder="分类" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup>
      </div>
    </NCard>

    <!-- 表格 -->
    <NDataTable
      :columns="columns"
      :data="dataList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
    />
  </div>
</template>
