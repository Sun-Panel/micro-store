<script lang="ts" setup>
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import { NButton, NCard, NDataTable, NInput, NInputGroup, NTag, useDialog, useMessage } from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'

import { deletes, getList as getListApi, setActive } from '@/api/admin/version'
import { buildTimeString } from '@/utils/cmn'
import Edit from './Edit/index.vue'
import SecretEdit from './SecretEdit/index.vue'

const defaultEditInfo: Version.Info = {
  version: '',
  type: 'beta',
  releaseTime: '',
  description: '',
  downloadURL: '',
  pageUrl: '',
  isRolledBack: false,
  aloneSecretKey: 0,
}

const message = useMessage()
const dialog = useDialog()
const tableIsLoading = ref<boolean>(false)
const editModalShow = ref<boolean>(false)
// const historyDialogShow = ref<boolean>(false)
const keyWord = ref<string>()
const editInfo = ref<Version.Info>(defaultEditInfo)
// const updateUserId = ref(0)

// 密钥弹窗相关
const secretModalShow = ref<boolean>(false)
const secretVersion = ref<string>('')
const secretIsCreate = ref<boolean>(true)

function createColumns({
  update,
}: {
  update: (row: Version.Info) => void
}): DataTableColumns<Version.Info> {
  const handleCreateSecret = (row: Version.Info) => {
    secretVersion.value = row.version
    secretIsCreate.value = true
    secretModalShow.value = true
  }

  const handleViewSecret = (row: Version.Info) => {
    secretVersion.value = row.version
    secretIsCreate.value = false
    secretModalShow.value = true
  }
  return [

    {
      title: '版本类型',
      key: 'type',
      render(row) {
        switch (row.type) {
          case 'release':
            return h(NTag, { type: 'success', size: 'small' }, '正式版')
          case 'beta':
            return h(NTag, { type: 'warning', size: 'small' }, 'Beta')
          case 'alpha':
            return h(NTag, { type: 'info', size: 'small' }, 'Alpha')
          case 'rc':
            return h(NTag, { type: 'default', size: 'small' }, 'RC')
          case 'dev':
            return h(NTag, { type: 'error', size: 'small' }, '开发版')
          default:
            return h(NTag, row.type)
        }
      },
    },
    {
      title: '版本号',
      key: 'version',
    },
    {
      title: '最新版本',
      key: 'isActive',
      render(row) {
        if (row.isActive) {
          return h(NTag, { type: 'success' }, '是')
        }
        else {
          return h(
            NButton,
            {
              size: 'tiny',
              onClick() {
                dialog.warning({
                  title: '警告',
                  content: `你确定将[ v${row.version} ]设置为 ${row.type} 最新版？`,
                  positiveText: '确定',
                  negativeText: '取消',
                  onPositiveClick: () => {
                    setLatestVersion(row.id as number, row.type)
                  },
                })
              },
            },
            '设为最新版',
          )
        }
      },
    },
    // {
    //   title: '描述',
    //   key: 'description',
    // },
    {
      title: '发布时间',
      key: 'releaseTime',
      render(row) {
        if (row.releaseTime)
          return buildTimeString(row.releaseTime, 'YYYY-MM-DD HH:mm')

        else
          return '-'
      },
    },
    {
      title: '独立秘钥',
      key: 'aloneSecretKey',
      render(row) {
        switch (row.aloneSecretKey) {
          case 0:
            return [
              h(NTag, { type: 'warning' }, '无'),
              h(NButton, { size: 'tiny', style: { marginLeft: '4px' }, type: 'primary', onClick: () => handleCreateSecret(row) }, '创建'),
            ]
          case 1:
            return [
              h(NTag, { type: 'success' }, '存在'),
              h(NButton, { size: 'tiny', style: { marginLeft: '4px' }, type: 'info', onClick: () => handleViewSecret(row) }, '详情'),
            ]
          case 2:
            return [
              h(NTag, { type: 'error' }, '已停用'),
              h(NButton, { size: 'tiny', style: { marginLeft: '4px' }, type: 'info', onClick: () => handleViewSecret(row) }, '详情'),
            ]
          default:
            return [
              h(NTag, { type: 'error' }, `未知状态：${row.aloneSecretKey}`),
              h(NButton, { size: 'tiny', style: { marginLeft: '4px' }, type: 'info', onClick: () => handleViewSecret(row) }, '详情'),
            ]
        }
      },
    },
    // {
    //   title: '加入时间',
    //   key: 'createTime',
    //   render(row) {
    //     return timeFormat(String(row.createTime))
    //   },
    // },
    {
      title: '',
      key: '',
      render(row) {
        return h(
          'div',
          {},
          [
            h(
              NButton,
              {
                size: 'tiny',
                type: 'info',
                onClick() {
                  console.log(row)
                  editInfo.value = { ...row }
                  editModalShow.value = true
                },
              },
              '编辑',
            ),
            h(
              NButton,
              {
                size: 'tiny',
                style: { marginLeft: '4px' },
                type: 'error',
                onClick() {
                  dialog.warning({
                    title: '警告',
                    content: `你确定将[ v${row.version} ]-(${row.type}) 删除？`,
                    positiveText: '确定',
                    negativeText: '取消',
                    onPositiveClick: () => {
                      deletes([row.id as number]).then(() => {
                        getList(null)
                        message.success('已删除')
                      }).catch(() => {
                        message.success('删除失败')
                      })
                    },
                  })
                },
              },
              '删除',
            ),
          ],
        )
      },
    },
  ]
}

const dataList = ref<Version.Info[]>()

const columns = createColumns({
  update(row: Version.Info) {
    editInfo.value = row
    editModalShow.value = true
  },
})
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

function handelDone() {
  getList(null)
}

function handleSecretDone() {
  secretModalShow.value = false
  getList(null)
}

function handleAdd() {
  editInfo.value = { ...defaultEditInfo }
  editModalShow.value = true
}

async function setLatestVersion(id: number, type: 'release' | 'beta' | 'alpha' | 'rc' | 'dev') {
  await setActive(id, type).then(() => {
    message.success('设置成功')
    getList(null)
  }).catch(() => {
    message.error('设置失败')
  })
}

async function getList(page: number | null) {
  tableIsLoading.value = true
  const req: AdminUserManage.GetListRequest = {
    page: page || pagination.page,
    limit: pagination.pageSize,
  }
  if (keyWord.value !== '')
    req.keyWord = keyWord.value

  const { data } = await getListApi<Common.ListResponse<Version.Info[]>>(req)
  pagination.itemCount = data.count
  if (data.list)
    dataList.value = data.list
  tableIsLoading.value = false
}

onMounted(() => {
  getList(null)
})
</script>

<template>
  <div>
    <NCard class="mb-[20px]">
      <div class="flex">
        <NInputGroup style="max-width:700px;">
          <NInput v-model:value="keyWord" :style="{ width: '50%' }" placeholder="请输入账号" @keyup.enter="handleSelect" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup>

        <NButton type="primary" style="margin-left: auto;" @click="handleAdd">
          添加一个版本
        </NButton>
      </div>
    </NCard>

    <NDataTable
      :columns="columns"
      :data="dataList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"

      @update:page="handlePageChange"
    />

    <Edit v-model:visible="editModalShow" :version="editInfo" @done="handelDone" />
    <SecretEdit v-model:visible="secretModalShow" :version="secretVersion" :is-create="secretIsCreate" @done="handleSecretDone" />
  </div>
</template>
