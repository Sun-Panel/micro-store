<script setup lang="ts">
import type { Language } from '@/store/modules/app/helper'
import { NButton, NCard, NForm, NFormItem, NGradientText, NInput, NSelect, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { login } from '@/api'
import { getLoginConfig } from '@/api/openness'
import { Captcha, SvgIcon } from '@/components/common'
import { t } from '@/locales'
import { router } from '@/router'
import { useAppStore, useAuthStore } from '@/store'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'
import { languageOptions } from '@/utils/defaultData'

// const userStore = useUserStore()
const authStore = useAuthStore()
const appStore = useAppStore()
// const route = useRoute()
const ms = useMessage()
const loading = ref(false)
const languageValue = ref<Language>(appStore.language)
const captchaRef = ref()

const isShowCaptcha = ref<boolean>(false)
const isShowRegister = ref<boolean>(false)

const form = ref<Login.LoginReqest>({
  username: '',
  password: '',
})

async function loginPost() {
  loading.value = true
  try {
    const res = await login<Login.LoginResponse>(form.value)
    authStore.setToken(res.data.token)
    authStore.setUserInfo(res.data)

    setTimeout(() => {
      ms.success(`Hi ${res.data.name},${t('login.welcomeMessage')}`)

      router.push('/')
      // if (authStore.userInfo?.role === 1)
      //   router.push('/')

      // else
      //   router.push('/')

      // if (encodeURIComponent(route.query.callback as string))
      //   // location.href = `#/${route.query.callback}` as string
      //   router.push(route.query.callback as string)

      // else
      //   router.push({ path: appStore.homeBase?.home_url as string })
    }, 500)

    loading.value = false
  }
  catch (error) {
    captchaRef.value.refresh()
    loading.value = false
    const { code, msg } = error as Common.Response
    if (code === -1) {
      ms.error(msg)
      return
    }

    apiRespErrMsg(error)
  }
}

function handleSubmit() {
  // 点击登录按钮触发
  loginPost()
}

async function getLoginVcodeApi() {
  const { data } = await getLoginConfig<Openness.open.LoginVcodeResponse>()
  isShowCaptcha.value = data.loginCaptcha
  isShowRegister.value = data.register.openRegister
}

function handleChangeLanuage(value: Language) {
  languageValue.value = value
  appStore.setLanguage(value)
}

onMounted(() => {
  getLoginVcodeApi()
  // getNotice(1)
})
</script>

<template>
  <div class="login-container">
    <NCard class="login-card shadow-2xl">
      <div class="mb-5 flex items-center justify-end">
        <div class="mr-2">
          <SvgIcon icon="ion-language" style="width: 20;height: 20;" />
        </div>
        <div class="min-w-[100px]">
          <NSelect v-model:value="languageValue" size="small" :options="languageOptions" @update-value="handleChangeLanuage" />
        </div>
      </div>

      <div class="login-title  ">
        <NGradientText :size="30" type="danger" class="!font-bold">
          {{ $t('common.appName') }} 微应用
        </NGradientText>
      </div>
      <NForm :model="form" label-width="100px" size="large" @keydown.enter="handleSubmit">
        <NFormItem>
          <NInput v-model:value="form.username" :placeholder="$t('login.usernamePlaceholder')">
            <template #prefix>
              <SvgIcon icon="ph:user-bold" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem>
          <NInput v-model:value="form.password" size="large" type="password" :placeholder="$t('login.passwordPlaceholder')">
            <template #prefix>
              <SvgIcon icon="mdi:password-outline" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem v-if="isShowCaptcha">
          <div class="w-[130px] h-[40px] mr-[20px] rounded border flex cursor-pointer">
            <Captcha ref="captchaRef" src="/api/captcha/getImage" />
          </div>
          <NInput v-model:value="form.vcode" size="large" type="text" clearable :placeholder="t('login.captchaCodePlaceholder')" />
        </NFormItem>
        <NFormItem style="margin-top: 10px">
          <NButton type="primary" block :loading="loading" @click="handleSubmit">
            {{ $t('login.loginButton') }}
          </NButton>
        </NFormItem>

        <!-- <div class="flex justify-end">
          <NButton v-if="isShowRegister" quaternary type="info" class="flex" @click="$router.push({ path: '/register' })">
            {{ t('login.register') }}
          </NButton>
          <NButton quaternary type="info" class="flex" @click="$router.push({ path: '/resetPassword' })">
            {{ t('login.forgetPassword') }}
          </NButton>
        </div> -->

        <div class="flex justify-center text-slate-300">
          Powered By <a href="https://github.com/hslr-s/sun-panel" target="_blank" class="ml-[5px] text-slate-500">Sun-Panel</a>
        </div>
      </NForm>
    </NCard>
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
