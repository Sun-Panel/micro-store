<script setup lang="ts">
import type { Language } from '@/store/modules/app/helper'
import { NResult, NSpin } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { oAuth2CodeBind, oAuth2CodeLogin } from '@/api/login'
import { t } from '@/locales'
import { router } from '@/router'
import { useAppStore, useAuthStore } from '@/store'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'

const route = useRoute()

const authStore = useAuthStore()
const appStore = useAppStore()

const code = route.query.code as string
const lang = route.query.lang as string
const callback = route.query.callback as string
const isBindStr = route.query.isBind as string
const loginLoading = ref(false)
const errMsg = ref('')

async function login() {
  appStore.setLanguage(lang as Language)
  loginLoading.value = true
  await oAuth2CodeLogin<Login.OAuth2CodeLoginResq>(code).then(({ data }) => {
    authStore.setToken(data.token)
    authStore.setUserInfo(data)

    if (callback)
      router.push(callback)
    else
      router.push({ path: appStore.homeBase?.home_url as string })
  }).catch((res) => {
    switch (res.code) {
      case -2:
        errMsg.value = t('oAuth2.reLogin')
        // 需要重新跳转到登录端点进行登录
        break
      case -3:
        // 需要重新跳转到登录端点进行登录
        errMsg.value = t('oAuth2.logoutAgainLogin')
        break
      default:
        apiRespErrMsg(res)
        break
    }
  })
  loginLoading.value = false
}

async function bind() {
  loginLoading.value = true
  await oAuth2CodeBind(code).then(() => {
    if (callback)
      router.push(callback)
    else
      router.push({ path: appStore.homeBase?.home_url as string })
  }).catch((res) => {
    switch (res.code) {
      case -2:
        // 需要重新跳转到登录端点进行登录
        errMsg.value = t('oAuth2.reLogin')
        break
      case -3:
        // 需要重新跳转到登录端点进行登录
        errMsg.value = t('oAuth2.logoutAgainLogin')
        break
      default:
        apiRespErrMsg(res)
        break
    }
  })
  loginLoading.value = false
}

onMounted(() => {
  if (isBindStr === 'true')
    bind()
  else
    login()
})
</script>

<template>
  <template v-if="loginLoading">
    <div class="flex justify-center items-center flex-col">
      <NSpin />
      <span class="text-xl mt-5 font-bold">
        Loading
      </span>
    <!-- 当前code：{{ code }}

    <NButton @click="router.push(callback)">
      假装通过直接直接跳转
    </NButton>

    <NButton @click="login">
      尝试获取token
    </NButton> -->
    </div>
  </template>
  <template v-else-if="errMsg">
    <NResult
      status="error"
    >
      <template #footer>
        {{ errMsg }}
      </template>
    </NResult>
  </template>
</template>
