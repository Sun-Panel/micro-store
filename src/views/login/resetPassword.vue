<script setup lang="ts">
import { NButton, NCard, NForm, NFormItem, NH2, NInput, NInputGroup, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { resetPasswordByVCode, sendResetPasswordVCode } from '@/api/login'
import { Captcha, SvgIcon, Verification } from '@/components/common'
import { router } from '@/router'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'

const verificationShow = ref<boolean>(false)
const verificationLoading = ref<boolean>(false)
const verificationId = ref<string>('')
const ms = useMessage()
const countdown = ref(0)
const captchaRef = ref()

const route = useRoute()

interface Form extends System.Register.SendRegisterVcodeRquest {
  passwordAgain: string
}

const form = ref<Form>({
  email: '',
  passwordAgain: '',
})

async function commit() {
  form.value.email = form.value.username
  await resetPasswordByVCode<System.Register.SendRegisterVcodeRquest>(form.value)
    .then(() => {
      ms.success('密码重置成功，请重新登录')
      router.push({ path: '/login' })
    }).catch((res) => {
      ms.error(res.msg)
    })
}

function handleSubmit() {
  if (form.value.password !== form.value.passwordAgain) {
    ms.error('两次密码输入不一致')
    return
  }
  commit()
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
  await sendResetPasswordVCode < System.Register.SendRegisterVcodeRquest > (
    form.value.username as string,
    form.value.vcode as string,
  ).then(() => {
    ms.success('验证码发送成功')
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value -= 1
      if (countdown.value === 0)
        clearInterval(timer)
    }, 1000)
    verificationLoading.value = false
  }).catch((res) => {
    verificationLoading.value = false
    captchaRef.value.refresh()
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

onMounted(() => {
  form.value.username = route.query?.u as string ?? ''
})
</script>

<template>
  <div class="login-container">
    <NCard class="login-card shadow-2xl">
      <div class="login-title">
        <NH2>重置密码</NH2>
      </div>
      <NForm :model="form" label-width="100px">
        <NFormItem label="邮箱">
          <NInput v-model:value="form.username" :disabled="route.query?.u ? true : false" placeholder="请输入邮箱">
            <template #prefix>
              <SvgIcon icon="ph:user-bold" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem label="新密码">
          <NInput v-model:value="form.password" type="password" placeholder="请输入新密码">
            <template #prefix>
              <SvgIcon icon="mdi:password-outline" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem label="再次输入新密码">
          <NInput v-model:value="form.passwordAgain" type="password" placeholder="请再次输入新密码">
            <template #prefix>
              <SvgIcon icon="mdi:password-outline" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem label="图形验证码">
          <div class="w-[130px] h-[40px] mr-[20px] rounded border flex cursor-pointer">
            <Captcha ref="captchaRef" src="/api/captcha/getImage" />
          </div>
          <NInput v-model:value="form.vcode" size="large" type="text" placeholder="请输入图像验证码" clearable />
        </NFormItem>

        <NFormItem label="邮箱验证码">
          <NInputGroup>
            <NInput v-model:value="form.emailVCode" placeholder="请输入邮箱的验证码" />
            <NButton type="primary" :loading="verificationLoading" :disabled="countdown > 0" ghost @click="startCountdown(form.vcode)">
              <span v-if="countdown > 0">{{ `${countdown}秒后重新获取` }} </span>
              <span v-else>获取邮箱验证码</span>
            </NButton>
          </NInputGroup>
        </NFormItem>

        <!-- <NFormItem>
          <span>
            <NImage
              class="w-[140px] h-[34px] flex"
              src="https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg"
            />
          </span>
          <NInput v-model:value="form.vcode" type="text" placeholder="图像验证码" />
        </NFormItem> -->
        <NFormItem style="margin-top: 10px">
          <NButton type="primary" block @click="handleSubmit">
            确定
          </NButton>
        </NFormItem>
      </NForm>
    </NCard>
    <Verification
      v-model:visible="verificationShow"
      v-model:loading="verificationLoading"
      :verification-id="verificationId"
      captcha-src="/api/captcha/getImage/200/60"
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
        /* background-color: #f2f6ff; */
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
