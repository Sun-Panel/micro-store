<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import {
  NAlert,
  NBreadcrumb,
  NBreadcrumbItem,
  NButton,
  NCard,
  NCheckbox,
  NCheckboxGroup,
  NCollapse,
  NCollapseItem,
  NDynamicTags,
  NForm,
  NFormItem,
  NImage,
  NImageGroup,
  NInput,
  NModal,
  NRadio,
  NRadioGroup,
  NSpace,
  NSwitch,
  NTag,
  useMessage,
} from 'naive-ui'
import { useRoute } from 'vue-router'
import { CodeBlock } from '@/components/common'
import { uploadPreviewImage, edit, getInfo, withdraw } from '@/api/admin/customCode'
import { apiRespErrMsg, getCurrentBaseUrlRoot } from '@/utils/cmn'
import { router } from '@/router'
import { genOnlyId, isValidOnlyId } from '@/utils/customCodeSnippet'

// 块本地状态：_uid 供 NCollapse 作为稳定的展开标识
interface BlockItem extends CustomCode.Block {
  _uid: string
  /** 块唯一标识（编辑态必填，用于复制粘贴去重） */
  onlyId: string
}

interface EditForm extends Omit<CustomCode.EditReq, 'blocks'> {
  blocks: BlockItem[]
}

const route = useRoute()
const message = useMessage()

const id = ref<number>(Number(route.query.id ?? 0))
const info = ref<CustomCode.Info>()
const saving = ref(false)

// 待审核期间只读，需先撤回
const readonly = computed(() => info.value?.canEdit === false)

// 前台公开地址
const publicUrl = computed(() => `${getCurrentBaseUrlRoot()}/customCode/${id.value}`)

function openPublicPage() {
  window.open(publicUrl.value, '_blank')
}

const versionOptions = [
  { label: 'v1', value: 1 },
  { label: 'v2', value: 2 },
]

const codeTypeOptions = [
  { label: 'JS', value: 1 },
  { label: 'CSS', value: 2 },
  { label: '页脚', value: 3 },
]

const codeTypeLabelMap: Record<number, string> = {
  1: 'JS',
  2: 'CSS',
  3: '页脚',
}

let uidSeed = 0

function createBlock(expanded = false): BlockItem {
  const uid = `block-${++uidSeed}`
  // 新增的片段默认展开，方便立即填写
  if (expanded)
    expandedNames.value.push(uid)

  return {
    codeType: 1,
    title: '',
    note: '',
    code: '',
    images: [],
    onlyId: genOnlyId(),
    _uid: uid,
  }
}

// 进入页面默认全部收缩
const expandedNames = ref<Array<string | number>>([])

const formRef = ref<FormInst | null>(null)
const model = ref<EditForm>({
  title: '',
  description: '',
  keywords: [],
  isOriginal: true,
  sourceNote: '',
  // 默认选中 v2
  versions: [2],
  blocks: [createBlock(true)],
  submit: false,
})

// 预览弹窗
const previewVisible = ref(false)
const previewBlock = ref<BlockItem | null>(null)

const rules: FormRules = {
  title: [
    { required: true, trigger: 'blur', message: '请输入标题' },
    { max: 100, trigger: 'blur', message: '标题不能超过100个字符' },
  ],
  description: [
    { required: true, trigger: 'blur', message: '请输入说明' },
    { max: 500, trigger: 'blur', message: '说明不能超过500个字符' },
  ],
  versions: {
    validator: () => {
      if (!model.value.versions.length)
        return new Error('请至少选择一个适用版本')
      return true
    },
    trigger: 'change',
  },
  sourceNote: {
    validator: () => {
      if (!model.value.isOriginal && !model.value.sourceNote.trim())
        return new Error('非原创必须填写来源说明')
      return true
    },
    trigger: 'blur',
  },
}

