<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NForm, NFormItem, NImage, NInput, NInputNumber, NModal, NSelect, NSpace, NTabPane, NTabs, NUpload } from 'naive-ui'
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import { create, update } from '@/api/admin/microAppDeveloper'
import { microAppChargeTypeMap } from '@/enums/panel'
import { t } from '@/locales'
import { useAuthStore } from '@/store'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'

interface Props {
  visible: boolean
  microAppInfo?: MicroApp.MicroAppInfo
  categoryOptions: { label: string, value: number }[]
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}

const authStore = useAuthStore()

// 表单初始值
const formInitValue = {
  id: 0,
  microAppId: '',
  appName: '',
  appDesc: '',
  remark: '',
  appIcon: '',
  categoryId: 0,
  chargeType: 0,
  price: 0,
}

const model = ref({ ...formInitValue })
const formRef = ref<FormInst | null>(null)

// 收费方式选项
const chargeTypeOptions = [
  { label: microAppChargeTypeMap[0], value: 0 },
  { label: microAppChargeTypeMap[1], value: 1 },
  { label: microAppChargeTypeMap[2], value: 2 },
]

// 语言标签映射
const langLabelMap: Record<string, string> = {
  'zh-CN': '中文',
  'en-US': 'English',
  'ja-JP': '日本語',
  'ko-KR': '한국어',
}

// 语言选项
const langOptions = [
  { label: '中文', value: 'zh-CN' },
  { label: 'English', value: 'en-US' },
  { label: '日本語', value: 'ja-JP' },
  { label: '한국어', value: 'ko-KR' },
]

// 已添加的语言列表
const addedLangs = ref<string[]>([])
const newLang = ref<string | null>(null)

// 本地语言数据
const localLangMap = ref<Record<string, { appName: string, appDesc: string }>>({})

// 表单验证规则 - 根据编辑或添加模式动态调整
const rules = computed<FormRules>(() => {
  const isEdit = !!props.microAppInfo?.id
  return {
    appName: [{ required: true, trigger: 'blur', message: '请输入应用名称' }],
    appIcon: isEdit ? [{ required: true, trigger: 'blur', message: '请上传应用图标' }] : [],
    microAppId: !isEdit ? [{ required: true, trigger: 'blur', message: '请输入应用标识' }] : [],
    categoryId: [{ required: true, type: 'number', trigger: 'change', message: '请选择分类' }],
  }
})

// 弹窗显示状态
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => emit('update:visible', visible),
})

// 用于 NUpload 显示的文件列表
const screenshotList = ref<any[]>([])

// 监听弹窗打开/关闭
watch(show, (newValue) => {
  if (newValue) {
    if (props.microAppInfo?.id) {
      // 编辑模式
      model.value = {
        id: props.microAppInfo.id,
        microAppId: props.microAppInfo.microAppId,
        appName: props.microAppInfo.appName || '',
        appDesc: props.microAppInfo.appDesc || '',
        remark: props.microAppInfo.remark || '',
        appIcon: props.microAppInfo.appIcon,
        categoryId: props.microAppInfo.categoryId,
        chargeType: props.microAppInfo.chargeType,
        price: props.microAppInfo.price,
      }
      // 初始化已有图片列表
      const screenshots = props.microAppInfo.screenshots ? props.microAppInfo.screenshots.split(',').filter(Boolean) : []
      screenshotList.value = screenshots.map((url: string, index: number) => ({
        id: String(index),
        name: url.split('/').pop() || `screenshot-${index}`,
        url,
        status: 'finished',
      }))
      // 初始化多语言数据
      initLangData()
    }
    else {
      // 创建模式
      model.value = { ...formInitValue }
      screenshotList.value = []
      // 创建模式默认添加中文
      addedLangs.value = ['zh-CN']
      localLangMap.value = { 'zh-CN': { appName: '', appDesc: '' } }
    }
  }
})

