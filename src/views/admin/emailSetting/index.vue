<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NAlert, NButton, NCard, NForm, NFormItem, NInput, NInputNumber, useMessage } from 'naive-ui'
import { AdminSystemSettingGetEmail, AdminSystemSettingSetEmail } from '@/api/admin'

const message = useMessage()

const formInitValue = {
  host: '',
  port: 465,
  mail: '',
  password: '',
}

const model = ref<AdminSystemSetting.Email>(formInitValue)
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  host: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  ],
  port: {
    required: true,
    trigger: 'blur',
    type: 'number',
    message: '必填',
  },
  mail: [{
    required: true,
    trigger: 'blur',
    message: '必填项且大于5个字符',
    min: 5,
  },
  {
    trigger: 'blur',
    message: '请输入邮箱格式',
    type: 'email',
  }],
  password: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项且大于5个字符',
      min: 5,
    },
  ],
}

const handleSave = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      const post = async () => {
        const { code, msg } = await AdminSystemSettingSetEmail(model.value)
        if (code === 0)
          message.success('保存成功')
        else
          message.warning(`保存失败，${msg}`)
      }
      post()
    }
  })
}

onMounted(() => {
  const post = async () => {
    const { code, data } = await AdminSystemSettingGetEmail<AdminSystemSetting.Email>()
    if (code === 0)
      model.value = data
  }
  post()
})
</script>

<template>
  <div>
    <NCard class="max-w-[500px]">
      <NAlert type="success" :show-icon="false">
        邮箱一般用于账号注册、密码找回等
      </NAlert>
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem path="host" label="服务器地址" style="margin-top: 20px;">
          <NInput v-model:value="model.host" type="text" placeholder="SMTP服务器地址" />
        </NFormItem>

        <NFormItem path="port" label="服务器端口" style="margin-top: 20px;">
          <NInputNumber v-model:value="model.port" placeholder="SMTP服务器端口" />
        </NFormItem>

        <NFormItem path="mail" label="发信邮箱账号" style="margin-top: 20px;">
          <NInput v-model:value="model.mail" type="text" placeholder="邮箱地址" />
        </NFormItem>

        <NFormItem path="password" label="邮箱密码" style="margin-top: 20px;">
          <NInput v-model:value="model.password" type="password" placeholder="邮箱密码或授权码" />
        </NFormItem>
      </NForm>
      <NButton type="success" @click="handleSave">
        保存
      </NButton>
    </NCard>
  </div>
</template>
