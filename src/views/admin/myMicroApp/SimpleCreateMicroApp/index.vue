<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NModal, NSpace } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { create } from '@/api/admin/microAppDeveloper'
import { getInfo } from '@/api/developer'
import { t } from '@/locales'
import { apiRespErrMsgAndCustomCodeNeg1Msg, message } from '@/utils/cmn/apiMessage'

interface Props {
  visible: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}

// 表单数据
const model = ref({
  microAppId: '',
  adminName: '',
})

const formRef = ref<FormInst | null>(null)
const developerName = ref<string>('')

// 表单验证规则
const rules = computed<FormRules>(() => ({
  microAppId: [
    { required: true, trigger: 'blur', message: '请输入微应用ID' },
    {
      validator: (_rule, value) => {
        if (!value) return true
        if (!developerName.value) {
          return new Error('获取开发者信息失败，请刷新重试')
        }
        if (!value.startsWith(`${developerName.value}-`)) {
          return new Error(`微应用ID必须以"${developerName.value}-"开头`)
        }
        if (value === `${developerName.value}-`) {
          return new Error('微应用ID不能仅包含开发者标识，请添加应用名称')
        }
        if (value.endsWith('-')) {
          return new Error('微应用ID不能以"-"结尾')
        }
        return true
      },
      trigger: 'blur',
    },
  ],
  appName: [{ required: true, trigger: 'blur', message: '请输入微应用名称' }],
}))

// 弹窗显示状态
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => emit('update:visible', visible),
})

// 获取开发者信息
async function fetchDeveloperInfo() {
  try {
    const res = await getInfo<any>()
    if (res.data?.developerName) {
      developerName.value = res.data.developerName
      // 自动填充前缀到微应用ID输入框
      if (!model.value.microAppId) {
        model.value.microAppId = `${developerName.value}-`
      }
    }
  }
  catch (error) {
    console.error('获取开发者信息失败:', error)
  }
}

// 监听弹窗打开/关闭
watch(show, async (newValue) => {
  if (newValue) {
    // 打开弹窗时清空表单并获取开发者信息
    model.value = {
      microAppId: '',
      adminName: '',
    }
    await fetchDeveloperInfo()
  }
})

// 提交表单
async function submit() {
  try {
    await create({
      microAppId: model.value.microAppId,
      appName: model.value.adminName,
    })
    emit('done')
    show.value = false
  }
  catch (error:any) {
    if(error.code===3000){
      message.error("微应用ID已存在",{duration:20000,closable:true})
      return 
    }
    apiRespErrMsgAndCustomCodeNeg1Msg(error, t('common.failed'))
  }
}

// 表单验证并提交
function handleValidateButtonClick(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      submit()
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 500px" title="创建微应用">
    <NForm ref="formRef" :model="model" :rules="rules">
      <!-- 微应用ID -->
      <NFormItem path="microAppId" label="微应用ID">
        <NInput v-model:value="model.microAppId" placeholder="请输入微应用ID（开发者标识-应用名称，不能以-结尾）" />
      </NFormItem>

      <!-- 微应用名称 -->
      <NFormItem path="appName" label="应用名称 (仅开发者页面可见)">
        <NInput v-model:value="model.adminName" placeholder="请输入微应用名称" />
      </NFormItem>
    </NForm>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="show = false">
          取消
        </NButton>
        <NButton type="success" @click="handleValidateButtonClick">
          创建
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
