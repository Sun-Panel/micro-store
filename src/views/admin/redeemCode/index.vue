<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import {
  NButton, NCard, NDataTable, NDatePicker, NForm, NFormItem, NInput, NInputGroup,
  NInputNumber, NModal, NSelect, NTable, useDialog, useMessage,
} from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'

import moment from 'moment'
import { create as createRedeemCodeApi, getList as getListApi, setInvalid } from '@/api/admin/redeemCode'
import { buildTimeString } from '@/utils/cmn'

interface Options {
  value?: number
  label?: string
  description?: string
}

const message = useMessage()
const dialog = useDialog()

const defaultModel: RedeemCode.CreateReq = {
  redeemType: 4,
  title: '',
  note: '',
  extendData: {
    days: 1,
  },
  prefix: 'QT',
  number: 1,
  expireTime: '',
}
const redeemTypeOptions = ref<Options[]>([
  {
    label: '微信支付',
    value: 1,
  },
  {
    label: '支付宝支付',
    value: 2,
  },
  {
    label: '贡献回馈',
    value: 3,
  },
  {
    label: '其他',
    value: 4,
  },
])

const tableIsLoading = ref<boolean>(false)
const updateUserAuthorizeDialogShow = ref<boolean>(false)
const createRedeemCodeModalShow = ref<boolean>(false)
const createRedeemCodeSuccessListModalShow = ref(false)
const keyWord = ref<string>()
const editUserUserInfo = ref<RedeemCode.Info>()
const model = ref<RedeemCode.CreateReq>({ ...defaultModel })
const expireTimeStamp = ref(Date.now() + (24 * 60 * 60 * 1000 * 30))
const createRedeemCodeSuccessList = ref<RedeemCode.Info[]>([])

const createColumns = ({
  update,
}: {
  update: (row: RedeemCode.Info) => void
}): DataTableColumns<RedeemCode.Info> => {
  return [
    {
      title: '标题',
      key: 'title',
    },
    {
      title: '兑换码',
      key: 'code',
    },
    {
      title: '类型',
      key: 'redeemType',
      render(row) {
        let text = ''
        switch (row.redeemType) {
          case 1:
            text = '微信支付'
            break
          case 2:
            text = '支付宝支付'
            break
          case 3:
            text = '贡献回馈'
            break
          case 4:
            text = '其他'
            break
          default:
            break
        }

        return text
      },
    },
    {
      title: '有效期',
      key: 'expireTime',
      render(row) {
        let createTime = '-'
        let expireTime = '-'

        if (row.createTime)
          createTime = buildTimeString(row.createTime, 'YYYY-MM-DD HH:mm') as ''

        if (row.expireTime)
          expireTime = buildTimeString(row.expireTime, 'YYYY-MM-DD HH:mm') as ''
        return h('div', { style: { fontSize: '13px' } }, [
          h('div', `创建：${createTime}`),
          h('div', `过期：${expireTime}`),
        ])
      },
    },
    {
      title: '备注',
      key: 'note',
    },
    {
      title: '授权天数',
      key: 'days',
      render(row) {
        if (row.extendData)
          return row.extendData.days
        return '-'
      },
    },
    {
      title: '作废/核销信息',
      key: '',
      render(row) {
        if (row.writeOffTime) {
          const html = []
          const wt = buildTimeString(row.writeOffTime, 'YYYY-MM-DD HH:mm')
          if (row.userInfo)
            html.push(h('div', `${row.userInfo.name}-[${row.userInfo.username}]`))

          html.push(h('div', `${wt}`))

          return h('div', { style: { fontSize: '13px' } }, html)
        }

        if (row.status === 4)
          return '已作废'

        if (row.status === 3)
          return '已过期'

        return h(
          NButton,
          {
            size: 'tiny',
            type: 'error',
            onClick() {
              dialog.warning({
                title: '警告',
                content: '你确定作废此兑换码？',
                positiveText: '确定',
                negativeText: '取消',
                onPositiveClick: () => {
                  setInvalid([row.code]).then(() => {
                    message.success('已删除')
                    getList(null)
                  }).catch(() => {
                    message.success('删除失败')
                    return false
                  })
                },
              })
            },
          },
          '作废',
        )
      },
    },
  ]
}

const userList = ref<RedeemCode.Info[]>()

const columns = createColumns({
  update(row: RedeemCode.Info) {
    editUserUserInfo.value = row
    updateUserAuthorizeDialogShow.value = true
  },
})
const pagination = reactive({
  page: 1,
  showSizePicker: true,
  pageSizes: [10, 30, 50, 100, 200],
  pageSize: 10,
  itemCount: 0,
  onChange: (page: number) => {
    pagination.page = page
    getList(null)
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    getList(null)
  },
  prefix(item: PaginationProps) {
    return `共 ${item.itemCount} 条记录`
  },
})

function handlePageChange(page: number) {
  getList(page)
}

// 查询
function handleSelect() {
  getList(null)
}

// function handelDone() {
//   updateUserAuthorizeDialogShow.value = false
//   message.success('操作成功')
//   getList(null)
// }

