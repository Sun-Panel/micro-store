<script setup lang="ts">
import { NCard, NGrid, NGridItem, NH4 } from 'naive-ui'
import { onMounted, onUnmounted, ref } from 'vue'
import moment from 'moment'
import type { PartOption } from './echarts/stackedLine.vue'
import StackedLine from './echarts/stackedLine.vue'
import Version from './Version.vue'
import { getClientLine, getStatistics, getUserLine } from '@/api/admin/dashboard'
import { NumberFormat } from '@/components/common'

interface ClientData {
  activeClientCount: number
  clientIncreaseCountToday: number
  installCount: number
  clientOnline24: number
  clientOnline48: number
  clientOnline72: number
  userCount: number
  userToday: number
}

const lineData = ref<number[]>()
const lineDates = ref<string[]>()
const statistics = ref<ClientData>()

const stackedLineOption = ref<PartOption | null>(null)
const clientLineOption = ref<PartOption | null>(null)
const StackedLineRef = ref()
const ClientLineRef = ref()

const handleResize = () => {
  StackedLineRef.value.resize()
  ClientLineRef.value.resize()
}

onMounted(() => {
  getData()
  getLineData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

// function getDaysDateTime(days: number): string[] {
//   const dates = []
//   for (let i = 0; i < days; i++) {
//     const date = moment().subtract(i, 'days').startOf('day').format('YYYY-MM-DD HH:SS:00')
//     dates.unshift(date)
//   }
//   return dates
// }

function getDaysZeroDateTime(days: number): string[] {
  const dates = []
  for (let i = 0; i < days; i++) {
    const date = moment().subtract(i, 'days').startOf('day').format('YYYY-MM-DD 00:00:00')
    dates.unshift(date)
  }
  return dates
}

function getDaysDate(days: number): string[] {
  const dates = []
  for (let i = 0; i < days; i++) {
    const date = moment().subtract(i, 'days').startOf('day').format('YYYY-MM-DD')
    dates.unshift(date)
  }
  return dates
}

function getLineData() {
  const days = 30
  const dates = getDaysZeroDateTime(days)
  lineDates.value = getDaysDate(days)
  getUserLine<number[]>(dates).then(({ data }) => {
    lineData.value = data
    stackedLineUpdate(days)
  })

  getClientLine<number[]>(dates).then(({ data }) => {
    lineData.value = data
    clientineUpdate(days)
  })
}

function stackedLineUpdate(days: number) {
  stackedLineOption.value = {
    // title: '用户折线',
    // legendData: ['USD', 'CNY'],
    xAxisData: lineDates.value || [], // ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
    datas: [
      {
        name: '增长',
        type: 'line',
        stack: 'Total',
        data: lineData.value || [],
      },
    ],
  }
  // StackedLineRef.value.update()
}

function clientineUpdate(days: number) {
  clientLineOption.value = {
    // title: '用户折线',
    // legendData: ['USD', 'CNY'],
    xAxisData: lineDates.value || [], // ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
    datas: [
      {
        name: '增长',
        type: 'line',
        stack: 'Total',
        data: lineData.value || [],
      },
    ],
  }
  // StackedLineRef.value.update()
}

async function getData() {
  const { data } = await getStatistics<ClientData>()
  statistics.value = data
}
</script>

<template>
  <div class="max-w-[1500px] mx-[auto]">
    <NH4 prefix="bar">
      概览
    </NH4>
    <NGrid cols="1 200:1 400:2 600:3" :x-gap="20" :y-gap="20">
      <NGridItem>
        <NCard>
          <div>
            总用户 / 今日新增
          </div>
          <div class="text-[24px]">
            <NumberFormat :num="statistics?.userCount || 0" />
            /
            <NumberFormat :num="statistics?.userToday || 0" />
          </div>
        </NCard>
      </NGridItem>
      <NGridItem>
        <NCard>
          <div>
            累计安装次数 / 活跃客户端（半年） / 今日新增
          </div>
          <div class="text-[24px]">
            <NumberFormat :num="statistics?.installCount || 0" />
            /
            <NumberFormat :num="statistics?.activeClientCount || 0" />
            /
            <NumberFormat :num="statistics?.clientIncreaseCountToday || 0" />
          </div>
        </NCard>
      </NGridItem>
      <NGridItem>
        <NCard>
          <div>
            在线数 24h / 48h / 72h
          </div>
          <div class="text-[24px]">
            <NumberFormat :num="statistics?.clientOnline24 || 0" />
            /
            <NumberFormat :num="statistics?.clientOnline48 || 0" />
            /
            <NumberFormat :num="statistics?.clientOnline72 || 0" />
          </div>
        </NCard>
      </NGridItem>
    </NGrid>

    <NH4 prefix="bar">
      各个版本数据统计
    </NH4>
    <NCard>
      <Version />
    </NCard>

    <NH4 prefix="bar">
      用户增长
    </NH4>

    <NGrid cols="1" :y-gap="20" class="mt-[20px]">
      <NGridItem>
        <NCard>
          <div class="h-[300px]">
            <StackedLine ref="StackedLineRef" :part-option="stackedLineOption" />
          </div>
        </NCard>
      </NGridItem>
    </NGrid>

    <NH4 prefix="bar">
      客户端增长
    </NH4>

    <NGrid cols="1" :y-gap="20" class="mt-[20px]">
      <NGridItem>
        <NCard>
          <div class="h-[300px]">
            <StackedLine ref="ClientLineRef" :part-option="clientLineOption" />
          </div>
        </NCard>
      </NGridItem>
    </NGrid>
  </div>
</template>
