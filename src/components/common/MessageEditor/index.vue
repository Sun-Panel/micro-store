<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { FormInst, FormRules, UploadFileInfo } from 'naive-ui'
import { NButton, NForm, NFormItem, NImage, NImageGroup, NInput, NPopconfirm, NSelect, NUpload, useMessage } from 'naive-ui'
import { useAuthStore } from '@/store'
import { getTemplateListByFlag as getTemplateListApi, send as sendApi } from '@/api/system/message'
import { t } from '@/locales'
import { router } from '@/router'

interface Props {
  flags: string[] // 列出的模板标识
  showSelect: boolean // 显示模板选择器
}

const props = defineProps<Props>()

interface TemplateInfo extends Message.TemplateInfo {
  mailToStr: string
}

const formInitValue: TemplateInfo = {
  name: '',
  title: '',
  mailToStr: '',
  content: '',
  note: '',
  args: [],
  toUsernames: '',
}

const authStore = useAuthStore()
const ms = useMessage()
const sending = ref(false)
const appendix = ref<string[]>([])

const templateInfoList = ref<Message.TemplateInfo[]>()
const activeTemplateInfo = ref<TemplateInfo>(formInitValue)
const currentSelectIndex = ref(0)

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
      message: '必填项',
    },
  ],
  toUsernames: {
    required: true,
    trigger: 'blur',
    type: 'string',
    message: '必填项',
  },
  content: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  ],
}

function handleSend(e: MouseEvent) {
  if (!activeTemplateInfo.value.args)
    activeTemplateInfo.value.args = []

  const argObj: { [key: string]: any } = {}
  for (let i = 0; i < activeTemplateInfo.value.args.length; i++) {
    const element = activeTemplateInfo.value.args[i]
    argObj[element.keyword as string] = element.value
  }
  const mailTo = splitEmailAddress(activeTemplateInfo.value.toUsernames as string)

  // 验证图片数量是否在限制内
  if (activeTemplateInfo.value.uploadAppendixShow) {
    if (activeTemplateInfo.value.uploadAppendixMin !== 0 && appendix.value.length < (activeTemplateInfo.value.uploadAppendixMin as number)) {
      ms.error(`请至少上传${activeTemplateInfo.value.uploadAppendixMin}张图片`)
      return
    }
  }

  // console.log('发送邮件的参数', activeTemplateInfo.value)

  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      const req: Message.SendReq = {
        title: activeTemplateInfo.value.title,
        content: activeTemplateInfo.value.content || '',
        templateArg: argObj,
        toUsernames: mailTo,
        appendixs: appendix.value,
      }
      send(req)
    }
    else {
      ms.error('请检查各项是否正确填写')
    }
  })
}

async function send(req: Message.SendReq) {
  if (sending.value === true) {
    ms.warning('正在发送请稍后')
    return
  }
  sending.value = true
  try {
    const res = await sendApi(req)
    console.log(res)
    sending.value = false
    ms.success('发送成功')
    if (activeTemplateInfo.value.saveAfterUrl)
      router.push(activeTemplateInfo.value.saveAfterUrl)
  }
  catch (error) {
    ms.error('发送失败')
    sending.value = false
  }
}

async function getTemplateList() {
  const { data } = await getTemplateListApi<Common.ListResponse<Message.TemplateInfo[]>>(props.flags)
  templateInfoList.value = data.list
  for (let i = 0; i < data.list.length; i++) {
    const element = data.list[i]
    templateOptions.value.push({
      label: element.name || '',
      value: i,
    })
  }

  if (templateInfoList.value.length >= 1) {
    currentSelectIndex.value = 0
    templateChange(currentSelectIndex.value)
  }
}