function validateBlocks(): string {
  const blocks = model.value.blocks
  if (!blocks.length)
    return '请至少添加一个代码片段'
  if (blocks.length > 10)
    return '代码片段最多10个'

  const onlyIdSet = new Set<string>()
  for (let i = 0; i < blocks.length; i++) {
    const block = blocks[i]
    const oid = block.onlyId.trim()
    if (!isValidOnlyId(oid))
      return `第${i + 1}个片段的唯一标识格式不正确（仅限字母、数字、下划线，最长40位）`
    if (onlyIdSet.has(oid))
      return `第${i + 1}个片段的唯一标识重复：${oid}`
    onlyIdSet.add(oid)
    if (!block.title.trim() && !block.note.trim())
      return `第${i + 1}个代码片段的标题和说明至少填写一个`
    if (!block.code.trim())
      return `第${i + 1}个代码片段的代码不能为空`
    if (block.images.length > 5)
      return `第${i + 1}个代码片段的预览图最多5张`
  }
  return ''
}

function addBlock() {
  if (model.value.blocks.length >= 10) {
    message.warning('代码片段最多10个')
    return
  }
  // 新增的块默认展开，方便立即填写
  model.value.blocks.push(createBlock(true))
}

function removeBlock(index: number) {
  if (model.value.blocks.length <= 1) {
    message.warning('至少保留一个代码片段')
    return
  }
  const [removed] = model.value.blocks.splice(index, 1)
  if (removed)
    expandedNames.value = expandedNames.value.filter(name => name !== removed._uid)
}

function moveBlock(index: number, offset: number) {
  const target = index + offset
  if (target < 0 || target >= model.value.blocks.length)
    return
  const blocks = model.value.blocks
  const temp = blocks[index]
  blocks[index] = blocks[target]
  blocks[target] = temp
}

function toggleExpandAll() {
  const allExpanded = model.value.blocks.length > 0
    && expandedNames.value.length >= model.value.blocks.length

  expandedNames.value = allExpanded
    ? []
    : model.value.blocks.map(block => block._uid)
}

function openPreview(block: BlockItem) {
  previewBlock.value = block
  previewVisible.value = true
}

async function handlePickImage(event: Event, block: BlockItem) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  // 置空以便重复选择同一张图
  input.value = ''

  if (!file)
    return
  if (file.size > 512 * 1024) {
    message.warning('预览图不能大于512K')
    return
  }
  if (block.images.length >= 5) {
    message.warning('每个代码片段最多5张预览图')
    return
  }

  try {
    const url = await uploadPreviewImage(file)
    block.images.push(url)
    message.success('上传成功')
  }
  catch (error) {
    message.warning(error instanceof Error ? error.message : '上传失败')
  }
}

function removeImage(block: BlockItem, index: number) {
  block.images.splice(index, 1)
}

async function submit(submitForReview: boolean) {
  if (readonly.value) {
    message.warning('审核中不可修改，请先撤回审核')
    return
  }

  const formErrors = await formRef.value?.validate().catch(err => err)
  if (formErrors)
    return

  const blockError = validateBlocks()
  if (blockError) {
    message.warning(blockError)
    return
  }

  saving.value = true
  try {
    // 去掉仅用于界面交互的展开状态
    const payload: CustomCode.EditReq = {
      ...model.value,
      submit: submitForReview,
      blocks: model.value.blocks.map(block => ({
        id: block.id,
        codeType: block.codeType,
        title: block.title,
        note: block.note,
        code: block.code,
        images: block.images,
      })),
    }
    const { code, msg } = await edit<Common.DataResponse>(payload)
    if (code === 0) {
      message.success(submitForReview ? '已提交审核' : '草稿已保存')
      router.push({ name: 'AdminCustomCodeMyList' })
    }
    else {
      message.warning(`操作失败: ${msg}`)
    }
  }
  catch {
    message.warning('保存失败，请稍后重试')
  }
  finally {
    saving.value = false
  }
}

async function handleWithdraw() {
  if (!id.value)
    return
  try {
    const { code, msg } = await withdraw<Common.DataResponse>(id.value)
    if (code === 0) {
      message.success('已撤回，可以重新编辑')
      await loadInfo()
    }
    else {
      message.warning(`操作失败: ${msg}`)
    }
  }
  catch {
    message.warning('撤回失败，请稍后重试')
  }
}

