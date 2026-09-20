<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NBreadcrumb, NBreadcrumbItem, NButton, NCard, NForm, NFormItem, NInput, NSelect, NSwitch, NTag, useMessage } from 'naive-ui'
import { useRoute } from 'vue-router'
import { edit, getInfo } from '@/api/admin/htmlPageManage'
import { TinymceEditor } from '@/components/common'
import { apiRespErrMsg, getCurrentBaseUrlRoot } from '@/utils/cmn'
import { router } from '@/router'

const route = useRoute()
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

const pageName = ref(route.query.pageName as string)
const info = ref<HtmlPage.ListItem>()

const rules: FormRules = {
  content:
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  pageName: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
    {
      // 与后端校验保持一致：页面名会作为预览页路径，不能含空白字符或斜杠
      pattern: /^[^\s/\\]+$/,
      trigger: 'blur',
      message: '不能包含空格或斜杠',
    },
  ],
  pageDescription:
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
}

const formInitValue: HtmlPage.ListItem = {
  pageDescription: '',
  pageName: '',
  isLogin: false,
  content: '',
  messageTemplatePosition: 'bottom',
  messageTemplateFlag: '',
}

const model = ref<HtmlPage.ListItem>({ ...formInitValue })
const formRef = ref<FormInst | null>(null)

const save = async () => {
  if (model.value.messageTemplateFlag && !model.value.messageTemplatePosition)
    model.value.messageTemplatePosition = 'bottom'

  try {
    const res = await edit(model.value as HtmlPage.EditReq)
    if (res.code === 0) {
      if (!pageName.value)
        router.back()

      message.success('保存成功')
    }
    else {
      message.warning(`操作失败: ${res.msg}`)
    }
  }
  catch {
    message.warning('保存失败，请稍后重试')
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

function getInfoByPageName(pageName: string) {
  getInfo<HtmlPage.EditReq>(pageName).then(({ data }) => {
    model.value = { ...data }
    info.value = { ...data }
  }).catch((res) => {
    apiRespErrMsg(res)
  })
}

onMounted(() => {
  if (pageName.value)
    getInfoByPageName(pageName.value)
})
</script>

<template>
  <div>
    <NCard size="small" class="mb-5">
      <NBreadcrumb>
        <NBreadcrumbItem @click="router.push({ name: 'AdminHtmlPageManage' })">
          HTML页面管理
        </NBreadcrumbItem>
        <NBreadcrumbItem>
          <template v-if="pageName">
            编辑
          </template>
          <template v-else>
            创建
          </template>
        </NBreadcrumbItem>
      </NBreadcrumb>
    </NCard>

    <NCard size="small" class="mb-5">
      页面地址
      <NTag type="info">
        {{ getCurrentBaseUrlRoot() }}/hPage/{{ model.pageName }}
      </NTag>
    </NCard>

    <NCard size="small" class="mb-5">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem path="pageName" label="英文唯一标识">
          <NInput v-model:value="model.pageName" type="text" placeholder="仅支持英文" :disabled="!!info?.pageName" />
        </NFormItem>

        <NFormItem path="pageDescription" label="描述">
          <NInput v-model:value="model.pageDescription" type="text" placeholder="描述信息" />
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

        <NFormItem path="content" label="页面内容">
          <TinymceEditor v-model:content="model.content" />
        </NFormItem>
      </NForm>

      <NButton type="success" @click="handleValidateButtonClick">
        保存
      </NButton>
    </NCard>
  </div>
</template>
