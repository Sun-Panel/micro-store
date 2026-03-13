<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NDropdown, NInput, NInputGroup, NSelect, NSwitch, useDialog, useMessage } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import EditCategory from './EditCategory/index.vue'
import { deletes, getList, updateStatus } from '@/api/admin/microAppCategory'
import { timeFormat } from '@/utils/cmn'
import { SvgIcon } from '@/components/common'

const message = useMessage()
const tableIsLoading = ref<boolean>(false)
const editCategoryDialogShow = ref<boolean>(false)
const keyWord = ref<string>()
const editCategoryInfo = ref<MicroAppCategory.CategoryInfo>()
const statusFilter = ref<number | null>(null)
const dialog = useDialog()

const statusOptions = [
  { label: '全部', value: null },
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 },
]

const createColumns = ({
  update,
  changeStatus,
}: {
  update: (row: MicroAppCategory.CategoryInfo) => void
  changeStatus: (row: MicroAppCategory.CategoryInfo) => void
}): DataTableColumns<MicroAppCategory.CategoryInfo> => {
  return [
    {
      title: 'ID',
      key: 'id',
      width: 80,
    },
    {
      title: '分类名称',
      key: 'name',
    },
    {
      title: '图标',
      key: 'icon',
      render(row) {
        return row.icon || '-'
      },
    },
    {
      title: '排序',
      key: 'sort',
      width: 80,
    },
    {
      title: '状态',
      key: 'status',
      width: 100,
      render(row) {
        return h(NSwitch, {
          value: row.status === 1,
          onUpdateValue: () => changeStatus(row),
        })
      },
    },
    {
      title: '创建时间',
      key: 'createTime',
      render(row) {
        return timeFormat(String(row.createTime))
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 100,
      render(row) {
        const btn = h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: 'small',
          },
          {
            default() {
              return h(SvgIcon, { icon: 'mingcute:more-1-fill' })
            },
          },
        )

        return h(NDropdown, {
          trigger: 'click',
          onSelect(key: string | number) {
            switch (key) {
              case 'update':
                update(row)
                break
              case 'delete':
                dialog.warning({
                  title: '警告',
                  content: `确定删除分类"${row.name}"吗？`,
                  positiveText: '确定',
                  negativeText: '取消',
                  onPositiveClick: () => {
                    handleDelete([row.id])
                  },
                })
                break
            }
          },
          options: [
            { label: '编辑', key: 'update' },
            { label: '删除', key: 'delete' },
          ],
        }, { default: () => btn })
      },
    },
  ]
}

const categoryList = ref<MicroAppCategory.CategoryInfo[]>([])

const columns = createColumns({
  update(row: MicroAppCategory.CategoryInfo) {
    editCategoryInfo.value = row
    editCategoryDialogShow.value = true
  },
  changeStatus(row: MicroAppCategory.CategoryInfo) {
    handleChangeStatus(row)
  },
})

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
    return `共 ${item.itemCount} 个分类`
  },
})

function handlePageChange(page: number) {
  fetchList()
}

function handleSelect() {
  pagination.page = 1
  fetchList()
}

function handleAdd() {
  editCategoryInfo.value = undefined
  editCategoryDialogShow.value = true
}

function handleDone() {
  editCategoryDialogShow.value = false
  message.success('操作成功')
  fetchList()
}

async function fetchList() {
  tableIsLoading.value = true
  const req: MicroAppCategory.GetListRequest = {
    page: pagination.page,
    limit: pagination.pageSize,
  }
  if (keyWord.value)
    req.keyWord = keyWord.value
  if (statusFilter.value !== null)
    req.status = statusFilter.value

  try {
    const { data } = await getList<Common.ListResponse<MicroAppCategory.CategoryInfo[]>>(req)
    pagination.itemCount = data.count
    categoryList.value = data.list || []
  }
  catch (error) {
    message.error('获取列表失败')
  }
  finally {
    tableIsLoading.value = false
  }
}

async function handleDelete(ids: number[]) {
  try {
    const { code } = await deletes({ ids })
    if (code === 0) {
      message.success('删除成功')
      fetchList()
    }
  }
  catch (error) {
    message.error('删除失败')
  }
}

async function handleChangeStatus(row: MicroAppCategory.CategoryInfo) {
  const newStatus = row.status === 1 ? 0 : 1
  try {
    const { code } = await updateStatus({ id: row.id, status: newStatus })
    if (code === 0) {
      message.success('状态更新成功')
      fetchList()
    }
  }
  catch (error) {
    message.error('状态更新失败')
  }
}

onMounted(() => {
  fetchList()
})
</script>

<template>
  <div>
    <NCard class="mb-[20px]">
      <div class="flex">
        <NInputGroup style="max-width: 700px;">
          <NInput v-model:value="keyWord" :style="{ width: '40%' }" placeholder="请输入分类名称搜索" />
          <NSelect
            v-model:value="statusFilter"
            :options="statusOptions"
            :style="{ width: '120px' }"
            placeholder="状态"
          />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup>
        <span class="flex ml-auto">
          <NButton type="primary" ghost @click="handleAdd">
            添加分类
          </NButton>
        </span>
      </div>
    </NCard>

    <NDataTable
      :columns="columns"
      :data="categoryList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
      @update:page="handlePageChange"
    />
    <EditCategory v-model:visible="editCategoryDialogShow" :category-info="editCategoryInfo" @done="handleDone" />
  </div>
</template>
