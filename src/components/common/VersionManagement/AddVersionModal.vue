<script lang="ts" setup>
import { NButton, NInput, NModal, NSelect, NSpace, NUpload, useMessage } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { createVersion, submitReview, uploadVersionPackage } from '@/api/admin/microAppVersion'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'

interface Props {
  visible: boolean
  appRecordId: number
}

const props = defineProps<Props>()

const emit = defineEmits<Emits>()

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'done'): void
}

const message = useMessage()
const addVersionLoading = ref(false)
const uploadLoading = ref(false)
const uploadedConfig = ref<MicroApp.VersionConfig | null>(null)
const uploadCacheId = ref('')

// 语言列表
const langOptions = [
  { label: '简体中文 (zh-CN)', value: 'zh-CN' },
  { label: 'English (en-US)', value: 'en-US' },
  // { label: '日本語 (ja-JP)', value: 'ja-JP' },
  // { label: '한국어 (ko-KR)', value: 'ko-KR' },
]

// 表单数据 - 版本说明改为多语言格式
const versionForm = ref({
  version: '',
  versionCode: 0,
  packageUrl: '',
  packageHash: '',
  versionDescMap: {} as Record<string, string>, // 多语言版本说明
})

// 当前选中的语言
const currentLang = ref('')

// 可添加的语言选项（排除已添加的）
const availableLangOptions = computed(() => {
  const usedLangs = Object.keys(versionForm.value.versionDescMap)
  return langOptions.filter(l => !usedLangs.includes(l.value))
})

// // 选中语言
// function selectLanguage(lang: string) {
//   currentLang.value = lang
// }

// 添加语言
function addLanguage() {
  if (!currentLang.value)
    return
  // 初始化为空字符串
  versionForm.value.versionDescMap[currentLang.value] = ''
  currentLang.value = ''
}

// 删除语言
function removeLanguage(lang: string) {
  const newMap = { ...versionForm.value.versionDescMap }
  delete newMap[lang]
  versionForm.value.versionDescMap = newMap
}

// 获取当前语言的版本说明
function getCurrentVersionDesc(lang: string): string {
  return versionForm.value.versionDescMap[lang] || ''
}

// 设置当前语言的版本说明
function setCurrentVersionDesc(lang: string, value: string) {
  versionForm.value.versionDescMap[lang] = value
}

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
  versionForm.value = { version: '', versionCode: 0, packageUrl: '', packageHash: '', versionDescMap: {} }
  uploadedConfig.value = null
  currentLang.value = ''
}

// 处理文件上传
async function handleUploadChange(options: { file: any }) {
  const file = options.file.file
  if (!file)
    return

  uploadLoading.value = true
  try {
    const res = await uploadVersionPackage<MicroApp.MicroAppVersionUploadResp>(file, String(props.appRecordId))
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
      uploadCacheId.value = res.data.uploadCacheId
      message.success('上传成功')
    }
    else {
      apiRespErrMsg(res)
    }
  }
  catch (error: any) {
    if (error.code === -2) {
      message.error('上传失败，请检查微应用名称等信息是否匹配', { duration: 50000, closable: true })
      return
    }

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

// 转换为多语言格式提交
function formatVersionDesc(): Record<string, { content: string }> {
  const result: Record<string, { content: string }> = {}
  for (const [lang, content] of Object.entries(versionForm.value.versionDescMap)) {
    if (content.trim()) {
      result[lang] = { content: content.trim() }
    }
  }
  return result
}

// 提交添加版本
async function handleAddVersion() {
  addVersionLoading.value = true
  try {
    // 创建版本
    const createRes = await createVersion<any>({
      versionDesc: formatVersionDesc(),
      uploadCacheId: uploadCacheId.value,
    })

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
  catch (error) {
    apiRespErrMsg(error)
  }
  finally {
    addVersionLoading.value = false
  }
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px" title="添加版本" :mask-closable="false">
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

      <!-- 版本说明（多语言） -->
      <div>
        <div class="mb-2">
          版本说明
        </div>

        <!-- 已添加的语言卡片 -->
        <div v-if="Object.keys(versionForm.versionDescMap).length > 0" class="space-y-3 mb-3">
          <div
            v-for="lang in Object.keys(versionForm.versionDescMap)"
            :key="lang"
            class="border rounded-lg p-3 transition-all"
            :class="currentLang === lang ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-gray-300 cursor-pointer'"
            @click="currentLang = lang"
          >
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center gap-2">
                <div
                  class="w-4 h-4 rounded-full border-2 flex items-center justify-center transition-all"
                  :class="currentLang === lang ? 'border-blue-500 bg-blue-500' : 'border-gray-300'"
                >
                  <div v-if="currentLang === lang" class="w-2 h-2 rounded-full bg-white" />
                </div>
                <span class="font-medium">{{ langOptions.find(l => l.value === lang)?.label || lang }}</span>
              </div>
              <NButton size="tiny" quaternary @click.stop="removeLanguage(lang)">
                删除
              </NButton>
            </div>
            <NInput
              :value="getCurrentVersionDesc(lang)"
              type="textarea"
              :placeholder="`请输入${langOptions.find(l => l.value === lang)?.label || lang}版本说明`"
              :rows="2"
              @update:value="(v: string) => setCurrentVersionDesc(lang, v)"
              @click.stop
            />
          </div>
        </div>

        <!-- 添加语言选择器 -->
        <div v-if="availableLangOptions.length > 0" class="flex gap-2">
          <NSelect
            v-model:value="currentLang"
            :options="availableLangOptions"
            placeholder="选择要添加的语言"
            style="width: 200px"
            filterable
            size="small"
          />
          <NButton type="primary" :disabled="!currentLang" size="small" @click="addLanguage">
            添加
          </NButton>
        </div>

        <!-- 未添加任何语言时的提示 -->
        <div v-if="Object.keys(versionForm.versionDescMap).length === 0" class="text-center py-4 text-gray-400">
          请先添加语言，然后填写版本说明
        </div>
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
