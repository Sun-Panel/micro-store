import type { App } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import { setupPageGuard } from './permission'
import AdminLayout from '@/views/admin/index.vue'
import HomeLayout from '@/views/home/layout/index.vue'

const adminRouter = {
  path: '/admin',
  name: 'Admin',
  component: AdminLayout,
  children: [
    // {
    //   path: '/admin/',
    //   name: 'AdminHome',
    //   component: () => import('@/views/admin/home/index.vue'),
    // },
    {
      path: '/admin/',
      name: 'Dashboard',
      component: () => import('@/views/admin/dashboard/index.vue'),
    },
    {
      path: '/admin/userManage',
      name: 'AdminUserManage',
      component: () => import('@/views/admin/userManage/index.vue'),
    },
    {
      path: '/admin/emailSetting',
      name: 'AdminSystemEmailSetting',
      component: () => import('@/views/admin/emailSetting/index.vue'),
    },
    {
      path: '/admin/websiteSetting',
      name: 'AdminSystemWebsiteSetting',
      component: () => import('@/views/admin/websiteSetting/index.vue'),
    },
    {
      path: '/admin/noticeMange',
      name: 'NoticeMange',
      component: () => import('@/views/admin/noticeMange/index.vue'),
    },
    {
      path: '/admin/aboutSetting',
      name: 'AdminSystemAboutSetting',
      component: () => import('@/views/admin/aboutSetting/index.vue'),
    },
    {
      path: '/admin/emailTemplate',
      name: 'AdminEmailTemplate',
      component: () => import('@/views/admin/emailTemplate/index.vue'),
    },
    {
      path: '/admin/email',
      name: 'AdminEmailSend',
      component: () => import('@/views/admin/email/send.vue'),
    },
    {
      path: '/admin/messge/template',
      name: 'MessageTemplate',
      component: () => import('@/views/admin/message/template.vue'),
    },
    {
      path: '/admin/goodsManage',
      name: 'AdminGoodsManage',
      component: () => import('@/views/admin/goodsManage/index.vue'),
    },
    {
      path: '/admin/proAuthorizeManage',
      name: 'AdminProAuthorizeManage',
      component: () => import('@/views/admin/proAuthorizeManage/index.vue'),
    },
    {
      path: '/admin/redeemCode',
      name: 'AdminRedeemCode',
      component: () => import('@/views/admin/redeemCode/index.vue'),
    },
    {
      path: '/admin/orderManage',
      name: 'AdminOrderManage',
      component: () => import('@/views/admin/orderManage/index.vue'),
    },
    {
      path: '/admin/systemVariable',
      name: 'AdminSystemVariable',
      component: () => import('@/views/admin/systemVariable/index.vue'),
    },
    {
      path: '/admin/mdPageManage',
      name: 'AdminMdPageManage',
      component: () => import('@/views/admin/mdPageManage/index.vue'),
    },
    {
      path: '/admin/versionManage',
      name: 'AdminVersionManage',
      component: () => import('@/views/admin/versionManage/index.vue'),
    },
    {
      path: '/admin/clientBlacklistIP',
      name: 'AdminClientBlacklistIP',
      component: () => import('@/views/admin/clientBlacklistIP/index.vue'),
    },
    {
      path: '/admin/clientCreateOnlineCache',
      name: 'AdminClientCreateOnlineCache',
      component: () => import('@/views/admin/clientCreateOnlineCache/index.vue'),
    },
    {
      path: '/admin/microAppCategory',
      name: 'AdminMicroAppCategory',
      component: () => import('@/views/admin/microAppCategory/index.vue'),
    },
    {
      path: '/admin/developer',
      name: 'AdminDeveloper',
      component: () => import('@/views/admin/developer/index.vue'),
    },
    {
      path: '/admin/myMicroApp',
      name: 'AdminMyMicroApp',
      component: () => import('@/views/admin/myMicroApp/index.vue'),
    },
    {
      path: '/admin/microAppManage',
      name: 'AdminMicroAppManage',
      component: () => import('@/views/admin/microAppManage/index.vue'),
    },
  ],
}

const platformLogin: RouteRecordRaw[] = [
  // {
  //   path: '/platform/order',
  //   name: 'PlatformOrder',
  //   component: () => import('@/views/platform/order.vue'),
  // },
  // {
  //   path: '/platform/pay',
  //   name: 'PlatformPay',
  //   component: () => import('@/views/platform/pay.vue'),
  // },
  // {
  //   path: '/platform/packageStore',
  //   name: 'PlatformPackageStore',
  //   component: () => import('@/views/packageStore/index.vue'),
  // },
  // {
  //   path: '/platform/orderInfo',
  //   name: 'PlatformOrderInfo',
  //   component: () => import('@/views/platform/orderInfo.vue'),
  // },
  {
    path: '/oAuth2/login',
    name: 'PlatformOrder',
    component: () => import('@/views/oAuth2/index.vue'),
  },
  {
    path: '/platform/proAuthorize',
    name: 'PlatformProAuthorize',
    component: () => import('@/views/platform/proAuthorize.vue'),
  },
  // {
  //   path: '/platform/message',
  //   name: 'PlatformMessage',
  //   component: () => import('@/views/platform/message/index.vue'),
  // },
  {
    path: '/platform/donateReward',
    name: 'DonateReward',
    component: () => import('@/views/platform/donateReward.vue'),
  },
  {
    path: '/platform/userInfo',
    name: 'PlatformUserInfo',
    component: () => import('@/views/platform/userInfo.vue'),
  },

  // {
  //   path: '/mdPage/:mdPageName',
  //   name: 'MdPage',
  //   component: () => import('@/views/mdPage/index.vue'),
  // },
  // {
  //   path: '/paypal',
  //   name: 'Paypal',
  //   component: () => import('@/views/packageStore/order/paypal.vue'),
  // },
  // {
  //   path: '/resetPassword',
  //   name: 'resetPassword',
  //   component: () => import('@/views/login/resetPassword.vue'),
  // },
  // {
  //   path: '/register',
  //   name: 'register',
  //   component: () => import('@/views/register/index.vue'),
  // },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
  },
  {
    path: '/authThirdAppLogin',
    name: 'authThirdAppLogin',
    component: () => import('@/views/thirdApp/oAuth2/index.vue'),
  },

]

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'HomeLayout',
    component: HomeLayout,
    children: [
      {
        path: '/',
        name: 'home',
        component: () => import('@/views/home/index.vue'),
      },
      {
        path: '/developer/register',
        name: 'DeveloperRegister',
        component: () => import('@/views/developer/register.vue'),
      },
      {
        path: '/pro',
        name: 'pro',
        component: () => import('@/views/home/pro.vue'),
      },
      {
        path: '/platformLogin',
        name: 'platformLogin',
        children: platformLogin,
      },
    ],
  },

  {
    path: '/404',
    name: '404',
    component: () => import('@/views/exception/404/index.vue'),
  },

  {
    path: '/500',
    name: '500',
    component: () => import('@/views/exception/500/index.vue'),
  },

  {
    path: '/:pathMatch(.*)*',
    name: 'notFound',
    redirect: '/404',
  },

  {
    path: '/test',
    name: 'test',
    component: () => import('@/views/exception/test/index.vue'),
  },

  adminRouter,
]

export const router = createRouter({
  // history: createWebHashHistory(),
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

setupPageGuard(router)

export async function setupRouter(app: App) {
  app.use(router)
  await router.isReady()
}
