<script setup lang="ts">
import { ref } from 'vue'
import type { Environments, Paddle, PaddleEventData } from '@paddle/paddle-js'
import { initializePaddle } from '@paddle/paddle-js'
import { NH3 } from 'naive-ui'

const props = defineProps<{
  clientToken: string
  transactionId: string
  environment: 'production' | 'sandbox' | string
  // quantity: number
}>()

const emit = defineEmits<{
  (e: 'eventCallback', event: PaddleEventData): void
}>()

// Create a local state to store Paddle instance
const paddle = ref<Paddle>()
const currencyCode = ref('')
const totalPrice = ref(0)
const loading = ref(true)
// const env = ref(props.environment === 'sandbox' ? 'sandbox' : 'production')

// Callback to open a checkout
const openCheckout = () => {
  paddle.value?.Checkout.open({
    transactionId: props.transactionId,
    settings: {
      displayMode: 'inline',
      frameTarget: 'checkout-container',
      frameInitialHeight: 450,
      frameStyle: 'width: 100%; min-width: 312px; background-color: transparent; border: none;',
    },

    // items: [{ priceId: 'pri_01hvg8rmef14pg62d73c037hf1', quantity: 1 }],
  })
}

// Download and initialize Paddle instance from CDN
// onMounted(async () => {
//   console.log('环境', env.value)
// })

async function Init() {
  initializePaddle({
    environment: props.environment as Environments,
    token: props.clientToken,
    eventCallback: (event: PaddleEventData) => {
      console.log(event)
      if (event.data) {
        currencyCode.value = event.data.currency_code
        totalPrice.value = event.data.totals.total
        loading.value = false
      }
      emit('eventCallback', event)
    },
  }).then(
    (paddleInstance: Paddle | undefined) => {
      if (paddleInstance) {
        paddle.value = paddleInstance
        openCheckout()
      }
    },
  )
}

// 暴露出去的内容
defineExpose({
  init() {
    console.log('初始化订单')
    Init()
  },
})
</script>

<template>
  <div>
    <div class="text-center">
      <span v-if="loading" class="font-extrabold">Loading...</span>
      <NH3 v-if="!loading">
        <span class="font-extrabold">
          <span>Price: </span>{{ currencyCode }} {{ totalPrice }}
        </span>
      </NH3>
    </div>
    <div class="checkout-container" />
    <!-- Your template code here -->
  </div>
</template>
