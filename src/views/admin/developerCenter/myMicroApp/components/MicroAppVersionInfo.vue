<script lang="ts" setup>
import type { DataTableColumns } from 'naive-ui'
import { NButton, NDataTable, NPopconfirm, NPopover, NSpace, NTag } from 'naive-ui'
import { computed, h } from 'vue'
import { SvgIcon } from '@/components/common'
import { MicroAppVersionStatus, microAppVersionStatusMap } from '@/enums/panel'
import { timeFormat } from '@/utils/cmn'

const props = defineProps<{
  versionList: MicroApp.VersionInfo[]
  loading?: boolean
  canAddVersion?: boolean
  canDeleteVersion?: boolean
  canSubmitReview?: boolean
  canOfflineVersion?: boolean
}>()

const emit = defineEmits<{
  (e: 'addVersion'): void
  (e: 'viewDetail', version: MicroApp.VersionInfo): void
  (e: 'submitReview', versionId: number): void
  (e: 'cancelReview', versionId: number): void
  (e: 'deleteVersion', ids: number[]): void
  (e: 'offlineVersion', version: MicroApp.VersionInfo): void
}>()

// 版本状态颜色
function getVersionStatusType(status: number): 'default' | 'success' | 'error' | 'warning' {
  if (status === MicroAppVersionStatus.APPROVED)
    return 'success'
  if (status === MicroAppVersionStatus.REJECTED)
    return 'error'
  if (status === MicroAppVersionStatus.OFFLINE)
    return 'warning'
  return 'default'
}

// 版本状态文本
function getVersionStatusText(status: number): string {
  return microAppVersionStatusMap[status] || '未知'
}

// 表格列配置
const columns = computed<DataTableColumns<MicroApp.VersionInfo>>(() => [
  {
    title: '版本号',
    key: 'version',
    width: 100,
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row) {
      const tag = h(NTag, { type: getVersionStatusType(row.status), size: 'small' }, {
        default: () => getVersionStatusText(row.status),
        ...(row.status === MicroAppVersionStatus.OFFLINE && row.offlineReason
          ? {
              icon: () => h(SvgIcon, { icon: 'lucide:info', style: 'color: #f0a020;' }),
            }
          : {}),
      })

      // 已下架状态且有下架原因时，用 NPopover 包裹整个 NTag
      if (row.status === MicroAppVersionStatus.OFFLINE && row.offlineReason) {
        return h(
          NPopover,
          { trigger: 'hover', placement: 'top' },
          {
            trigger: () => tag,
            default: () => [
              h('div', { style: 'font-weight: 500; margin-bottom: 4px;' }, '下架原因：'),
              h('div', { style: 'max-width: 280px; word-break: break-all;' }, row.offlineReason),
            ],
          },
        )
      }

      return tag
    },
  },
  {
    title: '上传时间',
    key: 'createTime',
    width: 160,
    render(row) {
      return timeFormat(String(row.createTime))
    },
  },
  {
    title: '审核时间',
    key: 'reviewTime',
    width: 160,
    render(row) {
      return row.reviewTime ? timeFormat(String(row.reviewTime)) : '-'
    },
  },
  {
    title: '审核备注',
    key: 'reviewNote',
    ellipsis: { tooltip: true },
  },
  {
    title: '操作',
    key: 'actions',
    width: 320,
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          // 查看详情
          h(NButton, { size: 'small', onClick: () => emit('viewDetail', row) }, {
            default: () => '查看',
          }),
          // 草稿/拒绝/下架状态，可以提交审核
          props.canSubmitReview && (row.status === MicroAppVersionStatus.DRAFT || row.status === MicroAppVersionStatus.REJECTED || row.status === MicroAppVersionStatus.OFFLINE)
            ? h(NButton, { size: 'small', type: 'primary', onClick: () => emit('submitReview', row.id) }, {
                default: () => '提交审核',
              })
            : null,
          // 待审核状态，可以撤销
          props.canSubmitReview && row.status === MicroAppVersionStatus.PENDING
            ? h(NButton, { size: 'small', type: 'warning', onClick: () => emit('cancelReview', row.id) }, {
                default: () => '撤销',
              })
            : null,

          // 已通过状态，可以下架
          props.canOfflineVersion && row.status === MicroAppVersionStatus.APPROVED
            ? h(NButton, { size: 'small', type: 'warning', onClick: () => emit('offlineVersion', row) }, {
                default: () => '下架',
              })
            : null,

          // 非通过且非下架状态，可以删除
          props.canDeleteVersion
            ? h(NPopconfirm, { onPositiveClick: () => emit('deleteVersion', [row.id]) }, {
                trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
                default: () => `确定删除版本 ${row.version} 吗？`,
              })
            : null,

        ],
      })
    },
  },
])
</script>

<template>
  <NCard title="版本管理">
    <template #header-extra>
      <NButton v-if="canAddVersion" type="primary" @click="emit('addVersion')">
        添加版本
      </NButton>
    </template>

    <NDataTable
      :columns="columns"
      :data="versionList"
      :loading="loading"
      :bordered="false"
    >
      <template #empty>
        <div class="text-center py-12 text-gray-400">
          {{ canAddVersion ? '暂无版本，点击"添加版本"上传第一个版本' : '暂无版本' }}
        </div>
      </template>
    </NDataTable>
  </NCard>
</template>
