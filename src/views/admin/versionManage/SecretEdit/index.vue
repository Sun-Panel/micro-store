<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import {
  NButton,
  NCheckbox,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NTag,
  useMessage,
} from 'naive-ui'
import { defineEmits, defineProps, ref, watch } from 'vue'
import { editSecret, getSecretByVersion } from '@/api/admin/version'

interface Props {
  visible: boolean
  version: string
  isCreate: boolean // true 表示创建，false 表示查看/编辑
}

interface Emit {
  (e: 'done'): void
  (e: 'update:visible', value: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

const defaultSecretInfo: Version.SecretInfo = {
  version: '',
  secretKey: '',
  status: true,
}

const model = ref<Version.SecretInfo>({ ...defaultSecretInfo })
const formRef = ref<FormInst | null>(null)
const show = ref(props.visible)
const message = useMessage()
const isLoading = ref(false)

watch(() => props.visible, async (newVal) => {
  show.value = newVal
  if (newVal) {
    if (props.version) {
      // 先重置为默认值
      model.value = { ...defaultSecretInfo }
      model.value.version = props.version

      // 如果不是创建模式，则加载现有密钥信息
      if (!props.isCreate) {
        isLoading.value = true
        try {
          const response = await getSecretByVersion<Version.SecretInfo>(props.version)
          model.value = { ...response.data }
        }
        catch (error) {
          message.error('加载密钥信息失败')
        }
        finally {
          isLoading.value = false
        }
      }
      else {
        // 创建模式，生成随机密钥
        model.value.version = props.version
        model.value.secretKey = generateSecretKey()
        model.value.status = true
      }
    }
  }
})

watch(show, (newVal) => {
  // 只有当 show 的值与 props.visible 不同时才 emit，避免循环触发
  if (newVal !== props.visible) {
    emit('update:visible', newVal)
  }
})

const rules: FormRules = {
  version: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  ],
  secretKey: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  ],
}

function generateSecretKey(): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < 32; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

function handleRegenerateKey() {
  model.value.secretKey = generateSecretKey()
}

function handleSave() {
  formRef.value?.validate((errors) => {
    if (!errors) {
      editSecret(model.value).then(() => {
        message.success('保存成功')
        show.value = false
        emit('done')
      }).catch(() => {
        message.error('保存失败')
      })
    }
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px" :title="`${isCreate ? '创建' : '编辑'}版本密钥`">
    <NForm v-if="!isLoading" ref="formRef" :model="model" :rules="rules">
      <NFormItem path="version" label="版本号">
        <NInput v-model:value="model.version" type="text" disabled placeholder="1.3.0" />
      </NFormItem>

      <NFormItem path="secretKey" label="密钥">
        <div class="flex w-full gap-2">
          <NInput v-model:value="model.secretKey" type="text" placeholder="请输入密钥" class="flex-1" />
          <NButton v-if="isCreate" type="info" @click="handleRegenerateKey">
            重新生成
          </NButton>
        </div>
      </NFormItem>

      <NFormItem label="状态">
        <NCheckbox v-model:checked="model.status">
          启用
        </NCheckbox>
      </NFormItem>

      <div v-if="!isCreate" class="mb-4">
        <div class="text-sm text-gray-500 mb-2">
          密钥提示：
        </div>
        <NTag type="info" size="small">
          请妥善保管密钥，不要泄露给他人
        </NTag>
      </div>
    </NForm>
    <div v-else class="flex justify-center py-8">
      加载中...
    </div>

    <template #footer>
      <div class="flex justify-end gap-2">
        <NButton @click="show = false">
          取消
        </NButton>
        <NButton type="primary" :disabled="isLoading" @click="handleSave">
          保存
        </NButton>
      </div>
    </template>
  </NModal>
</template>
