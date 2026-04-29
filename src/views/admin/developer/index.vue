<script lang="ts" setup>
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import { NButton, NCard, NDataTable, NDropdown, NInput, NInputGroup, NSelect, NSwitch, useDialog, useMessage } from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'
import { deletes, getList, updateStatus } from '@/api/admin/developer'
import { SvgIcon } from '@/components/common'
import { timeFormat } from '@/utils/cmn'
import EditDeveloper from './EditDeveloper/index.vue'

const message = useMessage()
const tableIsLoading = ref<boolean>(false)
const editDialogShow = ref<boolean>(false)
const keyWord = ref<string>()
const editInfo = ref<Developer.DeveloperInfo>()
const statusFilter = ref<number | null>(null)
const dialog = useDialog()

const statusOptions = [
  { label: '全部', value: null },
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 },
]

function createColumns({
  update,
  changeStatus,
}: {
  update: (row: Developer.DeveloperInfo) => void
  changeStatus: (row: Developer.DeveloperInfo) => void
}): DataTableColumns<Developer.DeveloperInfo> {
  return [
    { title: 'ID', key: 'id', width: 80 },
    {
      title: '用户',
      key: 'user',
      width: 200,
      render(row) {
        if (!row.user)
          return '-'
        return h('div', { class: 'flex items-center gap-2' }, [
          // h(NAvatar, { size: 'small', src: row.user.headImage, round: true }),
          h('div', { class: 'flex flex-col' }, [
            h('span', { class: 'text-sm' }, row.user.name || row.user.username || '-'),
            h('span', { class: 'text-xs text-gray-400' }, row.user.mail || '-'),
          ]),
        ])
      },
    },
    { title: '开发者标识', key: 'developerName' },
    { title: '作者昵称', key: 'name' },
    { title: '联系邮箱', key: 'contactMail' },
    { title: '收款方式', key: 'paymentMethod' },
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
        const btn = h(NButton, { strong: true, tertiary: true, size: 'small' }, {
          default: () => h(SvgIcon, { icon: 'mingcute:more-1-fill' }),
        })

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
                  content: `确定删除开发者"${row.developerName}"吗？`,
                  positiveText: '确定',
                  negativeText: '取消',
                  onPositiveClick: () => handleDelete([row.id]),
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

const dataList = ref<Developer.DeveloperInfo[]>([])

const columns = createColumns({
  update(row: Developer.DeveloperInfo) {
    editInfo.value = row
    editDialogShow.value = true
  },
  changeStatus(row: Developer.DeveloperInfo) {
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
    return `共 ${item.itemCount} 个开发者`
  },
})

function handleSelect() {
  pagination.page = 1
  fetchList()
}

// function handleAdd() {
//   editInfo.value = undefined
//   editDialogShow.value = true
// }

function handleDone() {
  editDialogShow.value = false
  message.success('操作成功')
  fetchList()
}

async function fetchList() {
  tableIsLoading.value = true
  const req: Developer.GetListRequest = {
    page: pagination.page,
    limit: pagination.pageSize,
  }
  if (keyWord.value)
    req.keyWord = keyWord.value
  if (statusFilter.value !== null)
    req.status = statusFilter.value

  try {
    const { data } = await getList<Common.ListResponse<Developer.DeveloperInfo[]>>(req)
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

async function handleChangeStatus(row: Developer.DeveloperInfo) {
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
          <NInput v-model:value="keyWord" :style="{ width: '40%' }" placeholder="请输入开发者标识或邮箱搜索" />
          <NSelect v-model:value="statusFilter" :options="statusOptions" :style="{ width: '120px' }" placeholder="状态" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup>
      </div>
    </NCard>

    <NDataTable
      :columns="columns"
      :data="dataList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
    />
    <EditDeveloper v-model:visible="editDialogShow" :developer-info="editInfo" @done="handleDone" />
  </div>
</template>
