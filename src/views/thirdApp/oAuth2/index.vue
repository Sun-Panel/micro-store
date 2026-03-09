<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NButton } from 'naive-ui'
import { useRoute } from 'vue-router'
import { getThirdAppInfo } from '@/api/thirdApp/thirdApp'
import { authLogin } from '@/api/thirdApp/oAuth2'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'
import { t } from '@/locales'

const route = useRoute()
const appid = route.query.appid as string
const response_type = route.query.response_type as string
const redirect_uri = route.query.redirect_uri as string
const appInfo = ref<ThirdApp.ThirdApp.GetThirdAppInfoResp>()

// 获取授权码并会跳到三方应用
// http://127.0.0.1:1003/#/authThirdAppLogin?appid=test_appid&redirect_uri=http://127.0.0.1&response_type=code
async function getCodeAndredirect() {
  await authLogin<{ code: string }>(appid, response_type, redirect_uri)
    .then(({ data }) => {
      const code = data.code
      // 获取到code，直接重定向
      const url = addQueryParamToUrl(redirect_uri, 'code', code)
      location.href = url
    })
    .catch((res) => {
      console.error(res)
      apiRespErrMsg(res)
    })
}

function addQueryParamToUrl(url: string, paramName: string, paramValue: string): string {
  // 创建一个 URL 对象
  const urlObj = new URL(url)
  // 设置新的查询参数
  urlObj.searchParams.set(paramName, paramValue)
  // 返回带有新参数的 URL 字符串
  return urlObj.toString()
}

onMounted(async () => {
  await getThirdAppInfo<ThirdApp.ThirdApp.GetThirdAppInfoResp>(appid).then(({ data }) => {
    appInfo.value = data

    // 自动授权，直接回调
    if (appInfo.value.isAutoAuth)
      getCodeAndredirect()
  }).catch((res) => {
    apiRespErrMsg(res)
  })
})
</script>

<template>
  <div v-if="appInfo">
    <NButton @click="getCodeAndredirect()">
      {{ t("thirdApp.oAuth2.agreeAuthLogin", { appName: appInfo?.appName }) }}
    </NButton>
  </div>
</template>
