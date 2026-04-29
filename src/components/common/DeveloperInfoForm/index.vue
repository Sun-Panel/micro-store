<script lang="ts" setup>
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NDivider, NForm, NFormItem, NImage, NInput, NSelect, NSpace, NTooltip, NUpload, useMessage } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { SvgIcon } from '@/components/common'
import { useAuthStore } from '@/store/modules/auth'

interface Props {
  /** 是否为编辑模式 */
  editMode?: boolean
  /** 是否已认证为开发者 */
  isDeveloper?: boolean
  /** 初始数据 */
  initialData?: Partial<Developer.RegisterRequest>
  /** Name 上次修改时间 */
  nameUpdatedAt?: string
}

const props = withDefaults(defineProps<Props>(), {
  editMode: false,
  isDeveloper: false,
  initialData: () => ({}),
})

const emit = defineEmits<{
  (e: 'submit', data: Developer.RegisterRequest, isEdit: boolean): void
}>()

const message = useMessage()
const authStore = useAuthStore()
const formRef = ref<FormInst | null>(null)
const loading = ref(false)

// 待上传的二维码图片文件
const pendingUploadFile = ref<File | null>(null)
// 二维码图片的本地预览URL
const qrcodePreviewUrl = ref('')

const model = ref<Developer.RegisterRequest>({
  developerName: '',
  contactMail: '',
  paymentName: '',
  paymentQrcode: '',
  paymentMethod: '',
  name: '',
})

const rules: FormRules = {
  developerName: [
    { required: true, trigger: 'blur', message: '必填项' },
    {
      pattern: /^[a-z][a-z0-9-]{2,}$/,
      trigger: 'blur',
      message: '只能包含小写字母、数字和连字符，以小写字母开头，且至少3个字符',
    },
  ],
  name: { required: true, trigger: 'blur', message: '必填项' },
  contactMail: [{ type: 'email', trigger: 'blur', message: '请输入有效的邮箱地址' }],
}

const cardTitle = computed(() => props.editMode ? '开发者信息' : '注册为开发者')
const buttonText = computed(() => props.editMode ? '保存修改' : '立即注册')

// Name 编辑冷却期判断（180天）
const isNameCooldown = computed(() => {
  if (!props.nameUpdatedAt) return false
  const updatedTime = new Date(props.nameUpdatedAt).getTime()
  const daysDiff = (Date.now() - updatedTime) / (1000 * 60 * 60 * 24)
  return daysDiff < 180
})

// Name 禁用提示
const nameDisabledTip = computed(() => {
  if (!props.nameUpdatedAt) return ''
  const updatedTime = new Date(props.nameUpdatedAt)
  const nextUpdateTime = new Date(updatedTime.getTime() + 180 * 24 * 60 * 60 * 1000)
  const daysRemaining = Math.ceil((nextUpdateTime.getTime() - Date.now()) / (1000 * 60 * 60 * 24))
  return `距离下次可修改还需 ${daysRemaining} 天（${nextUpdateTime.toLocaleDateString()}）`
})

// 收款方式选项
const paymentMethodOptions = [
  { label: '微信', value: 'wechat_qr' },
  { label: '支付宝', value: 'alipay_qr' },
]

// 监听 initialData 变化，更新 model
watch(() => props.initialData, (newData) => {
  if (newData && Object.keys(newData).length > 0) {
    model.value = {
      developerName: newData.developerName || '',
      contactMail: newData.contactMail || '',
      paymentName: newData.paymentName || '',
      paymentQrcode: newData.paymentQrcode || '',
      paymentMethod: newData.paymentMethod || '',
      name: newData.name || '',
    }
    // 如果有已上传的二维码URL，设置为预览URL
    if (newData.paymentQrcode)
      qrcodePreviewUrl.value = newData.paymentQrcode
  }
}, { immediate: true, deep: true })

// ========== 图片上传处理 ==========

// 上传单个文件到服务器
async function uploadFile(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const formData = new FormData()
    formData.append('imgfile', file)

    const xhr = new XMLHttpRequest()
    xhr.open('POST', '/api/file/uploadImg')
    xhr.setRequestHeader('token', authStore.token || '')
    xhr.onload = () => {
      if (xhr.status === 200) {
        try {
          const res = JSON.parse(xhr.responseText)
          if (res.code === 0 && res.data && res.data.imageUrl)
            resolve(res.data.imageUrl)

          else
            reject(new Error(res.message || '上传失败'))
        }
        catch (error) {
          reject(error)
        }
      }
      else {
        reject(new Error('上传失败'))
      }
    }
    xhr.onerror = () => reject(new Error('上传失败'))
    xhr.send(formData)
  })
}

// 处理文件选择变化 - 只创建本地预览，不上传
function handleQrcodeFileChange({ fileList }: { fileList: any[] }) {
  if (fileList && fileList.length > 0) {
    const file = fileList[0].file || fileList[0].originFileObj
    if (file) {
      pendingUploadFile.value = file as File
      // 创建本地预览URL
      qrcodePreviewUrl.value = URL.createObjectURL(file as File)
    }
  }
}

