<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NInput, NInputGroup, NModal, NTag } from 'naive-ui'
import type { PaginationProps } from 'naive-ui'
// import UpdateUserAuthorize from './UpdateUserAuthorize/index.vue'
import OrderInfo from './OrderInfo/index.vue'
import { getList as getListApi } from '@/api/admin/orderManage'
import { OrderStatus } from '@/enums/goodsOrder'
import { getCurrencySymbol, timeFormat } from '@/utils/cmn'

const tableIsLoading = ref<boolean>(false)
const orderInfoDialogShow = ref<boolean>(false)
const keyWord = ref<string>()

const activeOrderNo = ref('')

const orderList = ref< GoodsOrder.Info[]>()

const columns = [
  {
    title: '订单编号',
    key: 'orderNo',
  },

  {
    title: '账号',
    key: 'orderNo',
    render(row: GoodsOrder.Info) {
      return `${row.user.username} - [${row.user.name}]`
    },
  },

  {
    title: '订单状态',
    key: 'status',
    render(row: GoodsOrder.Info) {
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
    render(row: GoodsOrder.Info) {
      return row.goodsSnapshotInfo?.title
    },
  },
  {
    title: '金额',
    key: 'countPrice',
    render(row: GoodsOrder.Info) {
      const symbol = getCurrencySymbol(row.currencyCode) ?? '￥'
      return `(${row.currencyCode ? row.currencyCode : 'CNY'}) ${symbol}${row.countPrice}`
    },
  },
  {
    title: '下单时间',
    key: 'createTime',
    render(row: GoodsOrder.Info) {
      return timeFormat(String(row.createTime))
    },
  },
  {
    title: '支付时间',
    key: 'payTime',
    render(row: GoodsOrder.Info) {
      if (row.payTime === '')
        return row.payTime

      return timeFormat(String(row.payTime))
    },
  },
  {
    title: '',
    key: '',
    render(row: GoodsOrder.Info) {
      return h(
        NButton,
        {
          size: 'tiny',
          type: 'info',
          onClick() {
            activeOrderNo.value = row.orderNo
            orderInfoDialogShow.value = true
          },
        },
        '订单详情',
      )
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

async function getList(page: number | null) {
  tableIsLoading.value = true
  const req: AdminUserManage.GetListRequest = {
    page: page || pagination.page,
    limit: pagination.pageSize,
  }
  if (keyWord.value !== '')
    req.keyWord = keyWord.value

  const { data } = await getListApi<Common.ListResponse<GoodsOrder.Info[]>>(req)
  pagination.itemCount = data.count
  if (data.list)
    orderList.value = data.list
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
          <NInput v-model:value="keyWord" :style="{ width: '50%' }" placeholder="请输入订单号" @keyup.enter="handleSelect" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup>
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

    <NModal v-model:show="orderInfoDialogShow" preset="card" style="max-width: 800px" title="订单详情">
      <OrderInfo :order-no="activeOrderNo" />
    </NModal>
  </div>
</template>
