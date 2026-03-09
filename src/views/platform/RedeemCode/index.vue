<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NInput } from 'naive-ui'
import { redeemCodeWriteOff } from '@/api/redeemCode'
import { Captcha } from '@/components/common'
import { apiRespErrMsg, message } from '@/utils/cmn/apiMessage'
import { t } from '@/locales'

const emits = defineEmits<{
  (e: 'success'): void
}>()
const code = ref('')
const vcode = ref('')

const captchaRef = ref()

const handleCheck = async () => {
  if (code.value === '')
    return

  await redeemCodeWriteOff(code.value, vcode.value).then((res) => {
    emits('success')
    message.success(t('common.success'))
  }).catch((res) => {
    captchaRef.value.refresh()
    switch (res.code) {
      case -2:
        // 已过期
        message.error(t('common.redeemCodeExpired'))
        break
      case -3:
        // 已被使用
        message.error(t('common.redeemCodeUsed'))
        break
      case -4:
        console.log('error_code', res.code)
        message.error(t('common.redeemCodeNotExist'))
        break
      default:
        apiRespErrMsg(res)
        break
    }
  })
}
</script>

<template>
  <div>
    <NInput v-model:value="code" size="large" type="text" />

    <div class="mt-2 text-orange-500">
      {{ $t('proAuthorize.redeemDescription') }}
    </div>

    <div v-if="code !== ''" class="flex mt-5">
      <div class="w-[130px] h-[40px] mr-[20px] rounded flex cursor-pointer">
        <Captcha ref="captchaRef" src="/api/captcha/getImage" />
      </div>
      <NInput
        v-model:value="vcode" type="text" clearable size="large"
        :placeholder="$t('login.captchaCodePlaceholder')"
      />
    </div>

    <div class="mt-2 flex justify-end">
      <NButton type="primary" :disabled="code === '' || vcode === ''" @click="handleCheck">
        {{ $t('common.redeem') }}
      </NButton>
    </div>
  </div>
</template>
