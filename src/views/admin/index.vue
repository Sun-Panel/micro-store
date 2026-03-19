<script lang="ts" setup>
import type { AdminMenuItem } from '@/utils/role'
import { NLayout, NLayoutSider, NMenu, NSpace } from 'naive-ui'
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAdminStore, useAuthStore } from '@/store'
import { filterAndConvertMenu, hasRole, ROLE_ADMIN, ROLE_AUDITOR, ROLE_DEVELOPER } from '@/utils/role'
import Header from './layout/Header/index.vue'

const adminStore = useAdminStore()
const authStore = useAuthStore()
const inverted = ref(false)
const router = useRouter()
const route = useRoute()

// 用户角色
const userRole = computed(() => authStore.userInfo?.role || 0)

// 菜单配置（可在此处使用多语言 t() 函数）
const ADMIN_MENU_CONFIG: AdminMenuItem[] = [
  {
    label: '仪表盘',
    key: 'Dashboard',
    roles: ROLE_ADMIN,
  },
  {
    label: '用户管理',
    key: 'user-manage',
    roles: ROLE_ADMIN,
    children: [
      { label: '用户列表', key: 'AdminUserManage' },
    ],
  },
  {
    label: '微应用管理',
    key: 'micro-app-manage',
    roles: ROLE_ADMIN,
    children: [
      { label: '分类管理', key: 'AdminMicroAppCategory' },
      { label: '开发者管理', key: 'AdminDeveloper' },
      { label: '应用管理', key: 'AdminMicroAppManage' },
    ],
  },
  {
    label: '审核管理',
    key: 'review-manage',
    roles: ROLE_ADMIN | ROLE_AUDITOR,
    children: [
      { label: '应用审核', key: 'AdminMicroAppReview' },
      { label: '版本审核', key: 'AdminVersionReview' },
    ],
  },
  {
    label: '开发者中心',
    key: 'developer-center',
    roles: ROLE_DEVELOPER | ROLE_ADMIN,
    children: [
      { label: '我的微应用', key: 'AdminMyMicroApp' },
    ],
  },
  {
    label: '系统管理',
    key: 'system-manage',
    roles: ROLE_ADMIN,
    children: [
      { label: '邮箱设置', key: 'AdminSystemEmailSetting' },
      { label: '网站设置', key: 'AdminSystemWebsiteSetting' },
      { label: '系统变量', key: 'AdminSystemVariable' },
    ],
  },
]

// 根据用户角色过滤后的菜单
const menuOptions = computed(() => {
  // 管理员可以看到所有菜单（使用 hasRole 兼容多角色）
  if (hasRole(userRole.value, ROLE_ADMIN)) {
    return ADMIN_MENU_CONFIG
  }
  return filterAndConvertMenu(userRole.value, ADMIN_MENU_CONFIG)
})
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
            :default-expanded-keys="['user-manage', 'micro-app-manage', 'review-manage', 'developer-center', 'system-manage']"
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
