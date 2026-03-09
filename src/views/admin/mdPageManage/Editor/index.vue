<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NDivider, NForm, NFormItem, NInput, NModal, NSelect, NSwitch, useMessage } from 'naive-ui'
import { edit } from '@/api/admin/mdPageManage'

interface Props {
  visible: boolean
  info: MdPage.ListItem | null
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const message = useMessage()
const messageTemplatePositionOptions = [
  {
    label: '底部',
    value: 'bottom',
  },
  {
    label: '顶部',
    value: 'top',
  },
]

const rules: FormRules = {
  content:
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  mdPageName:
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  mdPageDescription:
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
}

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}

const formInitValue = ref< MdPage.ListItem> ({
  mdPageDescription: '',
  mdPageName: '',
  isLogin: false,
  content: '',
  messageTemplatePosition: 'bottom',
  messageTemplateFlag: '',
})

const model = ref({ ...formInitValue.value })
const formRef = ref<FormInst | null>(null)

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

watch(show, (newValue, oldValue) => {
  if (props.info?.mdPageName)
    model.value = { ...props.info } || {}

  else
    model.value = { ...formInitValue.value }
})

const save = async () => {
  if (model.value.messageTemplateFlag && !model.value.messageTemplatePosition)
    model.value.messageTemplatePosition = 'bottom'

  const res = await edit(model.value as MdPage.EditReq)
  if (res.code === 0) {
    emit('done')
    show.value = false
  }
  else {
    message.warning(`操作失败: ${res.msg}`)
  }
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      save()

    else
      console.log(errors)
  })
}

function goToLink(url: string) {
  window.open(url)
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px;" :mask-closable="false" :title="info?.mdPageName ? '编辑' : '添加'">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="mdPageName" label="页面标识">
        <NInput v-model:value="model.mdPageName" type="text" placeholder="仅支持英文" :disabled="info?.mdPageName" />
      </NFormItem>

      <NFormItem path="mdPageDescription" label="页面描述">
        <NInput v-model:value="model.mdPageDescription" type="text" placeholder="描述信息" />
      </NFormItem>

      <NFormItem label="是否需要登录">
        <NSwitch v-model:value="model.isLogin" />
      </NFormItem>

      <NFormItem label="站内信模板">
        <NInput v-model:value="model.messageTemplateFlag" type="text" placeholder="模板标识" />
      </NFormItem>
      <NFormItem v-if="model.messageTemplateFlag" label="站内信模板显示位置">
        <NSelect v-model:value="model.messageTemplatePosition" :options="messageTemplatePositionOptions" />
      </NFormItem>

      <NDivider />
      <NFormItem path="value" label="页面内容">
        <NInput v-model:value="model.content" placeholder="value" type="textarea" :min="5" />
      </NFormItem>
    </NForm>

    <div>
      <NButton type="info" size="tiny" @click="goToLink('https://b3log.org/vditor/demo/vue.html')">
        Markdown Editor
      </NButton>

      <NButton type="info" size="tiny" style="margin-left: 5px;" @click="goToLink('https://jsoneditoronline.org/')">
        JSON Editor
      </NButton>
    </div>

    <template #footer>
      <NButton type="success" @click="handleValidateButtonClick">
        保存
      </NButton>
    </template>
  </NModal>
</template>
