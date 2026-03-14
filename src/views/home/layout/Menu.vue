<script setup lang="ts">
import type { MenuOption } from 'naive-ui'
import { NMenu } from 'naive-ui'
import { h, ref } from 'vue'
import { SvgIconOnline } from '@/components/common'
import { t } from '@/locales'

withDefaults(defineProps<{
  isVertical: boolean
}>(), {
  isVertical: false,
})
const devDocLinks = 'https://doc.sun-panel.top/v2/zh_cn/micro_app_dev/'

const activeKey = ref('aaa')
const menuOptions: MenuOption[] = [

  {
    label: () => a('/developer/register', t('menu.becomeDeveloper')),
    key: 'becomeDeveloper',
  },

  {
    label: () => aBlank(devDocLinks, t('menu.devDoc')),
    key: 'community',
    // children: [
    //   {
    //     label: () => aBlank('https://github.com/hslr-s/sun-panel/discussions', 'Github社区'),
    //     key: 'community',
    //   },
    //   {
    //     label: () => aBlank(qqGroupLink, 'QQ 群聊'),
    //     key: 'narrator',
    //   },
    //   {
    //     label: () => aBlank('https://t.me/+bwOFXt6zXf43Njk1', 'TG 群聊'),
    //     key: 'sheep-man',
    //   },
    // ],
  },
]

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
