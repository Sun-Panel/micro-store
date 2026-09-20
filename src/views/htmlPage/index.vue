<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { HtmlRender, MessageEditor } from '@/components/common'
import { get } from '@/api/system/htmlPage'

const route = useRoute()
const pageConfig = ref<HtmlPage.Info>()
const loading = ref(false)

// 快速切换页面时，用递增序号丢弃已过期的请求，避免旧响应覆盖新内容
let latestRequestId = 0

async function getPageConfig() {
  const pageName = route.params.p as string
  if (!pageName) {
    pageConfig.value = { content: '404 Page does not exist.' }
    return
  }

  const requestId = ++latestRequestId
  loading.value = true
  try {
    const { data } = await get<HtmlPage.Info>(pageName)
    if (requestId === latestRequestId)
      pageConfig.value = data
  }
  catch {
    if (requestId === latestRequestId)
      pageConfig.value = { content: '404 Page does not exist.' }
  }
  finally {
    if (requestId === latestRequestId)
      loading.value = false
  }
}

watch(() => route.params.p, () => {
  getPageConfig()
})

onMounted(() => {
  getPageConfig()
})
</script>

<template>
  <div>
    <template v-if="!loading && pageConfig?.messageTemplateFlag && pageConfig.messageTemplatePosition === 'top'">
      <MessageEditor :flags="[pageConfig?.messageTemplateFlag]" :show-select="false" />
    </template>
    <HtmlRender :content="pageConfig?.content" :loading="loading" />
    <template v-if="!loading && pageConfig?.messageTemplateFlag && pageConfig.messageTemplatePosition === 'bottom'">
      <MessageEditor :flags="[pageConfig?.messageTemplateFlag]" :show-select="false" />
    </template>
  </div>
</template>
