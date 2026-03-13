<script setup lang="ts">
import type { DropdownOption } from 'naive-ui'
import type { RouteLocationRaw } from 'vue-router'
import { NBadge, NButton, NDrawer, NDrawerContent, NDropdown, NImage, useMessage } from 'naive-ui'

import { onMounted, onUnmounted, ref } from 'vue'
import { logout } from '@/api'
import { getHomeBase, getLoginConfig } from '@/api/openness'
import { SvgIconOnline } from '@/components/common'
import { t } from '@/locales'
import { router } from '@/router'
import { useAppStore, useAuthStore, useUserStore } from '@/store'
import Menu from './Menu.vue'

const userStore = useUserStore()
const authStore = useAuthStore()
const appStore = useAppStore()

const ms = useMessage()
const isShowRegister = ref(false)
const unReadCount = ref(0)
const homeBase = ref<Openness.open.HomeBase>()
const mobileDrawerShow = ref(false)
const isMobile = ref(false)

const timerId: any = 0

const options = [
  // {
  //   label: t('menu.myMessage'),
  //   key: 'PlatformMessage',
  //   type: 'render',
  //   render: myMessageRender,
  // },
  {
    key: 'header-divider',
    type: 'divider',
  },
  // {
  //   label: t('menu.userInfo'),
  //   key: 'PlatformUserInfo',
  // },
  // {
  //   label: t('menu.prAuthorizeInfo'),
  //   key: 'PlatformProAuthorize',
  // },
  // {
  //   label: t('menu.myOrder'),
  //   key: 'PlatformOrder',
  // },
  {
    label: t('common.logout'),
    key: 'logout',
  },
]

// function myMessageRender() {
//   return h(
//     'div',
//     {
//       style: { display: 'flex', aligItems: 'center', marginLeft: '14px', marginTop: '10px', marginBottom: '10px', cursor: 'pointer' },
//       onClick() {
//         router.push({ name: 'PlatformMessage' })
//       },
//     },
//     [
//       h('div', {}, t('menu.myMessage')),
//       h(NBadge, {
//         max: 99,
//         value: unReadCount.value,
//         style: { marginLeft: '4px' },
//       }),
//     ],
//   )
// }

async function logoutApi() {
  await logout()
  userStore.resetUserInfo()
  authStore.removeToken()
  router.push('/')
  ms.success(t('settingUserInfo.logoutSuccess'))
  setTimeout(() => {
    // appStore.removeToken()
    location.reload()// 强制刷新一下页面
  }, 200)
}

function handleSelect(key: string | number, option: DropdownOption) {
  switch (key) {
    case 'logout':
      logoutApi()
      break

    default:
      router.push({ name: key as string })
      break
  }
}

async function getLoginVcodeApi() {
  const { data } = await getLoginConfig<Openness.open.LoginVcodeResponse>()
  isShowRegister.value = data.register.openRegister
}

// async function getUnReadMessageCount() {
//   if (!authStore.userInfo)
//     return

//   try {
//     const { data } = await getUnReadCount<{ count: number }>()
//     unReadCount.value = data.count
//   }
//   catch (error) {
//   }
// }

async function getHomeBasePost() {
  try {
    const { data } = await getHomeBase<Openness.open.HomeBase>()
    homeBase.value = data
    document.title = data.logo_text as string
    appStore.setHomeBase(homeBase.value)
  }
  catch (error) {
    ms.error('服务器出错了')
  }
}

function handleBackHome() {
  location.href = homeBase.value?.logo_click_to_link as string
}

function handleGoToPage(option: RouteLocationRaw) {
  mobileDrawerShow.value = false
  router.push(option)
}

function handleResize() {
  if (window.innerWidth > 800)
    isMobile.value = false

  else
    isMobile.value = true
}

function goOAuth2(callbackUrl: string) {
  window.location.href = `/api/oAuth2/v1/login?callback=${callbackUrl}`
}

