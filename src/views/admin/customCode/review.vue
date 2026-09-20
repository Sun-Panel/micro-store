<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NH3, NInput, NInputGroup, NTag, NTooltip, useMessage } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'

import { getList as getListApi } from '@/api/admin/customCodeReview'
import { router } from '@/router'

const tableIsLoading = ref<boolean>(false)
const keyWord = ref<string>()
const list = ref<CustomCode.ReviewListItem[]>([])
const message = useMessage()

function codeTypeTexts(codeTypes?: number[]) {
  if (!codeTypes?.length)
    return '-'
  const map: Record<number, string> = { 1: 'JS', 2: 'CSS', 3: '页脚' }
  return codeTypes.map(item => map[item] ?? item).join(' / ')
}

const columns: DataTableColumns<CustomCode.ReviewListItem> = [
  { title: '标题', key: 'title' },
  { title: '作者', key: 'authorName', width: 140 },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  {
    title: '代码类型',
    key: 'codeTypes',
    width: 140,
    render(row) {
      return codeTypeTexts(row.codeTypes)
    },
  },
  {
    title: '类型',
    key: 'onlineStatus',
    width: 100,
    render(row) {
      return row.onlineStatus === 1
        ? h(NTag, { type: 'info', size: 'small' }, { default: () => '修改' })
        : h(NTag, { size: 'small' }, { default: () => '新提交' })
    },
  },
  { title: '提交时间', key: 'createTime', width: 180 },
  {
    title: '机器预审',
    key: 'machineAudit',
    width: 120,
    render(row) {
      const ma = row.machineAudit
      if (!ma)
        return h(NTag, { size: 'small' }, { default: () => '-' })
      if (ma.auto)
        return h(NTag, { type: 'success', size: 'small' }, { default: () => '自动通过' })
      const reasons = ma.reasons?.length ? ma.reasons : ['未受信任或安全扫描未通过']
      return h(
        NTooltip,
        { placement: 'left' },
        {
          trigger: () =>
            h(NTag, { type: 'warning', size: 'small' }, { default: () => '转人工' }),
          default: () =>
            h('div', reasons.map(r => h('div', { style: 'margin:2px 0;max-width:360px' }, r))),
        },
      )
    },
  },
  {
    title: '',
    key: 'actions',
    width: 100,
    render(row) {
      return h(
        NButton,
        {
          size: 'tiny',
          type: 'primary',
          onClick: () => router.push({ name: 'AdminCustomCodeReviewDetail', query: { id: row.id } }),
        },
        '审核',
      )
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

async function getList(page: number | null) {
  tableIsLoading.value = true
  try {
    const req: Common.ListRequest = {
      page: page || pagination.page,
      limit: pagination.pageSize,
    }
    if (keyWord.value)
      req.keyword = keyWord.value

    const { data } = await getListApi<Common.ListResponse<CustomCode.ReviewListItem[]>>(req)
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
    <NH3>自定义代码审核</NH3>
    <NCard class="mb-[20px]">
      <NInputGroup style="max-width:700px;">
        <NInput v-model:value="keyWord" clearable :style="{ width: '50%' }" placeholder="搜索标题或描述" @keyup.enter="getList(null)" />
        <NButton type="primary" @click="getList(null)">
          查询
        </NButton>
      </NInputGroup>
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
