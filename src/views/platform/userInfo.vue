<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { NButton, NCard, NH2, NInput, useLoadingBar, useMessage } from 'naive-ui'
import { useAuthStore } from '@/store'
import { buildTimeString } from '@/utils/cmn'
import { getAuthorize as getAuthorizeApi } from '@/api/proAuthorize'
import { getInfo as getUserInfoApi, unBindSunStore as unBindSunStoreApi, updateInfo } from '@/api/system/user'
import { t } from '@/locales'
import { apiRespErrMsg } from '@/utils/cmn/apiMessage'

const proAuthorizeExpiredTime = ref('')
const editName = ref(false)
const authStore = useAuthStore()
const nikeName = ref(authStore.userInfo?.name || '')
const ms = useMessage()
const loadingBar = useLoadingBar()
const userInfo = ref<User.Info>()

async function getAuthorize() {
  try {
    const { data } = await getAuthorizeApi<ProAuthorize.GetAuthorizeResp>()
    proAuthorizeExpiredTime.value = data.expiredTime
  }
  catch (error) {
    console.error(error)
  }
}

async function getUserInfo() {
  await getUserInfoApi<User.Info>().then(({ data }) => {
    userInfo.value = data
  }).catch((res) => {
    apiRespErrMsg(res)
  })
}

async function handleUpdateInfo() {
  loadingBar.start()
  await updateInfo(nikeName.value).then(() => {
    editName.value = false
    authStore.refreshUserInfo()
  }).catch(() => {
    ms.error(t('common.saveFail'))
  })
  loadingBar.finish()
}

function handleBindSunStore() {
  location.href = generateUrl()
}

function generateUrl() {
  const callbackUrl = encodeURIComponent('/platform/userInfo')
  console.log('回调地址', callbackUrl)
  return `/api/oAuth2/v1/login?callback=${callbackUrl}&isBind=true`
}

function unBindSunStore() {
  unBindSunStoreApi().then(() => {
    getUserInfo()
  }).catch((res) => {
    apiRespErrMsg(res)
  })
}

onMounted(() => {
  getAuthorize()
  getUserInfo()
})
</script>

<template>
  <div>
    <NH2 align-text prefix="bar">
      个人信息
    </NH2>

    <NCard>
      <div class="item-box">
        <div class="item">
          <span>
            昵称：
            <span v-if="!editName">
              {{ nikeName }}
            </span>
            <span v-if="editName">
              <NInput
                v-model:value="nikeName"
                size="small"
                style="max-width: 200px;"
                type="text"
                placeholder="请输入"
              />
            </span>
          </span>

          <span class="ml-5">
            <NButton v-if="!editName" size="tiny" tertiary type="info" @click="editName = true">
              修改
            </NButton>
            <NButton v-if="editName" size="tiny" type="info" @click="handleUpdateInfo">
              保存
            </NButton>
          </span>
        </div>

        <div class="item">
          账号/邮箱：{{ authStore.userInfo?.username }}
        </div>

        <div class="item">
          注册日期：{{ buildTimeString(authStore.userInfo?.createTime, "YYYY-MM-DD") }}
        </div>

        <div class="item">
          PRO到期时间：{{ buildTimeString(proAuthorizeExpiredTime, "YYYY-MM-DD") || "未授权" }}
        </div>

        <div class="item">
          SunStore平台：
          <template v-if="!userInfo?.isBindSunStore">
            <NButton type="warning" size="small" @click="handleBindSunStore">
              点击绑定
            </NButton>
          </template>
          <template v-else>
            已绑定
            <NButton quaternary type="error" size="small" @click="unBindSunStore">
              解绑
            </NButton>
          </template>
        </div>

        <div class="mt-5">
          <NButton type="warning" @click="ms.warning('请到主平台修改密码 联合平台不支持单独修改密码，请联系管理员')">
            重置密码
          </NButton>
        </div>
      </div>
    </NCard>
  </div>
</template>

<style scoped>
.item-box>.item {
    margin-top: 10px;
    font-size: 16px;
}
</style>
