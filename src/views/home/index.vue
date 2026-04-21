<script setup lang="ts">
import { NAvatar, NCard } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { getList as getListApi } from '@/api/microApp'
import { router } from '@/router'

interface MicroAppListItem extends MicroApp.Info {
  // developerName: string
  developer: MicroApp.DeveloperInfo
}

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
function getLangList(item: MicroAppListItem): string[] {
  const langList = item.langList || []
  if (langList.length > 0) {
    return langList.map(l => l.lang)
  }
  return ['zh-CN']
}

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
  <div class="home-container">
    <div class="grid-layout">
      <NCard
        v-for="item in list"
        :key="item.id"
        size="small"
        class="app-card"
        hoverable
        @click="handleCardClick(item)"
      >
        <div class="flex items-center gap-2">
          <NAvatar
            :size="50"
            :style="{
              backgroundColor: 'transparent',
            }"
            :src="item.appIcon || '-'"
          />
          <div class="flex flex-col">
            <div class="text-lg font-medium">
              {{ getAppName(item) || 'Unknown' }}
            </div>
            <div class="text-sm text-gray-500">
              作者：{{ item.developer?.name || 'Unknown' }}
            </div>
          </div>
        </div>
        <div class="text-sm mt-2 text-gray-600">
          {{ getAppDesc(item) || '-' }}
        </div>
      </NCard>
    </div>
  </div>
</template>

<style scoped>
.home-container {
  padding: 16px;
}

.grid-layout {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  justify-content: flex-start;
}

.app-card {
  flex: 0 0 auto;
  width: 280px;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.app-card:hover {
  transform: translateY(-2px);
}

/* 响应式：根据容器宽度自动调整 */
@media (min-width: 1200px) {
  .app-card {
    width: calc((100% - 48px) / 4); /* 4列 */
  }
}

@media (min-width: 900px) and (max-width: 1199px) {
  .app-card {
    width: calc((100% - 32px) / 3); /* 3列 */
  }
}

@media (min-width: 600px) and (max-width: 899px) {
  .app-card {
    width: calc((100% - 16px) / 2); /* 2列 */
  }
}

@media (max-width: 599px) {
  .app-card {
    width: 100%; /* 1列 */
  }
}
</style>
