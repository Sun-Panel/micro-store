<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal } from 'naive-ui'
import { edit as updateApi } from '@/api/admin/systemVariable'
import { t } from '@/locales'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'

interface Props {
  visible: boolean
  info: SystemVariable.SystemVariableEditReq | null
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}

const formInitValue = {
  description: '',
  name: '',
  value: '',
}

const model = ref<SystemVariable.SystemVariableEditReq>(formInitValue)
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  name: [
    {
      required: true,
      trigger: 'blur',
      message: t('form.required'),
    },
  ],
  description: [
    {
      required: true,
      trigger: 'blur',
      message: t('form.required'),
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
  await updateApi(model.value).then(() => {
    emit('done')
    show.value = false
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
  <NModal v-model:show="show" preset="card" style="width: 600px;" :title="model.id ? '修改' : '添加'">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="name" label="变量名">
        <NInput v-model:value="model.name" type="text" placeholder="仅支持英文" />
      </NFormItem>

      <NFormItem path="description" label="描述">
        <NInput v-model:value="model.description" type="text" placeholder="模板名称" />
      </NFormItem>
    </NForm>

    <template #footer>
      <NButton type="success" @click="handleValidateButtonClick">
        {{ t('common.save') }}
      </NButton>
    </template>
  </NModal>
</template>
