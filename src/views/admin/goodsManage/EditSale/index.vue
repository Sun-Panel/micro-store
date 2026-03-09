<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst } from 'naive-ui'
import { NButton, NForm, NFormItem, NInputNumber, NModal, NSwitch, useMessage } from 'naive-ui'
import { updateSale as updateSaleApi } from '@/api/admin/goodsManage'

interface Props {
  visible: boolean
  info: Admin.GoodsManage.GoodsInfo | null
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const message = useMessage()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', id: number): void// 创建完成
}

const model = ref<Admin.GoodsManage.GoodsInfo>({ })
const formRef = ref<FormInst | null>(null)

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

watch(show, (newValue, oldValue) => {
  if (props.info?.id)
    model.value = { ...props.info } || {}

  else
    message.error('出错了，该商品不存在')
})

async function update(data: Admin.GoodsManage.UpdateSaleReq) {
  try {
    const res = await updateSaleApi<Admin.GoodsManage.GoodsInfo>(data)
    if (res.code === 0)
      emit('done', 0)

    else if (res.code !== -1)
      message.warning('操作失败')
  }
  catch (error) {
    console.error(error)
    message.warning('操作失败')
  }
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    const saveInfo: Admin.GoodsManage.UpdateSaleReq = {
      status: model.value.status === 1,
      sort: model.value.sort,
      num: model.value.num,
      id: props.info?.id,
    }

    if (!errors && props.info?.id)
      update(saveInfo)
    else
      console.log(errors)
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px" title="修改销售信息">
    <NForm ref="formRef" :model="model">
      <NFormItem path="description" label="状态">
        <NSwitch
          v-model:value="model.status"
          :checked-value="1"
          :unchecked-value="2"
        >
          <template #checked>
            上架
          </template>
          <template #unchecked>
            下架
          </template>
        </NSwitch>
      </NFormItem>

      <NFormItem path="sort" label="排序">
        <NInputNumber
          v-model:value="model.sort"
        />
      </NFormItem>

      <NFormItem path="num" label="库存">
        <NInputNumber
          v-model:value="model.num"
          placeholder="数量"
          :min="0"
        />
      </NFormItem>
    </NForm>

    <template #footer>
      <NButton type="success" @click="handleValidateButtonClick">
        保存
      </NButton>
    </template>
  </NModal>
</template>
