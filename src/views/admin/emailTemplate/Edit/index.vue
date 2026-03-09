<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NForm, NFormItem, NInput, NModal, useMessage } from 'naive-ui'
import { edit as updateApi } from '@/api/admin/emailTemplate'

interface Props {
  visible: boolean
  userId?: number
  info?: EmailTemplate.Info
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const message = useMessage()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}

const formInitValue = {
  title: '',
  content: '',
  name: '',
  note: '',
  flag: '',
  isSeparateSend: true,
}

const templatesArgs = ref<EmailTemplate.Arg[]>([])
const model = ref<EmailTemplate.Info>(formInitValue)
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
  title: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  ],
  content: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  ],
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
    model.value = { ...props.info } || {}
    templatesArgs.value = model.value.args || []
  }
  else {
    model.value = { ...formInitValue }
  }
})

const add = async () => {
  model.value.args = templatesArgs.value
  const res = await updateApi(model.value)
  if (res.code === 0)
    emit('done')

  else
    message.warning(`操作失败: ${res.msg}`)
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      add()
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

      <NFormItem path="title" label="邮件标题">
        <NInput v-model:value="model.title" type="text" placeholder="标题" />
      </NFormItem>

      <NFormItem path="content" label="邮件内容">
        <NInput v-model:value="model.content" placeholder="请输入内容" type="textarea" />
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
              <NButton style="margin-left: 12px" size="small" @click="removeArgItem(index)">
                删除
              </NButton>
            </div>
          </div>
          <NButton attr-type="button" size="small" style="margin-top: 12px" @click="addArgItem">
            增加参数
          </NButton>
        </NCard>
      </NFormItem>

      <NFormItem path="url" label="备注">
        <NInput v-model:value="model.note" type="text" />
      </NFormItem>
    </NForm>

    <template #footer>
      <NButton type="success" @click="handleValidateButtonClick">
        保存
      </NButton>
    </template>
  </NModal>
</template>
