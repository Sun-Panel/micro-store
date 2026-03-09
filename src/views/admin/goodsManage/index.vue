<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NDropdown, NInput, NInputGroup, NTag, useDialog, useMessage } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import Edit from './Edit/index.vue'
import EditSale from './EditSale/index.vue'
import { createSnapshot as createSnapshotApi, deletes as deletesApi, getList as getListApi } from '@/api/admin/goodsManage'

import { timeFormat } from '@/utils/cmn'
import { SvgIcon } from '@/components/common'

const message = useMessage()
const tableIsLoading = ref<boolean>(false)
const editDialogShow = ref<boolean>(false)
const editSaleDialogShow = ref<boolean>(false)
const keyWord = ref<string>()
const editGoodsInfo = ref<Admin.GoodsManage.GoodsInfo | null>(null)
const dialog = useDialog()

const createColumns = ({
  updateBase,
}: {
  updateBase: (row: Admin.GoodsManage.GoodsInfo) => void
}): DataTableColumns<Admin.GoodsManage.GoodsInfo> => {
  return [
    {
      title: 'sort',
      key: 'sort',
      // render(row) {
      //   return timeFormat(String(row.updateTime))
      // },
    },
    {
      title: '标题',
      key: 'title',
    },
    {
      title: '库存数量',
      key: 'num',
    },
    {
      title: '更新时间',
      key: 'updateTime',
      render(row) {
        return timeFormat(String(row.updateTime))
      },
    },

    {
      title: '商品快照',
      key: 'lastSnapshotId',
      render(row) {
        if (row.lastSnapshotId) {
          return h(NTag, { type: 'success', bordered: false }, '有')
        }

        else {
          return h(NButton, {
            type: 'error',
            size: 'tiny',
            onClick: () => {
              handleCreateSnapshot(row)
            },
          }, '创建')
        }
      },
    },
    {
      title: '状态',
      key: 'status',
      render(row) {
        let msg = ''
        switch (row.status) {
          case 1:
            msg = '已上架'
            return h(NTag, { type: 'success' }, msg)
          case 2:
            msg = '已下架'
            return h(NTag, { type: 'error' }, msg)
          default:
            msg = '未知'
            break
        }
        return msg
      },
    },
    {
      title: '操作',
      key: '',
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
              return h(
                SvgIcon, {
                  icon: 'mingcute:more-1-fill',
                },
              )
            },
          },
        )

        return h(NDropdown, {
          trigger: 'click',
          onSelect(key: string | number) {
            console.log(key)
            switch (key) {
              case 'updateBase':
                updateBase(row)
                break
              case 'updateSale':
                updateSale(row)
                break
              case 'delete':
                dialog.warning({
                  title: '警告',
                  content: `你确定删除${row.title}？`,
                  positiveText: '确定',
                  negativeText: '取消',
                  onPositiveClick: () => {
                    deletes([row.id as number])
                  },
                })
                break

              default:
                break
            }
          },
          options: [
            {
              label: '修改基础信息',
              key: 'updateBase',
            },
            {
              label: '修改库存、排序、状态',
              key: 'updateSale',
            },
            {
              label: '删除',
              key: 'delete',
            },
          ],
        }, { default: () => btn })
      },
    },
  ]
}

const userList = ref<Admin.GoodsManage.GoodsInfo[]>()

const columns = createColumns({
  updateBase(row: Admin.GoodsManage.GoodsInfo) {
    editGoodsInfo.value = row
    editDialogShow.value = true
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
    return `共 ${item.itemCount} 件商品`
  },
})

function handlePageChange(page: number) {
  getList(page)
}

async function handleCreateSnapshot(row: Admin.GoodsManage.GoodsGetListItemResp) {
  console.log('创建快照', row.id)
  try {
    await createSnapshotApi(row?.id as number)
    getList(null)
  }
  catch (error) {
    console.log('创建快照错误', error)
  }
}

// 查询
function handleSelect() {
  getList(null)
}

// 添加
function handleAdd() {
  editDialogShow.value = true
  editGoodsInfo.value = null
}

function handelDone() {
  editDialogShow.value = false
  message.success('操作成功')
  getList(null)
}

function handelSaleDone() {
  editSaleDialogShow.value = false
  message.success('操作成功')
  getList(null)
}

async function getList(page: number | null) {
  tableIsLoading.value = true
  const req: AdminUserManage.GetListRequest = {
    page: page || pagination.page,
    limit: pagination.pageSize,
  }
  if (keyWord.value !== '')
    req.keyWord = keyWord.value

  const { data } = await getListApi<Common.ListResponse<Admin.GoodsManage.GoodsGetListItemResp[]>>(req)
  pagination.itemCount = data.count
  if (data.list)
    userList.value = data.list
  tableIsLoading.value = false
}

// 修改销售信息
function updateSale(row: Admin.GoodsManage.GoodsInfo) {
  editGoodsInfo.value = row
  editSaleDialogShow.value = true
}

async function deletes(ids: number[]) {
  const { code } = await deletesApi(ids)
  if (code === 0) {
    message.success('已删除')
    getList(null)
  }
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
          <NInput v-model:value="keyWord" :style="{ width: '50%' }" placeholder="请输入关键字来查询" />
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
      :data="userList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"

      @update:page="handlePageChange"
    />
    <Edit v-model:visible="editDialogShow" :info="editGoodsInfo" @done="handelDone" />
    <EditSale v-model:visible="editSaleDialogShow" :info="editGoodsInfo" @done="handelSaleDone" />
  </div>
</template>
