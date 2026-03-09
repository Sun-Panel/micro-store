<script setup lang="ts">
import { onMounted, ref } from 'vue'
import {
  NButton, NButtonGroup, NCard, NCheckbox, NCheckboxGroup, NCollapseTransition,
  NGrid, NGridItem, NSpace, NTabPane, NTabs,
} from 'naive-ui'
import VersionPie from './echarts/VersionPie.vue'
import type { PartOption } from './echarts/stackedLine.vue'
import StackedLine from './echarts/stackedLine.vue'
import { getActiveClientVersionStatistics, getVersionHistory, getVersions } from '@/api/admin/dashboard'
import { NumberFormat } from '@/components/common'
import { timeFormat } from '@/utils/cmn'

// import VersionBar from './echarts/VersionBar.vue'

interface VersionStatistics {
  version: string
  count: number
  onlineCount48: number
}

interface HistoryClientVersionStatistics {
  dateTime: string
  onlineNum24h: { [key: string]: number }
  onlineNum48h: { [key: string]: number }
  onlineNum72h: { [key: string]: number }
  activeClientNum: { [key: string]: number }
}

const allVersions = ref<Version.Info[]>([])
const filterVersions = ref<Version.Info[]>([])
const checkboxVersions = ref<string[]>([])
const activeClientNumUpdateTime = ref('')
const versionStatistic = ref<VersionStatistics[]>([])
const filtterVersionStatistic = ref<VersionStatistics[]>([])
const showFilter = ref(false)
const hisdtoryLine = ref<PartOption>({})
const versionsHistory = ref<HistoryClientVersionStatistics[]>([])

async function getClientVersionStatistics() {
  getActiveClientVersionStatistics<HistoryClientVersionStatistics>().then(({ data }) => {
    for (const key in data.activeClientNum) {
      if (Object.prototype.hasOwnProperty.call(data.activeClientNum, key)) {
        const element = data.activeClientNum[key]

        versionStatistic.value.push({
          version: key,
          count: element,
          onlineCount48: data.onlineNum48h?.[key] ?? 0,
        })

        // checkboxVersions.value.push(key)
      }
    }
    // console.log('统计数据', versionStatistic.value)
    activeClientNumUpdateTime.value = data.dateTime
    // console.log('shijian ', data.dateTime)
    filterVersions.value = filterAllVersion(versionStatistic.value, allVersions.value)
    versionStatistic.value = sortVersionStatistics(versionStatistic.value, filterVersions.value)
    //   filtterVersionStatistic.value = versionStatistic.value
    filtterVersionStatistic.value = [...versionStatistic.value]

    for (let i = 0; i < versionStatistic.value.length; i++) {
      const element = versionStatistic.value[i]
      checkboxVersions.value.push(element.version)
    }
  })
}

function sortVersionStatistics(
  versionStatistics: VersionStatistics[],
  allVersion: Version.Info[],
): VersionStatistics[] {
  // 将 allVersion 的版本顺序保存到一个 Map 中，方便快速查找排序顺序
  const versionOrder = new Map<string, number>()
  allVersion.forEach((version, index) => {
    versionOrder.set(version.version, index)
  })

  // 对 versionStatistics 按照 allVersion 中的顺序进行排序
  return versionStatistics.sort((a, b) => {
    const indexA = versionOrder.get(a.version) ?? Number.MAX_SAFE_INTEGER
    const indexB = versionOrder.get(b.version) ?? Number.MAX_SAFE_INTEGER
    return indexA - indexB
  })
}

function filterAllVersion(
  versionStatistics: VersionStatistics[],
  allVersion: Version.Info[],
): Version.Info[] {
  // 创建一个 Set，包含所有 versionStatistics 中的版本号
  const versionSet = new Set(versionStatistics.map(stat => stat.version))

  // 过滤 allVersion 中存在但不在 versionStatistics 中的项
  return allVersion.filter(version => versionSet.has(version.version))
}

function handleFiltterRefresh() {
  filtterVersionStatistic.value = []
  for (let i = 0; i < versionStatistic.value.length; i++) {
    const element = versionStatistic.value[i]

    if (checkboxVersions.value.includes(element.version)) {
      // 删除
      filtterVersionStatistic.value.push(element)
    }
  }
  buildVersionsHistoryLine()
}

// 版本过滤按钮被点击
function handleOnFiltterButtonClick(type: string) {
  switch (type) {
    case 'all':
      if (checkboxVersions.value.length > 0) {
        // 取消全选
        // filtterVersionStatistic.value = []
        checkboxVersions.value = []
      }
      else {
        // 全选
        for (let i = 0; i < versionStatistic.value.length; i++) {
          const element = versionStatistic.value[i]
          checkboxVersions.value.push(element.version)
        }
      }
      break
    case 'release':
      checkboxVersions.value = []
      for (let i = 0; i < filterVersions.value.length; i++) {
        const element = filterVersions.value[i]
        if (element.type === 'release')
          checkboxVersions.value.push(element.version)
      }

      break
    case 'beta':
      checkboxVersions.value = []
      for (let i = 0; i < filterVersions.value.length; i++) {
        const element = filterVersions.value[i]
        if (element.type === 'beta')
          checkboxVersions.value.push(element.version)
      }
      break

    default:
      break
  }
}

