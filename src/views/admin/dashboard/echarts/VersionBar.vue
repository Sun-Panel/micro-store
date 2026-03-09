<script setup lang="ts">
// 示例参考：https://echarts.apache.org/examples/zh/editor.html?c=bar-stack&lang=ts
import * as echarts from 'echarts/core'
import {
  GridComponent,
  LegendComponent,
  TooltipComponent,
} from 'echarts/components'
import { BarChart } from 'echarts/charts'
import { CanvasRenderer } from 'echarts/renderers'
import { nextTick, onMounted, watchEffect } from 'vue'

const props = defineProps<{
  versions: string[]
  onlineDatas: number[]
  countDatas: number[]
}>()

echarts.use([
  TooltipComponent,
  GridComponent,
  LegendComponent,
  BarChart,
  CanvasRenderer,
])

const elemId = 'versionBar'
let myChart: echarts.ECharts | null = null

function render() {
  const chartDom = document.getElementById(elemId)
  if (!chartDom)
    return

  if (!myChart)
    myChart = echarts.init(chartDom)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
    },
    legend: {},
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: [
      {
        type: 'category',
        data: props.versions,
        axisLabel: {
          // rotate: 45, // 旋转标签
          fontSize: 12, // 缩小字体大小
        },
      },
    ],
    yAxis: [
      {
        type: 'value',
      },
    ],
    series: [
      {
        name: '总数',
        type: 'bar',
        emphasis: {
          focus: 'series',
        },
        data: props.countDatas,
      },
      {
        name: '48h在线',
        type: 'bar',
        stack: 'Ad',
        emphasis: {
          focus: 'series',
        },
        data: props.onlineDatas,
      },
    ],
  }

  if (myChart) {
    myChart.setOption(option)
    myChart.resize()
  }
}

onMounted(() => {
  nextTick(() => {
    render()
  })
})

watchEffect(() => {
  if (props.versions.length && props.onlineDatas.length && props.countDatas.length) {
    nextTick(() => {
      render()
    })
  }
})
</script>

<template>
  <div :id="elemId" class="h-[100%]" />
</template>
