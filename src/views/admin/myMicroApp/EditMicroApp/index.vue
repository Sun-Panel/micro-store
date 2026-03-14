<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NImage, NInput, NInputNumber, NModal, NSelect, NSpace, NUpload } from 'naive-ui'
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import { create, update } from '@/api/admin/microApp'
import { microAppChargeTypeMap } from '@/enums/panel'
import { t } from '@/locales'
import { useAuthStore } from '@/store'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'

interface Props {
  visible: boolean
  microAppInfo?: MicroApp.MicroAppInfo
  authorId: number
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
  screenshots: [] as string[],
  authorId: 0,
}

const model = ref({ ...formInitValue })
const formRef = ref<FormInst | null>(null)

// 收费方式选项
const chargeTypeOptions = [
  { label: microAppChargeTypeMap[0], value: 0 },
  { label: microAppChargeTypeMap[1], value: 1 },
  { label: microAppChargeTypeMap[2], value: 2 },
]

// 表单验证规则
const rules: FormRules = {
  appName: [{ required: true, trigger: 'blur', message: '请输入应用名称' }],
  appIcon: [{ required: true, trigger: 'blur', message: '请上传应用图标' }],
  categoryId: [{ required: true, type: 'number', trigger: 'change', message: '请选择分类' }],
}

// 弹窗显示状态
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => emit('update:visible', visible),
})

// 上传图集列表
const screenshotList = ref<{ id: string, url: string, status?: string }[]>([])

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
        screenshots: props.microAppInfo.screenshots ? props.microAppInfo.screenshots.split(',').filter(Boolean) : [],
        authorId: props.authorId,
      }
      screenshotList.value = model.value.screenshots.map((url, index) => ({ id: String(index), url }))
    }
    else {
      // 创建模式
      model.value = { ...formInitValue, authorId: props.authorId }
      screenshotList.value = []
    }
  }
})

// 提交表单
async function submit() {
  try {
    const screenshotsStr = model.value.screenshots.join(',')

    if (model.value.id) {
      // 更新
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
        price: model.value.price,
        authorId: props.authorId,
        screenshots: screenshotsStr,
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
function handleScreenshotFinish(data: any) {
  if (data.event && data.event.target) {
    const xhr = data.event.target
    const response = xhr.response
    if (response) {
      const res = JSON.parse(response)
      if (res.code === 0 && res.data && res.data.imageUrl) {
        screenshotList.value.push({ id: data.file.id, url: res.data.imageUrl, status: 'finished' })
        model.value.screenshots = screenshotList.value.map((s: any) => s.url)
      }
    }
  }
}

function handleScreenshotRemove(data: any) {
  screenshotList.value = screenshotList.value.filter((s: any) => s.id !== data.file.id)
  model.value.screenshots = screenshotList.value.map((s: any) => s.url)
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 700px" :title="`${microAppInfo?.id ? '编辑' : '创建'}微应用`">
    <NForm ref="formRef" :model="model" :rules="rules">
      <!-- 应用唯一标识：仅创建时显示 -->
      <NFormItem v-if="!microAppInfo?.id" path="microAppId" label="应用唯一标识">
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
        <NUpload
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
