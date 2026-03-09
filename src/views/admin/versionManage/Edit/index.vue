<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import {
  NButton,
  NCheckbox,
  NDatePicker,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NSelect,
} from 'naive-ui'
import { defineEmits, defineProps, ref, watchEffect } from 'vue'
import { edit } from '@/api/admin/version'

interface Props {
  visible: boolean
  version: Version.Info
}

interface Emit {
  (e: 'done', id: number): void // 创建完成
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()

const defaultEditInfo: Version.Info = {
  version: '',
  type: 'beta',
  releaseTime: '',
  description: '',
  downloadURL: '',
  pageUrl: '',
  isRolledBack: false,
  aloneSecretKey: '',
}

const model = ref<Version.Info>({ ...defaultEditInfo })
const formRef = ref<FormInst | null>(null)
const show = ref(props.visible)

// const show = computed({
//   get: () => props.visible,
//   set: (visible: boolean) => {
//     emit('update:visible', visible)
//   },
// })

watchEffect(() => {
  show.value = props.visible
  if (props.visible) {
    if (props.version?.id)
      model.value = { ...props.version }
    else
      model.value = { ...defaultEditInfo }
  }
})

interface Options {
  value?: string
  label?: string
  description?: string
}

const options = ref<Options[]>([
  {
    label: '正式版本',
    value: 'release',
  },
  {
    label: 'Beta',
    value: 'beta',
  },
  {
    label: 'Alpha',
    value: 'alpha',
  },
  {
    label: 'RC',
    value: 'rc',
  },
  {
    label: '开发版本',
    value: 'dev',
  },
])

const rules: FormRules = {
  version: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  ],
  type: {
    required: true,
    trigger: 'blur',
    type: 'string',
    message: '必选项',
  },
  releaseTime: {
    required: true,
    trigger: 'blur',
    type: 'string',
    message: '必选项',
  },
}

function handleSave() {
  edit(model.value).then(() => {
    show.value = false
    emit('done', 0)
  })
}

function handleSelectReleaseTime(value: string | null, timestampValue: number | null) {
  if (value)
    model.value.releaseTime = value
}
</script>

<template>
  <NModal v-model:show="show " preset="card" style="width: 600px" :title="`${props.version?.id ? '编辑' : '添加'}`">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="type" label="版本类型">
        <NSelect
          v-model:value="model.type"
          :options="options"
        />
      </NFormItem>

      <NFormItem path="version" label="版本名称">
        <NInput v-model:value="model.version" type="text" placeholder="1.3.0" />
      </NFormItem>

      <NFormItem path="releaseTime" label="发布日期">
        <NDatePicker
          type="datetime" clearable
          :default-formatted-value="model.releaseTime || null"
          @update:formatted-value="handleSelectReleaseTime"
        />
      </NFormItem>

      <NFormItem path="description" label="描述">
        <NInput v-model:value="model.description" type="text" />
      </NFormItem>

      <NFormItem path="pageUrl" label="说明页地址">
        <NInput v-model:value="model.pageUrl" type="text" />
      </NFormItem>

      <NFormItem path="downloadURL" label="下载地址">
        <NInput v-model:value="model.downloadURL" type="text" />
      </NFormItem>

      <NFormItem label="是否撤包">
        <NCheckbox
          v-model:checked="model.isRolledBack"
        >
          是
        </NCheckbox>
      </NFormItem>
    </NForm>

    <template #footer>
      <div class="flex justify-end">
        <NButton type="success" @click="handleSave">
          保存
        </NButton>
      </div>
    </template>
  </NModal>
</template>
