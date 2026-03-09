<script setup lang="ts">
// 示例参考：https://echarts.apache.org/examples/zh/editor.html?c=line-stack&lang=ts
import * as echarts from 'echarts/core'

import { nextTick, onMounted, watch } from 'vue'
import type {
  GridComponentOption,
  LegendComponentOption,
  TitleComponentOption,
  ToolboxComponentOption,
  TooltipComponentOption,
} from 'echarts/components'
import {
  GridComponent,
  LegendComponent,
  TitleComponent,
  ToolboxComponent,
  TooltipComponent,
} from 'echarts/components'
import type { LineSeriesOption } from 'echarts/charts'
import { LineChart } from 'echarts/charts'
import { UniversalTransition } from 'echarts/features'
import { CanvasRenderer } from 'echarts/renderers'
import { randomCode } from '@/utils/cmn'

export interface PartOption {
  title?: string
  legendData?: string[] // 分类 ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  xAxisData?: string[] // X坐标标题
  datas?: LineSeriesOption[] | undefined
}

const props = defineProps<{
  partOption: PartOption | null
}>()

echarts.use([
  TitleComponent,
  ToolboxComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  LineChart,
  CanvasRenderer,
  UniversalTransition,
])

type EChartsOption = echarts.ComposeOption<
  | TitleComponentOption
  | ToolboxComponentOption
  | TooltipComponentOption
  | GridComponentOption
  | LegendComponentOption
  | LineSeriesOption
>

const elemId = `echartsStackedLine_${randomCode(10)}`
let myChart: any
// const data = computed(props.partOption)

function update() {
  if (!props.partOption) {
    console.error('PartOption 数据不完整，无法更新图表。')
    return
  }

  const option: EChartsOption = {
    title: {
      text: props.partOption.title || '',
    },
    tooltip: {
      trigger: 'axis',
    },
    legend: {
      data: props.partOption.legendData || [],
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    toolbox: {
      feature: {
        saveAsImage: {},
      },
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: props.partOption.xAxisData || [],
    },
    yAxis: {
      type: 'value',
    },
    series: props.partOption.datas || [],
  }

  option && myChart.setOption(option, true)
  myChart.resize()
}

watch(
  () => props.partOption,
  () => {
    console.log(props.partOption)
    nextTick(() => {
      update()
    })
  },
  { deep: true }, // 深度监听
)

onMounted(() => {
  const chartDom = document.getElementById(elemId)
  myChart = echarts.init(chartDom)
  nextTick(() => {
    update()
  })
})

defineExpose({
  update() {
    update()
  },
  resize() {
    myChart.resize()
  },
})
</script>

<template>
  <div :id="elemId" class="h-[100%]" />
</template>
