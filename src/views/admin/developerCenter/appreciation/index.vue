<script setup lang="ts">
import { NButton, NCard, NSpace, NAlert, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import Vditor from '@/components/common/Vditor/index.vue'
import { getInfo, updateMyInfo } from '@/api/developer'
import { useAuthStore } from '@/store/modules/auth'

const message = useMessage()
const authStore = useAuthStore()

// 当前开发者的完整信息（保存时需全量回传，避免覆盖其它字段）
const devInfo = ref<Developer.DeveloperInfo | null>(null)
const content = ref('')
const loading = ref(false)
const saving = ref(false)

// 自定义图片上传：复用项目的 /api/file/uploadImg（字段 imgfile，响应 data.imageUrl）
function uploadHandler(
  _event: Event,
  files: File[],
  success: (md: string) => void,
  failure: (msg: string) => void,
) {
  const file = files[0]
  if (!file) {
    failure('请选择图片')
    return
  }
  const formData = new FormData()
  formData.append('imgfile', file)

  const xhr = new XMLHttpRequest()
  xhr.open('POST', '/api/file/uploadImg')
  xhr.setRequestHeader('token', authStore.token || '')
  xhr.onload = () => {
    if (xhr.status === 200) {
      try {
        const res = JSON.parse(xhr.responseText)
        if (res.code === 0 && res.data?.imageUrl)
          success(`![](${res.data.imageUrl})`)
        else
          failure(res.message || '上传失败')
      }
      catch {
        failure('上传解析失败')
      }
    }
    else {
      failure('上传失败')
    }
  }
  xhr.onerror = () => failure('上传失败')
  xhr.send(formData)
}

const vditorOptions = {
  cache: { enable: false },
  height: 460,
  upload: {
    handler: uploadHandler,
  },
}

async function load() {
  loading.value = true
  try {
    const { data } = await getInfo<Developer.DeveloperInfo>()
    devInfo.value = data || null
    content.value = data?.rewardContent || ''
  }
  catch {
    message.error('加载赞赏信息失败')
  }
  finally {
    loading.value = false
  }
}

async function save() {
  if (!devInfo.value) {
    message.warning('请稍候，信息加载中…')
    return
  }
  saving.value = true
  try {
    await updateMyInfo({
      developerName: devInfo.value.developerName,
      name: devInfo.value.name,
      contactMail: devInfo.value.contactMail || '',
      paymentName: devInfo.value.paymentName || '',
      paymentQrcode: devInfo.value.paymentQrcode || '',
      paymentMethod: devInfo.value.paymentMethod || '',
      rewardContent: content.value,
    })
    message.success('保存成功')
  }
  catch (error: any) {
    message.error(error?.message || '保存失败')
  }
  finally {
    saving.value = false
  }
}

onMounted(() => {
  load()
})
</script>

<template>
  <div class="p-[20px] max-w-[900px] mx-auto">
    <NCard title="赞赏信息" :bordered="false">
      <NAlert type="info" class="mb-[16px]">
        这里的内容会在微应用 / 自定义代码详情页的「赞赏」弹窗中展示。支持 Markdown：可写感谢语、插入图片（点击编辑器工具栏的图片按钮上传收款码），链接会自动变为可点击形式。
      </NAlert>

      <Vditor
        v-model="content"
        :options="vditorOptions"
      />

      <NSpace class="mt-[16px]" justify="end">
        <NButton
          type="primary"
          :loading="saving"
          @click="save"
        >
          保存
        </NButton>
      </NSpace>
    </NCard>
  </div>
</template>
