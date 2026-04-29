<script lang="ts" setup>
import { useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { checkIsDeveloper, getInfo, register, updateMyInfo } from '@/api/developer'
import { DeveloperInfoForm } from '@/components/common'

const message = useMessage()
const formRef = ref<InstanceType<typeof DeveloperInfoForm> | null>(null)
const isDeveloper = ref(false)
const isEdit = ref(false)
const initialData = ref<Partial<Developer.RegisterRequest>>({})
const nameUpdatedAt = ref<string>('')

async function checkDeveloperStatus() {
  try {
    const res = await checkIsDeveloper<any>()
    isDeveloper.value = res.data?.isDeveloper

    if (isDeveloper.value) {
      const devRes = await getInfo<any>()
      if (devRes.data) {
        initialData.value = {
          developerName: devRes.data.developerName,
          contactMail: devRes.data.contactMail || '',
          paymentName: devRes.data.paymentName || '',
          paymentQrcode: devRes.data.paymentQrcode || '',
          paymentMethod: devRes.data.paymentMethod || '',
          name: devRes.data.name,
        }
        nameUpdatedAt.value = devRes.data.nameUpdatedAt || ''
        isEdit.value = true
      }
    }
  }
  catch (error) {
    console.error(error)
  }
}

async function handleSubmit(data: Developer.RegisterRequest, editMode: boolean) {
  formRef.value?.setLoading(true)
  try {
    if (editMode) {
      await updateMyInfo(data)
      message.success('更新成功')
      // 重新获取信息更新 nameUpdatedAt
      const devRes = await getInfo<any>()
      if (devRes.data?.nameUpdatedAt) {
        nameUpdatedAt.value = devRes.data.nameUpdatedAt
      }
    }
    else {
      await register(data)
      message.success('注册成功')
      isDeveloper.value = true
      isEdit.value = true
    }
  }
  catch (error: any) {
    message.error(error?.message || '操作失败')
  }
  finally {
    formRef.value?.setLoading(false)
  }
}

onMounted(() => {
  checkDeveloperStatus()
})
</script>

<template>
  <div class="min-h-[calc(100vh-120px)] flex items-center justify-center p-4">
    <DeveloperInfoForm
      ref="formRef"
      :edit-mode="isEdit"
      :is-developer="isDeveloper"
      :initial-data="initialData"
      :name-updated-at="nameUpdatedAt"
      @submit="handleSubmit"
    />
  </div>
</template>
