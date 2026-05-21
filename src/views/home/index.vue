<script setup lang="ts">
import { NButton, NCard, NEllipsis } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { getList as getListApi } from '@/api/microApp'
import defaultAppIcon from '@/assets/image_fail.png'

import { SvgIconOnline } from '@/components/common'
import { router } from '@/router'
import { useLocalAppStore } from '@/store'

import { isIframe } from '@/utils/cmn'

interface MicroAppListItem extends MicroApp.Info {
  // developerName: string
  developer: MicroApp.DeveloperInfo
}

const localAppStore = useLocalAppStore()

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

// 安装/打开按钮点击事件（占位函数）
function handleInstall(_item: MicroAppListItem) {
  // TODO: 实现安装/打开逻辑
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
  <div class="home-container">
    <div v-if="list.length === 0" class="empty-state">
      <img src="@/assets/image_fail.png" alt="empty" class="empty-icon">
      <p class="empty-text">
        {{ $t('home.emptyStateText') || '暂无应用' }}
      </p>
    </div>
    <div v-else class="grid-layout">
      <NCard
        v-for="item in list"
        :key="item.id"
        size="small"
        class="app-card"
        hoverable
        @click="handleCardClick(item)"
      >
        <div class="card-content">
          <div class="app-main">
            <div class="app-icon-wrapper">
              <img
                :src="getAppIcon(item)"
                :alt="getAppName(item)"
                class="app-icon-img"
              >
              <div class="app-stats">
                <span class="stat-item" :title="$t('common.downloadCount')">
                  <SvgIconOnline icon="grommet-icons:download" />
                  <span class="stat-value">{{ item.downloadCount || 0 }}</span>
                </span>
                <!-- <span class="stat-item" :title="$t('common.installCount')">
                  <SvgIconOnline icon="grommet-icons:install" />
                  <span class="stat-value">{{ item.installCount || 0 }}</span>
                </span> -->
              </div>
            </div>
            <div class="app-info">
              <NEllipsis class="app-name" :line-clamp="1">
                {{ getAppName(item) || 'Unknown' }}
              </NEllipsis>
              <p class="app-desc">
                {{ getAppDesc(item) || '暂无描述' }}
              </p>
            </div>
          </div>
          <div v-if="isIframe()" class="card-actions">
            <NButton
              size="tiny"
              :type="localAppStore.isInstalled(item.microAppId) ? 'default' : 'primary'"
              class="btn-install"
              @click.stop="handleInstall(item)"
            >
              {{ localAppStore.isInstalled(item.microAppId) ? '已安装' : '安装' }}
            </NButton>
            <!-- <NButton
              size="tiny"
              type="default"
              class="btn-action"
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
.home-container {
  padding: 16px;
}

.grid-layout {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 12px;
}

.app-card {
  border-radius: 8px;
  cursor: pointer;
  transition: box-shadow 0.2s ease;
}

.app-card:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.card-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-main {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.app-icon-wrapper {
  flex-shrink: 0;
  width: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.app-icon-img {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  object-fit: cover;
}

.app-info {
  flex: 1;
  min-width: 0;
}

.app-name {
  font-size: 15px;
  font-weight: 700;
  color: #333;
  margin-bottom: 4px;
}

.app-desc {
  font-size: 12px;
  color: #999;
  line-height: 1.4;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.app-stats {
  display: flex;
  justify-content: center;
  gap: 8px;
  font-size: 10px;
  color: #666;
  width: 100%;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 3px;
}

.stat-icon {
  font-size: 10px;
  color: #999;
}

.stat-value {
  font-weight: 500;
}

.card-actions {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex-shrink: 0;
}

.btn-install,
.btn-action {
  width: 60px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

.empty-icon {
  width: 64px;
  height: 64px;
  opacity: 0.4;
  margin-bottom: 16px;
}

.empty-text {
  color: #999;
  font-size: 14px;
}

@media (max-width: 600px) {
  .grid-layout {
    grid-template-columns: 1fr;
  }
}
</style>
