<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCard, NCheckbox, NForm, NFormItem, NInput, NPopconfirm, NSelect, useMessage } from 'naive-ui'
import { getList } from '@/api/admin/emailTemplate'
import { send as sendEmail } from '@/api/admin/email'
import { t } from '@/locales'

interface EmailInfo extends EmailTemplate.Info {
  mailToStr: string
  mailTo: string[]
  replaceArgs: boolean
}

const message = useMessage()
const sending = ref(false)
let templateInfos: { [key: number]: EmailTemplate.Info }

const formInitValue: EmailInfo = {
  name: '',
  title: '',
  mailToStr: '',
  mailTo: [''],
  replaceArgs: false,
  content: '',
  note: '',
  isSeparateSend: true,
  flag: '',
}

const templateArgs = ref<{
  name: string
  keyword: string
  value?: string
  isDefault?: boolean
}[]>([])
const model = ref<EmailInfo>(formInitValue)
const formRef = ref<FormInst | null>(null)
const templateOptions = ref<{
  label: string
  value: number | undefined
}[]>([])
const rules: FormRules = {
  name: [
    {
      required: true,
      trigger: 'blur',
      message: t('form.required'),
    },
  ],
  title: {
    required: true,
    trigger: 'blur',
    type: 'string',
    message: t('form.required'),
  },
  mailToStr: {
    required: true,
    trigger: 'blur',
    type: 'string',
    message: t('form.required'),
  },
  content: [
    {
      required: true,
      trigger: 'blur',
      message: t('form.required'),
    },
  // {
  //   trigger: 'blur',
  //   message: '请输入邮箱格式',
  //   type: 'email',
  // },
  ],
  password: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项且大于5个字符',
      min: 5,
    },
  ],
}

function handleSend(e: MouseEvent) {
  for (let i = 0; i < templateArgs.value.length; i++) {
    const element = templateArgs.value[i]
    if (!model.value.templateArg)
      model.value.templateArg = {}
    model.value.templateArg[element.keyword] = element.value || ''
  }
  model.value.replaceArgs = true
  model.value.mailTo = splitEmailAddress(model.value.mailToStr)
  // console.log('发送邮件的参数', model.value)
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      const req: Email.SendEmailReq = {
        body: model.value.content,
        replaceArg: true,
        title: model.value.title,
        mailTo: model.value.mailTo,
        templateArg: model.value.templateArg,
        isSeparateSend: model.value.isSeparateSend,
      }

      send(req)
    }
  })
}

async function send(req: Email.SendEmailReq) {
  if (sending.value === true) {
    message.warning('正在发送请稍后')
    return
  }
  sending.value = true
  try {
    await sendEmail(req)
    sending.value = false
    message.success('发送成功')
  }
  catch (error) {
    message.error('发送失败')
    sending.value = false
  }
}

async function getTemplateList() {
  const { data } = await getList<Common.ListResponse<EmailTemplate.Info[]>>()
  templateInfos = {}
  for (let i = 0; i < data.list.length; i++) {
    const element = data.list[i]
    templateInfos[element.id || 0] = (element)
    templateOptions.value.push({
      label: element.name,
      value: element.id,
    })
  }
}

function removeArgItem(index: number) {
  templateArgs.value.splice(index, 1)
}

function addArgItem(index: number) {
  templateArgs.value.push({
    keyword: '',
    name: '',
    value: '',
  })
}

function templateChange(value: number) {
  const info = templateInfos[value]
  templateArgs.value = []
  if (info.args) {
    for (let i = 0; i < info.args?.length; i++) {
      const element = info.args[i]
      templateArgs.value.push({
        name: element.name,
        keyword: element.keyword,
        value: '',
        isDefault: true,
      })
    }
    model.value.content = info.content
    model.value.title = info.title
  }
}

function splitEmailAddress(input: string): string[] {
  // 使用正则表达式匹配逗号或空格进行分割
  const regex = /[,\s]+/
  const splittedStrings = input.split(regex)

  // 过滤出符合邮箱格式的项
  const validEmails = splittedStrings.filter((item) => {
    // 使用简单的邮箱正则表达式进行验证
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(item)
  })

  return validEmails
}

onMounted(() => {
  getTemplateList()
})
</script>

<template>
  <div>
    <NCard class="min-w-[300px]">
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem label="邮件模板">
          <NSelect v-model:value="model.id" clearable style="max-width: 300px;" :options="templateOptions" placeholder="可以选择一个模板" @update:value="templateChange" />
        </NFormItem>

        <NCard size="small">
          模板参数
          <div
            v-for="(item, index) in templateArgs"
            :key="index"
            class="mt-2"
          >
            <div class="flex items-center">
              <NInput v-if="item.isDefault" :value="`(预设) ${item.name}`" size="small" style="margin-right: 5px" readonly disabled title="预设参数说明" />
              <NInput v-model:value="item.keyword" size="small" clearable placeholder="变量名-示例：{name}" />
              <NInput v-model:value="item.value" size="small" style="margin-left: 5px" clearable placeholder="变量值-示例：张三" />
              <NButton style="margin-left: 10px" type="error" size="small" @click="removeArgItem(index)">
                删除
              </NButton>
            </div>
          </div>
          <NButton attr-type="button" size="small" type="info" style="margin-top: 12px" @click="addArgItem">
            增加自定义参数
          </NButton>
        </NCard>

        <NFormItem path="mailToStr" label="收件人" style="margin-top: 20px;">
          <NInput
            v-model:value="model.mailToStr"
            type="textarea"
            :autosize="{
              minRows: 1,
              maxRows: 3,
            }"
            placeholder="邮箱地址，多个使用英文逗号隔开"
          />
        </NFormItem>

        <NFormItem path="title" label="标题">
          <NInput
            v-model:value="model.title"
            type="textarea"
            :autosize="{
              minRows: 1,
              maxRows: 3,
            }"
            placeholder="标题"
          />
        </NFormItem>

        <NFormItem path="content" label="内容">
          <NInput
            v-model:value="model.content"
            placeholder="请输入内容"
            :autosize="{
              minRows: 3,
            }"
            type="textarea"
          />
        </NFormItem>
        <NFormItem label="群发方式">
          <NCheckbox v-model:checked="model.isSeparateSend">
            单独发送
          </NCheckbox>
        </NFormItem>
      </NForm>
      <NPopconfirm
        @positive-click="handleSend"
      >
        <template #trigger>
          <NButton type="success">
            发送邮件
          </NButton>
        </template>
        请确实无误继续发送吗
      </NPopconfirm>
    </NCard>
  </div>
</template>