async function loadInfo() {
  try {
    const { data } = await getInfo<CustomCode.Info>(id.value)
    info.value = data
    model.value = {
      id: data.id,
      title: data.title,
      description: data.description,
      keywords: data.keywords ?? [],
      isOriginal: data.isOriginal,
      sourceNote: data.sourceNote ?? '',
      versions: data.versions?.length ? data.versions : [2],
      blocks: data.blocks?.length
        ? data.blocks.map((block) => {
            return {
              id: block.id,
              codeType: block.codeType,
              title: block.title ?? '',
              note: block.note ?? '',
              code: block.code ?? '',
              images: block.images ?? [],
              onlyId: block.onlyId ?? '',
              _uid: `block-${++uidSeed}`,
            }
          })
        : [createBlock(true)],
      submit: false,
    }

    // 进入页面默认全部收缩
    expandedNames.value = []
  }
  catch (res) {
    apiRespErrMsg(res)
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
        <NBreadcrumbItem @click="router.push({ name: 'AdminCustomCodeMyList' })">
          我的自定义代码
        </NBreadcrumbItem>
        <NBreadcrumbItem>
          {{ id ? '编辑' : '新建' }}
        </NBreadcrumbItem>
      </NBreadcrumb>
    </NCard>

    <NCard v-if="id" size="small" class="mb-5">
      公开地址
      <NTag type="info" style="cursor: pointer;" @click="openPublicPage">
        {{ publicUrl }}
      </NTag>
    </NCard>

    <NAlert v-if="readonly" type="warning" class="mb-5" title="审核中">
      内容正在审核，暂时不可编辑。如需修改请先撤回审核。
      <template #footer>
        <NButton size="small" type="warning" @click="handleWithdraw">
          撤回审核
        </NButton>
      </template>
    </NAlert>

    <NAlert v-if="info?.reviewStatus === 2" type="error" class="mb-5" title="审核未通过">
      {{ info.reviewNote || '请查看审核意见并修改后重新提交' }}
    </NAlert>

    <NCard size="small" class="mb-5">
      <NForm ref="formRef" :model="model" :rules="rules" :disabled="readonly">
        <NFormItem path="title" label="标题">
          <NInput v-model:value="model.title" :maxlength="100" show-count placeholder="最多100个字符" />
        </NFormItem>

        <NFormItem path="description" label="说明">
          <NInput v-model:value="model.description" :maxlength="500" type="textarea" show-count placeholder="最多500个字符" />
        </NFormItem>

        <NFormItem label="关键词">
          <NDynamicTags v-model:value="model.keywords" :max="10" />
        </NFormItem>

        <NFormItem label="是否原创">
          <NSwitch v-model:value="model.isOriginal" />
        </NFormItem>

        <NFormItem v-if="!model.isOriginal" path="sourceNote" label="来源说明">
          <NInput v-model:value="model.sourceNote" placeholder="可以是地址，也可以是其他说明" />
        </NFormItem>

        <NFormItem path="versions" label="适用客户端版本">
          <NCheckboxGroup v-model:value="model.versions">
            <NSpace>
              <NCheckbox v-for="item in versionOptions" :key="item.value" :value="item.value" :label="item.label" />
            </NSpace>
          </NCheckboxGroup>
        </NFormItem>
      </NForm>

      <div class="mb-[10px] flex items-center">
      <span class="font-bold">代码片段（{{ model.blocks.length }}/10）</span>
      <NButton size="tiny" class="ml-auto" @click="toggleExpandAll">
        {{ expandedNames.length >= model.blocks.length && model.blocks.length > 0 ? '全部收起' : '全部展开' }}
      </NButton>
    </div>

    <NCollapse v-model:expanded-names="expandedNames" class="mb-[20px]">
      <NCollapseItem
        v-for="(block, index) in model.blocks"
        :key="block._uid"
        :name="block._uid"
        :title="`${codeTypeLabelMap[block.codeType] ?? block.codeType}-代码片段: ${block.title || '未命名'}`"
      >
        <template #header-extra>
          <!-- header-extra 的点击会连带触发折叠，这里阻止冒泡 -->
          <div @click.stop>
            <NSpace :size="6">
              <NButton size="tiny" :disabled="index === 0 || readonly" @click="moveBlock(index, -1)">
                上移
              </NButton>
              <NButton size="tiny" :disabled="index === model.blocks.length - 1 || readonly" @click="moveBlock(index, 1)">
                下移
              </NButton>
              <NButton size="tiny" type="error" :disabled="readonly" @click="removeBlock(index)">
                删除
              </NButton>
            </NSpace>
          </div>
        </template>
        <NForm :disabled="readonly">
          <NFormItem label="代码类型">
            <NRadioGroup v-model:value="block.codeType">
              <NSpace>
                <NRadio v-for="item in codeTypeOptions" :key="item.value" :value="item.value">
                  {{ item.label }}
                </NRadio>
              </NSpace>
            </NRadioGroup>
          </NFormItem>

          <NFormItem label="片段标题">
            <NInput v-model:value="block.title" :maxlength="100" placeholder="标题与说明至少填写一个" />
          </NFormItem>

          <NFormItem label="唯一标识">
            <div class="flex w-full items-center gap-2">
              <NInput v-model:value="block.onlyId" :maxlength="40" placeholder="用于跨项目粘贴去重，可自定义" />
              <NButton size="tiny" type="default" @click="block.onlyId = genOnlyId()">
                重新生成
              </NButton>
            </div>
          </NFormItem>

          <NFormItem label="说明">
            <NInput v-model:value="block.note" type="textarea" :rows="3" placeholder="多行说明，与标题至少填写一个" />
          </NFormItem>

          <NFormItem label="预览图（≤512K，最多5张）">
            <div>
              <input
                type="file"
                accept="image/png,image/jpeg,image/gif,image/webp"
                :disabled="readonly || block.images.length >= 5"
                @change="handlePickImage($event, block)"
              >
              <div v-if="block.images.length" class="mt-[8px]">
                <NImageGroup>
                  <div class="flex flex-wrap gap-[8px]">
                    <div v-for="(url, imgIndex) in block.images" :key="url" class="flex flex-col items-start gap-[4px]">
                      <NImage :src="url" width="100" />
                      <NButton size="tiny" type="error" :disabled="readonly" @click="removeImage(block, imgIndex)">
                        移除
                      </NButton>
                    </div>
                  </div>
                </NImageGroup>
              </div>
            </div>
          </NFormItem>

          <NFormItem label="代码">
            <NInput v-model:value="block.code" type="textarea" :rows="6" placeholder="粘贴代码，支持预览高亮" />
          </NFormItem>
        </NForm>

        <NButton size="small" :disabled="!block.code.trim()" @click="openPreview(block)">
          预览
        </NButton>
      </NCollapseItem>
    </NCollapse>

    <NSpace class="mb-[20px]">
      <NButton :disabled="readonly" @click="addBlock">
        添加代码片段
      </NButton>
    </NSpace>
    </NCard>

    

    <NSpace>
      <NButton :disabled="readonly" :loading="saving" @click="submit(false)">
        保存草稿
      </NButton>
      <NButton type="primary" :disabled="readonly" :loading="saving" @click="submit(true)">
        提交审核
      </NButton>
    </NSpace>

    <NModal v-model:show="previewVisible" preset="card" style="width: 720px;" title="代码预览" :mask-closable="true">
      <div v-if="previewBlock">
        <div class="mb-[8px]">
         {{ codeTypeLabelMap[previewBlock.codeType] ?? previewBlock.codeType }}代码片段: {{ previewBlock.title || '未命名' }}
        </div>
        <CodeBlock :code="previewBlock.code" :code-type="previewBlock.codeType" :default-collapsed="false" />
      </div>
    </NModal>
  </div>
</template>
