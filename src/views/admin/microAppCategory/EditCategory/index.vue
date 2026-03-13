<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NSelect } from 'naive-ui'
import { create, update } from '@/api/admin/microAppCategory'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'
import { t } from '@/locales'

interface Props {
  visible: boolean
  categoryInfo?: MicroAppCategory.CategoryInfo
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}

const formInitValue: MicroAppCategory.CreateRequest = {
  name: '',
  icon: '',
  sort: 0,
  status: 1,
}

const model = ref<MicroAppCategory.CreateRequest & { id?: number }>({ ...formInitValue })
const formRef = ref<FormInst | null>(null)

const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 },
]

const rules: FormRules = {
  name: [
    {
      required: true,
      trigger: 'blur',
      message: '请输入分类名称',
    },
  ],
}

const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

watch(show, (newValue) => {
  if (newValue) {
    if (props.categoryInfo?.id) {
      model.value = {
        id: props.categoryInfo.id,
        name: props.categoryInfo.name,
        icon: props.categoryInfo.icon,
        sort: props.categoryInfo.sort,
        status: props.categoryInfo.status,
      }
    }
    else {
      model.value = { ...formInitValue }
    }
  }
})

const submit = async () => {
  try {
    if (model.value.id) {
      await update({
        id: model.value.id,
        name: model.value.name,
        icon: model.value.icon,
        sort: model.value.sort,
        status: model.value.status,
      })
    }
    else {
      await create({
        name: model.value.name,
        icon: model.value.icon,
        sort: model.value.sort,
        status: model.value.status,
      })
    }
    emit('done')
  }
  catch (error) {
    apiRespErrMsgAndCustomCodeNeg1Msg(error, t('common.failed'))
  }
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      submit()
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 500px" :title="`${categoryInfo?.id ? '编辑' : '添加'}分类`">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="name" label="分类名称">
        <NInput v-model:value="model.name" type="text" placeholder="请输入分类名称" />
      </NFormItem>

      <NFormItem path="icon" label="分类图标">
        <NInput v-model:value="model.icon" type="text" placeholder="请输入图标URL或图标名称" />
      </NFormItem>

      <NFormItem path="sort" label="排序">
        <NInputNumber v-model:value="model.sort" :min="0" placeholder="数字越大越靠前" style="width: 100%;" />
      </NFormItem>

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