onMounted(() => {
  getHomeBasePost()
  getLoginVcodeApi()
  appStore.setTheme('light')
  // getUnReadMessageCount()
  // timerId = setInterval(getUnReadMessageCount, 5000)
  // handleResize()
  // window.addEventListener('resize', handleResize)
  // getNotice(1)
})

onUnmounted(() => {
  clearInterval(timerId)
  window.removeEventListener('resize', handleResize)
})
</script>

<template>
  <div class="border-b">
    <!-- header -->
    <div class="flex max-w-[1200px] mx-auto px-2">
      <div class="w-full flex items-center h-[60px] ">
        <div class="flex cursor-pointer" @click="handleBackHome">
          <div class="flex w-[50px]">
            <NImage
              :width="40"
              preview-disabled
              :src="homeBase?.logo_url"
            />
          </div>
          <div class="flex items-center text-xl font-bold whitespace-nowrap dark:text-white">
            <!-- {{ homeBase?.logo_text }} -->
            Sun-Panel | {{ t('common.microAppStore') }}
          </div>
        </div>

        <div v-if="!isMobile" class="w-full flex justify-end">
          <span>
            <Menu />
          </span>
        </div>

        <div v-if="!isMobile" class="min-w-[150px] mx-5">
          <template v-if="!authStore.userInfo?.token">
            <!-- <span v-if="isShowRegister" class="mr-4">
              <NButton type="info" size="small" ghost @click="handleGoToPage({ name: 'register' })">
                {{ t('login.register') }}
              </NButton>
            </span> -->
            <span class="mr-4">
              <NButton size="small" @click="handleGoToPage({ name: 'login' })">
                {{ t('login.login') }}
              </NButton>
            </span>
            <span class="mr-4">
              <NButton type="success" size="small" @click="goOAuth2('/')">
                授权登录
              </NButton>
            </span>
          </template>
          <template v-else>
            <div class="flex justify-end">
              <NDropdown :options="options" @select="handleSelect">
                <NBadge :value="unReadCount" :max="99">
                  <NButton type="info" size="small" ghost>
                    {{ authStore.userInfo?.name }}
                  </NButton>
                </NBadge>
              </NDropdown>

              <!-- <NButton v-if="authStore.userInfo?.role === 1" type="info" size="small" class="!mx-5" @click="handleGoToPage('/admin')">
                后台管理
              </NButton> -->
            </div>
          </template>
        </div>

        <div v-if="isMobile" class="w-full flex justify-end">
          <SvgIconOnline
            icon="material-symbols:menu"
            class="w-8 h-8 cursor-pointer"
            @click="mobileDrawerShow = true"
          />
        </div>
      </div>
    </div>
    <NDrawer v-model:show="mobileDrawerShow" weight="100%" placement="right">
      <NDrawerContent closable>
        <template #header>
          <div class="min-w-[150px] ">
            <template v-if="!authStore.userInfo?.id">
              <!-- <span v-if="isShowRegister" class="mr-4">
                <NButton type="info" size="small" ghost @click="handleGoToPage({ name: 'register' })">
                  {{ t('login.register') }}
                </NButton>
              </span> -->
              <span>
                <NButton type="success" size="small" @click="handleGoToPage({ name: 'login' })">
                  {{ t('login.login') }}
                </NButton>
              </span>
              <span>
                <NButton type="success" size="small" @click="handleGoToPage({ name: 'login' })">
                  授权登录
                </NButton>
              </span>
            </template>
            <template v-else>
              <NDropdown :options="options" @select="handleSelect">
                <NBadge :value="unReadCount" :max="99">
                  <NButton type="info" size="small" ghost>
                    {{ authStore.userInfo?.name }}
                  </NButton>
                </NBadge>
              </NDropdown>
            </template>
          </div>
        </template>
        <!-- <Menu :is-vertical="true" /> -->
      </NDrawerContent>
    </NDrawer>
  </div>
</template>