function templateChange(value: number) {
  if (templateInfoList.value)
    activeTemplateInfo.value = templateInfoList.value[value] as TemplateInfo

  else
    ms.error('此模板损坏')
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

const handleUploadFinish = ({
  file,
  event,
}: {
  file: UploadFileInfo
  event?: ProgressEvent
}) => {
  const res = JSON.parse((event?.target as XMLHttpRequest).response)
  if (res.code === 0) {
    const imageUrl = res.data.imageUrl
    appendix.value.push(imageUrl)
    // emit('update:itemIcon', itemIconInfo.value || null)
  }
  else {
    // apiRespErrMsg(res)
    ms.error(`${t('common.uploadFail')}:${res.msg}`)
  }

  return file
}

function handleDeleteAppendix(index: number) {
  appendix.value.splice(index, 1)
}

function handleUploadBefore(): boolean {
  if (activeTemplateInfo.value.uploadAppendixMax !== 0 && appendix.value.length >= (activeTemplateInfo.value.uploadAppendixMax as number)) {
    ms.error(`最多只能上传${activeTemplateInfo.value.uploadAppendixMax}张图片`)
    return false
  }

  return true
}

onMounted(() => {
  getTemplateList()
})
</script>

<template>
  <div>
    <NForm ref="formRef" :model="activeTemplateInfo" :rules="rules">
      <NFormItem v-if="showSelect" label="模板">
        <NSelect v-model:value="currentSelectIndex" style="max-width: 300px;" :options="templateOptions" placeholder="可以选择一个模板" @update:value="templateChange" />
      </NFormItem>

      <NFormItem v-if="activeTemplateInfo?.titleShow === 1" path="title" label="标题">
        <NInput
          v-model:value="activeTemplateInfo.title"
          type="textarea"
          :autosize="{
            minRows: 1,
            maxRows: 3,
          }"
          :disabled="activeTemplateInfo?.titleEdit === 0 "
          placeholder="标题(可选)"
        />
      </NFormItem>

      <NFormItem v-if="activeTemplateInfo?.toUsernameShow === 1" path="toUsernames" label="收信人">
        <NInput
          v-model:value="activeTemplateInfo.toUsernames"
          type="textarea"
          :autosize="{
            minRows: 1,
            maxRows: 3,
          }"
          :disabled="activeTemplateInfo?.toUsernameEdit === 0 "
          placeholder="(账号)邮箱地址，多个使用英文逗号隔开"
        />
      </NFormItem>

      <div
        v-for="(item, index) in activeTemplateInfo.args"
        :key="index"
        class="mt-2"
      >
        <NFormItem
          :path="`args[${index}].value`" :label="item.name" :rule="{
            required: item.required,
            message: `必填项`,
            trigger: ['input', 'blur'],
          }"
        >
          <NInput
            v-model:value="item.value"
            style="margin-left: 5px"
            clearable
            :placeholder="item.placeholder"
          />
        </NFormItem>
      </div>

      <NFormItem v-if="activeTemplateInfo?.contentShow === 1" path="content" label="内容">
        <NInput
          v-model:value="activeTemplateInfo.content"
          placeholder="请输入内容"
          :disabled="activeTemplateInfo?.contentEdit === 0 "
          :autosize="{
            minRows: 3,
          }"
          type="textarea"
        />
      </NFormItem>

      <NFormItem v-if="activeTemplateInfo.uploadAppendixShow === 1" label="附加图片">
        <div>
          <NUpload
            action="/api/file/uploadImg"
            accept="image/png, image/jpeg,image/gif,image/webp"
            name="imgfile"
            :headers="{
              token: authStore.token as string,
            }"
            :show-file-list="false"
            @finish="handleUploadFinish"
            @before-upload="handleUploadBefore"
          >
            <NButton size="small">
              上传图片
            </NButton>
          </NUpload>
          <div class="flex">
            <NImageGroup>
              <div
                v-for="item, index in appendix"
                :key="index"
                class="flex items-center"
              >
                <div class="flex flex-col items-center m-2 ">
                  <div class="w-[100px] h-[100px] flex items-center">
                    <NImage
                      width="100"
                      :src="item"
                    />
                  </div>
                  <div class="mt-2">
                    <NButton size="small" type="error" @click="handleDeleteAppendix(index)">
                      删除
                    </NButton>
                  </div>
                </div>
              </div>
            </NImageGroup>
          </div>
        </div>
      </NFormItem>
    </NForm>

    <NPopconfirm
      @positive-click="handleSend"
    >
      <template #trigger>
        <NButton type="success">
          {{ activeTemplateInfo.saveButtonName ? activeTemplateInfo.saveButtonName : "发送" }}
        </NButton>
      </template>
      确实无误要继续吗
    </NPopconfirm>
  </div>
</template>