// 初始化多语言数据
function initLangData() {
  if (props.microAppInfo?.langList && props.microAppInfo.langList.length > 0) {
    const map: Record<string, { appName: string, appDesc: string }> = {}
    const langs: string[] = []
    props.microAppInfo.langList.forEach((lang) => {
      map[lang.lang] = {
        appName: lang.appName,
        appDesc: lang.appDesc,
      }
      langs.push(lang.lang)
    })
    localLangMap.value = map
    addedLangs.value = langs
  }
  else {
    // 如果没有语言，默认添加中文
    addedLangs.value = ['zh-CN']
    localLangMap.value = { 'zh-CN': { appName: '', appDesc: '' } }
  }
}

// 添加语言
function addLang(lang: string) {
  if (lang && !addedLangs.value.includes(lang)) {
    addedLangs.value.push(lang)
    localLangMap.value[lang] = { appName: '', appDesc: '' }
  }
  newLang.value = null
}

// 移除语言
function removeLang(lang: string) {
  if (addedLangs.value.length <= 1)
    return
  const index = addedLangs.value.indexOf(lang)
  if (index > -1) {
    addedLangs.value.splice(index, 1)
    delete localLangMap.value[lang]
  }
}

// 获取可用语言选项
const availableLangOptions = computed(() => {
  return langOptions.filter(l => !addedLangs.value.includes(l.value))
})

