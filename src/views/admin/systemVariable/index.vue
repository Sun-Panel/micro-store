<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NEllipsis, NH3, NInput, NInputGroup, useDialog } from 'naive-ui'
import type { PaginationProps } from 'naive-ui'
import AddOrEdit from './AddOrEdit/index.vue'
import Editor from './Editor/index.vue'
import { getMultiple } from '@/api/system/systemVariable'

import { deleteByName as deleteByNameApi, getList as getListApi,getByCache,clearCache as clearCacheOne } from '@/api/admin/systemVariable'
import { apiRespErrMsg, message } from '@/utils/cmn/apiMessage'

interface PageConfig {
  system_variable_is_create: string
  system_variable_is_edit: string
  system_variable_is_delete: string
}

const pageConfig = ref<PageConfig>()
const tableIsLoading = ref<boolean>(false)
const editorShow = ref<boolean>(false)
const addOrEditShow = ref(false)
const activeInfo = ref<SystemVariable.SystemVariableEditReq | null>(null)
const keyWord = ref<string>()

const dialog = useDialog()

const orderList = ref< SystemVariable.SystemVariableListItem[]>()

const columns = [
  {
    title: '描述',
    key: 'description',
  },

  {
    title: '变量名称',
    key: 'configName',
  },

  {
    title: '变量值',
    key: 'configValue',
    render(row: SystemVariable.SystemVariableListItem) {
      return h(NEllipsis,
        {
          style: { maxWidth: '300px' },
          lineClamp: 1,
          expandTrigger: 'click',
        },
        row.configValue,
      )
    },
  },

  {
    title: '',
    key: '',
    render(row: SystemVariable.SystemVariableListItem) {
      const updateKeyButton = h(
        NButton,
        {
          size: 'tiny',
          type: 'warning',
          style: { marginLeft: '5px' },
          onClick() {
            activeInfo.value = {
              id: row.id,
              description: row.description,
              name: row.configName,
              value: row.configValue,
            }
            addOrEditShow.value = true
          },
        },
        '修改',
      )

      const deleteButton = h(
        NButton,
        {
          size: 'tiny',
          type: 'error',
          style: { marginLeft: '5px' },
          onClick() {
            dialog.warning({
              title: '警告',
              content: '你确定要删除这个变量吗，如果删除错可能会导致很严重的后果？',
              positiveText: '确定',
              negativeText: '取消',
              onPositiveClick: () => {
                deleteByName(row.configName)
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
              activeInfo.value = {
                id: row.id,
                description: row.description,
                name: row.configName,
                value: row.configValue,
              }
              editorShow.value = true
            },
          },
          '修改值',
        ),
        h(
          NButton,
          {
            size: 'tiny',
            type: 'warning',
            style: { marginLeft: '5px' },
            onClick() {
              forcedFlusheCache(row.configName,row.configValue)
            },
          },
          '更新缓存',
        )
      ]

      if (pageConfig.value?.system_variable_is_edit === 'true')
        btns.push(updateKeyButton)

      if (pageConfig.value?.system_variable_is_delete === 'true')
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
  addOrEditShow.value = true
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

  const { data } = await getListApi<Common.ListResponse<SystemVariable.SystemVariableListItem[]>>(req)
  pagination.itemCount = data.count
  if (data.list)
    orderList.value = data.list
  tableIsLoading.value = false
}

async function getPageSystemVariableConfig() {
  const keys = [
    'system_variable_is_delete', 'system_variable_is_edit', 'system_variable_is_create',
  ]
  try {
    const { data } = await getMultiple<PageConfig>(keys)
    pageConfig.value = data
  }
  catch (error) {

  }
}

async function forcedFlusheCache(name: string,lastestValue:string) {
  await clearCacheOne(name).then(() => {
    getByCache<string>( name ).then(({ data }) => {
      if(data===lastestValue){
        message.success("缓存更新成功")

      }else{
        dialog.error({
          title: '错误',
          content: '缓存更新失败',
        })

      }
    }).catch(() => {
      dialog.error({
        title: '错误',
        content: '数据获取失败',
      })
    })
  }).catch((error) => {
    
    apiRespErrMsg(error)
  })
}

function handleRefreshList() {
  getList(null)
}

onMounted(() => {
  getPageSystemVariableConfig()
  getList(null)
})
</script>

<template>
  <div>
    <NH3>系统变量管理</NH3>

    <NCard class="mb-[20px]">
      <div class="flex">
        <NInputGroup style="max-width:700px;">
          <NInput v-model:value="keyWord" clearable :style="{ width: '50%' }" placeholder="请输入描述或者变量名" @keyup.enter="handleSelect" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup>

        <span v-if="pageConfig?.system_variable_is_create === 'true'" class="flex ml-auto">
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

    <Editor v-model:visible="editorShow" :info="activeInfo" @done="handleRefreshList" />
    <AddOrEdit v-model:visible="addOrEditShow" :info="activeInfo" @done="handleRefreshList" />
  </div>
</template>
