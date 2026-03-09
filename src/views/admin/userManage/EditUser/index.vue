<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal, NSelect } from 'naive-ui'
import { AdminUserManageEdit } from '@/api/admin'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'
import { t } from '@/locales'

interface Props {
  visible: boolean
  userId?: number
  userInfo?: User.Info
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', id: number): void// 创建完成
}

const formInitValue = {
  name: '',
  username: '',
  role: 2,
  status: 3,
}

const model = ref<User.Info>(formInitValue)
const formRef = ref<FormInst | null>(null)

interface Options {
  value?: number
  label?: string
  description?: string
}

const options = ref<Options[]>([
  {
    label: '管理员',
    value: 1,
  },
  {
    label: '普通用户',
    value: 2,
  },
])

const statusOptions = ref<Options[]>([
  {
    label: '启用',
    value: 1,
  },
  {
    label: '停用',
    value: 2,
  },
  {
    label: '未激活',
    value: 3,
  },
])

const rules: FormRules = {
  username: [
    {
      required: true,
      trigger: 'blur',
      message: '请输入账号且大于5个字符',
      min: 5,
    },
    {
      trigger: 'blur',
      message: '请输入邮箱作为账号',
      type: 'email',
    },
  ],
  role: {
    required: true,
    trigger: 'blur',
    type: 'number',
    message: '请选择角色',
  },
  status: {
    required: true,
    trigger: 'blur',
    type: 'number',
    message: '请选择账号状态',
  },
}

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

watch(show, (newValue, oldValue) => {
  if (props.userInfo?.id)
    model.value = props.userInfo || {}

  else
    model.value = formInitValue
})

const add = async () => {
  await AdminUserManageEdit<User.Info>(model.value).then((res) => {
    emit('done', res.data.id as number)
  }).catch((error) => {
    apiRespErrMsgAndCustomCodeNeg1Msg(error, t('common.failed'))
  })
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
  <NModal v-model:show="show" preset="card" style="width: 600px" :title="`${userInfo?.id ? '编辑' : '添加'}用户`">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="username" label="账号" style="margin-top: 20px;">
        <NInput v-model:value="model.username" type="text" placeholder="邮箱地址作为账号" />
      </NFormItem>

      <NFormItem path="name" label="昵称" style="margin-top: 20px;">
        <NInput v-model:value="model.name" type="text" placeholder="请输入昵称" />
      </NFormItem>

      <NFormItem path="role" label="角色">
        <NSelect
          v-model:value="model.role"
          :options="options"
        />
      </NFormItem>

      <NFormItem path="status" label="状态">
        <NSelect
          v-model:value="model.status"
          :options="statusOptions"
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
