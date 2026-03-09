<script setup lang="ts">
import { onMounted, ref } from 'vue'
import List from './list.vue'
import { getAuthorizeHistoryRecordByUserId } from '@/api/admin/proAuthorize'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'

interface Props {
  userId: number
}
const props = defineProps<Props>()
const list = ref<ProAuthorize.GetAuthorizeHistoryRecordResp[]>([])

async function getList() {
  await getAuthorizeHistoryRecordByUserId<Common.ListResponse<ProAuthorize.GetAuthorizeHistoryRecordResp[]>>(props.userId).then(({ data }) => {
    list.value = data.list
  }).then((error) => {
    apiRespErrMsgAndCustomCodeNeg1Msg(error)
  })
}

onMounted(() => {
  // console.log(props.userId)
  getList()
})
</script>

<template>
  <div>
    <!-- <NH2>授权历史</NH2> -->
    <List :data-list="list" />
  </div>
</template>
