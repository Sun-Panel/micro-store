<script setup lang="ts">
import { defineProps, ref } from 'vue'

import { NButton, NCard, NH4, NModal, NPopover, NTag } from 'naive-ui'
import { getPayPlatformText } from '@/utils/cmn'
import { OrderStatus } from '@/enums/goodsOrder'
import { GoodsInfo } from '@/views/components'
import { SvgIconOnline } from '@/components/common'
import { t } from '@/locales'

interface Props {
  orderInfo: GoodsOrder.Info | null
}

const props = defineProps<Props>()
const goodsInfoModalShow = ref(false)

function getPayStatusTagTypeName(status: number) {
  switch (status) {
    case OrderStatus.PAY_SUCCESS:
      return 'info'

    case OrderStatus.FINISH:
      return 'success'

    default:
      return ''
  }
}

function getPayStatusText(code: number): string {
  switch (code) {
    case OrderStatus.CLOSE:
      return t('order.status.close')

    case OrderStatus.PAY_WAIT:
      return t('order.status.payWait')

    case OrderStatus.PAY_SUCCESS:
      return t('order.status.paySuccess')

    case OrderStatus.FINISH:
      return t('order.status.complete')

    default:
      return t('order.status.unknown')
  }
}

function openPayLink() {
  window.open(props.orderInfo?.payUrl)
}
</script>

<template>
  <div>
    <div v-if="orderInfo?.user" class="mt-5">
      <NH4 prefix="bar" style="margin-bottom: 5px;" align-text>
        {{ t('order.userInfo') }}
      </NH4>
      <div>
        <div>
          {{ t('common.username') }} : {{ orderInfo?.user.username }}
        </div>
        <div>
          {{ t('common.nikeName') }} : {{ orderInfo?.user.name }}
        </div>
      </div>
    </div>
    <div class="mt-5">
      <NH4 prefix="bar" style="margin-bottom: 5px;" align-text>
        {{ t('order.goodsInfo') }}
      </NH4>
      <NCard>
        <div>
          <span>
            {{ t('order.goodsTitle') }} : {{ orderInfo?.goodsSnapshotInfo?.title }}
          </span>
          <span class="ml-2">

            <NButton type="tertiary" size="tiny" @click="goodsInfoModalShow = true">
              {{ t('order.goodsSnapshot') }}
            </NButton>
          </span>
        </div>
        <div v-if="orderInfo?.goodsSnapshotInfo?.discount">
          {{ t('order.goodsActivityLabel') }} :
          <NTag size="small" :bordered="false" type="error">
            {{ orderInfo?.goodsSnapshotInfo.discount }}
          </NTag>
        </div>
        <div>
          {{ t('common.amount') }} : (CNY) {{ orderInfo?.goodsSnapshotInfo?.price }}
        </div>
      </NCard>
    </div>

    <div class="mt-5">
      <NH4 prefix="bar" style="margin-bottom: 5px;" align-text>
        {{ t('order.orderInfo') }}
      </NH4>
      <NCard>
        <div class="detail-box">
          <div>
            {{ t('order.orderStatus') }} :
            <NTag size="small" :type="getPayStatusTagTypeName(orderInfo?.status || 0)">
              {{ getPayStatusText(orderInfo?.status || 0) }}
            </NTag>
            <span v-if="orderInfo?.status === OrderStatus.PAY_WAIT && orderInfo.payUrl" class="ml-5">
              <NButton size="tiny" type="info" @click="openPayLink">{{ t('order.openThePaymentPage') }}</NButton>
            </span>
          </div>

          <div class="flex items-center">
            <span>
              {{ t('order.payAmount') }}  : ({{ orderInfo?.currencyCode || 'CNY' }}) {{ orderInfo?.countPrice }}
            </span>
            <span class="ml-2">
              <NPopover trigger="hover">
                <template #trigger>
                  <SvgIconOnline icon="akar-icons:question" />
                </template>
                <div class="text-xs">{{ t("order.serviceFeeTip") }}</div>
              </NPopover>
            </span>
          </div>

          <div>
            {{ t('order.createOrderTime') }} : {{ orderInfo?.createTime }}
          </div>

          <div>
            {{ t('order.paySuccessTime') }} : {{ orderInfo?.payTime }}
          </div>

          <div>
            {{ t('order.orderNo') }} : {{ orderInfo?.orderNo }}
          </div>

          <div>
            {{ t('order.payPlatform') }} :
            <NTag>
              {{ getPayPlatformText(orderInfo?.payPlatform || 0) }}
            </NTag>
          </div>

          <div>
            {{ t('order.payPlatform') }} -  {{ t('order.orderNo') }} : {{ orderInfo?.payPlatformOrderNo }}
          </div>
        </div>
      </NCard>
    </div>

    <NModal v-model:show="goodsInfoModalShow" preset="card" style="width: 600px" :title="t('order.goodsSnapshot')">
      <GoodsInfo :goods-info="orderInfo?.goodsSnapshotInfo as Goods.GoodsInfo" />
    </NModal>
  </div>
</template>

<style scoped>
.detail-box > div {
  margin-top: 10px;
}
</style>
