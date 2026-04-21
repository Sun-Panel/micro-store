<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NDivider, NForm, NFormItem, NImage, NInput, NInputNumber, NModal, NSelect, NSpace, NSwitch, NUpload, useMessage } from 'naive-ui'
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import { create, update } from '@/api/admin/microAppDeveloper'
import { microAppChargeTypeMap, microAppThirdChargeTypeMap } from '@/enums/panel'
import { t } from '@/locales'
import { useAppStore, useAuthStore } from '@/store'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'

interface Props {
  visible: boolean
  microAppInfo?: MicroApp.MicroAppReviewInfo
  categoryOptions: Category.Info[]
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}

const authStore = useAuthStore()
const appStore = useAppStore()
const locale = appStore.language
const message = useMessage()

// 表单初始值
const formInitValue: MicroApp.Info = {
  id: 0,
  adminName: '',
  microAppId: '',
  appName: '',
  appDesc: '',
  remark: '',
  appIcon: '',
  categoryId: 0 as any, // 使用 any 类型以便设置为 null
  chargeType: 0,
  points: 0,
  status: 0,
  thirdCharge: 0,
  haveIframe: false,
}

const model = ref({ ...formInitValue })
const formRef = ref<FormInst | null>(null)

// 积分数量默认值
const pointsDefaultValue = ref(1)

// 收费方式选项
const chargeTypeOptions = [
  { label: microAppChargeTypeMap[0], value: 0 },
  { label: `${microAppChargeTypeMap[1]}-(开发中)`, value: 1, disabled: true },
  { label: `${microAppChargeTypeMap[2]}-(开发中)`, value: 2, disabled: true },
]

// 第三方收费方式选项
const thirdChargeOptions = [
  { label: microAppThirdChargeTypeMap[0], value: 0 },
  { label: microAppThirdChargeTypeMap[1], value: 1 },
  { label: microAppThirdChargeTypeMap[2], value: 2 },
]

// 语言标签映射
const langLabelMap: Record<string, string> = {
  'zh-CN': '中文',
  'en-US': 'English',
  // 'ja-JP': '日本語',
  // 'ko-KR': '한국어',
}

// 语言选项
const langOptions = [
  { label: '中文', value: 'zh-CN' },
  { label: 'English', value: 'en-US' },
  // { label: '日本語', value: 'ja-JP' },
  // { label: '한국어', value: 'ko-KR' },
]

// 已添加的语言列表
const addedLangs = ref<string[]>([])
const newLang = ref<string | null>(null)

// 本地语言数据
const localLangMap = ref<Record<string, { appName: string, appDesc: string }>>({})

// 每个语言的折叠状态
const collapsedLangs = ref<Record<string, boolean>>({})

// 切换语言折叠状态
function toggleLangCollapsed(lang: string) {
  collapsedLangs.value[lang] = !collapsedLangs.value[lang]
}

// 获取文本摘要（前50个字符）
function getTextSummary(text: string, maxLength = 50) {
  if (!text)
    return '未填写'
  return text.length > maxLength ? `${text.slice(0, maxLength)}...` : text
}

