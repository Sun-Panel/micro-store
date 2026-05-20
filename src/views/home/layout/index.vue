<script setup lang="ts">
import { NLoadingBarProvider } from 'naive-ui'
import { onMounted } from 'vue'
import { useAuthStore } from '@/store'
import { useIframe } from '../composables/useIframe'
import Header from './Header.vue'

const authStore = useAuthStore()
const { sendLoginEvent, sendGetInstalledApps } = useIframe()
// const loadingBar = useLoadingBar()

function isIframe() {
  return window !== window.top
}

onMounted(() => {
  // loadingBar.start()
  if (authStore.userInfo?.id)
    authStore.refreshUserInfo()

  // 在 iframe 中发送登录事件
  if (isIframe()) {
    sendGetInstalledApps()
    sendLoginEvent()
  }
})
</script>

<template>
  <div>
    <div class="fixed top-0 w-full z-10 bg-slate-100">
      <Header />
    </div>
    <NLoadingBarProvider>
      <div class="max-w-[1200px] mx-auto" :class="isIframe() ? '' : 'mt-[60px]'">
        <router-view class="p-[20px]" />
      </div>
    </NLoadingBarProvider>
  </div>
</template>
