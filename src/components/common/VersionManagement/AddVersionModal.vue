<script lang="ts" setup>
import { NButton, NInput, NModal, NSpace, NUpload, useMessage } from 'naive-ui'
import { ref, watch } from 'vue'
import { createVersion, submitReview, uploadVersionPackage } from '@/api/admin/microAppVersion'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'

interface Props {
  visible: boolean
  appRecordId: number
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:visible': [value: boolean]
  'done': []
}>()

const message = useMessage()
const addVersionLoading = ref(false)
const uploadLoading = ref(false)
const uploadedConfig = ref<MicroApp.VersionConfig | null>(null)

// 表单数据
const versionForm = ref({
  version: '',
  versionCode: 0,
  packageUrl: '',
  packageHash: '',
  versionDesc: '',
})

// 双向绑定
const show = ref(props.visible)
watch(() => props.visible, (val) => {
  show.value = val
  if (!val) {
    resetForm()
  }
})
watch(show, (val) => {
  emit('update:visible', val)
})

// 重置表单
function resetForm() {
  versionForm.value = { version: '', versionCode: 0, packageUrl: '', packageHash: '', versionDesc: '' }
  uploadedConfig.value = null
}

// 处理文件上传
async function handleUploadChange(options: { file: any }) {
  const file = options.file.file
  if (!file)
    return

  uploadLoading.value = true
  try {
    const res = await uploadVersionPackage<any>(file)
    if (res.code === 0 && res.data) {
      versionForm.value.packageUrl = res.data.url
      versionForm.value.packageHash = res.data.hash || ''
      // 保存上传的配置信息
      if (res.data.config) {
        res.data.config.icon = res.data.iconURL || res.data.config.icon
      }
      uploadedConfig.value = res.data.config || null
      // 自动填充版本号
      if (res.data.config?.version) {
        versionForm.value.version = res.data.config.version
        handleVersionInput(res.data.config.version)
      }
      message.success('上传成功')
    }
    else {
      apiRespErrMsg(res)
    }
  }
  catch (error) {
    apiRespErrMsg(error)
  }
  finally {
    uploadLoading.value = false
  }
}

// 处理版本号输入
function handleVersionInput(value: string) {
  versionForm.value.version = value
  const parts = value.split('.')
  let code = 0
  if (parts.length >= 1)
    code += Number(parts[0]) * 100
  if (parts.length >= 2)
    code += Number(parts[1]) * 10
  if (parts.length >= 3)
    code += Number(parts[2])
  versionForm.value.versionCode = code
}

// 提交添加版本
async function handleAddVersion() {
  if (!versionForm.value.packageUrl || !versionForm.value.version) {
    message.warning('请上传版本包并填写版本号')
    return
  }

  addVersionLoading.value = true
  try {
    // 创建版本
    const createRes = await createVersion<any>({
      appRecordId: props.appRecordId,
      version: versionForm.value.version,
      versionCode: versionForm.value.versionCode,
      packageUrl: versionForm.value.packageUrl,
      packageHash: versionForm.value.packageHash || '',
      versionDesc: versionForm.value.versionDesc,
      config: uploadedConfig.value || undefined,
    })

    if (createRes.code === 0) {
      // 自动提交审核
      if (createRes.data?.id) {
        const reviewRes = await submitReview<any>({ versionId: createRes.data.id })
        if (reviewRes.code !== 0) {
          apiRespErrMsg(reviewRes)
          return
        }
      }

      message.success('版本添加成功，已提交审核')
      show.value = false
      emit('done')
    }
    else {
      apiRespErrMsg(createRes)
    }
  }
  catch (error) {
    apiRespErrMsg(error)
  }
  finally {
    addVersionLoading.value = false
  }
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 500px" title="添加版本" :mask-closable="false">
    <div class="space-y-4">
      <div>
        <div class="mb-2">
          选择版本包 <span class="text-red-500">*</span>
        </div>
        <NUpload
          accept=".zip"
          :max="1"
          :custom-request="(options: any) => handleUploadChange({ file: options.file })"
          :show-file-list="false"
        >
          <NButton :loading="uploadLoading">
            {{ uploadLoading ? '上传中...' : '选择文件' }}
          </NButton>
        </NUpload>
        <div v-if="uploadedConfig !== null" class="text-xs mt-1">
          <div v-if="uploadedConfig.icon" class="mb-2">
            <img :src="uploadedConfig.icon" class="w-12 h-12 object-contain border rounded">
          </div>
          <div v-if="uploadedConfig.appInfo?.['zh-CN']?.appName || uploadedConfig.appInfo?.['en-US']?.appName" class="text-lg font-bold">
            应用名称：{{ uploadedConfig.appInfo?.['zh-CN']?.appName || uploadedConfig.appInfo?.['en-US']?.appName }}
          </div>
          <div v-if="uploadedConfig.microAppId" class="text-gray-500">
            应用ID：{{ uploadedConfig.microAppId }}
          </div>
          <div v-if="uploadedConfig.author" class="text-gray-500">
            作者：{{ uploadedConfig.author }}
          </div>
          <div v-if="uploadedConfig.version" class="text-blue-500">
            版本号：{{ uploadedConfig.version }}
          </div>
          <div v-else class="text-orange-500">
            未检测到版本号，请手动填写
          </div>
        </div>
        <div class="text-xs text-gray-400 mt-1">
          支持 .zip 格式的微应用包
        </div>
      </div>
      <div>
        <div class="mb-2">
          版本号 <span class="text-red-500">*</span>
        </div>
        <NInput v-model:value="versionForm.version" placeholder="如：1.0.0" @update:value="handleVersionInput" />
      </div>
      <div>
        <div class="mb-2">
          版本说明
        </div>
        <NInput v-model:value="versionForm.versionDesc" type="textarea" placeholder="请输入版本说明" :rows="3" />
      </div>
    </div>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="show = false">
          取消
        </NButton>
        <NButton type="primary" :loading="addVersionLoading" @click="handleAddVersion">
          添加
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
