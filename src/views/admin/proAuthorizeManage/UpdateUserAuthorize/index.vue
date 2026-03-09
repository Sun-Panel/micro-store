<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, useMessage } from 'naive-ui'
import { updateUserExpiredTimeByDay } from '@/api/admin/proAuthorize'

interface Props {
  visible: boolean
  userId: number
}

interface Form {
  dayNum: number
  note: string
  userId?: number
  adminNote?: string
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const message = useMessage()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', id: number): void// 创建完成
}

const model = ref<Form>({
  userId: 0,
  dayNum: 0,
  note: '',
  adminNote: '',
})
const formRef = ref<FormInst | null>(null)

const rules: FormRules = {
  dayNum: {
    required: true,
    trigger: 'blur',
    message: '必填项',
    type: 'number',
  },
  note: {
    required: true,
    trigger: 'blur',
    message: '必填项',
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
  if (props.userId) {
    model.value = {
      dayNum: 0,
      note: '',
    }
  }
  else {
    message.error('出错了，该商品不存在')
  }
})

async function handleUpdateUserExpiredTimeByDay() {
  const saveData: Admin.ProAuthorize.ProAuthorizeUpdateUserExpiredTimeByDayReq = {
    userId: props.userId,
    dayNum: model.value.dayNum,
    note: model.value.note,
    adminNote: model.value.adminNote,
  }
  try {
    const res = await updateUserExpiredTimeByDay(saveData)
    if (res.code === 0)
      emit('done', 0)

    else if (res.code !== -1)
      message.warning('操作失败')
  }
  catch (error) {
    console.error(error)
    message.warning('操作失败')
  }
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors && props.userId)
      handleUpdateUserExpiredTimeByDay()
    else
      console.log(errors)
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px" title="修改">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="dayNum" label="天数">
        <NInputNumber
          v-model:value="model.dayNum"
        />
      </NFormItem>

      <NFormItem path="note" label="备注">
        <NInput
          v-model:value="model.note"
          placeholder="改动的原因"
        />
      </NFormItem>

      <NFormItem path="adminNote" label="管理员备注">
        <NInput
          v-model:value="model.adminNote"
          placeholder="管理员备注，对用户不可见"
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
