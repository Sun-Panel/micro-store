<script setup lang="ts">
import type { AppStatusInfo } from './composables/useAppInstallStatus'
import { NAlert, NButton, NCard, NTooltip, useMessage } from 'naive-ui'
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { getList as getListApi } from '@/api/microApp'
import defaultAppIcon from '@/assets/image_fail.png'
import { SvgIconOnline } from '@/components/common'
import { router } from '@/router'
import { isIframe } from '@/utils/cmn'
import { getAppDescByLang, getAppNameByLang, getBrowserLang, getLangMapFromAppInfo } from '@/utils/functions/lang'
import { useAppInstallStatus } from './composables/useAppInstallStatus'
import { getDownloadUrl } from './composables/useDownload'
import { useIframe } from './composables/useIframe'
import { useRouterHelper } from './composables/useRouterHelper'

interface MicroAppListItem extends MicroApp.Info {
  developer: MicroApp.DeveloperInfo
}

const message = useMessage()
const { sendInstallApp } = useIframe()
const { getDetailPath } = useRouterHelper('v1')
const { getAppButtonStatus } = useAppInstallStatus()

const list = ref<MicroAppListItem[]>([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const hasMore = computed(() => list.value.length < total.value)
const sentinelRef = ref<HTMLElement | null>(null)
const req = ref<MicroApp.GetListRequest>({
  page: 1,
  limit: 30,
  onlyWithVersion: true,
})

// 当前语言（浏览器语言不会变化，只计算一次）
const currentLang = getBrowserLang()

// 缓存每个应用的语言 Map，避免在模板中重复构建
function getLangMap(item: MicroAppListItem): Record<string, MicroApp.LangInfo> {
  return getLangMapFromAppInfo(item) as Record<string, MicroApp.LangInfo>
}

function getAppName(item: MicroAppListItem): string {
  return getAppNameByLang(getLangMap(item), currentLang, item.appName)
}

function getAppDesc(item: MicroAppListItem): string {
  return getAppDescByLang(getLangMap(item), currentLang, item.appDesc)
}

// 获取有效的应用图标，无效时返回默认图标
function getAppIcon(item: MicroAppListItem): string {
  return item.appIcon || defaultAppIcon
}

// 按钮状态缓存：避免模板中同一个 item 调用 getAppButtonStatus 多达 6 次
const buttonStatusCache = computed(() => {
  const cache = new Map<string, AppStatusInfo>()
  for (const item of list.value) {
    cache.set(item.microAppId, getAppButtonStatus(item.microAppId, item.latestVersion, item.latestLowVersion))
  }
  return cache
})

function getButtonStatus(item: MicroAppListItem): AppStatusInfo {
  return buttonStatusCache.value.get(item.microAppId) || { text: '安装', type: 'primary', disabled: false, isInstalled: false, hasUpdate: false, incompatible: false }
}

// 安装/打开按钮点击事件
async function handleInstall(item: MicroAppListItem) {
  if (isIframe()) {
    const url = await getDownloadUrl(item.microAppId)
    sendInstallApp({ microAppId: item.microAppId, url })
  }
  else {
    message.info('请在私有部署项目中打开微应用商店安装')
  }
}

// // 购买/卸载按钮点击事件（占位函数）
// function handleAction(_item: MicroAppListItem) {
//   // TODO: 实现购买/卸载逻辑
// }

// // 模拟模式控制：URL参数 ?mock=true 或 localStorage mock=true
// const isMockMode = ref(
//   new URLSearchParams(window.location.search).get('mock') === 'true'
//     || localStorage.getItem('mock') === 'true',
// )
//
// function toggleMockMode() {
//   isMockMode.value = !isMockMode.value
//   localStorage.setItem('mock', String(isMockMode.value))
//   // 切换模式后重新加载数据
//   page.value = 1
//   list.value = []
//   total.value = 0
//   nextTick(() => getList())
// }

// 模拟数据生成函数：生成足够的数据测试翻页（默认60条，每页30条可翻两页）
// function generateMockData(count = 60): MicroAppListItem[] {
//   const appNames = [
//     '智能日历', '天气助手', '记账本', '待办事项', '备忘录',
//     '计算器', '番茄时钟', '汇率换算', '二维码生成', '密码管理',
//     '文件管理', '音乐播放器', '视频播放器', '图片编辑', 'PDF阅读器',
//     '翻译工具', '笔记应用', '习惯追踪', '运动记录', '睡眠监测',
//     '食谱大全', '读书笔记', '新闻聚合', '天气预警', '倒数日',
//     '时间记录', '剪贴板', '屏幕录制', '截图工具', '网络检测',
//     '代码编辑器', 'Markdown编辑', '思维导图', '流程图', '白板工具',
//     '密码生成', '单位换算', '颜色选择', '字体管理', '图标库',
//     'API测试', 'JSON格式化', '正则测试', 'Base64编码', '时间戳转换',
//     '文本差异', 'CSV编辑器', '数据可视化', '图表工具', '仪表盘',
//     '任务看板', '日程管理', '会议安排', '团队协作', '项目管理',
//     '代码片段', 'API文档', '接口管理', '版本控制', '部署工具',
//   ]
//   const descriptions = [
//     '高效便捷的工具应用，提升工作效率',
//     '简洁实用的生活助手，让生活更简单',
//     '功能强大的专业工具，满足各种需求',
//     '轻松管理日常事务，随时随地访问',
//     '安全可靠的数据管理，保护您的隐私',
//   ]
//   const developers = [
//     { id: 1, name: '张三', developerName: '张三工作室', avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg' },
//     { id: 2, name: '李四', developerName: '李四科技', avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg' },
//     { id: 3, name: '王五', developerName: '王五软件', avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg' },
//     { id: 4, name: '赵六', developerName: '赵六开发', avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg' },
//     { id: 5, name: '孙七', developerName: '孙七创意', avatar: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg' },
//   ]
//   const categories = [1, 2, 3]
//   const chargeTypes = [0, 1, 2]

//   return Array.from({ length: count }, (_, i) => ({
//     id: i + 1,
//     microAppId: `app-${String(i + 1).padStart(3, '0')}`,
//     appName: appNames[i % appNames.length],
//     appIcon: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
//     appDesc: descriptions[i % descriptions.length],
//     developer: developers[i % developers.length],
//     developerId: developers[i % developers.length].id,
//     categoryId: categories[i % categories.length],
//     chargeType: chargeTypes[i % chargeTypes.length],
//     points: i % 3 === 0 ? 100 : 0,
//     status: 1,
//     downloadCount: Math.floor(Math.random() * 10000),
//     installCount: Math.floor(Math.random() * 5000),
//     latestVersion: '1.0.0',
//   }))
// }

// // 模拟分页数据
// const allMockData = generateMockData()

// function getMockPageData(page: number, limit: number) {
//   const start = (page - 1) * limit
//   const end = start + limit
//   return {
//     list: allMockData.slice(start, end),
//     count: allMockData.length,
//   }
// }

let observer: IntersectionObserver | null = null

async function getList(append = false) {
  if (loading.value)
    return
  loading.value = true
  const targetPage = append ? page.value + 1 : 1

  // // 模拟模式：直接使用本地模拟数据
  // if (isMockMode.value) {
  //   // 模拟网络延迟
  //   await new Promise(resolve => setTimeout(resolve, 300 + Math.random() * 500))
  //   const mockResult = getMockPageData(targetPage, req.value.limit)
  //   list.value = append ? [...list.value, ...mockResult.list] : mockResult.list
  //   total.value = mockResult.count
  //   page.value = targetPage
  //   loading.value = false
  //   updateObserver()
  //   return
  // }

  try {
    const { data } = await getListApi<Common.ListResponse<MicroAppListItem[]>>({ ...req.value, page: targetPage })
    list.value = append ? [...list.value, ...data.list] : data.list
    total.value = data.count
    page.value = targetPage
  }
  catch {
    // // API 调用失败时回退到模拟数据
    // console.warn('API 调用失败，使用模拟数据')
    // const mockResult = getMockPageData(targetPage, req.value.limit)
    // list.value = append ? [...list.value, ...mockResult.list] : mockResult.list
    // total.value = mockResult.count
    // page.value = targetPage
  }
  finally {
    loading.value = false
    updateObserver()
  }
}

// 滚动到底自动加载更多（观察底部哨兵元素，兼容 iframe 场景）
function updateObserver() {
  nextTick(() => {
    const el = sentinelRef.value
    if (!el)
      return
    if (observer)
      observer.disconnect()
    observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting && hasMore.value && !loading.value)
        getList(true)
    }, { rootMargin: '200px 0px' })
    observer.observe(el)
  })
}

