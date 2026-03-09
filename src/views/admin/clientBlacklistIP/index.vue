<script lang="ts" setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NCard, NDataTable, NForm, NFormItem, NInput, NModal, NSelect, useDialog, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import moment from 'moment'
import { deletes as deletesApi, getAll as getAllApi, set as setApi } from '@/api/admin/clientBlacklistIP'

interface BlacklistIPCacheParams {
  LastVisitTimestamp: number // 最后访问时间戳
  ExpirationTimestamp: number // 过期时间戳
}

interface BlacklistItem {
  ip: string
  lastTime: string
  expirationTime: string
}

const addModalShow = ref(false)
const banTimeOptions = [
  {
    label: '1小时',
    value: 3600,
  },
  {
    label: '5小时',
    value: 18000,
  },
  {
    label: '1天',
    value: 86400,
  },
  {
    label: '2天',
    value: 86400 * 2,
  },
  {
    label: '3天',
    value: 86400 * 3,
  },
  {
    label: '7天',
    value: 86400 * 7,
  },
  {
    label: '30天',
    value: 2592000,
  },
  {
    label: '180天',
    value: 6048000,
  },
  {
    label: '1年',
    value: 31536000,
  },
  {
    label: '5年',
    value: 5 * 31536000,
  },
  {
    label: '10年',
    value: 10 * 31536000,
  },
]

const addModalForm = ref<{
  ips: string
  banTime: number
}>({
  ips: '',
  banTime: 86400,
})

const addModalRules = {
  ips: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
    },
  ],
  banTime: [
    {
      required: true,
      trigger: 'blur',
      message: '必填项',
      type: 'nubmer',
    },
  ],
}
const blacklist = ref<BlacklistItem[]>([])
const message = useMessage()
const dialog = useDialog()
const tableIsLoading = ref<boolean>(false)

const createColumns = (): DataTableColumns<BlacklistItem> => {
  return [
    {
      title: 'IP地址',
      key: 'ip',
    },
    {
      title: '解除时间',
      key: 'expirationTime',
    },
    {
      title: '最后拦截时间',
      key: 'lastTime',
    //   render(row) {
    //     if (row.isActive) {
    //       return h(NTag, { type: 'success' }, '是')
    //     }
    //     else {
    //       return h(
    //         NButton,
    //         {
    //           size: 'tiny',
    //           onClick() {
    //             dialog.warning({
    //               title: '警告',
    //               content: `你确定将[ v${row.version} ]设置为 ${row.type} 最新版？`,
    //               positiveText: '确定',
    //               negativeText: '取消',
    //               onPositiveClick: () => {
    //                 setLatestVersion(row.id as number, row.type)
    //               },
    //             })
    //           },
    //         },
    //         '设为最新版',
    //       )
    //     }
    //   },
    },

    {
      title: '',
      key: '',
      render(row) {
        return h(
          'div',
          {},
          [
            h(
              NButton,
              {
                size: 'tiny',
                type: 'info',
                onClick() {
                  handleEdit(row.ip)
                },
              },
              '编辑',
            ),
            h(
              NButton,
              {
                size: 'tiny',
                style: { marginLeft: '4px' },
                type: 'error',
                onClick() {
                  dialog.warning({
                    title: '警告',
                    content: `你确定要将IP [ ${row.ip} ] 移出黑名单？`,
                    positiveText: '确定',
                    negativeText: '取消',
                    onPositiveClick: () => {
                      deletesApi([row.ip]).then(() => {
                        getList()
                        message.success('已删除')
                      }).catch(() => {
                        message.success('删除失败')
                      })
                    },
                  })
                },
              },
              '删除',
            ),
          ],
        )
      },
    },
  ]
}

const columns = createColumns()

function handlePageChange() {
  getList()
}

function handleAdd() {
  addModalForm.value.ips = ''
  addModalShow.value = true
}

function handleEdit(ip: string) {
  addModalForm.value.ips = ip
  addModalShow.value = true
}

function handleSetSave() {
  const ips = splitAndTrimString(addModalForm.value.ips, '\n')
  const banTime = addModalForm.value.banTime
  setApi(ips, banTime).then(() => {
    message.success('操作成功')
    addModalShow.value = false
    getList()
  }).catch(() => {
    message.error('操作失败')
  })
}

function splitAndTrimString(inputString: string, delimiter: string): string[] {
  return inputString.split(delimiter).map(item => item.trim())
}

async function getList() {
  tableIsLoading.value = true

  await getAllApi<{ [key: string]: BlacklistIPCacheParams }>().then(({ data }) => {
    blacklist.value = []
    for (const ip in data) {
      if (Object.prototype.hasOwnProperty.call(data, ip)) {
        const element = data[ip]
        blacklist.value.push({
          ip,
          lastTime: element.LastVisitTimestamp ? timestampToFormattedString(element.LastVisitTimestamp, 'YYYY-MM-DD HH:mm') : '-',
          expirationTime: timestampToFormattedString(element.ExpirationTimestamp, 'YYYY-MM-DD HH:mm'),
        })
      }
    }
    return data
  }).catch(() => {
    return {
      list: [],
      count: 0,
    }
  })

  tableIsLoading.value = false
}

function timestampToFormattedString(timestampInSeconds: number, format?: string) {
  if (isNaN(timestampInSeconds))
    throw new Error('Invalid timestamp')

  if (!format)
    format = 'YYYY-MM-DD HH:mm:ss'

  const date = moment.unix(timestampInSeconds)
  return date.format(format)
}

onMounted(() => {
  getList()
})
</script>

<template>
  <div>
    <NCard class="mb-[20px]">
      <div class="flex">
        <!-- <NInputGroup style="max-width:700px;">
          <NInput v-model:value="keyWord" :style="{ width: '50%' }" placeholder="请输入账号" @keyup.enter="handleSelect" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup> -->

        <NButton type="primary" style="margin-left: auto;" @click="handleAdd">
          添加
        </NButton>
      </div>
    </NCard>

    <NDataTable
      :columns="columns"
      :data="blacklist"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
      size="small"

      @update:page="handlePageChange"
    />

    <NModal v-model:show="addModalShow" preset="card" style="width: 600px;">
      <NForm ref="formRef" :model="addModalForm" :rules="addModalRules">
        <NFormItem path="ips" label="IP 地址 （每行一个）">
          <NInput
            v-model:value="addModalForm.ips"
            type="textarea"
            placeholder="请输入IP地址"
          />
        </NFormItem>

        <NFormItem label="封禁时长" path="banTime">
          <NSelect v-model:value="addModalForm.banTime" :options="banTimeOptions" />
        </NFormItem>
      </NForm>

      <template #footer>
        <div class="flex justify-end">
          <NButton type="success" @click="handleSetSave">
            保存
          </NButton>
        </div>
      </template>
    </NModal>
  </div>
</template>
