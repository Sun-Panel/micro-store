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

// 图片上传：复用项目的 /api/file/uploadImg（字段 imgfile，响应 data.imageUrl）
// 注意：vditor 3.x 的 upload 契约为 url + format 适配响应，旧版
// (event, files, success, failure) 的 handler 写法已失效，会导致上传静默失败。
// 工具栏使用与默认一致的结构，仅移除录音（record）按钮。
const vditorOptions = {
  cache: { enable: false },
  height: 460,
  toolbar: [
    'emoji', 'headings', 'bold', 'italic', 'strike', 'link', '|',
    'list', 'ordered-list', 'check', 'outdent', 'indent', '|',
    'quote', 'line', 'code', 'inline-code', 'insert-before', 'insert-after', '|',
    'upload', 'table', '|',
    'undo', 'redo', '|',
    'fullscreen', 'edit-mode',
    { name: 'more', toolbar: ['both', 'code-theme', 'content-theme', 'export', 'outline', 'preview', 'devtools', 'info', 'help'] },
  ],
  upload: {
    url: '/api/file/uploadImg',
    fieldName: 'imgfile',
    max: 500 * 1024, // 500KB，上传尺寸限制
    multiple: false,
    // 动态注入登录 token（后端通过 token 请求头鉴权）
    setHeaders: () => ({ token: authStore.token || '' }),
    // 将后端返回 {code, data:{imageUrl}} 适配为 vditor 期望的 {code, data:{succMap}}
    format: (files: File[], responseText: string) => {
      try {
        const res = JSON.parse(responseText)
        if (res.code === 0 && res.data?.imageUrl) {
          const name = files?.[0]?.name || 'image'
          return JSON.stringify({
            code: 0,
            data: { succMap: { [name]: res.data.imageUrl }, errFiles: [] },
          })
        }
        return JSON.stringify({
          code: 1,
          msg: res.msg || '上传失败',
          data: { succMap: {}, errFiles: [] },
        })
      }
      catch {
        return JSON.stringify({
          code: 1,
          msg: '上传失败',
          data: { succMap: {}, errFiles: [] },
        })
      }
    },
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
