<script lang="ts" setup>
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NDivider, NForm, NFormItem, NInput, NSpace, useMessage } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { SvgIcon } from '@/components/common'

interface Props {
  /** 是否为编辑模式 */
  editMode?: boolean
  /** 是否已认证为开发者 */
  isDeveloper?: boolean
  /** 初始数据 */
  initialData?: Partial<Developer.RegisterRequest>
}

const props = withDefaults(defineProps<Props>(), {
  editMode: false,
  isDeveloper: false,
  initialData: () => ({}),
})

const emit = defineEmits<{
  (e: 'submit', data: Developer.RegisterRequest, isEdit: boolean): void
}>()

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const loading = ref(false)

const model = ref<Developer.RegisterRequest>({
  developerName: '',
  contactMail: '',
  paymentName: '',
  paymentQrcode: '',
  paymentMethod: '',
  name: '',
})

const rules: FormRules = {
  developerName: [
    { required: true, trigger: 'blur', message: '仅支持小写英文字母、“-”（中划线），确定后不可随意修改' },
    {
      pattern: /^[a-z][a-z0-9-]*$/,
      trigger: 'blur',
      message: '只能包含小写字母、数字和中划线，且以字母开头',
    },
  ],
  name: { required: true, trigger: 'blur', message: '必填项' },
  contactMail: [{ type: 'email', trigger: 'blur', message: '请输入有效的邮箱地址' }],
}

const cardTitle = computed(() => props.editMode ? '开发者信息' : '成为开发者')
const buttonText = computed(() => props.editMode ? '保存修改' : '立即注册')

// 监听 initialData 变化，更新 model
watch(() => props.initialData, (newData) => {
  if (newData && Object.keys(newData).length > 0) {
    model.value = {
      developerName: newData.developerName || '',
      contactMail: newData.contactMail || '',
      paymentName: newData.paymentName || '',
      paymentQrcode: newData.paymentQrcode || '',
      paymentMethod: newData.paymentMethod || '',
      name: newData.name || '',
    }
  }
}, { immediate: true, deep: true })

async function handleSubmit() {
  await formRef.value?.validate()

  loading.value = true
  try {
    emit('submit', model.value, props.editMode)
  }
  catch (error: any) {
    message.error(error?.message || '操作失败')
  }
  finally {
    loading.value = false
  }
}

function setLoading(value: boolean) {
  loading.value = value
}

defineExpose({
  setLoading,
})
</script>

<template>
  <NCard :title="cardTitle" style="max-width: 600px;">
    <template #header-extra>
      <NSpace align="center">
        <span v-if="isDeveloper" class="text-green-500 text-sm">
          <SvgIcon icon="mdi:check-circle" class="mr-1" />
          已认证为开发者
        </span>
      </NSpace>
    </template>

    <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" label-width="100">
      <NDivider title-placement="left">
        开发者基础信息
      </NDivider>
      <NFormItem path="developerName" label="开发者标识">
        <NInput
          v-model:value="model.developerName"
          :disabled="editMode"
          placeholder="纯英文，多词用-分割，如：my-team"
        />
      </NFormItem>

      <NFormItem path="name" label="开发者名称">
        <NInput v-model:value="model.name" placeholder="公开显示的开发者名称" />
      </NFormItem>

      <NFormItem path="contactMail" label="联系邮箱">
        <NInput v-model:value="model.contactMail" placeholder="用于接收重要通知" />
      </NFormItem>

      <NDivider title-placement="left">
        收款信息
      </NDivider>

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
          {{ buttonText }}
        </NButton>
      </NSpace>
    </template>
  </NCard>
</template>
