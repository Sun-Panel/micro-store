<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { useRoute, useRouter } from 'vue-router'
import { NButton, NH2, useMessage } from 'naive-ui'
import { getOrderInfo as getOrderInfoApi, queryPayStatus as queryPayStatusApi } from '@/api/goodsOrder'
import { OrderInfo } from '@/views/components'
import { OrderStatus } from '@/enums/goodsOrder'
import { t } from '@/locales'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'

const route = useRoute()
const router = useRouter()
const ms = useMessage()
const orderNo = route.query.orderNo as string
const orderInfo = ref<GoodsOrder.Info | null>(null)

async function getOrderInfo() {
  await getOrderInfoApi<GoodsOrder.Info>(orderNo).then(({ data }) => {
    orderInfo.value = data
  }).catch((error) => {
    apiRespErrMsgAndCustomCodeNeg1Msg(error, t('common.queryFailedWaitRetry'))
  })
}

function handleClick() {
  router.push({ name: 'PlatformOrder' })
}

async function queryPayStatus(isMessage?: boolean) {
  if (orderInfo.value?.status === OrderStatus.PAY_WAIT) {
    await queryPayStatusApi(orderInfo.value?.orderNo).then(() => {
      getOrderInfo()
      if (isMessage)
        ms.success(t('order.queryOrderStatusSuccess'))
    }).catch((error) => {
      apiRespErrMsgAndCustomCodeNeg1Msg(error, t('order.queryOrderStatusFailed'))
    })
  }
}

onMounted(async () => {
  await getOrderInfo()
  queryPayStatus()
})
</script>

<template>
  <div>
    <NH2>{{ t('order.orderDetail') }}</NH2>
    <OrderInfo :order-info="orderInfo as GoodsOrder.Info" />

    <NButton v-if="orderInfo?.status === OrderStatus.PAY_WAIT" style="margin-top:5px ;" @click="queryPayStatus(true)">
      {{ t('order.queryOrderStatus') }}
    </NButton>

    <NButton type="info" style="margin-top:5px ;margin-left:5px" @click="handleClick">
      {{ t('order.goToOrderHome') }}
    </NButton>
  </div>
</template>
