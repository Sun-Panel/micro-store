<script lang="ts" setup>
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NForm, NFormItem, NInput, NSpace, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { checkIsDeveloper, getInfo, register, updateMyInfo } from '@/api/developer'
import { SvgIcon } from '@/components/common'

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const isDeveloper = ref(false)
const isEdit = ref(false)

const model = ref<Developer.RegisterRequest>({
  developerName: '',
  contactMail: '',
  paymentName: '',
  paymentQrcode: '',
  paymentMethod: '',
})

const rules: FormRules = {
  developerName: [
    { required: true, trigger: 'blur', message: '请输入开发者标识' },
    {
      pattern: /^[a-z][a-z0-9-]*$/,
      trigger: 'blur',
      message: '只能包含小写字母、数字和中划线，且以字母开头',
    },
  ],
  contactMail: [{ type: 'email', trigger: 'blur', message: '请输入有效的邮箱地址' }],
}

async function checkDeveloperStatus() {
  try {
    const res = await checkIsDeveloper<any>()
    isDeveloper.value = res.data?.isDeveloper

    if (isDeveloper.value) {
      const devRes = await getInfo<any>()
      if (devRes.data) {
        model.value = {
          developerName: devRes.data.developerName,
          contactMail: devRes.data.contactMail || '',
          paymentName: devRes.data.paymentName || '',
          paymentQrcode: devRes.data.paymentQrcode || '',
          paymentMethod: devRes.data.paymentMethod || '',
        }
        isEdit.value = true
      }
    }
  }
  catch (error) {
    console.error(error)
  }
}

async function handleSubmit() {
  await formRef.value?.validate()

  loading.value = true
  try {
    if (isEdit.value) {
      await updateMyInfo(model.value)
      message.success('更新成功')
    }
    else {
      await register(model.value)
      message.success('注册成功')
      isDeveloper.value = true
      isEdit.value = true
    }
  }
  catch (error: any) {
    message.error(error?.message || '操作失败')
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  checkDeveloperStatus()
})
</script>

<template>
  <div class="min-h-[calc(100vh-120px)] flex items-center justify-center p-4">
    <NCard :title="isEdit ? '开发者信息' : '成为开发者'" style="max-width: 600px;">
      <template #header-extra>
        <NSpace align="center">
          <span v-if="isDeveloper" class="text-green-500 text-sm">
            <SvgIcon icon="mdi:check-circle" class="mr-1" />
            已认证
          </span>
        </NSpace>
      </template>

      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" label-width="100">
        <NFormItem path="developerName" label="开发者标识">
          <NInput
            v-model:value="model.developerName"
            :disabled="isEdit"
            placeholder="纯英文，多词用-分割，如：my-team"
          />
        </NFormItem>

        <NFormItem path="contactMail" label="联系邮箱">
          <NInput v-model:value="model.contactMail" placeholder="用于接收重要通知" />
        </NFormItem>

        <NFormItem path="paymentName" label="收款人姓名">
          <NInput v-model:value="model.paymentName" placeholder="请输入收款人真实姓名" />
        </NFormItem>

        <NFormItem path="paymentQrcode" label="收款二维码">
          <NInput v-model:value="model.paymentQrcode" placeholder="收款二维码图片URL" />
        </NFormItem>

        <NFormItem path="paymentMethod" label="收款方式">
          <NInput v-model:value="model.paymentMethod" placeholder="如：支付宝、微信等" />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end">
          <NButton type="primary" :loading="loading" @click="handleSubmit">
            {{ isEdit ? '保存修改' : '立即注册' }}
          </NButton>
        </NSpace>
      </template>
    </NCard>
  </div>
</template>
