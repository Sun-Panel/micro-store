<script setup lang="ts">
import type { FormInst } from 'naive-ui'
import { NButton, NCard, NForm, NFormItem, NH2, NInput, NInputGroup, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { commit as registerCommitApi, sendRegisterVcode } from '@/api/register'
import { router } from '@/router'
import { Captcha, SvgIcon, Verification } from '@/components/common'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'
import { t } from '@/locales'

const verificationShow = ref<boolean>(false)
const verificationLoading = ref<boolean>(false)
const verificationId = ref<string>('')
const ms = useMessage()
const countdown = ref(0)
const formRef = ref<FormInst | null>(null)
const captchaRef = ref()

const route = useRoute()

interface Form extends System.Register.SendRegisterVcodeRquest {
  passwordAgain: string
}

const form = ref<Form>({
  username: '',
  password: '',
  passwordAgain: '',
  referralCode: '',
})

const rules = {
  username: [{
    required: true,
    trigger: 'blur',
    message: t('form.required'),
  },
  {
    message: '请填入邮箱格式作为账号',
    type: 'email',
    trigger: 'blur',
  },
  ],
  password: {
    required: true,
    trigger: 'blur',
    message: t('form.required'),

  },
  passwordAgain: {
    required: true,
    trigger: 'blur',
    message: t('form.required'),
  },

  vcode: {
    required: true,
    trigger: 'blur',
    message: t('form.required'),
  },
}

async function registerCommit() {
  await registerCommitApi<System.Register.SendRegisterVcodeRquest>(form.value).then(() => {
    ms.success('注册完成，准备跳转登录页面')
    router.push({ path: '/login' })
  }).catch((res) => {
    ms.error(res.msg)
  })
}

async function startCountdown(getCerificationId?: string, vCode?: string, obj?: { refresh: () => void }) {
  if (!form.value.vcode || form.value.vcode === '') {
    ms.error('请先输入图形验证码')
    return
  }
  if (form.value.password !== form.value.passwordAgain) {
    ms.error('两次密码输入不一致')
    return
  }
  verificationLoading.value = true
  const verification: Common.VerificationRequest = {
    codeId: getCerificationId,
    vCode,
  }
  const data: System.Register.SendRegisterVcodeRquest = { ...form.value }
  data.verification = verification
  await sendRegisterVcode<System.Register.SendRegisterVcodeRquest>(data).then(() => {
    ms.success(t('common.emailVerificationCodeSent'))
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value -= 1
      if (countdown.value === 0)
        clearInterval(timer)
    }, 1000)
    verificationLoading.value = false
  }).catch((res) => {
    captchaRef.value.refresh()
    verificationLoading.value = false
    if (res.code === -1) {
      ms.error(res.msg)
      return
    }
    else if (res.code === 1101) {
      // verificationShow.value = true
      verificationId.value = res.data.verification?.codeId || ''
      ms.warning(res.msg)
      return
    }
    else if (res.code === 1102) {
    // 验证码错误
      ms.warning(res.msg)
      return
    }

    apiRespErrMsg(res)
  })
}

function handleSubmit() {
  if (form.value.password !== form.value.passwordAgain) {
    ms.error('两次密码输入不一致')
    return
  }
  // 点击登录按钮触发
  registerCommit()
}

function handleSendEmail(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      startCountdown()
    }
    else {
      console.log(errors)
      ms.error('请检查表单各项是否存在什么问题')
    }
  })
}

onMounted(() => {
  form.value.referralCode = route.query?.referralCode as string ?? ''
})
</script>

<template>
  <div class="login-container">
    <NCard class="login-card shadow-2xl">
      <div class="login-title">
        <NH2>注册</NH2>
      </div>
      <NForm ref="formRef" :model="form" label-width="100px" :rules="rules">
        <NFormItem label="邮箱" path="username">
          <NInput v-model:value="form.username" placeholder="请输入邮箱作为登录账号">
            <template #prefix>
              <SvgIcon icon="ph:user-bold" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem label="密码" path="password">
          <NInput v-model:value="form.password" type="password" placeholder="请输入密码">
            <template #prefix>
              <SvgIcon icon="mdi:password-outline" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem label="确认密码" path="passwordAgain">
          <NInput v-model:value="form.passwordAgain" type="password" placeholder="请再次输入密码">
            <template #prefix>
              <SvgIcon icon="mdi:password-outline" />
            </template>
          </NInput>
        </NFormItem>

        <!-- <NFormItem label="邀请码">
          <NInput v-model:value="form.referralCode" type="text" placeholder="邀请码(可选)">
            <template #prefix>
              <SvgIcon icon="solar:password-minimalistic-broken" />
            </template>
          </NInput>
        </NFormItem> -->

        <NFormItem label="图形验证码" path="vcode">
          <div class="w-[130px] h-[40px] mr-[20px] rounded border flex cursor-pointer">
            <Captcha ref="captchaRef" src="/api/captcha/getImage" />
          </div>
          <NInput v-model:value="form.vcode" size="large" type="text" placeholder="请输入图像验证码" clearable />
        </NFormItem>

        <NFormItem label="邮箱验证码">
          <NInputGroup>
            <NInput v-model:value="form.emailVCode" placeholder="请输入邮箱的验证码" />
            <NButton type="primary" :loading="verificationLoading" :disabled="countdown > 0" ghost @click="handleSendEmail">
              <span v-if="countdown > 0">{{ `${countdown}秒后重新获取` }} </span>
              <span v-else>获取邮箱验证码</span>
            </NButton>
          </NInputGroup>
        </NFormItem>

        <NButton type="primary" block @click="handleSubmit">
          注册
        </NButton>
        <NButton type="primary" block quaternary @click="$router.push({ path: '/login' })">
          已有账号，去登录
        </NButton>
      </NForm>
    </NCard>

    <Verification
      v-model:visible="verificationShow"
      v-model:loading="verificationLoading"
      :verification-id="verificationId"
      @on-submit="startCountdown"
    />
  </div>
</template>

  <style>
    .login-container {
        padding: 20px;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
    }

    /* 夜间模式 */
    .dark .login-container{
      background-color: rgb(43, 43, 43);
    }

    @media (min-width: 600px) {
        .login-card {
            width: auto;
            margin: 0px 10px;
        }
        .login-button {
            width: 100%;
        }
    }

    .login-card {
        margin: 20px;
        min-width:400px;
    }

  .login-title{
    text-align: center;
    margin: 20px;
  }
  </style>
