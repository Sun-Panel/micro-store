<script setup lang="ts">
import { computed, defineProps, ref } from 'vue'

import type { DataTableColumns } from 'naive-ui'
import { NDataTable } from 'naive-ui'
import { buildTimeString } from '@/utils/cmn'

const props = defineProps<{
  dataList: ProAuthorize.GetAuthorizeHistoryRecordResp[]
}>()

const proAuthorizeHistoryList = computed(() => props.dataList)
const tableIsLoading = ref(false)

const columns: DataTableColumns<ProAuthorize.GetAuthorizeHistoryRecordResp> = [
  {
    title: '发生时间',
    key: 'changeTime',
    render(row) {
      return buildTimeString(row.changeTime, 'YYYY-MM-DD HH:mm')
    },
  },
  {
    title: '变动天数',
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
    title: '备注',
    key: 'note',
    render(row) {
      if (row.orderNo !== '')
        return '- 付费购买 -'

      return row.note
    },
  },
  {
    title: '过期时间',
    key: 'expiredTime',
    render(row) {
      return buildTimeString(row.expiredTime, 'YYYY-MM-DD HH:mm')
    },
  },
  {
    title: '订单编号',
    key: 'orderNo',
  },
  {
    title: '管理员备注',
    key: 'adminNote',
  },
]
</script>

<template>
  <div>
    <NDataTable
      :columns="columns"
      :data="proAuthorizeHistoryList"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
    />
  </div>
</template>
