<script setup lang="ts">
import { onMounted } from 'vue'
import { NButton, NH3 } from 'naive-ui'
import { router } from '@/router'
import { useAuthStore } from '@/store'

const appStore = useAuthStore()

// interface RootPageDescription {
//   root_page_description: string
// }

// const rootPageDescription = ref<string>('')
// const loading = ref(false)
// // async function getVar() {
// //   loading.value = true
// //   await getRootPageDescription<RootPageDescription>().then(({ data }) => {
// //     rootPageDescription.value = data.root_page_description
// //   // loading.value = false
// //   }).catch(() => {
// //     rootPageDescription.value = 'load fail'
// //     loading.value = false
// //   })
// // }

onMounted(() => {
  // getVar()
})
</script>

<template>
  <div class="flex justify-center">
    <NH3 v-if="!appStore.userInfo?.username">
      普通账号请使用授权登录，管理员请使用账号密码登录
    </NH3>
    <div v-else>
      <NButton @click="router.push('/platform/proAuthorize')">
        查看我的授权信息
      </NButton>
    </div>
  </div>
</template>
