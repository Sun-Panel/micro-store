<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NH3, NInput, NInputGroup, NTag, useDialog, useMessage } from 'naive-ui'
import type { PaginationProps, DataTableColumns } from 'naive-ui'

import { deleteById as deleteByIdApi, getMyList as getMyListApi, offline as offlineApi, withdraw as withdrawApi } from '@/api/admin/customCode'
import { router } from '@/router'

const tableIsLoading = ref<boolean>(false)
const keyWord = ref<string>()
const list = ref<CustomCode.ListItem[]>([])

const dialog = useDialog()
const message = useMessage()

function mainStatusTag(status: number) {
  switch (status) {
    case 1:
      return { text: '已上线', type: 'success' as const }
    case 0:
      return { text: '待审核', type: 'warning' as const }
    case -1:
      return { text: '草稿', type: 'default' as const }
    case 2:
      return { text: '已拒绝', type: 'error' as const }
    case 3:
      return { text: '已下架', type: 'default' as const }
    default:
      return { text: '未知', type: 'default' as const }
  }
}

function codeTypeTexts(codeTypes?: number[]) {
  if (!codeTypes?.length)
    return '-'
  const map: Record<number, string> = { 1: 'JS', 2: 'CSS', 3: '页脚' }
  return codeTypes.map(item => map[item] ?? item).join(' / ')
}

const columns: DataTableColumns<CustomCode.ListItem> = [
  { title: '标题', key: 'title' },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  {
    title: '代码类型',
    key: 'codeTypes',
    render(row) {
      return codeTypeTexts(row.codeTypes)
    },
  },
  {
    title: '状态',
    key: 'status',
    render(row) {
      const tag = mainStatusTag(row.status)
      return h('div', { style: 'display:flex;gap:6px;align-items:center' }, [
        h(NTag, { type: tag.type, size: 'small' }, { default: () => tag.text }),
        row.reviewStatus === 0 ? h(NTag, { type: 'warning', size: 'small' }, { default: () => '修改审核中' }) : null,
      ])
    },
  },
  { title: '阅读量', key: 'readCount', width: 90 },
  {
    title: '',
    key: 'actions',
    render(row) {
      const buttons = [
        h(
          NButton,
          {
            size: 'tiny',
            type: 'info',
            onClick: () => router.push({ name: 'AdminCustomCodeEdit', query: { id: row.id } }),
          },
          '编辑',
        ),
      ]

      if (row.reviewStatus === 0) {
        buttons.push(
          h(
            NButton,
            {
              size: 'tiny',
              style: { marginLeft: '5px' },
              onClick: () => handleWithdraw(row),
            },
            '撤回审核',
          ),
        )
      }

      if (row.status === 1) {
        buttons.push(
          h(
            NButton,
            {
              size: 'tiny',
              style: { marginLeft: '5px' },
              onClick: () => handleOffline(row),
            },
            '下架',
          ),
        )
      }

      if (row.status !== 1) {
        buttons.push(
          h(
            NButton,
            {
              size: 'tiny',
              type: 'error',
              style: { marginLeft: '5px' },
              onClick: () => handleDelete(row),
            },
            '删除',
          ),
        )
      }

      return buttons
    },
  },
]

const pagination = reactive({
  page: 1,
  showSizePicker: true,
  pageSizes: [10, 30, 50],
  pageSize: 10,
  itemCount: 0,
  onChange: (page: number) => {
    pagination.page = page
    getList(null)
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    getList(null)
  },
  prefix(item: PaginationProps) {
    return `共 ${item.itemCount} 条`
  },
})

function handlePageChange(page: number) {
  getList(page)
}

function handleSelect() {
  getList(null)
}

function handleAdd() {
  router.push({ name: 'AdminCustomCodeEdit' })
}

async function handleWithdraw(row: CustomCode.ListItem) {
  dialog.warning({
    title: '撤回审核',
    content: '撤回后可以重新编辑并提交审核，已上线的版本不受影响。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      const { code, msg } = await withdrawApi<Common.DataResponse>(row.id)
      if (code === 0) {
        message.success('已撤回')
        getList(null)
      }
      else {
        message.warning(`操作失败: ${msg}`)
      }
    },
  })
}

function handleOffline(row: CustomCode.ListItem) {
  dialog.warning({
    title: '下架',
    content: '下架后前台将不再展示，确定继续吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      const { code, msg } = await offlineApi<Common.DataResponse>(row.id)
      if (code === 0) {
        message.success('已下架')
        getList(null)
      }
      else {
        message.warning(`操作失败: ${msg}`)
      }
    },
  })
}

function handleDelete(row: CustomCode.ListItem) {
  dialog.warning({
    title: '警告',
    content: '删除后不可恢复，确定删除吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      const { code, msg } = await deleteByIdApi<Common.DataResponse>(row.id)
      if (code === 0) {
        message.success('已删除')
        getList(null)
      }
      else {
        message.warning(`删除失败: ${msg}`)
      }
    },
  })
}

async function getList(page: number | null) {
  tableIsLoading.value = true
  try {
    const req: Common.ListRequest = {
      page: page || pagination.page,
      limit: pagination.pageSize,
    }
    if (keyWord.value)
      req.keyword = keyWord.value

    const { data } = await getMyListApi<Common.ListResponse<CustomCode.ListItem[]>>(req)
    pagination.itemCount = data.count
    list.value = data.list ?? []
  }
  catch {
    message.warning('列表加载失败，请稍后重试')
  }
  finally {
    tableIsLoading.value = false
  }
}

onMounted(() => {
  getList(null)
})
</script>

<template>
  <div>
    <NH3>我的自定义代码</NH3>
    <NCard class="mb-[20px]">
      <div class="flex">
        <NInputGroup style="max-width:700px;">
          <NInput v-model:value="keyWord" clearable :style="{ width: '50%' }" placeholder="搜索标题或描述" @keyup.enter="handleSelect" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup>

        <span class="flex ml-auto">
          <NButton type="primary" ghost @click="handleAdd">
            新建
          </NButton>
        </span>
      </div>
    </NCard>

    <NDataTable
      :columns="columns"
      :data="list"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
      @update:page="handlePageChange"
    />
  </div>
</template>
