<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { MarkdownRender, MessageEditor } from '@/components/common'
import { getMultiple } from '@/api/system/systemVariable'

interface Data {
  donate_reward_description: string
}

const describeData = ref<Data>()
const loading = ref(false)

async function getDescribeData() {
  loading.value = true
  try {
    const { data } = await getMultiple<Data>(['donate_reward_description'])
    describeData.value = data

    // loading.value = false
  }
  catch (error) {
    loading.value = false
  }
}

onMounted(() => {
  getDescribeData()
})
</script>

<template>
  <div class="m-10">
    <div class="mb-10">
      <MarkdownRender v-model:loading="loading" :content="describeData?.donate_reward_description as string" mode="light" />
    </div>

    <div>
      <MessageEditor :flags="['sponsor_reward_cn']" :show-select="false" />
    </div>
  </div>
</template>
