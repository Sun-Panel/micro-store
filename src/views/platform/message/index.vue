<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NButton, NEmpty, NH2, NTabPane, NTabs } from 'naive-ui'
import ListItemCard from './ListItemCard/index.vue'
import { getReceiveList as getReceiveListApi, getSendList as getSendListApi } from '@/api/system/message'

const receiveList = ref<Message.MessageInfo[]>([])
const sendList = ref<Message.MessageInfo[]>([])
const isLoading = ref(false)

async function getReceiveList() {
  try {
    const { data } = await getReceiveListApi <Common.ListResponse<Message.MessageInfo[]>>()
    receiveList.value = []
    receiveList.value = data.list
  }
  catch (error) {

  }
}

async function getSendList() {
  try {
    const { data } = await getSendListApi <Common.ListResponse<Message.MessageInfo[]>>()
    sendList.value = data.list
  }
  catch (error) {

  }
}

async function handleRefresh() {
  isLoading.value = true
  await getReceiveList()
  await getSendList()
  isLoading.value = false
}

onMounted(() => {
  handleRefresh()
})
</script>

<template>
  <div>
    <NH2 prefix="bar">
      <div class="flex items-center">
        <span>站内信</span>
        <span class="ml-5 flex justify-end">
          <NButton size="small" :loading="isLoading" tertiary type="info" @click="handleRefresh">
            刷新
          </NButton>
        </span>
      </div>
    </NH2>
    <NTabs type="line" animated>
      <NTabPane name="receive" tab="收信箱">
        <template v-if="receiveList.length === 0">
          <NEmpty description="这里什么都没有" />
        </template>
        <template v-else>
          <div v-for="item, index in receiveList" :key="index" class="mb-2">
            <ListItemCard :is-send="false" :message-info="item" @delete-done="handleRefresh" />
          </div>
        </template>
      </NTabPane>
      <NTabPane name="send" tab="发信箱">
        <template v-if="sendList.length === 0">
          <NEmpty description="这里什么都没有" />
        </template>
        <template v-else>
          <div v-for="item, index in sendList" :key="index" class="mb-2">
            <ListItemCard :is-send="true" :message-info="item" @delete-done="handleRefresh" />
          </div>
        </template>
      </NTabPane>
    </NTabs>
  </div>
</template>
