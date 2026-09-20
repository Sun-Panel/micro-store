<script setup lang="ts">
import type { MenuOption } from 'naive-ui'
import { NMenu } from 'naive-ui'
import { computed, h, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
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
const route = useRoute()
const router = useRouter()

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

// 发布自定义代码片段：与"发布微应用"逻辑一致，开发者进入管理页，非开发者引导注册
const publishCustomCodeOption = computed<MenuOption>(() => ({
  label: () => a(
    hasDeveloperPermission.value
      ? '/admin/customCode/myList'
      : '/developer/register',
    t('menu.publishCustomCode'),
  ),
  key: 'publishCustomCode',
}))

// 注册开发者：单一入口，无子菜单；仅对非开发者显示，点击前往注册页
const registerDeveloperOption = computed<MenuOption>(() => ({
  label: () => a('/developer/register', t('menu.registerDeveloper')),
  key: 'registerDeveloper',
}))

// 顶部"微应用 / 自定义代码"浏览切换，控制首页展示的内容类型
const browseValue = computed(() =>
  route.path.startsWith('/customCode') ? 'customCode' : 'microApp',
)

const menuOptions = computed<MenuOption[]>(() => {
  const options: MenuOption[] = []

  // // 微应用 / 自定义代码 浏览切换（二级菜单，父级标签显示当前选中项）
  // options.push({
  //   label: () => browseValue.value === 'customCode' ? t('menu.customCode') : t('menu.microApp'),
  //   key: 'browse',
  //   children: [
  //     { label: () => t('menu.microApp'), key: 'browse-microApp' },
  //     { label: () => t('menu.customCode'), key: 'browse-customCode' },
  //   ],
  // })

  // // 微应用 / 自定义代码 浏览切换（二级菜单，父级标签显示当前选中项）
  // options.push({
  //   label: () => t('menu.microApp'),
  //   key: 'browse-microApp',
  // })

  // 自定义代码 浏览切换（二级菜单，父级标签显示当前选中项）
  options.push({
    label: () => t('menu.customCode'),
    key: 'browse-customCode',
  })

  // 发布（二级菜单）：微应用 / 自定义代码 两个发布入口
  options.push({
    label: () => t('menu.publish'),
    key: 'publish',
    children: [
      publishMicroAppOption.value,
      publishCustomCodeOption.value,
    ],
  })

  // 注册开发者入口（仅非开发者可见，单一入口无二级菜单）
  if (!hasDeveloperPermission.value)
    options.push(registerDeveloperOption.value)

  options.push({
    label: () => aBlank(devDocLinks, t('menu.devDoc')),
    key: 'community',
  })

  return options
})

function handleSelect(key: string) {
  if (key === 'browse-microApp')
    router.push('/')
  else if (key === 'browse-customCode')
    router.push('/customCode')
}

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
    @select="handleSelect"
  />
</template>
