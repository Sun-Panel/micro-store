<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NDropdown, useDialog, useMessage } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import Edit from './Edit/index.vue'
import { deletes as deletesRequest, getList as getListRequest } from '@/api/admin/notice'
import { SvgIcon } from '@/components/common'
import { timeFormat } from '@/utils/cmn'

const message = useMessage()
const tableIsLoading = ref<boolean>(false)
const editDialogShow = ref<boolean>(false)
// const keyWord = ref<string>()
const info = ref<Notice.NoticeInfo>()
const dialog = useDialog()

const createColumns = ({
  update,
}: {
  update: (row: Notice.NoticeInfo) => void
}): DataTableColumns<Notice.NoticeInfo> => {
  return [
    {
      title: '展示类型',
      key: 'displayType',
      render(row) {
        let msg = ''
        switch (row.displayType) {
          case 1:
            msg = '登录页面'
            break
          case 2:
            msg = '首页'
            break
          default:
            break
        }
        return msg
      },
    },
    {
      title: '平台',
      key: 'title',
    },
    {
      title: '内容',
      key: 'content',
    },
    {
      title: '允许已读',
      key: 'oneRead',
      render(row) {
        let msg = ''
        switch (row.oneRead) {
          case 0:
            msg = '否'
            break
          case 1:
            msg = '是'
            break
          default:
            break
        }
        return msg
      },
    },
    {
      title: '登录可见',
      key: 'isLogin',
      render(row) {
        let msg = ''
        switch (row.isLogin) {
          case 0:
            msg = '否'
            break
          case 1:
            msg = '是'
            break
          default:
            break
        }
        return msg
      },
    },
    {
      title: '创建/更新时间',
      key: 'createTime',
      render(row) {
        return `${timeFormat(String(row.createTime))} / ${timeFormat(String(row.updateTime))}`
      },
    },
    {
      title: '操作',
      key: '',
      render(row) {
        const btn = h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: 'small',
          },
          {
            default() {
              return h(
                SvgIcon, {
                  icon: 'mingcute:more-1-fill',
                },
              )
            },
          },
        )

        return h(NDropdown, {
          trigger: 'click',
          onSelect(key: string | number) {
            console.log(key)
            switch (key) {
              case 'update':
                update(row)
                break
              case 'delete':
                dialog.warning({
                  title: '警告',
                  content: `你确定删除${row.title}？`,
                  positiveText: '确定',
                  negativeText: '取消',
                  onPositiveClick: () => {
                    deletes([row.id as number])
                  },
                })
                break

              default:
                break
            }
          },
          options: [
            {
              label: '修改信息',
              key: 'update',
            },
            {
              label: '删除',
              key: 'delete',
            },
          ],
        }, { default: () => btn })
      },
    },
  ]
}

const list = ref<Notice.NoticeInfo[]>()

const columns = createColumns({
  update(row: Notice.NoticeInfo) {
    handleEdit(row)
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

function handleEdit(row?: Notice.NoticeInfo) {
  if (row)
    info.value = row

  editDialogShow.value = true
}

function handlePageChange(page: number) {
  getList(page)
}

// 查询
// function handleSelect() {
//   getList(null)
// }

function handelEditDone() {
  editDialogShow.value = false
  message.success('操作完成')
  getList(null)
}

async function getList(page: number | null) {
  tableIsLoading.value = true
  //   const req: AdminUserManage.GetListRequest = {
  //     page: page || pagination.page,
  //     limit: pagination.pageSize,
  //   }
  //   if (keyWord.value !== '')
  //     req.keyWord = keyWord.value

  const { data } = await getListRequest<Common.ListResponse<Notice.NoticeInfo[]>>()
  pagination.itemCount = data.count
  if (data.list)
    list.value = data.list
  tableIsLoading.value = false
}

async function deletes(ids: number[]) {
  const { code } = await deletesRequest(ids)
  if (code === 0) {
    message.success('已删除')
    getList(null)
  }
}

onMounted(() => {
  getList(null)
})
</script>

<template>
  <div>
    <NCard class="mb-[20px]">
      <div class="flex">
        <!-- <NInputGroup style="max-width:700px;">
          <NInput v-model:value="keyWord" :style="{ width: '50%' }" placeholder="请输入账号或者昵称来查询" />
          <NButton type="primary" @click="handleSelect">
            查询
          </NButton>
        </NInputGroup> -->
        <!-- <span class="flex ml-auto"> -->
        <span class="flex">
          <NButton type="primary" ghost @click="handleEdit()">
            添加
          </NButton>
        </span>
      </div>
    </NCard>

    <NDataTable
      :columns="columns"
      :data="list"

      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"

      @update:page="handlePageChange"
    />
    <Edit v-model:visible="editDialogShow" :info="info" @done="handelEditDone" />
  </div>
</template>