// 删除图片
function handleQrcodeRemove() {
  model.value.paymentQrcode = ''
  pendingUploadFile.value = null
  // 释放本地预览URL
  if (qrcodePreviewUrl.value && qrcodePreviewUrl.value.startsWith('blob:'))
    URL.revokeObjectURL(qrcodePreviewUrl.value)

  qrcodePreviewUrl.value = ''
}

async function handleSubmit() {
  await formRef.value?.validate()

  loading.value = true
  try {
    // 如果有待上传的文件，先上传
    if (pendingUploadFile.value) {
      try {
        const imageUrl = await uploadFile(pendingUploadFile.value)
        model.value.paymentQrcode = imageUrl
        qrcodePreviewUrl.value = imageUrl
        // 上传成功后清理待上传文件
        pendingUploadFile.value = null
      }
      catch (error: any) {
        message.error(error?.message || '图片上传失败')
        return
      }
    }

    emit('submit', model.value, props.editMode)
  }
  catch (error: any) {
    message.error(error?.message || '操作失败')
  }
  finally {
    loading.value = false
  }
}

function setLoading(value: boolean) {
  loading.value = value
}

defineExpose({
  setLoading,
})
</script>

<template>
  <NCard
    :title="cardTitle"
    style="max-width: 600px;"
  >
    <template #header-extra>
      <NSpace align="center">
        <span
          v-if="isDeveloper"
          class="text-green-500 text-sm"
        >
          <SvgIcon
            icon="mdi:check-circle"
            class="mr-1"
          />
          已认证为开发者
        </span>
      </NSpace>
    </template>

    <NForm
      ref="formRef"
      :model="model"
      :rules="rules"
      label-placement="left"
      label-width="100"
    >
      <NDivider title-placement="left">
        开发者基础信息
      </NDivider>
      <NFormItem
        path="developerName"
        label="开发者标识"
      >
        <NInput
          v-model:value="model.developerName"
          :disabled="editMode"
          placeholder="支持小写字母、数字，至少3个字符且首个字符不能为数字"
        />
      </NFormItem>

      <NFormItem
        path="name"
        label="开发者名称"
      >
        <NTooltip
          v-if="isNameCooldown"
          trigger="hover"
        >
          <template #trigger>
            <div class="w-full">
              <NInput
                v-model:value="model.name"
                :disabled="true"
                placeholder="公开显示的开发者名称，每180天可以修改1次"
              />
            </div>
          </template>
          {{ nameDisabledTip }}
        </NTooltip>
        <NInput
          v-else
          v-model:value="model.name"
          placeholder="公开显示的开发者名称，每180天可以修改1次"
        />
      </NFormItem>

      <NFormItem
        path="contactMail"
        label="联系邮箱"
      >
        <NInput
          v-model:value="model.contactMail"
          placeholder="用于接收重要通知"
        />
      </NFormItem>

      <NDivider title-placement="left">
        收款信息
      </NDivider>

      <NFormItem
        path="paymentName"
        label="收款人姓名"
      >
        <NInput
          v-model:value="model.paymentName"
          placeholder="请输入收款人真实姓名"
        />
      </NFormItem>

      <NFormItem
        path="paymentMethod"
        label="收款方式"
      >
        <NSelect
          v-model:value="model.paymentMethod"
          :options="paymentMethodOptions"
          placeholder="请选择收款方式"
        />
      </NFormItem>

      <NFormItem
        path="paymentQrcode"
        label="收款二维码"
      >
        <div class="w-full">
          <!-- 图片预览 -->
          <div
            v-if="qrcodePreviewUrl"
            class="mb-3"
          >
            <NImage
              :src="qrcodePreviewUrl"
              width="150"
              height="150"
              object-fit="contain"
              :img-props="{ style: 'border: 1px solid #e5e7eb; border-radius: 4px; cursor: pointer;' }"
            />
          </div>
          <div class="flex items-center gap-2">
            <!-- 隐藏的输入框，用于保存URL -->
            <NInput
              v-show="false"
              v-model:value="model.paymentQrcode"
            />
            <!-- 上传按钮 -->
            <NUpload
              :show-file-list="false"
              accept="image/*"
              :custom-request="() => { }"
              @change="handleQrcodeFileChange"
            >
              <!-- 上传收款码 -->
              <div class="flex gap-2">
                <NButton>
                  {{ qrcodePreviewUrl ? '更换二维码' : '上传二维码' }}
                </NButton>
                <!-- 删除按钮 -->
                <NButton
                  v-if="qrcodePreviewUrl"
                  type="error"
                  @click.stop="handleQrcodeRemove"
                >
                  删除
                </NButton>
              </div>
            </NUpload>
          </div>
        </div>
      </NFormItem>
    </NForm>

    <template #footer>
      <NSpace justify="end">
        <NButton
          type="primary"
          :loading="loading"
          @click="handleSubmit"
        >
          {{ buttonText }}
        </NButton>
      </NSpace>
    </template>
  </NCard>
</template>
