<script lang="ts" setup>
import { computed, nextTick, onMounted, ref } from 'vue'
import { NButton, useMessage } from 'naive-ui'
import { SvgIconOnline } from '@/components/common'
import { useAdminStore } from '@/store'
import { getHomeBase } from '@/api/openness'
import { router } from '@/router'
// interface Props {
//   // usingContext: boolean
// }

// interface Emit {
//   (ev: 'export'): void
//   (ev: 'toggleUsingContext'): void
// }

// defineProps<Props>()

// const emit = defineEmits<Emit>()

const adminStore = useAdminStore()
const ms = useMessage()
const collapsed = computed(() => adminStore.siderCollapsed)
const homeBase = ref<Openness.open.HomeBase>()

function handleUpdateCollapsed() {
  adminStore.setSiderCollapsed(!collapsed.value)
}

function onScrollToTop() {
  const scrollRef = document.querySelector('#scrollRef')
  if (scrollRef)
    nextTick(() => scrollRef.scrollTop = 0)
}

async function getHomeBasePost() {
  try {
    const { data } = await getHomeBase<Openness.open.HomeBase>()
    homeBase.value = data
    document.title = data.logo_text as string
  }
  catch (error) {
    ms.error('服务器出错了')
  }
}
// function handleExport() {
//   emit('export')
// }

// function toggleUsingContext() {
//   emit('toggleUsingContext')
// }

onMounted(() => {
  getHomeBasePost()
})
</script>

<template>
  <header
    class="sticky top-0 left-0 right-0 z-30 border-b dark:border-neutral-800 bg-white/80 dark:bg-black/20 backdrop-blur"
  >
    <div class="relative flex items-center justify-between min-w-0 overflow-hidden h-14">
      <div class="flex items-center">
        <button
          class="flex items-center justify-center w-11 h-11"
          @click="handleUpdateCollapsed"
        >
          <SvgIconOnline v-if="collapsed" class="text-2xl" icon="ri:align-justify" />
          <SvgIconOnline v-else class="text-2xl" icon="ri:align-right" />
        </button>
      </div>
      <h1
        class="flex-1 px-4 pr-6 overflow-hidden cursor-pointer select-none text-ellipsis whitespace-nowrap"
        @dblclick="onScrollToTop"
      >
        <div class="text-[18px] font-bold">
          {{ homeBase?.logo_text }}后台
        </div>
      </h1>

      <div class="mx-5" @click="router.push({ path: '/' })">
        <NButton size="small" strong secondary type="info">
          返回平台
        </NButton>
      </div>
    </div>
  </header>
</template>
