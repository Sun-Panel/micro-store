<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NH3, NInput, NInputGroup, useDialog } from 'naive-ui'
import type { PaginationProps } from 'naive-ui'
import Editor from './Editor/index.vue'

import { deleteByMdPageName as deleteByNameApi, getList as getListApi } from '@/api/admin/mdPageManage'

const tableIsLoading = ref<boolean>(false)
const editorShow = ref<boolean>(false)
const activeInfo = ref<MdPage.ListItem | null>()
const keyWord = ref<string>()

const dialog = useDialog()

const orderList = ref< MdPage.ListItem[]>()

const columns = [
  {
    title: '页面描述',
    key: 'mdPageDescription',
  },

  {
    title: '页面名称',
    key: 'mdPageName',
  },

  {
    title: '是否需要登录',
    key: 'isLogin',
    render(row: MdPage.ListItem) {
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
    render(row: MdPage.ListItem) {
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
    render(row: MdPage.ListItem) {
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
                deleteByName(row.mdPageName)
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
              activeInfo.value = row
              editorShow.value = true
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
  activeInfo.value = null
  editorShow.value = true
}

async function deleteByName(name: string) {
  try {
    const { data } = await deleteByNameApi(name)
    console.log(data)
    handleRefreshList()
  }
  catch (error) {

  }
}

async function getList(page: number | null) {
  tableIsLoading.value = true
  const req: AdminUserManage.GetListRequest = {
    page: page || pagination.page,
    limit: pagination.pageSize,
  }
  if (keyWord.value !== '')
    req.keyWord = keyWord.value

  const { data } = await getListApi<Common.ListResponse<MdPage.ListItem[]>>(req)
  pagination.itemCount = data.count
  if (data.list)
    orderList.value = data.list
  tableIsLoading.value = false
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
    <NH3>Markdown页面管理</NH3>
    <NCard class="mb-[20px]">
      <div class="flex">
        <NInputGroup style="max-width:700px;">
          <NInput v-model:value="keyWord" clearable :style="{ width: '50%' }" placeholder="请输入描述或者变量名" @keyup.enter="handleSelect" />
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

    <Editor v-model:visible="editorShow" :info="activeInfo as MdPage.ListItem" @done="handleRefreshList" />
  </div>
</template>
