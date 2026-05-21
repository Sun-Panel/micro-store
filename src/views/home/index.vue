<script setup lang="ts">
import { NButton, NCard, NEllipsis, useMessage } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { getList as getListApi } from '@/api/microApp'
import defaultAppIcon from '@/assets/image_fail.png'

import { SvgIconOnline } from '@/components/common'
import { router } from '@/router'
import { useLocalAppStore } from '@/store'
import { isIframe } from '@/utils/cmn'
import { useIframe } from './composables/useIframe'

interface MicroAppListItem extends MicroApp.Info {
  // developerName: string
  developer: MicroApp.DeveloperInfo
}

const localAppStore = useLocalAppStore()
const message = useMessage()
const { sendInstallApp } = useIframe()

const list = ref<MicroAppListItem[]>([])
const req = ref<MicroApp.GetListRequest>({
  page: 1,
  limit: 10,
})

// ==================== 多语言处理 ====================
// 浏览器语言检测
function getBrowserLang(): string {
  const lang = navigator.language || (navigator as any).userLanguage || 'zh-CN'
  if (lang.startsWith('zh'))
    return 'zh-CN'
  if (lang.startsWith('en'))
    return 'en-US'
  if (lang.startsWith('ja'))
    return 'ja-JP'
  if (lang.startsWith('ko'))
    return 'ko-KR'
  return 'zh-CN'
}

// 当前语言
const currentLang = computed(() => getBrowserLang())

// 获取应用的多语言列表
// function getLangList(item: MicroAppListItem): string[] {
//   const langList = item.langList || []
//   if (langList.length > 0) {
//     return langList.map(l => l.lang)
//   }
//   return ['zh-CN']
// }

// 获取指定语言下的应用名称
function getAppName(item: MicroAppListItem): string {
  const langMap: Record<string, MicroApp.LangInfo> = {}
  const langList = item.langList || []
  langList.forEach((l) => {
    langMap[l.lang] = l
  })
  return langMap[currentLang.value]?.appName
    || langMap['zh-CN']?.appName
    || item.appName
    || ''
}

// 获取指定语言下的应用描述
function getAppDesc(item: MicroAppListItem): string {
  const langMap: Record<string, MicroApp.LangInfo> = {}
  const langList = item.langList || []
  langList.forEach((l) => {
    langMap[l.lang] = l
  })
  return langMap[currentLang.value]?.appDesc
    || langMap['zh-CN']?.appDesc
    || item.appDesc
    || ''
}

// 获取有效的应用图标，无效时返回默认图标
function getAppIcon(item: MicroAppListItem): string {
  return item.appIcon || defaultAppIcon
}

// 安装/打开按钮点击事件
function handleInstall(item: MicroAppListItem) {
  if (isIframe()) {
    // 在 iframe 中，发送安装消息给父窗口
    sendInstallApp({ microAppId: item.microAppId, authCode: '5555' })
  }
  else {
    // 非 iframe 环境，显示提示
    message.info('请在私有部署项目中打开微应用商店安装')
  }
}

// 购买/卸载按钮点击事件（占位函数）
function handleAction(_item: MicroAppListItem) {
  // TODO: 实现购买/卸载逻辑
}

// 模拟10条数据（当API调用失败时使用）
// const mockData: MicroApp.Info[] = [
//   {
//     id: 1,
//     microAppId: 'app-001',
//     appName: '智能日历',
//     appIcon: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     appDesc: '智能日历管理工具，支持日程提醒和智能排程',
//     developer: {
//       id: 1,
//       name: '张三',
//       avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     },
//     developerId: 1,
//     categoryId: 1,
//     chargeType: 0,
//     points: 0,
//     status: 1,
//   },
//   {
//     id: 2,
//     microAppId: 'app-002',
//     appName: '天气助手',
//     appIcon: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     appDesc: '实时天气查询，未来7天天气预报',
//     developer: {
//       id: 2,
//       name: '李四',
//       avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     },
//     developerId: 2,
//     categoryId: 1,
//     chargeType: 0,
//     points: 0,
//     status: 1,
//   },
//   {
//     id: 3,
//     microAppId: 'app-003',
//     appName: '记账本',
//     appIcon: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     appDesc: '简洁实用的个人记账应用，支持多种账本分类',
//     developer: {
//       id: 1,
//       name: '张三',
//       avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     },
//     developerId: 1,
//     categoryId: 2,
//     chargeType: 1,
//     points: 100,
//     status: 1,
//   },
//   {
//     id: 4,
//     microAppId: 'app-004',
//     appName: '待办事项',
//     appIcon: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     appDesc: '高效的待办事项管理工具，支持标签和优先级',
//     developer: {
//       id: 3,
//       name: '王五',
//       avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     },
//     developerId: 3,
//     categoryId: 2,
//     chargeType: 0,
//     points: 0,
//     status: 1,
//   },
//   {
//     id: 5,
//     microAppId: 'app-005',
//     appName: '备忘录',
//     appIcon: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     appDesc: '快速记录笔记和想法，支持富文本编辑',
//     developer: {
//       id: 2,
//       name: '李四',
//       avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     },
//     developerId: 2,
//     categoryId: 3,
//     chargeType: 2,
//     points: 0,
//     status: 1,
//   },
//   {
//     id: 6,
//     microAppId: 'app-006',
//     appName: '计算器',
//     appIcon: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     appDesc: '功能强大的科学计算器，支持历史记录',
//     developer: {
//       id: 4,
//       name: '赵六',
//       avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     },
//     developerId: 4,
//     categoryId: 1,
//     chargeType: 0,
//     points: 0,
//     status: 1,
//   },
//   {
//     id: 7,
//     microAppId: 'app-007',
//     appName: '番茄时钟',
//     appIcon: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     appDesc: '专注时间管理，番茄工作法工具',
//     developer: {
//       id: 3,
//       name: '王五',
//       avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     },
//     developerId: 3,
//     categoryId: 3,
//     chargeType: 0,
//     points: 0,
//     status: 1,
//   },
//   {
//     id: 8,
//     microAppId: 'app-008',
//     appName: '汇率换算',
//     appIcon: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     appDesc: '实时汇率查询，支持多种货币换算',
//     developer: {
//       id: 4,
//       name: '赵六',
//       avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     },
//     developerId: 4,
//     categoryId: 2,
//     chargeType: 0,
//     points: 0,
//     status: 1,
//   },
//   {
//     id: 9,
//     microAppId: 'app-009',
//     appName: '二维码生成',
//     appIcon: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     appDesc: '快速生成各类二维码，支持文本、链接等',
//     developer: {
//       id: 1,
//       name: '张三',
//       avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     },
//     developerId: 1,
//     categoryId: 1,
//     chargeType: 1,
//     points: 50,
//     status: 1,
//   },
//   {
//     id: 10,
//     microAppId: 'app-010',
//     appName: '密码管理',
//     appIcon: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     appDesc: '安全的密码管理工具，支持多平台同步',
//     developer: {
//       id: 2,
//       name: '李四',
//       avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     },
//     developerId: 2,
//     categoryId: 3,
//     chargeType: 2,
//     points: 0,
//     status: 1,
//   },
// ]

