<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getMultiple } from '@/api/system/systemVariable'
import { MarkdownRender } from '@/components/common'

interface Props {
  varName: string
}

const props = defineProps<Props>()
const content = ref('')
const loading = ref(false)

async function getPageSystemVariableConfig() {
  loading.value = true
  const keys = [props.varName]
  try {
    const { data } = await getMultiple<any>(keys)
    if (data[props.varName]) {
      content.value = data[props.varName]
    }
    else {
      content.value = 'load fail'
      loading.value = false
    }
  }
  catch (error) {
    loading.value = false
  }
}

onMounted(() => {
  getPageSystemVariableConfig()
})
</script>

<template>
  <MarkdownRender v-model:loading="loading" :content="content" mode="dark" />
</template>
