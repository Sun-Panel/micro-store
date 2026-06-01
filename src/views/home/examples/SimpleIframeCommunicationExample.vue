<script setup lang="ts">
/**
 * 简单的 iframe 通信示例
 * 展示如何在微应用商店 iframe 中使用 useIframeCommunication composable
 */

import { onMounted, ref } from 'vue'
import { MICRO_APP_EVENTS, MICRO_APP_STORE_EVENTS } from '../composables/iframeEvents'
import { useIframeCommunication } from '../composables/useIframeCommunication'

const communicationStatus = ref('未连接')
const lastMessage = ref<any>(null)

// 初始化 iframe 通信
const { on, once, onQuickRequest, sendMessage, sendRequest, reinit, destroy } = useIframeCommunication({
  targetOrigin: '*', // 生产环境建议指定具体 origin
  sourceId: 'microAppStore',
  isIframe: true,
  debug: true,
  autoDestroy: true,
})

// 监听主面板发送的应用安装事件
on(MICRO_APP_STORE_EVENTS.INSTALL_APP, (data) => {
  communicationStatus.value = '已连接'
  lastMessage.value = { type: 'install', data }
  // console.log('收到安装应用事件:', data)
})

// 监听一次性的通信就绪事件
once(MICRO_APP_EVENTS.COMMUNICATION_READY, () => {
  communicationStatus.value = '通信就绪'
  // console.log('通信已就绪')
})

// 注册快捷回复处理器
onQuickRequest(MICRO_APP_STORE_EVENTS.GET_APP_LIST, (data, ctx) => {
  // console.log('收到获取应用列表请求:', data)
  // 这里可以调用本地 API 获取应用列表
  ctx.reply({
    list: [
      { id: 1, appName: '示例应用1' },
      { id: 2, appName: '示例应用2' },
    ],
  })
})

// 向主面板发送通信就绪事件
onMounted(() => {
  sendMessage('main', MICRO_APP_EVENTS.COMMUNICATION_READY)
})

// 发送请求示例
async function requestUserInfo() {
  try {
    await sendRequest('main', 'getUserInfo', { userId: 123 })
    // console.log('获取用户信息成功:', userInfo)
  }
  catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

// 重新初始化示例（例如更换 postMessageKey）
function reinitialize() {
  reinit({
    targetOrigin: '*',
    sourceId: 'microAppStore',
    isIframe: true,
    debug: true,
    autoDestroy: true,
  })
}
</script>

<template>
  <div class="p-4">
    <h2 class="text-lg font-bold mb-4">
      iframe 通信示例
    </h2>
    <div class="mb-4">
      <p>通信状态: <span class="font-mono">{{ communicationStatus }}</span></p>
      <p v-if="lastMessage">
        最后消息: <span class="font-mono">{{ JSON.stringify(lastMessage) }}</span>
      </p>
    </div>
    <div class="flex gap-2">
      <button
        class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        @click="requestUserInfo"
      >
        请求用户信息
      </button>
      <button
        class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
        @click="reinitialize"
      >
        重新初始化
      </button>
      <button
        class="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
        @click="destroy"
      >
        销毁通信
      </button>
    </div>
  </div>
</template>
