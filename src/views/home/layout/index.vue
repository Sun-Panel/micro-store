<script setup lang="ts">
import { onMounted } from 'vue'
import { NLoadingBarProvider } from 'naive-ui'
import Header from './Header.vue'
import { useAuthStore } from '@/store'

const authStore = useAuthStore()
// const loadingBar = useLoadingBar()

function isIframe() {
  return window !== window.top
}

onMounted(() => {
  // loadingBar.start()
  if (authStore.userInfo?.username)
    authStore.refreshUserInfo()
})
</script>

<template>
  <div>
    <div v-show="!isIframe()" class="fixed top-0 w-full z-10 bg-slate-100">
      <Header />
    </div>
    <NLoadingBarProvider>
      <div class="max-w-[1200px] mx-auto" :class="isIframe() ? '' : 'mt-[60px]'">
        <router-view class="p-[20px]" />
      </div>
    </NLoadingBarProvider>
  </div>
</template>
