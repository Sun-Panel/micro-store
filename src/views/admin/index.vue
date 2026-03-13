<script lang="ts" setup>
import { NLayout, NLayoutSider, NMenu, NSpace } from 'naive-ui'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAdminStore } from '@/store'
import Header from './layout/Header/index.vue'
// import {
//   BookOutline as BookIcon,
//   PersonOutline as PersonIcon,
//   WineOutline as WineIcon,
// } from '@vicons/ionicons5'
const adminStore = useAdminStore()
const inverted = ref(false)
const router = useRouter()
const route = useRoute()

const menuOptions = [

  // {
  //   label: '后台首页',
  //   key: 'AdminHome',
  // },
  {
    label: '仪表盘',
    key: 'Dashboard',
  },
  {
    label: '用户管理',
    key: 'user-manage',
    children: [
      {
        label: '用户列表',
        key: 'AdminUserManage',
      },
    ],
  },
  {
    label: '微应用管理',
    key: 'micro-app-manage',
    children: [
      {
        label: '分类管理',
        key: 'AdminMicroAppCategory',
      },
    ],
  },

  // {
  //   label: '兑换码',
  //   key: 'AdminRedeemCode',
  // },

  // {
  //   label: '公告管理',
  //   key: 'NoticeMange',
  // },

  {
    label: '系统管理',
    key: 'system-manage',
    children: [
      {
        label: '邮箱设置',
        key: 'AdminSystemEmailSetting',
      },
      {
        label: '网站设置',
        key: 'AdminSystemWebsiteSetting',
      },
      // {
      //   label: '关于设置',
      //   key: 'AdminSystemAboutSetting',
      // },
      {
        label: '系统变量',
        key: 'AdminSystemVariable',
      },
      // {
      //   label: 'Markdown页面',
      //   key: 'AdminMdPageManage',
      // },

    ],
  },
]

//   export default defineComponent({
//     setup () {
//       return {
//         inverted: ref(false),
//         menuOptions
//       }
//     }
//   })
</script>

<template>
  <NSpace vertical>
    <NLayout>
      <Header />
      <NLayout has-sider>
        <NLayoutSider
          bordered
          collapse-mode="width"
          :collapsed-width="0"
          :width="240"
          :native-scrollbar="false"
          :inverted="inverted"
          :collapsed="adminStore.siderCollapsed"
          class="h-[calc(100vh-60px)]"
        >
          <NMenu
            :value="route.name as string"
            :collapsed="adminStore.siderCollapsed"
            :inverted="inverted"
            :collapsed-width="0"
            :collapsed-icon-size="0"
            :options="menuOptions"
            :default-expanded-keys="['user-manage', 'EmailNotify']"
            @update:value="(key: string) => router.push({ name: key })"
          />
        </NLayoutSider>
        <NLayout class="h-[calc(100vh-60px)] dark:bg-slate-900 bg-slate-100 ">
          <router-view class=" h-[calc(100vh-60px)] p-[20px] bg-slate-100 dark:bg-slate-900" />
          <!-- <RouterView v-slot="{ Component, route }">
            <component :is="Component" :key="route.fullPath" />
          </RouterView> -->
        </NLayout>
      </NLayout>
      <!-- <NLayoutFooter :inverted="inverted" bordered>
        Footer Footer Footer
      </NLayoutFooter> -->
    </NLayout>
  </NSpace>
</template>
