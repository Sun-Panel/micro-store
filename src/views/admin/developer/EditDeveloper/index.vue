<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal, NSelect } from 'naive-ui'
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import { update } from '@/api/admin/developer'
import { t } from '@/locales'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'

interface Props {
  visible: boolean
  developerInfo?: Developer.DeveloperInfo
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}

const formInitValue = {
  id: 0,
  developerName: '',
  contactMail: '',
  paymentName: '',
  name: '',
  // paymentQrcode: '',
  // paymentMethod: '',
  status: 1,
}

const model = ref<Developer.UpdateRequest>({ ...formInitValue })
const formRef = ref<FormInst | null>(null)

const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 },
]

const rules: FormRules = {
  developerName: [{ required: true, trigger: 'blur', message: '请输入开发者标识' }],
}

const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => emit('update:visible', visible),
})

watch(show, (newValue) => {
  if (newValue && props.developerInfo?.id) {
    model.value = {
      id: props.developerInfo.id,
      developerName: props.developerInfo.developerName,
      contactMail: props.developerInfo.contactMail || '',
      name: props.developerInfo.name || '',
      // paymentName: props.developerInfo.paymentName || '',
      // paymentQrcode: props.developerInfo.paymentQrcode || '',
      // paymentMethod: props.developerInfo.paymentMethod || '',
      status: props.developerInfo.status,
    }
  }
  else {
    model.value = { ...formInitValue }
  }
})

async function submit() {
  try {
    await update(model.value)
    emit('done')
  }
  catch (error) {
    apiRespErrMsgAndCustomCodeNeg1Msg(error, t('common.failed'))
  }
}

function handleValidateButtonClick(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      submit()
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 500px" title="编辑开发者">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="developerName" label="开发者标识">
        <NInput v-model:value="model.developerName" placeholder="纯英文，多词用-分割" />
      </NFormItem>

      <NFormItem path="paymentName" label="名字">
        <NInput v-model:value="model.name" placeholder="作者名字" />
      </NFormItem>

      <NFormItem path="contactMail" label="联系邮箱">
        <NInput v-model:value="model.contactMail" placeholder="请输入联系邮箱" />
      </NFormItem>

      <!-- <NFormItem path="paymentName" label="收款人姓名">
        <NInput v-model:value="model.paymentName" placeholder="请输入收款人真实姓名" />
      </NFormItem>

      <NFormItem path="paymentQrcode" label="收款二维码">
        <NInput v-model:value="model.paymentQrcode" placeholder="收款二维码图片URL" />
      </NFormItem>

      <NFormItem path="paymentMethod" label="收款方式">
        <NInput v-model:value="model.paymentMethod" placeholder="如：支付宝、微信" />
      </NFormItem> -->

      <NFormItem path="status" label="状态">
        <NSelect v-model:value="model.status" :options="statusOptions" />
      </NFormItem>
    </NForm>

    <template #footer>
      <NButton type="success" @click="handleValidateButtonClick">
        保存
      </NButton>
    </template>
  </NModal>
</template>