// 表单验证规则 - 根据编辑或添加模式动态调整
const rules = computed<FormRules>(() => {
  const isEdit = !!props.microAppInfo?.id
  return {
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

// 监听收费方式变化，当切换到积分收费时设置默认值
watch(() => model.value.chargeType, (newType) => {
  if (newType === 1 && (!model.value.points || model.value.points === 0))
    model.value.points = pointsDefaultValue.value
})

// 监听弹窗打开/关闭
watch(show, (newValue) => {
  if (newValue) {
    if (props.microAppInfo?.id) {
      // 编辑模式
      model.value = {
        ...formInitValue,
        id: props.microAppInfo.id,
        microAppId: props.microAppInfo.microAppId,
        appName: props.microAppInfo.appName || '',
        appDesc: props.microAppInfo.appDesc || '',
        remark: props.microAppInfo.remark || '',
        appIcon: props.microAppInfo.appIcon,
        categoryId: props.microAppInfo.categoryId,
        chargeType: props.microAppInfo.chargeType,
        points: props.microAppInfo.points,
        adminName: props.microAppInfo.adminName,
        thirdCharge: props.microAppInfo.thirdCharge || 0,
        haveIframe: props.microAppInfo.haveIframe || false,
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
      // 设置默认积分数量
      model.value.points = pointsDefaultValue.value
      screenshotList.value = []
      // 创建模式：语言列表为空，会在 initLangData() 中自动添加系统语言
      addedLangs.value = []
      localLangMap.value = {}
      initLangData()
    }
  }
})

// 初始化多语言数据
function initLangData() {
  if (props.microAppInfo?.langMap && Object.keys(props.microAppInfo.langMap).length > 0) {
    // 编辑模式：使用已有的多语言数据
    localLangMap.value = props.microAppInfo.langMap
    addedLangs.value = Object.keys(props.microAppInfo.langMap)
    // 同步主语言数据
    if (addedLangs.value.length > 0) {
      const firstLang = addedLangs.value[0]
      model.value.appName = localLangMap.value[firstLang]?.appName || ''
      model.value.appDesc = localLangMap.value[firstLang]?.appDesc || ''
      // 默认所有语言都收起
      addedLangs.value.forEach((lang) => {
        collapsedLangs.value[lang] = true
      })
    }
  }
  else {
    // 如果没有语言，根据当前系统语言自动添加
    let defaultLang = locale || 'zh-CN'
    // 如果系统语言不在支持的语言列表中，默认使用中文
    if (!langOptions.find(l => l.value === defaultLang))
      defaultLang = 'zh-CN'

    addedLangs.value = [defaultLang]
    localLangMap.value = { [defaultLang]: { appName: '', appDesc: '' } }
    // 创建模式：第一项都为空时展开，否则收起
    const firstLangData = localLangMap.value[defaultLang]
    collapsedLangs.value[defaultLang] = !!firstLangData?.appName || !!firstLangData?.appDesc
  }
}

// 添加语言
function addLang(lang: string) {
  if (lang && !addedLangs.value.includes(lang)) {
    addedLangs.value.push(lang)
    localLangMap.value[lang] = { appName: '', appDesc: '' }
    // 新添加的语言默认展开
    collapsedLangs.value[lang] = false
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
    delete collapsedLangs.value[lang]
  }
}

// 获取可用语言选项
const availableLangOptions = computed(() => {
  return langOptions.filter(l => !addedLangs.value.includes(l.value))
})

// 分类选项转换
const categorySelectOptions = computed(() => {
  return props.categoryOptions?.map((item: Category.Info) => ({
    label: item.name,
    value: item.id,
  })) || []
})

// 监听分类选项，当创建时默认选中第一个分类
watch(categorySelectOptions, (options) => {
  if (options.length > 0 && (!model.value.id && (model.value.categoryId === 0 || !model.value.categoryId)))
    model.value.categoryId = options[0].value!
}, { immediate: true })

// 提交表单
async function submit() {
  try {
    // 从 screenshotList 获取所有有效的 URL
    const screenshotsStr = screenshotList.value
      .map(f => f.url)
      .filter((url: string) => url)
      .join(',')

    // 使用第一个语言的数据作为主语言数据
    const primaryLang = addedLangs.value[0]
    const primaryData = localLangMap.value[primaryLang] || { appName: '', appDesc: '' }

    if (model.value.id) {
      // 更新：同时提交主信息和多语言信息
      await update({
        id: model.value.id,
        appName: primaryData.appName,
        appIcon: model.value.appIcon,
        appDesc: primaryData.appDesc,
        remark: model.value.remark,
        categoryId: model.value.categoryId,
        chargeType: model.value.chargeType,
        points: model.value.points,
        screenshots: screenshotsStr,
        langMap: localLangMap.value,
        adminName: model.value.adminName,
        thirdCharge: model.value.thirdCharge,
        haveIframe: model.value.haveIframe,
      })
    }
    else {
      // 创建
      await create({
        microAppId: model.value.microAppId || undefined,
        appName: primaryData.appName,
        appIcon: model.value.appIcon,
        appDesc: primaryData.appDesc,
        remark: model.value.remark,
        categoryId: model.value.categoryId,
        chargeType: model.value.chargeType,
        points: model.value.points,
        screenshots: screenshotsStr,
        langMap: localLangMap.value,
        thirdCharge: model.value.thirdCharge,
        haveIframe: model.value.haveIframe,
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
function handleValidateButtonClick() {
  // 验证第一个语言的应用名称是否为空
  const primaryLang = addedLangs.value[0]
  if (primaryLang && (!localLangMap.value[primaryLang]?.appName || localLangMap.value[primaryLang]?.appName.trim() === '')) {
    message.error('请输入应用名称')
    return
  }

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
      else {
        // ...
        message.error(`图标上传失败：${res.msg}`)
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

// 上传完成后处理 - 设置文件 URL 或处理错误
function handleScreenshotFinish({ file, event }: { file: any, event?: any }) {
  // 从 event.target.response 中解析响应（参考图标上传的处理方式）
  let res: any = null
  if (event?.target?.response) {
    try {
      res = JSON.parse(event.target.response)
    }
    catch {
      // 如果解析失败，尝试从 file.response 获取
      res = file.response
    }
  }
  else {
    res = file.response
  }

  // 检查上传是否成功（后端成功响应为 code === 0）
  const isSuccess = res?.code === 0

  if (!isSuccess) {
    // 上传失败，从列表中移除该文件
    const index = screenshotList.value.findIndex(f => f.id === file.id)
    if (index !== -1)
      screenshotList.value.splice(index, 1)

    // 显示错误提示
    const errorMsg = res?.msg || res?.message || '上传失败'
    window.$message?.error(errorMsg)
    return
  }

  // 获取上传后的 URL
  const imageUrl = res?.data?.imageUrl || res?.data?.url
  if (imageUrl)
    file.url = imageUrl
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 800px" title="编辑基本信息">
    <NForm ref="formRef" :model="model" :rules="rules" require-mark-placement="left">
      <!-- 应用唯一标识：仅创建时显示 -->
      <!-- <NFormItem path="microAppId" label="应用ID (MicroAppID)">
        <NInput v-model:value="model.microAppId" placeholder="请输入应用标识（作者唯一标识-应用名称）" :disabled="microAppInfo?.id" />
      </NFormItem> -->

      <!-- 应用名称 -->
      <NFormItem path="adminName" label="名称 (仅开发者页面可见)">
        <NInput v-model:value="model.adminName" placeholder="请输入应用名称" size="small" />
      </NFormItem>

      <NDivider style="margin-top: 0px;" />

      <!-- 多语言设置 -->
      <div style="margin-bottom: 16px;">
        <div style="font-weight: bold; margin-bottom: 12px;">
          本土化设置
        </div>
        <NCard v-for="lang in addedLangs" :key="`lang-${lang}`" size="small" style="margin-bottom: 12px;">
          <!-- <template #header> -->
          <div class="flex justify-between items-center">
            <span class="font-bold">
              {{ langLabelMap[lang] || lang }}
            </span>
            <div class="flex items-center gap-2">
              <NButton text size="small" @click="toggleLangCollapsed(lang)">
                {{ collapsedLangs[lang] ? '编辑' : '收起' }}
              </NButton>
              <NButton v-if="addedLangs.length > 1" text size="small" type="error" @click="removeLang(lang)">
                删除
              </NButton>
            </div>
          </div>
          <!-- </template> -->

          <!-- 折叠状态：显示摘要 -->
          <div v-if="collapsedLangs[lang]" class="py-1">
            <div style="font-size: 13px; color: #666; margin-bottom: 4px;">
              <strong>名称：</strong>{{ getTextSummary(localLangMap[lang].appName) }}
            </div>
            <div style="font-size: 13px; color: #666;">
              <strong>简介：</strong>{{ getTextSummary(localLangMap[lang].appDesc) }}
            </div>
          </div>

          <!-- 展开状态：显示完整表单 -->
          <div v-else class="py-2">
            <NFormItem label="应用名称">
              <NInput
                v-model:value="localLangMap[lang].appName"
                :placeholder="`请输入应用名称 (${lang})`"
              />
            </NFormItem>
            <NFormItem label="应用简介">
              <NInput
                v-model:value="localLangMap[lang].appDesc"
                type="textarea"
                :placeholder="`请输入应用简介 (${lang})`"
                :rows="3"
              />
            </NFormItem>
          </div>
        </NCard>
      </div>

      <!-- 添加更多语言 -->
      <div v-if="availableLangOptions.length > 0" class="flex items-center gap-2 mb-4">
        <span style="font-size: 12px; color: #999;">添加其他语言：</span>
        <NSelect
          v-model:value="newLang"
          :options="availableLangOptions"
          placeholder="选择语言"
          style="width: 150px;"
          @update:value="addLang"
        />
      </div>

      <NDivider style="margin: 16px 0;" />

      <!-- 应用图标 -->
      <NFormItem path="appIcon" label="应用图标">
        <div>
          <div class="text-xs text-gray-400 flex items-center gap-3">
            <span>尺寸要求：1:1 正方形，最大 512KB</span>
          </div>
          <div class="flex items-center gap-2">
            <div v-if="model.appIcon" style="margin-top: 10px;width: 100%;">
              <NImage :src="model.appIcon" width="80" height="80" />
            </div>
            <div class="flex items-center" style="width: 100%;">
              <NInput v-show="false" v-model:value="model.appIcon" placeholder="图标URL或上传" style="flex: 1;" />
              <NUpload
                action="/api/admin/developer/myMicroApp/uploadIcon"
                :headers="{ token: authStore.token }"
                name="iconfile"
                :show-file-list="false"
                accept=".png,.jpg,.jpeg,.svg,.ico"
                style="margin-left: 10px;"
                @finish="handleIconFinish"
                @error="handleIconError"
              >
                <NButton size="small">
                  上传图片
                </NButton>
              </NUpload>
            </div>
          </div>
        </div>
      </NFormItem>

      <!-- 所属分类 -->
      <NFormItem path="categoryId" label="所属分类">
        <NSelect v-model:value="model.categoryId" :options="categorySelectOptions" placeholder="请选择应用分类" clearable />
      </NFormItem>

      <!-- 收费方式 -->
      <NFormItem path="chargeType" label="收费方式">
        <NSelect v-model:value="model.chargeType" :options="chargeTypeOptions" placeholder="请选择收费方式" />
      </NFormItem>

      <!-- 积分数量：仅收费方式为积分时显示 -->
      <NFormItem v-if="model.chargeType === 1" path="price" label="积分数量">
        <NInputNumber v-model:value="model.points" :min="1" :precision="0" style="width: 100%;" placeholder="请输入积分数量" />
      </NFormItem>

      <!-- 第三方收费方式 -->
      <NFormItem :label="t('microApp.thirdCharge')">
        <NSelect v-model:value="model.thirdCharge" :options="thirdChargeOptions" placeholder="请选择第三方收费方式" />
      </NFormItem>

      <!-- 是否包含iframe -->
      <NFormItem label="是否包含iframe">
        <NSwitch v-model:value="model.haveIframe" />
      </NFormItem>

      <!-- 应用备注 -->
      <!-- <NFormItem label="应用备注">
        <NInput v-model:value="model.remark" type="textarea" placeholder="请输入应用备注" :rows="2" />
      </NFormItem> -->

      <!-- 应用图集 -->
      <NFormItem label="应用图集">
        <template #label>
          应用图集
          <div class="text-xs text-gray-400 flex items-center gap-3">
            <span>尺寸要求：4:3 比例，最大 2MB</span>
            <span>推荐尺寸：</span>
            <span class="bg-gray-100 px-2 rounded">1920×1440</span>
            <span class="bg-gray-100 px-2 rounded">1600×1200</span>
            <span class="bg-gray-100 px-2 rounded">1280×960</span>
            <span class="bg-gray-100 px-2 rounded">800×600</span>
          </div>
        </template>

        <!-- 使用 v-model:file-list 管理所有图片 -->
        <NUpload
          v-model:file-list="screenshotList"
          action="/api/admin/developer/myMicroApp/uploadScreenshot"
          :headers="{ token: authStore.token }"
          name="screenshotfile"
          accept=".png,.jpg,.jpeg,.webp,.gif"
          list-type="image-card"
          :max="5"
          @finish="handleScreenshotFinish"
          @remove="handleScreenshotRemove"
        >
          <div>上传截图</div>
        </NUpload>
      </NFormItem>
    </NForm>

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
