<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NH3, NInput, NInputGroup, useDialog, useMessage } from 'naive-ui'
import type { PaginationProps } from 'naive-ui'

import { deleteByPageName as deleteByNameApi, getList as getListApi } from '@/api/admin/htmlPageManage'
import { router } from '@/router'

const tableIsLoading = ref<boolean>(false)
const keyWord = ref<string>()

const dialog = useDialog()
const message = useMessage()

const orderList = ref<HtmlPage.ListItem[]>()

const columns = [
  {
    title: '页面描述',
    key: 'pageDescription',
  },

  {
    title: '页面名称',
    key: 'pageName',
  },

  {
    title: '是否需要登录',
    key: 'isLogin',
    render(row: HtmlPage.ListItem) {
      return row.isLogin ? '是' : '否'
    },
  },

  {
    title: '站内信模板',
    key: 'messageTemplateFlag',
  },

  {
    title: '站内信模板定位',
    key: 'messageTemplatePosition',
    render(row: HtmlPage.ListItem) {
      let text = ''
      switch (row.messageTemplatePosition) {
        case 'top':
          text = '顶部'
          break
        case 'bottom':
          text = '底部'
          break
        default:
          break
      }
      return text
    },
  },

  {
    title: '',
    key: '',
    render(row: HtmlPage.ListItem) {
      const deleteButton = h(
        NButton,
        {
          size: 'tiny',
          type: 'error',
          style: { marginLeft: '5px' },
          onClick() {
            dialog.warning({
              title: '警告',
              content: '你确定要删除这个页面吗，删除后不可以恢复？',
              positiveText: '确定',
              negativeText: '取消',
              onPositiveClick: () => {
                deleteByName(row.pageName)
              },

            })
          },
        },
        '删除',
      )

      const btns = [
        h(
          NButton,
          {
            size: 'tiny',
            type: 'info',
            onClick() {
              router.push({ name: 'AdminHtmlPageManageEdit', query: { pageName: row.pageName } })
            },
          },
          '修改',
        ),
      ]

      btns.push(deleteButton)

      return btns
    },
  },
]
const pagination = reactive({
  page: 1,
  showSizePicker: true,
  pageSizes: [10, 30, 50, 100, 200],
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

// 查询
function handleSelect() {
  getList(null)
}

function handleAdd() {
  router.push({ name: 'AdminHtmlPageManageEdit' })
}

async function deleteByName(name: string) {
  try {
    const { code, msg } = await deleteByNameApi<Common.DataResponse>(name)
    if (code !== 0) {
      message.warning(`删除失败: ${msg}`)
      return
    }
    handleRefreshList()
  }
  catch (error) {
    message.warning('删除失败，请稍后重试')
  }
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

    const { data } = await getListApi<Common.ListResponse<HtmlPage.ListItem[]>>(req)
    pagination.itemCount = data.count
    if (data.list)
      orderList.value = data.list
  }
  catch {
    message.warning('列表加载失败，请稍后重试')
  }
  finally {
    tableIsLoading.value = false
  }
}

function handleRefreshList() {
  getList(null)
}

onMounted(() => {
  getList(null)
})
</script>

<template>
  <div>
    <NH3>HTML页面管理</NH3>
    <NCard class="mb-[20px]">
      <div class="flex">
        <NInputGroup style="max-width:700px;">
          <NInput v-model:value="keyWord" clearable :style="{ width: '50%' }" placeholder="请输入描述或者页面名称" @keyup.enter="handleSelect" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup>

        <span class="flex ml-auto">
          <NButton type="primary" ghost @click="handleAdd">
            添加
          </NButton>
        </span>
      </div>
    </NCard>

    <NDataTable
      :columns="columns"
      :data="orderList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"

      @update:page="handlePageChange"
    />
  </div>
</template>