function buildVersionsHistoryLine() {
  console.log('选中的版本', checkboxVersions.value)
  const updatedOption: PartOption = {
    legendData: [],
    xAxisData: [],
    datas: [],
    title: '30天-历史记录',
  }

  updatedOption.legendData = checkboxVersions.value
  updatedOption.xAxisData = []
  updatedOption.datas = []

  const versionHistoryLineData: { [key: string]: number[] } = {}

  for (let i = 0; i < versionsHistory.value.length; i++) {
    const element = versionsHistory.value[i]
    updatedOption.xAxisData.push(timeFormat(element.dateTime))

    for (const versionName in element.activeClientNum) {
      if (Object.prototype.hasOwnProperty.call(element.activeClientNum, versionName)) {
        const elementVersionNum = element.activeClientNum[versionName]
        if (!versionHistoryLineData[versionName])
          versionHistoryLineData[versionName] = []

        if (checkboxVersions.value.includes(versionName))
          versionHistoryLineData[versionName].push(elementVersionNum)
        else
          versionHistoryLineData[versionName].push(0)
      }
    }
  }

  // console.log('versionHistoryLineData', versionHistoryLineData)

  updatedOption.datas = []
  for (const key in versionHistoryLineData) {
    if (Object.prototype.hasOwnProperty.call(versionHistoryLineData, key)) {
      const element = versionHistoryLineData[key]
      if (checkboxVersions.value.includes(key)) {
        // console.log('正在插入', key, updatedOption.datas)
        updatedOption.datas.push({
          name: key,
          type: 'line',
          data: element,
        })
        // console.log('已插入', key, updatedOption.datas)
      }
    }
  }

  // 重新赋值以触发视图更新
  hisdtoryLine.value = updatedOption
}

onMounted(async () => {
  await getVersions<Common.ListResponse<Version.Info[]>>().then(({ data }) => {
    allVersions.value = data.list
  })

  await getClientVersionStatistics()
  await getVersionHistory<Common.ListResponse<HistoryClientVersionStatistics[]>>().then(({ data }) => {
    versionsHistory.value = data.list
  })
  buildVersionsHistoryLine()
})
</script>

<template>
  <div>
    <div>
      <div class="flex items-center my-2">
        <div class="mr-auto">
          <NCheckbox v-model:checked="showFilter" label="展开版本过滤" />
        </div>
        <div class="text-xs text-gray-400">
          最后更新时间：{{ timeFormat(activeClientNumUpdateTime) }}
        </div>
      </div>

      <NCollapseTransition :show="showFilter">
        <NButtonGroup size="tiny" class="my-1">
          <NButton type="info" @click="handleOnFiltterButtonClick('all')">
            全选/取消
          </NButton>
          <NButton type="success" @click="handleOnFiltterButtonClick('release')">
            仅正式版本
          </NButton>
          <NButton type="warning" @click="handleOnFiltterButtonClick('beta')">
            仅Beta版本
          </NButton>
        </NButtonGroup>
        <NFlex>
          <NCard size="small" class="my-2">
            <NCheckboxGroup v-model:value="checkboxVersions">
              <NSpace item-style="display: flex;">
                <NCheckbox
                  v-for="item, index in versionStatistic" :key="index" :value="item.version"
                  :label="`v${item.version}`"
                />
              </NSpace>
            </NCheckboxGroup>
          </NCard>

          <NButton type="info" size="small" @click="handleFiltterRefresh">
            查看筛选结果
          </NButton>
        </NFlex>
      </NCollapseTransition>
    </div>
    <NTabs type="line" animated>
      <NTabPane name="pie" tab="版本占比">
        <NGrid cols="1" :y-gap="20" class="mt-[20px]">
          <NGridItem>
            <!-- <NCard> -->
            <div class="h-[300px]">
              <VersionPie :version-client-count="filtterVersionStatistic || []" />
            </div>
            <!-- </NCard> -->
          </NGridItem>
        </NGrid>
      </NTabPane>

      <NTabPane name="card" tab="详细">
        <NGrid cols="1 200:1 400:2 600:3 800:4" :x-gap="20" :y-gap="20">
          <NGridItem v-for="item, index in filtterVersionStatistic" :key="index">
            <NCard>
              <div>
                v{{ item.version }}
              </div>
              <div>
                <span class="mr-5">
                  总数：
                  <NumberFormat :num="item?.count || 0" />
                </span>

                在线数(48h)：
                <NumberFormat :num="item?.onlineCount48 || 0" />
              </div>
            </NCard>
          </NGridItem>
        </NGrid>
      </NTabPane>

      <NTabPane name="history" tab="历史记录">
        <div class="h-[500px]">
          <StackedLine ref="ClientLineRef" :part-option="hisdtoryLine" />
        </div>
      </NTabPane>

      <!-- <NTabPane name="bar" tab="柱形图">
          <NGrid cols="1" :y-gap="20" class="mt-[20px]">
            <NGridItem>
              <NCard>
                <div class="h-[300px]">
                  <VersionBar
                    :versions="versionDatas.versions"
                    :count-datas="versionDatas.countDatas"
                    :online-datas="versionDatas.onlineDatas"
                  />
                </div>
              </NCard>
            </NGridItem>
          </NGrid>
        </NTabPane> -->
    </NTabs>
  </div>
</template>
