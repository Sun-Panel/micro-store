<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NButton, NH4, NInput, useDialog, useMessage } from 'naive-ui'
import { OrderInfo } from '@/views/components'
import { getOrderInfo, orderManageUpdateStatusByOrderNo, updateAdminNoteByOrderNo as updateAdminNoteByOrderNoApi } from '@/api/admin/orderManage'
import { OrderStatus } from '@/enums/goodsOrder'

interface Props {
  orderNo: string
}

const props = defineProps<Props>()
const ms = useMessage()
const dialog = useDialog()
const orderInfo = ref<GoodsOrder.Info | null>(null)
const isUpdateAdminNoteEdit = ref(false)
const adminNote = ref('')

async function getinfo() {
  try {
    const { data } = await getOrderInfo<GoodsOrder.Info>(props.orderNo)
    orderInfo.value = data
    adminNote.value = data.adminNote
  }
  catch (error) {
    ms.error(`修改详情失败：${String(error)}`)
  }
}

async function updateAdminNoteByOrderNo() {
  console.log(adminNote.value)
  try {
    await updateAdminNoteByOrderNoApi(props.orderNo, adminNote.value)
    isUpdateAdminNoteEdit.value = false
  }
  catch (error) {
    ms.error(`修改备注失败：${String(error)}`)
  }
}

function handleUpdateStatus() {
  dialog.warning({
    title: '警告',
    content: '你确定修改此订单的状态吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      updateFinishStatus()
    },

  })
}

async function updateFinishStatus() {
  try {
    await orderManageUpdateStatusByOrderNo(props.orderNo, OrderStatus.FINISH)
    getinfo()
  }
  catch (error) {
    ms.error(`修改状态失败：${String(error)}`)
  }
}

onMounted(() => {
  getinfo()
})
</script>

<template>
  <div>
    <div>
      <OrderInfo :order-info="orderInfo" />
    </div>

    <div>
      <div class="flex mt-4 items-center">
        <NH4 prefix="bar" style="margin-bottom: 5px;" type="error" align-text>
          <div class="flex">
            <div>
              管理员备注
            </div>
            <div class="ml-2">
              <NButton v-if="isUpdateAdminNoteEdit" size="tiny" type="success" @click="updateAdminNoteByOrderNo">
                保存
              </NButton>
              <NButton v-else size="tiny" @click="isUpdateAdminNoteEdit = true">
                修改备注
              </NButton>
            </div>
          </div>
        </NH4>
      </div>
      <div class="flex">
        <NInput v-if=" isUpdateAdminNoteEdit" v-model:value="adminNote" />
        <div v-else>
          {{ adminNote ?? "-" }}
        </div>
      </div>
    </div>
    <div v-if="orderInfo?.status === OrderStatus.PAY_SUCCESS " class="mt-5">
      <NButton type="success" @click="handleUpdateStatus">
        修改订单为已完成的状态
      </NButton>
    </div>
  </div>
</template>