function getList() {
  getListApi<Common.ListResponse<MicroAppListItem[]>>(req.value).then(({ data }) => {
    list.value = data.list
  }).catch(() => {
    // API调用失败时使用模拟数据
    // list.value = mockData
  })
}

function handleCardClick(item: MicroAppListItem) {
  // console.log('item', item)
  // 可以在这里添加跳转逻辑
  router.push(`/microApp/${item.id}`)
}

onMounted(() => {
  getList()
  // list.value = mockData
})
</script>

<template>
  <div class="p-4">
    <div v-if="list.length === 0" class="flex flex-col items-center justify-center py-15 px-5">
      <img src="@/assets/image_fail.png" alt="empty" class="w-16 h-16 opacity-40 mb-4">
      <p class="text-gray-400 text-sm">
        {{ $t('home.emptyStateText') || '暂无应用' }}
      </p>
    </div>
    <div v-else class="grid-layout">
      <NCard
        v-for="item in list"
        :key="item.id"
        size="small"
        class="!rounded-lg cursor-pointer transition-shadow duration-200 hover:shadow-md"
        hoverable
        @click="handleCardClick(item)"
      >
        <div class="flex items-center gap-3">
          <div class="flex items-start gap-2.5 flex-1 min-w-0">
            <div class="flex-shrink-0 w-10 flex flex-col items-center gap-1">
              <img
                :src="getAppIcon(item)"
                :alt="getAppName(item)"
                class="w-10 h-10 rounded-lg object-cover"
              >
              <div class="flex justify-center gap-2 text-[10px] text-gray-500 w-full">
                <span class="flex items-center gap-[3px]" :title="$t('common.downloadCount')">
                  <SvgIconOnline icon="grommet-icons:download" />
                  <span class="font-medium">{{ item.downloadCount || 0 }}</span>
                </span>
                <!-- <span class="flex items-center gap-[3px]" :title="$t('common.installCount')">
                  <SvgIconOnline icon="grommet-icons:install" />
                  <span class="font-medium">{{ item.installCount || 0 }}</span>
                </span> -->
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="text-base font-bold m-0 leading-[1.4] truncate block w-full" :title="getAppName(item)">
                {{ getAppName(item) || 'Unknown' }}
              </div>
              <div class="flex items-center gap-1 text-[11px] text-gray-500">
                <SvgIconOnline icon="boxicons:user" class="w-[15px] h-[15px]" />
                <span class="truncate">{{ item.developer?.name || item.developer?.developerName || '未知' }}</span>
              </div>
              <p class="text-xs text-gray-400 leading-[1.4] m-0 line-clamp-2">
                {{ getAppDesc(item) || '暂无描述' }}
              </p>
            </div>
          </div>
          <div v-if="isIframe()" class="flex flex-col gap-1 flex-shrink-0">
            <NButton
              size="tiny"
              :type="localAppStore.isInstalled(item.microAppId) ? 'default' : 'primary'"
              class="w-[60px]"
              @click.stop="handleInstall(item)"
            >
              {{ localAppStore.isInstalled(item.microAppId) ? '已安装' : '安装' }}
            </NButton>
            <!-- <NButton
              size="tiny"
              type="default"
              class="w-[60px]"
              @click.stop="handleAction(item)"
            >
              {{ localAppStore.isInstalled(item.microAppId) ? '卸载' : '详情' }}
            </NButton> -->
          </div>
        </div>
      </NCard>
    </div>
  </div>
</template>

<style scoped>
.grid-layout {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 12px;
}

@media (max-width: 600px) {
  .grid-layout {
    grid-template-columns: 1fr;
  }
}
</style>
