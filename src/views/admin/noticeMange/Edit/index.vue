<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal, NSelect, useMessage } from 'naive-ui'
import { edit as updateApi } from '@/api/admin/notice'

interface Props {
  visible: boolean
  userId?: number
  info?: Notice.NoticeInfo
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
  displayType: 1,
  oneRead: 0,
  url: '',
  isLogin: 1,
}

const model = ref<Notice.NoticeInfo>(formInitValue)
const formRef = ref<FormInst | null>(null)

interface Options {
  value?: number
  label?: string
  description?: string
}

const displayTypeOptions = ref<Options[]>([
  {
    label: '登录页',
    value: 1,
  },
  {
    label: '首页',
    value: 2,
  },
])

const boolOptions = ref<Options[]>([
  {
    label: '否',
    value: 0,
  },
  {
    label: '是',
    value: 1,
  },
])

const rules: FormRules = {
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
  if (props.info?.id)
    model.value = { ...props.info } || {}

  else
    model.value = formInitValue
})

const add = async () => {
  const res = await updateApi(model.value)
  if (res.code === 0)
    emit('done')

  else if (res.code !== -1)
    message.warning('操作失败')
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
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px;" :title="model.id ? '修改' : '添加'">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="displayType" label="展示类型">
        <NSelect
          v-model:value="model.displayType"
          :options="displayTypeOptions"
        />
      </NFormItem>

      <NFormItem path="title" label="标题">
        <NInput v-model:value="model.title" type="text" placeholder="标题" />
      </NFormItem>

      <NFormItem path="content" label="内容">
        <NInput v-model:value="model.content" placeholder="请输入内容" type="textarea" />
      </NFormItem>

      <NFormItem path="url" label="地址">
        <NInput v-model:value="model.url" type="text" placeholder="跳转地址:http(s)://" />
      </NFormItem>

      <NFormItem path="oneRead" label="允许读取">
        <NSelect
          v-model:value="model.oneRead"
          :options="boolOptions"
        />
      </NFormItem>

      <NFormItem path="isLogin" label="登录可见">
        <NSelect
          v-model:value="model.isLogin"
          :options="boolOptions"
        />
      </NFormItem>
    </NForm>

    <template #footer>
      <NButton type="success" @click="handleValidateButtonClick">
        保存
      </NButton>
    </template>
  </NModal>
</template>
