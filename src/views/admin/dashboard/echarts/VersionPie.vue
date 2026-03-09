<script setup lang="ts">
// 示例参考：https://echarts.apache.org/examples/zh/editor.html?c=bar-stack&lang=ts
import * as echarts from 'echarts/core'
import {
  LegendComponent,
  TitleComponent,
  TooltipComponent,
} from 'echarts/components'
import { PieChart } from 'echarts/charts'
import { LabelLayout } from 'echarts/features'
import { CanvasRenderer } from 'echarts/renderers'
import { nextTick, onMounted, watchEffect } from 'vue'

interface VersionClientCount {
  version: string
  count: number
  onlineCount48: number
}

const props = defineProps<{
  versionClientCount: VersionClientCount[]
}>()

echarts.use([
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  PieChart,
  CanvasRenderer,
  LabelLayout,
])

const elemId = 'versionPie'
let myChart: echarts.ECharts | null = null

function render() {
  const chartDom = document.getElementById(elemId)
  if (!chartDom)
    return

  if (!myChart)
    myChart = echarts.init(chartDom)

  const data = []
  let allVersionCount = 0
  for (let i = 0; i < props.versionClientCount.length; i++) {
    const element = props.versionClientCount[i]

    allVersionCount += element.count
  }

  for (let i = 0; i < props.versionClientCount.length; i++) {
    const element = props.versionClientCount[i]

    data.push({
      value: element.count,
      name: element.version,
      onlineCount48: element.onlineCount48,
      percentage: `${(element.count / allVersionCount * 100).toFixed(2)}%`,
    })
  }

  const option = {
    title: {
      text: `版本统计 : ${allVersionCount}`,
      // subtext: `总客户端数据：${allVersionCount}`,
      left: 'center',
    },
    tooltip: {
      trigger: 'item',
      formatter(params: any) {
      // 自定义悬浮窗显示内容
        return `<b>版本号:</b> v${params.name} <br> <b>在用数:</b> ${params.value}  <br>  <b>占比:</b> ${params.data.percentage} <br> <b>48h在线数:</b> ${params.data.onlineCount48}`
      },
    },
    legend: {
      orient: 'vertical',
      left: 'left',
    },
    series: [
      {
        // name: 'Access From',
        type: 'pie',
        radius: '60%',
        data,
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
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
  if (props.versionClientCount) {
    nextTick(() => {
      render()
    })
  }
})
</script>

<template>
  <div :id="elemId" class="h-[100%]" />
</template>
