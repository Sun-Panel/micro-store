<script setup lang="ts">
import type { MenuOption } from 'naive-ui'
import { NMenu } from 'naive-ui'
import { computed, h, ref } from 'vue'
import { SvgIconOnline } from '@/components/common'
import { t } from '@/locales'
import { useAuthStore } from '@/store'
import { hasRole, ROLE_DEVELOPER } from '@/utils/role'

withDefaults(defineProps<{
  isVertical: boolean
}>(), {
  isVertical: false,
})

const authStore = useAuthStore()

// 判断是否有开发者权限
const hasDeveloperPermission = computed(() => hasRole(authStore.userInfo!.role || 0, ROLE_DEVELOPER))

const devDocLinks = 'https://doc.sun-panel.top/v2/zh_cn/micro_app_dev/'

const activeKey = ref('aaa')
const becomeDeveloperOption: MenuOption = {
  label: () => a('/developer/register', t('menu.becomeDeveloper')),
  key: 'becomeDeveloper',
}
const myMicroAppOption: MenuOption = {
  label: () => a('/admin/developerCenter/myMicroApp', t('menu.myMicroApp')),
  key: 'appManagement',
}

const menuOptions = computed<MenuOption[]>(() => {
  const options: MenuOption[] = []

  if (hasDeveloperPermission.value) {
    // 有开发者权限：显示"我的微应用"
    options.push(myMicroAppOption)
  }
  else {
    // 没有开发者权限：显示"成为开发者"
    options.push(becomeDeveloperOption)
  }

  options.push({
    label: () => aBlank(devDocLinks, t('menu.devDoc')),
    key: 'community',
  })

  return options
})

function a(url: string, text: string) {
  return h(
    'a',
    {
      href: url,
      // rel: 'noopenner noreferrer',
    },
    text,
  )
}

function aBlank(url: string, text: string) {
  return h(
    'a',
    {
      href: url,
      target: '_blank',
      style: { display: 'flex', alignItems: 'center' },
      // rel: 'noopenner noreferrer',
    },
    [
      text,
      h(SvgIconOnline, { icon: 'ion:open-outline', style: { marginLeft: '4px' } }),
    ],
  )
}
</script>

<template>
  <NMenu
    v-model:value="activeKey"
    :mode="isVertical ? 'vertical' : 'horizontal'"
    :options="menuOptions"
    responsive
  />
</template>
