<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst } from 'naive-ui'
import { NButton, NDivider, NForm, NFormItem, NH4, NInput, NModal, useDialog } from 'naive-ui'
import { set as setApi } from '@/api/admin/systemVariable'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'
import { t } from '@/locales'

interface Props {
  visible: boolean
  info: SystemVariable.SystemVariableEditReq | null
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const dialog = useDialog()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}

const formInitValue = {
  description: '',
  name: '',
  value: '',
}

const model = ref<SystemVariable.SystemVariableEditReq>(formInitValue)
const formRef = ref<FormInst | null>(null)

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

watch(show, (newValue, oldValue) => {
  if (props.info?.id)
    model.value = { ...props.info } || {}

  else
    model.value = formInitValue
})

const save = async () => {
  await setApi(model.value.name, model.value.value).then((res) => {
    emit('done')
    show.value = false
  }).catch((error) => {
    apiRespErrMsgAndCustomCodeNeg1Msg(error, `${t('common.failed')}:${error.msg}`)
  })
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      dialog.warning({
        title: '警告',
        content: '修改请慎重，如果是JSON等格式内容请自行校验一下，你确定要保存吗',
        positiveText: '确定',
        negativeText: '取消',
        onPositiveClick: () => {
          save()
        },
      })
    }
    else {
      console.log(errors)
    }
  })
}

function goToLink(url: string) {
  window.open(url)
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px;" :mask-closable="false" title="Modify the Variable Value">
    <NForm ref="formRef" :model="model">
      <!-- https://b3log.org/vditor/demo/vue.html -->
      <!-- https://jsoneditoronline.org/ -->

      <div>
        <NH4 prefix="bar">
          Name: {{ info?.name }}
        </NH4>
        Description: {{ info?.description }}
      </div>
      <NDivider />
      <NFormItem path="value" label="Value">
        <NInput v-model:value="model.value" placeholder="value" type="textarea" :min="5" />
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
        Save
      </NButton>
    </template>
  </NModal>
</template>
