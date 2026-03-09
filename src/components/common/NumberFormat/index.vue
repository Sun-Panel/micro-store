<script setup lang="ts">
import { NTooltip } from 'naive-ui'

withDefaults(defineProps<{
  num: number
  isAmount?: boolean
}>(), {
  isAmount: false,
})

function format(Num: number, isAmount: boolean): string {
  // 定义单位后缀
  const suffixes = ['', 'k', 'M', 'B', 'T']

  // 寻找合适的单位后缀
  let suffixIndex = 0
  while (Num >= 1000 && suffixIndex < suffixes.length - 1) {
    Num /= 1000
    suffixIndex++
  }

  // 格式化金额
  let formattedNum = Num.toFixed(0)

  if (isAmount) {
    if (suffixIndex === 0)
      formattedNum = Num.toFixed(2) // 如果没有单位后缀，则保留两位小数
    else
      formattedNum = Num.toFixed(1) // 使用一个小数点来表示k、M、B、T等单位
  }
  else {
    if (suffixIndex === 0)
      formattedNum = Num.toFixed(0)
    else
      formattedNum = Num.toFixed(1) // 使用一个小数点来表示k、M、B、T等单位
  }

  // 添加单位后缀
  formattedNum += suffixes[suffixIndex]

  return formattedNum
}
</script>

<template>
  <NTooltip placement="bottom" trigger="hover">
    <template #trigger>
      {{ format(num, isAmount) }}
    </template>
    <span> {{ num }} </span>
  </NTooltip>
</template>
