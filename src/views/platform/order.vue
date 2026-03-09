<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NDataTable, NDropdown, NH2, NTag, useDialog } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import { getOrderList, queryPayStatus as queryPayStatusApi } from '@/api/goodsOrder'
import { getCurrencySymbol, timeFormat } from '@/utils/cmn'
import { SvgIcon } from '@/components/common'
import { router } from '@/router'
import { OrderStatus } from '@/enums/goodsOrder'

const tableIsLoading = ref<boolean>(false)
const editUserDialogShow = ref<boolean>(false)
const editUserUserInfo = ref<User.Info>()
const dialog = useDialog()

const createColumns = ({
  update,
}: {
  update: (row: GoodsOrder.Info) => void
}): DataTableColumns<GoodsOrder.Info> => {
  return [
    {
      title: '订单编号',
      key: 'orderNo',
    },

    {
      title: '支付平台',
      key: 'payPlatform',
      render(row) {
        let msg = ''
        switch (row.payPlatform) {
          case 1:
            msg = '支付宝'
            break
          case 2:
            msg = '微信'
            break
          case 3:
            msg = 'PayPal/信用卡'
            break
          default:
            msg = '未知'
            break
        }
        return msg
      },
    },
    {
      title: '订单状态',
      key: 'status',
      render(row) {
        let msg = ''
        switch (row.status) {
          case OrderStatus.PAY_WAIT:
            msg = '未付款'
            break
          case OrderStatus.PAY_SUCCESS:
            msg = '已付款'
            return h(NTag, { type: 'info' }, msg)
          case OrderStatus.CLOSE:
            msg = '订单过期'
            break
          case OrderStatus.CANCEL:
            msg = '取消'
            break
          case OrderStatus.FINISH:
            msg = '已完成'
            return h(NTag, { type: 'success' }, msg)
          default:
            break
        }
        return msg
      },
    },
    {
      title: '商品标题',
      key: 'goodsTitle',
      render(row) {
        return row.goodsSnapshotInfo?.title
      },
    },
    {
      title: '金额',
      key: 'countPrice',
      render(row) {
        const symbol = getCurrencySymbol(row.currencyCode) ?? '￥'
        return `(${row.currencyCode ? row.currencyCode : 'CNY'}) ${symbol}${row.countPrice}`
      },
    },
    {
      title: '下单时间',
      key: 'createTime',
      render(row) {
        return timeFormat(String(row.createTime))
      },
    },
    {
      title: '支付时间',
      key: 'payTime',
      render(row) {
        if (row.payTime === '')
          return row.payTime

        return timeFormat(String(row.payTime))
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
              case 'queryPayStatus':
                queryPayStatus(row)
                break
              case 'orderInfo':
                handleToOrderInfo(row)
                break
              case 'delete':
                dialog.warning({
                  title: '警告',
                  content: `你确定删除${row.orderNo}？`,
                  positiveText: '确定',
                  negativeText: '取消',
                  onPositiveClick: () => {
                    // deletes([row.OrderNo])
                  },
                })
                break

              default:
                break
            }
          },
          options: [
            {
              label: '查询支付状态',
              key: 'queryPayStatus',
            },
            {
              label: '订单详情',
              key: 'orderInfo',
            },
            // {
            //   label: '删除',
            //   key: 'delete',
            // },
          ],
        }, { default: () => btn })
      },
    },
  ]
}

const goodsList = ref<GoodsOrder.Info[]>()

const columns = createColumns({
  update(row: User.Info) {
    editUserUserInfo.value = row
    editUserDialogShow.value = true
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
    return `共 ${item.itemCount} 条记录`
  },
})

function handlePageChange(page: number) {
  getList(page)
}

// // 查询
// function handleSelect() {
//   getList(null)
// }

// // 添加
// function handleAdd() {
//   editUserDialogShow.value = true
//   editUserUserInfo.value = {}
// }

// function handelDone() {
//   editUserDialogShow.value = false
//   message.success('操作成功')
//   getList(null)
// }

async function getList(page: number | null) {
  tableIsLoading.value = true
  // const req: AdminUserManage.GetListRequest = {
  //   page: page || pagination.page,
  //   limit: pagination.pageSize,
  // }
  // if (keyWord.value !== '')
  //   req.keyWord = keyWord.value

  const { data } = await getOrderList<Common.ListResponse<GoodsOrder.Info[]>>()
  pagination.itemCount = data.count
  console.log(data.list)
  if (data.list)
    goodsList.value = data.list
  tableIsLoading.value = false
}

// async function deletes(ids: string[]) {
//   // const { code } = await AdminUserManageDelete(ids)
//   // if (code === 0) {
//   //   message.success('已删除')
//   //   getList(null)
//   // }
// }

async function queryPayStatus(row: GoodsOrder.Info) {
  try {
    await queryPayStatusApi(row.orderNo)
    getList(null)
  }
  catch (error) {

  }
}

function handleClick(keys: Array<string | number>, rows: object[], meta: { row: object | undefined; action: 'check' | 'uncheck' | 'checkAll' | 'uncheckAll' }) {
  console.log(keys, rows)
}

function handleToOrderInfo(row: GoodsOrder.Info) {
  router.push({ name: 'PlatformOrderInfo', query: { orderNo: row.orderNo } })
}

onMounted(() => {
  getList(null)
})
</script>

<template>
  <div>
    <NH2 align-text prefix="bar">
      订单记录
    </NH2>
    <!-- <NCard class="mb-[20px]">
      <div class="flex">
        <NInputGroup style="max-width:700px;">
          <NInput v-model:value="keyWord" :style="{ width: '50%' }" placeholder="请输入账号或者昵称来查询" />
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
    </NCard> -->

    <NDataTable
      :columns="columns"
      :data="goodsList"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
      @update:checked-row-keys="handleClick"
      @update:page="handlePageChange"
    />
  </div>
</template>
