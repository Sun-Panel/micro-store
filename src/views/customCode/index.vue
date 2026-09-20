<script setup lang="ts">
import { NButton, NButtonGroup, NEmpty, NInput, NSelect, NSpin, NTag } from 'naive-ui'
import { onMounted, reactive, ref } from 'vue'
import { getList as getListApi } from '@/api/system/customCode'
import { router } from '@/router'

const loading = ref(false)
const list = ref<CustomCode.ListItem[]>([])
const keyWord = ref('')
const sortBy = ref('published_at')
const version = ref(0)
const codeType = ref(0)

const versionOptions = [
  { label: '全部版本', value: 0 },
  { label: 'v1', value: 1 },
  { label: 'v2', value: 2 },
]

const codeTypeOptions = [
  { label: '全部类型', value: 0 },
  { label: 'JS', value: 1 },
  { label: 'CSS', value: 2 },
  { label: '页脚', value: 3 },
]

const pagination = reactive({
  page: 1,
  pageSize: 12,
  itemCount: 0,
})

const codeTypeMap: Record<number, string> = { 1: 'JS', 2: 'CSS', 3: '页脚' }

function codeTypeTexts(codeTypes?: number[]) {
  if (!codeTypes?.length)
    return ''
  return codeTypes.map(item => codeTypeMap[item] ?? item).join(' / ')
}

function formatDate(value?: string) {
  if (!value)
    return ''
  return value.slice(0, 10)
}

async function getList(page = pagination.page) {
  loading.value = true
  try {
    const { data } = await getListApi<Common.ListResponse<CustomCode.ListItem[]>>({
      page,
      limit: pagination.pageSize,
      keyword: keyWord.value,
      sortBy: sortBy.value,
      sortOrder: 'desc',
      version: version.value,
      codeType: codeType.value,
    })
    pagination.itemCount = data.count
    pagination.page = page
    list.value = data.list ?? []
  }
  finally {
    loading.value = false
  }
}

function handleSearch() {
  getList(1)
}

function handleSort(value: string) {
  sortBy.value = value
  getList(1)
}

function goDetail(id: number) {
  router.push({ name: 'CustomCodeDetail', params: { id } })
}

onMounted(() => {
  getList(1)
})
</script>

<template>
  <div class="p-[20px] max-w-[1100px] mx-auto">
    <div class="flex items-center gap-[10px] mb-[16px]">
      <NInput v-model:value="keyWord" placeholder="搜索标题、描述或关键词" clearable style="max-width: 360px;" @keyup.enter="handleSearch" />
      <NButton type="primary" @click="handleSearch">
        搜索
      </NButton>

      <NSelect v-model:value="version" :options="versionOptions" style="width: 120px;" @update:value="getList(1)" />
      <NSelect v-model:value="codeType" :options="codeTypeOptions" style="width: 130px;" @update:value="getList(1)" />

      <div class="ml-auto flex gap-[8px]">
        <NButtonGroup>
          <NButton :type="sortBy === 'published_at' ? 'primary' : 'default'" size="small" @click="handleSort('published_at')">
            最新
          </NButton>
          <NButton :type="sortBy === 'read_count' ? 'primary' : 'default'" size="small" @click="handleSort('read_count')">
            最热
          </NButton>
        </NButtonGroup>
      </div>
    </div>

    <NSpin :show="loading">
      <div v-if="list.length" class="grid gap-[16px] sm:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="item in list"
          :key="item.id"
          class="p-[16px] rounded-[10px] border border-slate-200 dark:border-slate-700 cursor-pointer hover:shadow-md transition"
          @click="goDetail(item.id)"
        >
          <div class="flex items-center gap-[6px] mb-[8px]">
            <div class="font-bold text-[16px] truncate">
              {{ item.title }}
            </div>
            <NTag v-if="item.codeTypes?.length" size="small">
              {{ codeTypeTexts(item.codeTypes) }}
            </NTag>
          </div>
          <p class="text-slate-500 text-[13px] desc-clamp mb-[10px]">
            {{ item.description }}
          </p>
          <div v-if="item.keywords?.length" class="mb-[10px]">
            <NTag v-for="kw in item.keywords.slice(0, 3)" :key="kw" size="small" class="mr-[4px]">
              {{ kw }}
            </NTag>
          </div>
          <div class="flex items-center text-[12px] text-slate-400 gap-[10px]">
            <span>{{ item.authorName || '匿名' }}</span>
            <span>阅读 {{ item.readCount }}</span>
            <span class="ml-auto">{{ formatDate(item.publishedAt) }}</span>
          </div>
        </div>
      </div>
      <NEmpty v-else description="暂无内容" />

      <div v-if="pagination.itemCount > pagination.pageSize" class="flex justify-center mt-[20px] gap-[8px]">
        <NButton size="small" :disabled="pagination.page <= 1" @click="getList(pagination.page - 1)">
          上一页
        </NButton>
        <span class="leading-[28px] text-[13px] text-slate-500">
          {{ pagination.page }} / {{ Math.ceil(pagination.itemCount / pagination.pageSize) }}
        </span>
        <NButton
          size="small"
          :disabled="pagination.page >= Math.ceil(pagination.itemCount / pagination.pageSize)"
          @click="getList(pagination.page + 1)"
        >
          下一页
        </NButton>
      </div>
    </NSpin>
  </div>
</template>

<style scoped>
/* tailwind 3.2 尚未内置 line-clamp，这里自行实现两行截断 */
.desc-clamp {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