// 提交表单
async function submit() {
  try {
    // 从 screenshotList 获取所有有效的 URL
    const screenshotsStr = screenshotList.value
      .map(f => f.url)
      .filter((url: string) => url)
      .join(',')

    if (model.value.id) {
      // 更新：同时提交主信息和多语言信息
      await update({
        id: model.value.id,
        appName: model.value.appName,
        appIcon: model.value.appIcon,
        appDesc: model.value.appDesc,
        remark: model.value.remark,
        categoryId: model.value.categoryId,
        chargeType: model.value.chargeType,
        price: model.value.price,
        screenshots: screenshotsStr,
        langMap: localLangMap.value,
      })
    }
    else {
      // 创建
      await create({
        microAppId: model.value.microAppId || undefined,
        appName: model.value.appName,
        appIcon: model.value.appIcon,
        appDesc: model.value.appDesc,
        remark: model.value.remark,
        categoryId: model.value.categoryId,
        chargeType: model.value.chargeType,
        points: model.value.price,
        screenshots: screenshotsStr,
        langMap: localLangMap.value,
      })
    }
    emit('done')
    show.value = false
  }
  catch (error) {
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

// ========== 图标上传处理 ==========
function handleIconFinish(data: any) {
  if (data.event && data.event.target) {
    const xhr = data.event.target
    const response = xhr.response
    if (response) {
      const res = JSON.parse(response)
      if (res.code === 0 && res.data && res.data.imageUrl) {
        model.value.appIcon = res.data.imageUrl
      }
    }
  }
}

function handleIconError(data: any) {
  console.error('Upload error:', data)
}

// ========== 图集上传处理 ==========

// 删除图片 - NUpload 自动更新 file-list，不需要手动处理
function handleScreenshotRemove() {
  // v-model:file-list 会自动更新 screenshotList，提交时直接使用即可
}

// 上传完成后处理 - 设置文件 URL
function handleScreenshotFinish({ file, event }: { file: any, event?: any }) {
  // 尝试从响应中获取上传后的 URL
  let imageUrl = file.response?.data?.imageUrl || file.response?.data?.url
  // 如果没有，尝试从 event 获取
  if (!imageUrl && event?.target?.response) {
    try {
      const res = JSON.parse(event.target.response)
      imageUrl = res.data?.imageUrl || res.data?.url
    }
    catch (e) {}
  }
  if (imageUrl) {
    file.url = imageUrl
  }
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 800px" :title="`${microAppInfo?.id ? '编辑' : '创建'}微应用`">
    <NTabs type="line" animated>
      <!-- 基本信息标签页 -->
      <NTabPane name="basic" tab="基本信息">
        <NForm ref="formRef" :model="model" :rules="rules">
          <!-- 应用唯一标识：仅创建时显示 -->
          <NFormItem v-if="!microAppInfo?.id" path="microAppId" label="应用ID (MicroAppID)">
            <NInput v-model:value="model.microAppId" placeholder="请输入应用标识（作者唯一标识-应用名称）" />
          </NFormItem>

          <!-- 应用名称 -->
          <NFormItem path="appName" label="应用名称">
            <NInput v-model:value="model.appName" placeholder="请输入应用名称" />
          </NFormItem>

          <!-- 应用图标 -->
          <NFormItem path="appIcon" label="应用图标">
            <div v-if="model.appIcon" style="margin-top: 10px;">
              <NImage :src="model.appIcon" width="80" height="80" />
            </div>
            <div class="flex items-center" style="width: 100%;">
              <NInput v-show="false" v-model:value="model.appIcon" placeholder="图标URL或上传" style="flex: 1;" />
              <NUpload
                action="/api/file/uploadImg"
                :headers="{ token: authStore.token }"
                name="imgfile"
                :show-file-list="false"
                accept="image/*"
                style="margin-left: 10px;"
                @finish="handleIconFinish"
                @error="handleIconError"
              >
                <NButton size="small">
                  上传图片
                </NButton>
              </NUpload>
            </div>
          </NFormItem>

          <!-- 所属分类 -->
          <NFormItem path="categoryId" label="所属分类">
            <NSelect v-model:value="model.categoryId" :options="categoryOptions" placeholder="请选择分类" />
          </NFormItem>

          <!-- 收费方式 -->
          <NFormItem path="chargeType" label="收费方式">
            <NSelect v-model:value="model.chargeType" :options="chargeTypeOptions" />
          </NFormItem>

          <!-- 积分数量：仅收费方式为积分时显示 -->
          <NFormItem v-if="model.chargeType === 1" path="price" label="积分数量">
            <NInputNumber v-model:value="model.price" :min="1" :precision="0" style="width: 100%;" />
          </NFormItem>

          <!-- 应用备注 -->
          <NFormItem label="应用备注">
            <NInput v-model:value="model.remark" type="textarea" placeholder="请输入应用备注" :rows="2" />
          </NFormItem>

          <!-- 应用图集 -->
          <NFormItem label="应用图集">
            <!-- 使用 v-model:file-list 管理所有图片 -->
            <NUpload
              v-model:file-list="screenshotList"
              action="/api/file/uploadImg"
              :headers="{ token: authStore.token }"
              name="imgfile"
              accept="image/*"
              list-type="image-card"
              :max="5"
              @finish="handleScreenshotFinish"
              @remove="handleScreenshotRemove"
            >
              <div>上传截图</div>
            </NUpload>
          </NFormItem>
        </NForm>
      </NTabPane>

      <!-- 多语言设置标签页 -->
      <NTabPane name="lang" tab="多语言设置">
        <div class="space-y-4">
          <!-- 添加语言 -->
          <div class="flex items-center gap-2 mb-4">
            <span>添加语言：</span>
            <NSelect
              v-model:value="newLang"
              :options="availableLangOptions"
              placeholder="选择语言"
              style="width: 150px;"
              @update:value="addLang"
            />
          </div>

          <!-- 语言列表 -->
          <NCard v-for="lang in addedLangs" :key="lang" size="small">
            <template #header>
              <div class="flex justify-between items-center">
                <span class="font-bold">{{ langLabelMap[lang] || lang }}</span>
                <NButton v-if="addedLangs.length > 1" text size="small" type="error" @click="removeLang(lang)">
                  删除
                </NButton>
              </div>
            </template>

            <NForm label-placement="top">
              <NFormItem :label="`${langLabelMap[lang] || lang} 应用名称`">
                <NInput
                  v-model:value="localLangMap[lang].appName"
                  :placeholder="`请输入应用名称 (${lang})`"
                />
              </NFormItem>
              <NFormItem :label="`${langLabelMap[lang] || lang} 应用简介`">
                <NInput
                  v-model:value="localLangMap[lang].appDesc"
                  type="textarea"
                  :placeholder="`请输入应用简介 (${lang})`"
                  :rows="2"
                />
              </NFormItem>
            </NForm>
          </NCard>
        </div>
      </NTabPane>
    </NTabs>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="show = false">
          取消
        </NButton>
        <NButton type="success" @click="handleValidateButtonClick">
          保存
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
