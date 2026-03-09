<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NCheckbox, NDivider, NForm, NFormItem, NInput, NInputNumber, NModal, NSwitch, useMessage } from 'naive-ui'
import { templateEdit as templateEditApi } from '@/api/admin/message'

interface Props {
  visible: boolean
  userId?: number
  info?: Message.TemplateInfo
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const message = useMessage()

interface MessageInfo {
  id: number
  toUsernames?: string
  name?: string
  flag?: string
  title?: string
  titleShow?: number
  titleEdit?: number
  content?: string
  contentShow?: number
  contentEdit?: number
  note?: string
  args?: Message.Arg[]
  toUsernameShow?: number
  toUsernameEdit?: number
  uploadAppendixShow?: number
  uploadAppendixMax?: number
  uploadAppendixMin?: number
  saveAfterUrl?: string
  saveButtonName?: string
}

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}

const formInitValue: MessageInfo = {
  id: props?.info?.id || 0,
  titleShow: 1,
  titleEdit: 1,
  content: '',
  contentShow: 1,
  contentEdit: 1,
  toUsernames: '',
  toUsernameShow: 1,
  toUsernameEdit: 1,
  saveButtonName: '发送',
  saveAfterUrl: '/',
}

const templatesArgs = ref<Message.Arg[]>([])
const model = ref<MessageInfo>(formInitValue)
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  flag: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  ],
  name: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  ],
  // title: [
  //   {
  //     required: true,
  //     trigger: 'blur',
  //     message: '必填项',
  //   },
  // ],
  // content: [
  //   {
  //     required: true,
  //     trigger: 'blur',
  //     message: '必填项',
  //   },
  // ],
}

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

watch(show, (newValue, oldValue) => {
  if (props.info?.id) {
    model.value = {
      id: props.info.id,
      toUsernames: props.info.toUsernames as string,
      name: props.info.name,
      flag: props.info.flag,
      title: props.info.title,
      titleEdit: props.info.titleEdit,
      titleShow: props.info.titleShow,
      content: props.info.content,
      contentShow: props.info.contentShow,
      contentEdit: props.info.contentEdit,
      note: props.info.note,
      toUsernameShow: props.info.toUsernameShow,
      toUsernameEdit: props.info.toUsernameEdit,
      args: props.info.args,
      uploadAppendixShow: props.info.uploadAppendixShow,
      uploadAppendixMax: props.info.uploadAppendixMax,
      uploadAppendixMin: props.info.uploadAppendixMin,
      saveAfterUrl: props.info.saveAfterUrl,
      saveButtonName: props.info.saveButtonName,
    } || {}
    templatesArgs.value = model.value.args || []
  }
  else {
    model.value = { ...formInitValue }
    templatesArgs.value = []
  }
})

async function edit() {
  const postArg: Message.TemplateInfo = {
    id: model.value.id,
    toUsernames: splitEmailAddress(model.value.toUsernames || ''),
    name: model.value.name || '',
    flag: model.value.flag || '',
    title: model.value.title,
    titleShow: model.value.titleShow,
    titleEdit: model.value.titleEdit,
    content: model.value.content,
    contentShow: model.value.contentShow,
    contentEdit: model.value.contentEdit,
    toUsernameShow: model.value.toUsernameShow,
    toUsernameEdit: model.value.toUsernameEdit,
    uploadAppendixShow: model.value.uploadAppendixShow,
    uploadAppendixMax: model.value.uploadAppendixMax,
    uploadAppendixMin: model.value.uploadAppendixMin,
    saveAfterUrl: model.value.saveAfterUrl,
    saveButtonName: model.value.saveButtonName,
  }

  if (postArg.content === '') {
    postArg.contentShow = 1
    postArg.contentEdit = 1
  }

  if (postArg.toUsernames.length === 0) {
    postArg.toUsernameShow = 1
    postArg.toUsernameEdit = 1
  }

  postArg.args = templatesArgs.value
  console.log('要保存的内容', postArg)
  const res = await templateEditApi(postArg)
  if (res.code === 0)
    emit('done')

  else
    message.warning(`操作失败: ${res.msg}`)
}

function splitEmailAddress(input: string): string[] {
  // 使用正则表达式匹配逗号或空格进行分割
  const regex = /[,\s]+/
  const splittedStrings = input.split(regex)

  // 过滤出符合邮箱格式的项
  const validEmails = splittedStrings.filter((item) => {
    // 使用简单的邮箱正则表达式进行验证
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(item)
  })

  return validEmails
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      edit()
    else
      console.log(errors)
  })
}

function removeArgItem(index: number) {
  templatesArgs.value.splice(index, 1)
}

