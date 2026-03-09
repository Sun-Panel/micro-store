<script setup lang="ts">
// 参考教程：https://www.npmjs.com/package/@paddle/paddle-js
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import type { Environments, PaddleEventData } from '@paddle/paddle-js'
import { useDialog } from 'naive-ui'
import { Paddle } from '@/components/common'
import { getOrderInfo as getOrderInfoApi, queryPayStatus } from '@/api/goodsOrder'
import { router } from '@/router'

// 获取订单详情和商品的priceId

const route = useRoute()
const dialog = useDialog()
const transactionId = ref('')
const clientToken = ref(route.query.tk as string)
const environment = ref(route.query.env as Environments)
const paddleRef = ref()
const orderNo = ref('')

async function getOrderInfo() {
  await getOrderInfoApi<GoodsOrder.Info>(route.query.orderNo as string).then(({ data }) => {
    transactionId.value = data.payPlatformOrderNo
    orderNo.value = data.orderNo
    paddleRef.value.init()
  }).catch(() => {

  })
}

function payEvent(event: PaddleEventData) {
  if (event.name === 'checkout.completed') {
    // 查询订单,更新订单状态
    queryPayStatus(orderNo.value)
    dialog.success({
      title: 'Success',
      content: 'You have successfully paid, do you go to the order detail page?',
      positiveText: 'YES',
      negativeText: 'NO',
      onPositiveClick: () => {
        router.push({ name: 'PlatformOrderInfo', query: { orderNo: orderNo.value } })
      },
    })
  }
}

onMounted(() => {
  getOrderInfo()
})
</script>

<template>
  <div>
    <Paddle ref="paddleRef" :transaction-id="transactionId" :client-token="clientToken" :environment="environment" @event-callback="payEvent" />
  </div>
</template>
