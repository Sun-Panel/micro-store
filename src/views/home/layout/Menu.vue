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
const hasDeveloperPermission = computed(() => {
  const role = authStore.userInfo?.role || 0
  return hasRole(role, ROLE_DEVELOPER)
})

const devDocLinks = 'https://doc.sun-panel.top/v2/zh_cn/micro_app_dev/'

const activeKey = ref('aaa')

const publishMicroAppOption = computed<MenuOption>(() => ({
  label: () => a(
    hasDeveloperPermission.value
      ? '/admin/developerCenter/myMicroApp'
      : '/developer/register',
    t('menu.publishMicroApp'),
  ),
  key: 'publishMicroApp',
}))

const menuOptions = computed<MenuOption[]>(() => {
  const options: MenuOption[] = []

  // 始终显示"发布微应用"按钮，根据权限跳转不同页面
  options.push(publishMicroAppOption.value)

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
