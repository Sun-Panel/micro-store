<script setup lang='ts'>
import { NButton, NCard, NImage, NRadio, NSpin, NTag, useDialog, useMessage } from 'naive-ui'
import { h, ref } from 'vue'
import alipayLogoImg from '@/assets/pay_platform/alipayLogo.png'
import { pay } from '@/api/goodsOrder'
import { t } from '@/locales'
import { router } from '@/router'
import { ErrorCode } from '@/enums/errorCode'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'

const props = defineProps<{
  price: number
  goodsOrderNo: string
  title?: string
  goodsTitle?: string
  goodsDiscount?: string
}>()
const emit = defineEmits<{
  (event: 'pay'): void
}>()
const getPayUrlLoading = ref(false)
const dialog = useDialog()
const ms = useMessage()
enum PayPlatform {
  AliPay = 'alipay',
  WeiXinPay = 'weixinpay',
  PayPal = 'paypal',
  Paddle = 'paddle',
}
const payPlatform = ref<PayPlatform>(PayPlatform.AliPay)

async function handlePay() {
  getPayUrlLoading.value = true
  let payPlatformNum = 0
  if (payPlatform.value === PayPlatform.PayPal)
    payPlatformNum = 3

  switch (payPlatform.value) {
    case PayPlatform.AliPay:
      payPlatformNum = 1
      break
    case PayPlatform.Paddle:
      payPlatformNum = 3
      break
    default:
      break
  }

  await pay<GoodsOrder.PayResponse>(props.goodsOrderNo, payPlatformNum).then(({ data }) => {
    getPayUrlLoading.value = false
    window.open(data.payUrl)
    askPayStatus(data.payUrl)
    emit('pay')
  }).catch((error) => {
    getPayUrlLoading.value = false
    if (error.code === ErrorCode.OrderCreateFailed) {
      ms.error(t('packageStore.createPlatformOrderFail'), { closable: true, duration: 20000 })
      return
    }
    else if ((error.code === ErrorCode.GoodsNoUsePayPlatform)) {
      ms.error(t('packageStore.goodsNoUsePayPlatform'))
      return
    }
    apiRespErrMsg(error)
  })
}
function askPayStatus(payUrl: string) {
  const d = dialog.warning({
    showIcon: false,
    content: () => t('components.pay.helpPayText'),
    action: () => {
      return [
        h(NButton, {
          type: 'info',
          size: 'small',
          onClick() {
            window.open(payUrl)
            d.destroy()
          },
        }, t('components.pay.openPayPage')),
        h(NButton, {
          type: 'info',
          size: 'small',
          onClick() {
            router.push({
              name: 'PlatformOrderInfo',
              query: {
                orderNo: props.goodsOrderNo,
              },
            })
            d.destroy()
          },
        }, t('components.pay.goToOrderInfoPage')),
      ]
    },
  })
}
</script>

<template>
  <div class="max-w-[1200px] mx-[auto]">
    <NSpin :show="getPayUrlLoading">
      <NCard>
        <div class="text-lg font-semibold">
          {{ title }}
        </div>

        <div v-if="goodsTitle">
          {{ goodsTitle }}
        </div>

        <div v-if="goodsDiscount">
          <NTag size="small" :bordered="false" type="error">
            {{ goodsDiscount }}
          </NTag>
        </div>

        <div class="text-lg font-semibold">
          {{ t('common.amount') }}: CNY<span class="text-red-500 text-xl"> ￥{{ price }}</span>
        </div>
        <div v-if="payPlatform === PayPlatform.Paddle" class="mt-2">
          <span class="text-orange-500">
            {{ t('components.pay.serviceFeeWarning') }}
          </span>
        </div>
      </NCard>

      <div class="mt-5">
        {{ t('components.pay.selectPayPlatform') }}
      </div>

      <div class="overflow-hidden gap-5 grid grid-cols-2">
        <NCard>
          <NRadio size="large" :checked="payPlatform === PayPlatform.AliPay" @click="payPlatform = PayPlatform.AliPay">
            <div class="flex">
              <NImage
                class="w-[140px] h-[30px] ml-[10px]"
                preview-disabled
                :src="alipayLogoImg"
                alt="AliPay[支付宝]"
              />
            </div>
          </NRadio>
        </NCard>

        <NCard>
          <NRadio size="large" :checked="payPlatform === PayPlatform.Paddle" @click="payPlatform = PayPlatform.Paddle">
            <div class="flex">
              {{ t('common.paypalOrCreditCard') }}
            </div>
          </NRadio>
        </NCard>

        <!-- <NCard>
        <NRadio size="large" :checked="payPlatform === PayPlatform.WeiXinPay" @click="payPlatform = PayPlatform.WeiXinPay">
          <div class="flex">
            微信
          </div>
        </NRadio>
      </NCard> -->

        <NCard>
          {{ t('components.pay.waitOther') }}
        </NCard>
      </div>

      <div class="flex mt-[50px] justify-end">
        <div>
          <NButton type="success" :loading="getPayUrlLoading" @click="handlePay">
            {{ t('goToPay') }}
          </NButton>
        </div>
      </div>
    </NSpin>
  </div>
</template>
