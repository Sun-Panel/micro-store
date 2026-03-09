<script setup lang="ts">
import { onMounted, ref } from 'vue'

import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NH2, NModal, NSpin, NTag } from 'naive-ui'
import RedeemCode from './RedeemCode/index.vue'
import { buildTimeString } from '@/utils/cmn'
import { getAuthorize as getAuthorizeApi, getAuthorizeHistoryRecord as getAuthorizeHistoryRecordApi } from '@/api/proAuthorize'
import { t } from '@/locales'

const proAuthorizeExpiredTime = ref('')
const proAuthorizeHistoryList = ref<ProAuthorize.GetAuthorizeHistoryRecordResp[]>([])
const tableIsLoading = ref(false)
const showHistoryList = ref(false)
const infoLoading = ref(false)

const redeemCodeModalShow = ref(false)

const columns: DataTableColumns<ProAuthorize.GetAuthorizeHistoryRecordResp> = [
  {
    title: t('common.time'),
    key: 'changeTime',
    render(row) {
      return buildTimeString(row.changeTime, 'YYYY-MM-DD HH:mm')
    },
  },
  {
    title: t('authorize.changeDays'),
    key: 'dayNum',
    render(row) {
      let day: string
      if (row.dayNum > 0)
        day = `+${row.dayNum}`
      else
        day = `${row.dayNum}`
      return day
    },
  },
  {
    title: t('common.note'),
    key: 'note',
    render(row) {
      if (row.orderNo !== '')
        return `- ${t('common.PayBuy')} -`

      return row.note
    },
  },
  {
    title: t('common.expiredTime'),
    key: 'expiredTime',
    render(row) {
      return buildTimeString(row.expiredTime, 'YYYY-MM-DD HH:mm')
    },
  },
  {
    title: t('order.orderNo'),
    key: 'orderNo',
  },
]

async function getAuthorize() {
  infoLoading.value = true
  try {
    const { data } = await getAuthorizeApi<ProAuthorize.GetAuthorizeResp>()
    proAuthorizeExpiredTime.value = data.expiredTime
    infoLoading.value = false // 避免闪烁 等完整加载出来在显示
  }
  catch (error) {
    console.error(error)
  }
}

async function getAuthorizeHistoryRecord() {
  tableIsLoading.value = true
  try {
    const { data } = await getAuthorizeHistoryRecordApi<Common.ListResponse<ProAuthorize.GetAuthorizeHistoryRecordResp[]>>()
    proAuthorizeHistoryList.value = data.list
  }
  catch (error) {
    console.error(error)
  }
  tableIsLoading.value = false
}

function handleLoadHistroy() {
  showHistoryList.value = true
  getAuthorizeHistoryRecord()
}

function handleClickRedeemCode() {
  redeemCodeModalShow.value = true
}

function handleRedeemCodeSuccess() {
  getAuthorize()
  handleLoadHistroy()
  redeemCodeModalShow.value = false
}

onMounted(() => {
  getAuthorize()
})
</script>

<template>
  <div>
    <template v-if="infoLoading">
      <div class="flex justify-center">
        <NSpin />
      </div>
    </template>
    <template v-else>
      <NH2 align-text prefix="bar">
        <div class="flex">
          {{ t('authorize.authorizeInfo') }}

          <div class="ml-5">
            <NButton size="small" strong secondary type="info" @click="handleClickRedeemCode">
              {{ t('authorize.redemptionCodeEntrance') }}
            </NButton>
          </div>
        </div>
      </NH2>

      <NCard>
        <div>
          <NTag type="warning">
            PRO
          </NTag>
          {{ t('common.expiredTime') }}: {{ proAuthorizeExpiredTime ? buildTimeString(proAuthorizeExpiredTime,
                                                                                      'YYYY-MM-DD HH:mm') : "-" }}
        </div>
      </NCard>

      <NH2 align-text prefix="bar">
        {{ t('common.historyRecord') }}
      </NH2>

      <template v-if="!showHistoryList">
        <NButton @click="handleLoadHistroy">
          {{ t('authorize.loadHistoryRecord') }}
        </NButton>
      </template>
      <template v-if="showHistoryList">
        <NDataTable
          :columns="columns" :data="proAuthorizeHistoryList" :bordered="false" :loading="tableIsLoading"
          :remote="true"
        />
      </template>
    </template>

    <NModal
      v-model:show="redeemCodeModalShow" preset="card" style="max-width: 500px"
      :title="t('authorize.proRedemptionCode') "
    >
      <RedeemCode @success="handleRedeemCodeSuccess" />
    </NModal>
  </div>
</template>
