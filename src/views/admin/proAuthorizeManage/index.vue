<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NInput, NInputGroup, NModal, useMessage } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import UpdateUserAuthorize from './UpdateUserAuthorize/index.vue'
import AuthorizeHistory from './AuthorizeHistory/index.vue'
import { getUserProAuthorizeList } from '@/api/admin/proAuthorize'
import { buildTimeString } from '@/utils/cmn'

const message = useMessage()
const tableIsLoading = ref<boolean>(false)
const updateUserAuthorizeDialogShow = ref<boolean>(false)
const historyDialogShow = ref<boolean>(false)
const keyWord = ref<string>()
const editUserUserInfo = ref<Admin.ProAuthorize.ProAuthorizeGetUserProAuthorizeListItemResp>()
const updateUserId = ref(0)

const createColumns = ({
  update,
}: {
  update: (row: Admin.ProAuthorize.ProAuthorizeGetUserProAuthorizeListItemResp) => void
}): DataTableColumns<Admin.ProAuthorize.ProAuthorizeGetUserProAuthorizeListItemResp> => {
  return [
    {
      title: '账户',
      key: 'username',
    },
    {
      title: '昵称',
      key: 'name',
    },
    {
      title: '过期时间',
      key: 'expiredTime',
      render(row) {
        if (row.expiredTime)
          return buildTimeString(row.expiredTime, 'YYYY-MM-DD HH:mm')

        else
          return '-'
      },
    },
    // {
    //   title: '加入时间',
    //   key: 'createTime',
    //   render(row) {
    //     return timeFormat(String(row.createTime))
    //   },
    // },
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
                  updateUserId.value = row.userId
                  updateUserAuthorizeDialogShow.value = true
                },
              },
              '增减授权',
            ),
            h(
              NButton,
              {
                size: 'tiny',
                style: { marginLeft: '4px' },
                onClick() {
                  updateUserId.value = row.userId
                  historyDialogShow.value = true
                  console.log('点击', updateUserId.value)
                },
              },
              '历史',
            ),
          ],
        )
      },
    },
    // {
    //   title: '操作',
    //   key: '',
    //   render(row) {
    //     const btn = h(
    //       NButton,
    //       {
    //         strong: true,
    //         tertiary: true,
    //         size: 'small',
    //       },
    //       {
    //         default() {
    //           return h(
    //             SvgIcon, {
    //               icon: 'mingcute:more-1-fill',
    //             },
    //           )
    //         },
    //       },
    //     )

    //     return h(NDropdown, {
    //       trigger: 'click',
    //       onSelect(key: string | number) {
    //         console.log(key)
    //         switch (key) {
    //           case 'update':
    //             update(row)
    //             break

    //           default:
    //             break
    //         }
    //       },
    //       options: [
    //         {
    //           label: '修改信息',
    //           key: 'update',
    //         },
    //       ],
    //     }, { default: () => btn })
    //   },
    // },
  ]
}

const userList = ref<Admin.ProAuthorize.ProAuthorizeGetUserProAuthorizeListItemResp[]>()

const columns = createColumns({
  update(row: Admin.ProAuthorize.ProAuthorizeGetUserProAuthorizeListItemResp) {
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
    return `共 ${item.itemCount} 位用户`
  },
})

function handlePageChange(page: number) {
  getList(page)
}

// 查询
function handleSelect() {
  getList(null)
}

function handelDone() {
  updateUserAuthorizeDialogShow.value = false
  message.success('操作成功')
  getList(null)
}

async function getList(page: number | null) {
  tableIsLoading.value = true
  const req: AdminUserManage.GetListRequest = {
    page: page || pagination.page,
    limit: pagination.pageSize,
  }
  if (keyWord.value !== '')
    req.keyWord = keyWord.value

  const { data } = await getUserProAuthorizeList<Common.ListResponse<Admin.ProAuthorize.ProAuthorizeGetUserProAuthorizeListItemResp[]>>(req)
  pagination.itemCount = data.count
  if (data.list)
    userList.value = data.list
  tableIsLoading.value = false
}

onMounted(() => {
  getList(null)
})
</script>

<template>
  <div>
    <NCard class="mb-[20px]">
      <div class="flex">
        <NInputGroup style="max-width:700px;">
          <NInput v-model:value="keyWord" :style="{ width: '50%' }" placeholder="请输入账号" @keyup.enter="handleSelect" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup>
      </div>
    </NCard>

    <NDataTable
      :columns="columns"
      :data="userList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"

      @update:page="handlePageChange"
    />
    <UpdateUserAuthorize v-model:visible="updateUserAuthorizeDialogShow" :user-id="updateUserId" @done="handelDone" />

    <NModal v-model:show="historyDialogShow" preset="card" style="max-width: 1200px" title="授权历史记录">
      <AuthorizeHistory :user-id="updateUserId" />
    </NModal>
  </div>
</template>
