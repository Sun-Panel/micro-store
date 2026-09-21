<script setup lang="ts">
import { NAlert, NButton, NCard, NImage, NImageGroup, NModal, NSpace, NSpin, NTag } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { get } from '@/api/system/customCode'
import { CodeBlock, MarkdownRender, SvgIconOnline } from '@/components/common'
import { router } from '@/router'
import { getCurrentBaseUrlRoot } from '@/utils/cmn'
import { useAppStore } from '@/store/modules/app'

const route = useRoute()
const loading = ref(false)
const detail = ref<CustomCode.Detail>()
const appStore = useAppStore()
// 赞赏弹窗
const showReward = ref(false)
const markdownMode = computed(() => (appStore.theme === 'dark' ? 'dark' : 'light'))

// 全部展开 / 收起
const expandAll = ref(false)

const id = computed(() => Number(route.params.id ?? 0))

const codeTypeLabelMap: Record<number, string> = { 1: 'JS', 2: 'CSS', 3: '页脚' }
const versionLabelMap: Record<number, string> = { 1: 'v1', 2: 'v2' }

function versionTexts(versions?: number[]) {
  if (!versions?.length)
    return '-'
  return versions.map(item => versionLabelMap[item] ?? item).join(' / ')
}

// 复制到代码注释里的时间：作者发布/修改时间
function formatUpdateTime(value?: string) {
  if (!value)
    return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime()))
    return ''
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

// 复制到代码注释里的时间：作者发布/修改时间（Unix 秒，用于复制标识）
function publishedAtUnix(): number {
  const v = detail.value?.publishedAt
  if (!v)
    return 0
  const t = new Date(v).getTime()
  return Number.isNaN(t) ? 0 : Math.floor(t / 1000)
}

// 每个块单独构造复制元信息（含唯一标识，用于跨项目粘贴去重）
function copyMetaFor(block: CustomCode.Block): {
  title: string
  author: string
  updateTime: string
  url: string
  id: number
  onlyId: string
  updateTimeUnix: number
} {
  return {
    title: detail.value?.title ?? '',
    author: detail.value?.authorName ?? '',
    updateTime: formatUpdateTime(detail.value?.publishedAt),
    url: `${getCurrentBaseUrlRoot()}/customCode/${id.value}`,
    id: id.value,
    onlyId: block.onlyId ?? '',
    updateTimeUnix: publishedAtUnix(),
  }
}

async function loadDetail() {
  loading.value = true
  try {
    const { data } = await get<CustomCode.Detail>(id.value)
    detail.value = data
  }
  catch {
    detail.value = undefined
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  if (id.value)
    loadDetail()
})
</script>

<template>
  <div class="p-[20px] max-w-[900px] mx-auto">
    <NSpin :show="loading">
      <template v-if="detail">
        <h1 class="text-[22px] font-bold mb-[8px]">
          {{ detail.title }}
        </h1>

        <div class="flex items-center gap-[10px] text-[13px] text-slate-500 mb-[12px]">
          <span>{{ detail.authorName || '匿名' }}</span>
          <span>更新于 {{ formatUpdateTime(detail.publishedAt) }}</span>
          <span>阅读 {{ detail.readCount }}</span>
          <NTag size="small">
            适用版本 {{ versionTexts(detail.versions) }}
          </NTag>
          <NTag v-if="!detail.isOriginal" size="small" type="warning">
            非原创
          </NTag>
        </div>

        <div v-if="detail.keywords?.length" class="mb-[12px]">
          <NTag v-for="kw in detail.keywords" :key="kw" size="small" class="mr-[6px]">
            {{ kw }}
          </NTag>
        </div>

        <NAlert v-if="!detail.isOriginal" type="info" class="mb-[12px]">
          来源说明：{{ detail.sourceNote }}
        </NAlert>

        <NAlert type="warning" class="mb-[16px]" title="使用前请注意">
          以下内容由用户自行发布，第三方 JS / CSS 可能影响面板安全与稳定性，请自行判断后再使用。
        </NAlert>

        <p class="text-slate-600 mb-[12px]">
          {{ detail.description }}
        </p>

        <div class="flex items-center mb-[10px]">
          <span class="font-bold">代码片段（{{ detail.blocks?.length ?? 0 }}）</span>
          <NButton size="tiny" class="ml-auto" @click="expandAll = !expandAll">
            {{ expandAll ? '全部收起' : '全部展开' }}
          </NButton>
        </div>

        <div v-if="!detail.blocks?.length" class="text-slate-400">
          暂无代码片段
        </div>

        <NCard
          v-for="(block, index) in detail.blocks"
          :key="index"
          size="small"
          class="mb-[16px]"
        >
          <template #header>
            <!-- 友好标题：xxx代码片段（类型） -->
            <div class="font-bold">
              {{ codeTypeLabelMap[block.codeType] ?? block.codeType }}-代码片段{{ block.title ? `: ${block.title}` : '' }}
            </div>
          </template>

          <div v-if="block.note" class="text-slate-600 text-[14px] mb-[10px] note-pre">
            {{ block.note }}
          </div>

          <div v-if="block.images?.length" class="mb-[10px]">
            <NImageGroup>
              <div class="flex flex-wrap gap-[8px]">
                <NImage v-for="url in block.images" :key="url" :src="url" width="160" />
              </div>
            </NImageGroup>
          </div>

          <!-- 代码区默认折叠；复制按钮常驻，折叠时同样可用 -->
          <CodeBlock
            :key="`${index}-${expandAll}`"
            :code="block.code"
            :code-type="block.codeType"
            :default-collapsed="!expandAll"
            :copy-meta="copyMetaFor(block)"
          />
        </NCard>

        <NSpace class="mt-[20px]">
          <NButton @click="router.push({ name: 'CustomCodeList' })">
            返回列表
          </NButton>
          <NButton @click="showReward = true">
            <template #icon>
              <SvgIconOnline icon="ph:heart" />
            </template>
            赞赏
          </NButton>
        </NSpace>
      </template>
      <NAlert v-else-if="!loading" type="error">
        内容不存在或已下架
      </NAlert>
    </NSpin>
  </div>

  <!-- 赞赏弹窗 -->
  <NModal
    v-model:show="showReward"
    preset="card"
    title="赞赏作者"
    style="max-width: 640px;"
  >
    <div v-if="detail?.authorRewardContent" class="appreciate-content">
      <MarkdownRender
        :content="detail.authorRewardContent"
        :mode="markdownMode"
      />
    </div>
    <div v-else class="py-8 text-center text-slate-400 dark:text-slate-500">
      作者暂未设置赞赏信息
    </div>
  </NModal>
</template>

<style scoped>
.note-pre {
  white-space: pre-wrap;
}
</style>