function addArgItem() {
  if (!templatesArgs.value)
    templatesArgs.value = []

  templatesArgs.value?.push({
    keyword: '',
    name: '',
    required: false,
    placeholder: '',
    formType: '',
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px;" :title="model.id ? '修改' : '添加'">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="flag" label="模板唯一标识">
        <NInput v-model:value="model.flag" type="text" placeholder="模板唯一标识，例如：template_cn_register" />
      </NFormItem>

      <NFormItem path="name" label="模板名称">
        <NInput v-model:value="model.name" type="text" placeholder="模板名称" />
      </NFormItem>

      <NFormItem path="title" label="标题">
        <div class="w-full">
          <div class="mb-2">
            <NInput v-model:value="model.title" type="text" placeholder="标题" />
          </div>
          <div class="flex">
            <div class="flex items-center gap-2">
              <span>允许编辑</span>
              <NSwitch
                v-model:value="model.titleEdit"
                :checked-value="1"
                :unchecked-value="0"
                :disabled="model.title === ''"
              />
            </div>
            <div class="flex items-center gap-2 ml-5">
              <span>显示</span>
              <NSwitch
                v-model:value="model.titleShow"
                :checked-value="1"
                :unchecked-value="0"
                :disabled="model.title === ''"
              />
            </div>
          </div>
        </div>
      </NFormItem>

      <NFormItem path="title" label="默认收信人">
        <div class="w-full">
          <div class="mb-2">
            <NInput v-model:value="model.toUsernames" type="text" placeholder="邮箱地址，多个用逗号隔开" />
          </div>

          <div class="flex">
            <div class="flex items-center gap-2">
              <span>允许编辑</span>
              <NSwitch
                v-model:value="model.toUsernameEdit"
                :checked-value="1"
                :unchecked-value="0"
                :disabled="model.toUsernames === ''"
              />
            </div>
            <div class="flex items-center gap-2 ml-5">
              <span>显示</span>
              <NSwitch
                v-model:value="model.toUsernameShow"
                :checked-value="1"
                :unchecked-value="0"
                :disabled="model.toUsernames === ''"
              />
            </div>
          </div>
        </div>
      </NFormItem>

      <NFormItem path="content" label="内容">
        <div class="w-full">
          <div class="mb-2">
            <NInput v-model:value="model.content" placeholder="请输入内容" type="textarea" />
          </div>

          <div class="flex">
            <div class="flex items-center gap-2">
              <span>允许编辑</span>
              <NSwitch
                v-model:value="model.contentEdit"
                :checked-value="1"
                :unchecked-value="0"
                :disabled="model.content === ''"
              />
            </div>
            <div class="flex items-center gap-2 ml-5">
              <span>显示</span>
              <NSwitch
                v-model:value="model.contentShow"
                :checked-value="1"
                :unchecked-value="0"
                :disabled="model.content === ''"
              />
            </div>
          </div>
        </div>
      </NFormItem>

      <NFormItem label="模板参数">
        <NCard size="small">
          <div
            v-for="(item, index) in templatesArgs"
            :key="index"
            class="mt-2"
          >
            <div class="flex">
              <NInput v-model:value="item.keyword" clearable placeholder="变量名-示例：({name})" />
              <NInput v-model:value="item.name" style="margin-left: 5px" clearable placeholder="变量说明-示例：(姓名)" />
            </div>
            <div class="mt-1">
              <NInput v-model:value="item.placeholder" clearable placeholder="占位提示" />
            </div>
            <div class="mt-1">
              <NCheckbox v-model:checked="item.required">
                是否为必填项
              </NCheckbox>
              <NButton size="tiny" type="error" @click="removeArgItem(index)">
                删除参数
              </NButton>
            </div>
            <NDivider />
          </div>

          <NButton attr-type="button" size="small" style="margin-top: 12px" @click="addArgItem">
            增加参数
          </NButton>
        </NCard>
      </NFormItem>

      <NFormItem label="显示上传附件框">
        <div class="w-full">
          <div>
            <NSwitch
              v-model:value="model.uploadAppendixShow"
              :checked-value="1"
              :unchecked-value="0"
            />
          </div>
          <div v-if="model.uploadAppendixShow === 1" class="my-2 flex">
            <div>
              最小/个
              <NInputNumber v-model:value="model.uploadAppendixMin" :min="0" placeholder="最小上传限制" />
            </div>
            <div class="ml-5">
              最大/个
              <NInputNumber v-model:value="model.uploadAppendixMax" :min="0" placeholder="最大上传限制" />
            </div>
          </div>
        </div>
      </NFormItem>

      <NFormItem label="发送按钮的名字">
        <NInput v-model:value="model.saveButtonName" type="text" placeholder="按钮名称" />
      </NFormItem>

      <NFormItem label="发送成功跳转地址">
        <NInput v-model:value="model.saveAfterUrl" type="text" placeholder="跳转地址" />
      </NFormItem>

      <!-- <NFormItem path="url" label="备注">
        <NInput v-model:value="model.note" type="text" />
      </NFormItem> -->
    </NForm>

    <template #footer>
      <NButton type="success" @click="handleValidateButtonClick">
        保存
      </NButton>
    </template>
  </NModal>
</template>
