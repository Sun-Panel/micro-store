<script setup lang="ts">
import { NButton } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { router } from '@/router'
import { getProDescription } from '@/api/openness'

import { MarkdownRender } from '@/components/common'

interface ProDescription {
  pro_func_description: string
  pro_before_buy_description: string
}

const proDescription = ref<ProDescription>()
const loading = ref(false)

// 立即升级
function handleUpdatePro() {
  router.push({ name: 'PlatformPackageStore' })
}

function handleGoToDonateReward() {
  router.push({ name: 'DonateReward' })
}

async function getProDescriptionPost() {
  loading.value = true
  try {
    const { data } = await getProDescription<ProDescription>()
    proDescription.value = data
    // loading.value = false
  }
  catch (error) {
    loading.value = false
  }
}

onMounted(() => {
  getProDescriptionPost()
})
</script>

<template>
  <div>
    <div>
      <MarkdownRender v-model:loading="loading" :content="proDescription?.pro_before_buy_description as string" mode="light" />
    </div>

    <div class="my-5">
      <NButton color="#8a2be2" @click="handleUpdatePro">
        立即升级PRO
      </NButton>

      <NButton style="margin-left: 5px;" @click="handleGoToDonateReward">
        对曾经打赏的用户回馈授权
      </NButton>
    </div>

    <div>
      <MarkdownRender v-model:loading="loading" :content="proDescription?.pro_func_description as string" mode="light" />
    </div>

    <div class="my-5">
      <NButton color="#8a2be2" @click="handleUpdatePro">
        立即升级PRO
      </NButton>
    </div>
  </div>
</template>