function handleCardClick(item: MicroAppListItem) {
  router.push(getDetailPath(item.id))
}

onMounted(() => {
  getList()
})

onBeforeUnmount(() => {
  observer?.disconnect()
})
</script>

<template>
  <div class="min-h-screen p-2 bg-gradient-to-br from-slate-50 via-white to-slate-100 dark:from-slate-900 dark:via-slate-800 dark:to-slate-900">
    <NAlert type="info" :show-icon="false" class="mb-4 !rounded-lg">
      📣 感谢各位小伙伴一直以来的支持和关注，现官方版v2已经加入微应用框架，目前还在开发测试期间，每位开发者都可以开发应用提交供大家使用。待v2正式版上线应用商店将会有奖励机制，但暂未完全确定奖励形式。
      <span v-if="isIframe()" style="display: inline">
        <NButton tag="a" href="https://doc.sun-panel.top/v2/zh_cn/micro_app_dev" style="margin-left: 5px;font-weight: 900;" type="primary" target="_blank" text>
          官方微应用开发文档
        </NButton>

        <NButton tag="a" href="/" type="primary" target="_blank" text style="margin-left: 15px;font-weight: 900;">
          发布微应用
        </NButton>
      </span>
    </NAlert>

    <!-- 开发调试：模拟模式控制 -->
    <!-- <div class="mb-4 flex items-center gap-2">
      <NButton size="small" :type="isMockMode ? 'warning' : 'default'" @click="toggleMockMode">
        {{ isMockMode ? '模拟模式 ON' : '模拟模式 OFF' }}
      </NButton>
      <span v-if="isMockMode" class="text-xs text-orange-500">
        使用本地模拟数据（{{ allMockData.length }} 条，每页 {{ req.limit }} 条）
      </span>
    </div> -->

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
        <div class="flex items-start gap-3">
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
            </div>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <div class="font-bold m-0 leading-[1.4] truncate flex-1 min-w-0" :title="getAppName(item)">
                {{ getAppName(item) || 'Unknown' }}
              </div>
              <NTooltip v-if="isIframe() && getButtonStatus(item).incompatible" trigger="hover">
                <template #trigger>
                  <NButton
                    size="tiny"
                    type="default"
                    disabled
                    class="w-[60px] !h-[24px] flex-shrink-0"
                  >
                    {{ getButtonStatus(item).text }}
                  </NButton>
                </template>
                {{ getButtonStatus(item).incompatibleMsg }}
              </NTooltip>
              <NButton
                v-else-if="isIframe()"
                size="tiny"
                :type="getButtonStatus(item).type"
                :disabled="getButtonStatus(item).disabled"
                class="w-[60px] !h-[24px] flex-shrink-0"
                @click.stop="handleInstall(item)"
              >
                {{ getButtonStatus(item).text }}
              </NButton>
            </div>
            <div class="flex items-center gap-1 text-[11px] text-gray-500">
              <SvgIconOnline icon="boxicons:user" class="w-[15px] h-[15px]" />
              <span class="truncate">{{ item.developer?.name || item.developer?.developerName || '未知' }}</span>
            </div>
            <p class="text-xs text-gray-400 leading-[1.4] m-0 description-clamp" :title="getAppDesc(item)">
              {{ getAppDesc(item) || '-' }}
            </p>
          </div>
        </div>
      </NCard>
    </div>
    <div v-if="list.length > 0" ref="sentinelRef" class="flex flex-col items-center gap-2 py-4">
      <NButton v-if="hasMore" size="small" :loading="loading" @click="getList(true)">
        加载更多
      </NButton>
      <span v-else class="text-xs text-gray-400">— END —</span>
    </div>
  </div>
</template>

<style scoped>
.grid-layout {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 12px;
}

.description-clamp {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

@media (max-width: 600px) {
  .grid-layout {
    grid-template-columns: 1fr;
  }
}
</style>
