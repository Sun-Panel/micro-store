<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { FormInst } from 'naive-ui'
import { NButton, NCard, NDivider, NForm, NFormItem, NInput, NSwitch, useMessage } from 'naive-ui'
import { AdminSystemSettingGetWebsiteSetting, AdminSystemSettingSetWebsiteSetting } from '@/api/admin'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'
import { t } from '@/locales'

const message = useMessage()

const formInitValue = {
  loginCaptcha: false,
  openRegister: false,
  webSiteUrl: '',
  emailSuffix: '',
}

const model = ref<AdminSystemSetting.Website>(formInitValue)
const formRef = ref<FormInst | null>(null)

const handleSave = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      const post = async () => {
        await AdminSystemSettingSetWebsiteSetting(model.value).then(() => {
          message.success(t('common.saveSuccess'))
        }).catch((error) => {
          apiRespErrMsgAndCustomCodeNeg1Msg(error, t('common.saveFail'))
        })
      }
      post()
    }
  })
}

function getCurrentDomain() {
  return `${location.protocol}//${location.hostname}${location.port === '' ? '' : (`:${location.port}`)}`
}

onMounted(() => {
  const post = async () => {
    const { code, data } = await AdminSystemSettingGetWebsiteSetting<AdminSystemSetting.Website>()
    if (code === 0)
      model.value = data
  }
  post()
})
</script>

<template>
  <div>
    <NCard class="max-w-[500px]">
      <NDivider title-placement="left">
        登录设置
      </NDivider>
      <NForm ref="formRef" :model="model" label-placement="left">
        <NFormItem label="登录验证码" style="margin-top: 20px;">
          <NSwitch v-model:value="model.loginCaptcha" />
        </NFormItem>
        <NDivider title-placement="left">
          注册设置
        </NDivider>
        <NFormItem label="新用户注册" style="margin-top: 20px;">
          <NSwitch v-model:value="model.openRegister" />
        </NFormItem>

        <!-- <NFormItem label="邮箱后缀" style="margin-top: 20px;">
          <NInput v-model:value="model.emailSuffix" type="text" placeholder="一般用于企业邮箱，例:@abc.com" />
        </NFormItem> -->
        <NDivider title-placement="left">
          其他
        </NDivider>
        <NFormItem label="网站地址" style="margin-top: 20px;">
          <div>
            <NInput v-model:value="model.webSiteUrl" type="text" placeholder="用于外站回跳本站" />
            <NButton quaternary type="info" @click="model.webSiteUrl = getCurrentDomain()">
              获取当前地址
            </NButton>
          </div>
        </NFormItem>
      </NForm>
      <NButton type="success" @click="handleSave">
        {{ t('common.save') }}
      </NButton>
    </NCard>
  </div>
</template>