function handleAdd() {
  createRedeemCodeModalShow.value = true
}

function handleUploadRedeemType(v: number) {
  switch (v) {
    case 1:
      model.value.prefix = 'WX'
      break
    case 2:
      model.value.prefix = 'AL'
      break
    case 3:
      model.value.prefix = 'HK'
      break
    case 4:
      model.value.prefix = 'QT'
      break
    default:
      break
  }
}

async function getList(page: number | null) {
  tableIsLoading.value = true
  const req: Admin.RedeemCode.GetListRequest = {
    page: page || pagination.page,
    limit: pagination.pageSize,
  }
  if (keyWord.value !== '')
    req.keyWord = keyWord.value

  const { data } = await getListApi<Common.ListResponse<RedeemCode.Info[]>>(req)
  pagination.itemCount = data.count
  if (data.list)
    userList.value = data.list
  tableIsLoading.value = false
}

function handleStartBuildRedeemCode() {
  if (model.value.expireTime === '')
    model.value.expireTime = timestampToString(expireTimeStamp.value)

  createRedeemCodeApi<Common.ListResponse<RedeemCode.Info[]>>(model.value).then(({ data }) => {
    console.log(data)
    message.success('创建成功')
    createRedeemCodeModalShow.value = false

    createRedeemCodeSuccessListModalShow.value = true
    createRedeemCodeSuccessList.value = data.list
    getList(null)
  }).catch((res) => {
    message.error(`创建失败:${res.msg}`)
  })
}

function handleUpdateExpireTime(value: number) {
  console.log(value)
  if (value)
    model.value.expireTime = timestampToString(value)
  else
    model.value.expireTime = ''
}

function timestampToString(timestamp: number) {
  return moment(timestamp).format('YYYY-MM-DD HH:mm:ss')
}

// 将格式为 YYYY-MM-DD HH:mm:ss 的字符串转换为毫秒级时间戳
// function stringToTimestamp(datetimeString: string) {
//   return moment(datetimeString, 'YYYY-MM-DD HH:mm:ss').valueOf()
// }

onMounted(() => {
  getList(null)
})
</script>

<template>
  <div>
    <NCard class="mb-[20px]">
      <div class="flex">
        <NInputGroup style="max-width:700px;">
          <NInput v-model:value="keyWord" :style="{ width: '50%' }" placeholder="请输入兑换码" @keyup.enter="handleSelect" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup>

        <span class="flex ml-auto">
          <NButton type="primary" ghost @click="handleAdd">
            创建PRO兑换码
          </NButton>
        </span>
      </div>
    </NCard>

    <NDataTable
      :columns="columns" :data="userList" :pagination="pagination" :bordered="false" :loading="tableIsLoading"
      :remote="true" @update:page="handlePageChange"
    />

    <NModal v-model:show="createRedeemCodeSuccessListModalShow" preset="card" style="max-width: 500px" title="已创建成功的兑换码">
      <NTable :bordered="false" :single-line="false">
        <thead>
          <tr>
            <th>兑换码</th>
            <th>过期时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item, index in createRedeemCodeSuccessList" :key="index">
            <td>{{ item.code }}</td>
            <td>{{ item.expireTime }}</td>
          </tr>
        </tbody>
      </NTable>
    </NModal>

    <NModal v-model:show="createRedeemCodeModalShow" preset="card" style="max-width: 500px" title="创建PRO兑换码">
      <NForm ref="formRef" :model="model" size="small">
        <NFormItem path="number" label="兑换码数量">
          <NInputNumber v-model:value="model.number" :min="1" placeholder="数量" />
        </NFormItem>

        <NFormItem path="redeemType" label="兑换码类型">
          <NSelect
            v-model:value="model.redeemType" :options="redeemTypeOptions"
            @update:value="handleUploadRedeemType"
          />
        </NFormItem>

        <NFormItem path="extendData.days" label="授权天数">
          <NInputNumber v-model:value="model.extendData.days" type="text" placeholder="授权天数" />
        </NFormItem>
        <NFormItem path="extendData.days" label="过期时间">
          <NDatePicker
            v-model:value="expireTimeStamp" type="datetime" format="yyyy-MM-dd HH:mm:ss" placeholder="选择日期时间"
            @update:value="handleUpdateExpireTime"
          />
        </NFormItem>

        <NFormItem path="prefix" label="前缀 (全拼大写例如：WX （微信）尽量保持在两个字母)">
          <NInput v-model:value="model.prefix" type="text" placeholder="前缀" />
        </NFormItem>

        <NFormItem path="title" label="标题">
          <NInput v-model:value="model.title" type="text" placeholder="标题" />
        </NFormItem>

        <NFormItem path="note" label="备注" style="margin-top: 20px;">
          <NInput v-model:value="model.note" type="text" placeholder="备注信息，兑换用户不可见" />
        </NFormItem>
      </NForm>

      <template #footer>
        <NButton type="success" @click="handleStartBuildRedeemCode">
          开始生成
        </NButton>
      </template>
    </NModal>
  </div>
</template>
