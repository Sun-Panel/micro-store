<script setup lang="ts">
import { onMounted, ref } from 'vue'
import {
  NAlert,
  NBreadcrumb,
  NBreadcrumbItem,
  NButton,
  NCard,
  NDivider,
  NImage,
  NImageGroup,
  NInput,
  NModal,
  NSpace,
  NTag,
  useMessage,
} from 'naive-ui'
import { useRoute } from 'vue-router'
import { CodeBlock } from '@/components/common'
import { approve, getInfo, reject } from '@/api/admin/customCodeReview'
import { apiRespErrMsg } from '@/utils/cmn'
import { router } from '@/router'

const route = useRoute()
const message = useMessage()

const id = ref<number>(Number(route.query.id ?? 0))
const detail = ref<CustomCode.ReviewDetail>()
const rejectVisible = ref(false)
const rejectNote = ref('')
const submitting = ref(false)

const codeTypeLabelMap: Record<number, string> = { 1: 'JS', 2: 'CSS', 3: '页脚' }
const versionLabelMap: Record<number, string> = { 1: 'v1', 2: 'v2' }

function versionTexts(versions?: number[]) {
  if (!versions?.length)
    return '-'
  return versions.map(item => versionLabelMap[item] ?? item).join(' / ')
}

async function loadInfo() {
  try {
    const { data } = await getInfo<CustomCode.ReviewDetail>(id.value)
    detail.value = data
  }
  catch (res) {
    apiRespErrMsg(res)
  }
}

async function handleApprove() {
  submitting.value = true
  try {
    const { code, msg } = await approve<Common.DataResponse>(id.value)
    if (code === 0) {
      message.success('审核通过')
      router.push({ name: 'AdminCustomCodeReview' })
    }
    else {
      message.warning(`操作失败: ${msg}`)
    }
  }
  catch {
    message.warning('操作失败，请稍后重试')
  }
  finally {
    submitting.value = false
  }
}

async function handleReject() {
  if (!rejectNote.value.trim()) {
    message.warning('请填写拒绝原因')
    return
  }

  submitting.value = true
  try {
    const { code, msg } = await reject<Common.DataResponse>(id.value, rejectNote.value)
    if (code === 0) {
      message.success('已拒绝')
      rejectVisible.value = false
      router.push({ name: 'AdminCustomCodeReview' })
    }
    else {
      message.warning(`操作失败: ${msg}`)
    }
  }
  catch {
    message.warning('操作失败，请稍后重试')
  }
  finally {
    submitting.value = false
  }
}

onMounted(() => {
  if (id.value)
    loadInfo()
})
</script>

<template>
  <div>
    <NCard size="small" class="mb-5">
      <NBreadcrumb>
        <NBreadcrumbItem @click="router.push({ name: 'AdminCustomCodeReview' })">
          自定义代码审核
        </NBreadcrumbItem>
        <NBreadcrumbItem>
          审核详情
        </NBreadcrumbItem>
      </NBreadcrumb>
    </NCard>

    <NCard v-if="detail" size="small">
      <NSpace align="center" class="mb-[10px]">
        <span class="text-[16px] font-bold">{{ detail.review.title }}</span>
        <NTag size="small">
          适用版本 {{ versionTexts(detail.review.versions) }}
        </NTag>
        <NTag v-if="!detail.review.isOriginal" size="small" type="warning">
          非原创
        </NTag>
        <NTag v-if="detail.online" size="small" type="info">
          修改审核
        </NTag>
        <NTag v-else size="small">
          新提交
        </NTag>
      </NSpace>

      <div class="text-slate-500 mb-[10px]">
        {{ detail.review.description }}
      </div>

      <div v-if="detail.review.keywords?.length" class="mb-[10px]">
        <NTag v-for="item in detail.review.keywords" :key="item" size="small" class="mr-[5px]">
          {{ item }}
        </NTag>
      </div>

      <NAlert
        v-if="detail.review.machineAudit"
        :type="detail.review.machineAudit.auto ? 'success' : 'warning'"
        class="mb-[12px]"
      >
        <div class="font-bold mb-[6px]">
          {{ detail.review.machineAudit.auto ? '机器预审：自动通过' : '机器预审：转人工审核' }}
        </div>
        <div class="text-[13px] text-slate-500 mb-[4px]">
          受信任作者：{{ detail.review.machineAudit.trusted ? '是' : '否' }}
          ｜ 安全扫描：{{ detail.review.machineAudit.scanPass ? '通过' : '未通过' }}
        </div>
        <ul v-if="detail.review.machineAudit.reasons?.length" class="list-disc pl-[18px] text-[13px]">
          <li v-for="(r, i) in detail.review.machineAudit.reasons" :key="i">
            {{ r }}
          </li>
        </ul>
      </NAlert>

      <NAlert v-if="!detail.review.isOriginal" type="info" class="mb-[12px]">
        来源说明：{{ detail.review.sourceNote }}
      </NAlert>

      <NDivider title-placement="left">
        待审核代码片段（{{ detail.review.blocks?.length ?? 0 }}）
      </NDivider>

      <div v-if="!detail.review.blocks?.length" class="text-slate-400 mb-[12px]">
        无代码片块
      </div>
      <div v-for="(block, index) in detail.review.blocks" :key="index" class="mb-[14px]">
        <div class="font-bold mb-[6px]">
          {{ block.title || '未命名' }}代码片段（{{ codeTypeLabelMap[block.codeType] ?? block.codeType }}）
        </div>
        <div v-if="block.note" class="text-slate-500 text-[13px] mb-[6px] note-pre">
          {{ block.note }}
        </div>
        <div v-if="block.images?.length" class="mb-[6px]">
          <NImageGroup>
            <div class="flex flex-wrap gap-[6px]">
              <NImage v-for="url in block.images" :key="url" :src="url" width="100" />
            </div>
          </NImageGroup>
        </div>
        <CodeBlock :code="block.code" :code-type="block.codeType" :default-collapsed="false" />
      </div>

      <template v-if="detail.online">
        <NDivider title-placement="left">
          当前线上版本（{{ detail.online.blocks?.length ?? 0 }} 个代码片段）
        </NDivider>

        <div v-if="!detail.online.blocks?.length" class="text-slate-400 mb-[12px]">
          无线上块
        </div>
        <div v-for="(block, index) in detail.online.blocks" :key="index" class="mb-[14px]">
          <div class="font-bold mb-[6px]">
            {{ block.title || '未命名' }}代码片段（{{ codeTypeLabelMap[block.codeType] ?? block.codeType }}）
          </div>
          <div v-if="block.note" class="text-slate-500 text-[13px] mb-[6px] note-pre">
            {{ block.note }}
          </div>
          <CodeBlock :code="block.code" :code-type="block.codeType" />
        </div>
      </template>

      <NDivider />

      <NSpace>
        <NButton type="primary" :loading="submitting" @click="handleApprove">
          通过
        </NButton>
        <NButton type="error" :loading="submitting" @click="rejectVisible = true">
          拒绝
        </NButton>
      </NSpace>
    </NCard>

    <NModal v-model:show="rejectVisible" preset="card" style="width: 500px;" title="拒绝原因" :mask-closable="false">
      <NInput v-model:value="rejectNote" type="textarea" :rows="4" placeholder="请填写拒绝原因，作者可见" />
      <template #footer>
        <NSpace justify="end">
          <NButton @click="rejectVisible = false">
            取消
          </NButton>
          <NButton type="error" :loading="submitting" @click="handleReject">
            确定拒绝
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
.note-pre {
  white-space: pre-wrap;
}
</style>
