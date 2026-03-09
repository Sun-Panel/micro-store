<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, useMessage } from 'naive-ui'
import { add as addApi, update as updateApi } from '@/api/admin/goodsManage'

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

const rules: FormRules = {
  title: [
    {
      required: true,
      trigger: 'blur',
      message: '必须大于2个字符',
      min: 2,
    },
  ],
  description: {
    required: true,
    trigger: 'blur',
    message: '必填项',
  },
  param: {
    required: true,
    trigger: 'blur',
    message: '必填项',
  },

}

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
    model.value = {}
})

async function add(goods: Admin.GoodsManage.GoodsInfo) {
  try {
    const res = await addApi<Admin.GoodsManage.GoodsInfo>(goods)
    if (res.code === 0)
      emit('done', res.data?.id as number)

    else if (res.code !== -1)
      message.warning('操作失败')
  }
  catch (error) {
    message.warning('操作失败')
  }
}

async function update(goods: Admin.GoodsManage.GoodsInfo) {
  try {
    const res = await updateApi<Admin.GoodsManage.GoodsInfo>(goods)
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

function isJSON(str: string) {
  try {
    JSON.parse(str)
    return true
  }
  catch (e) {
    console.log(e)
    return false
  }
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    const saveInfo = { ...model.value }
    if (isJSON(model.value?.param as string)) {
      saveInfo.param = JSON.parse(model.value?.param as string)
    }
    else {
      message.error('Param 必须为json格式')
      return
    }

    if (!errors) {
      if (props.info?.id)
        update(saveInfo as Admin.GoodsManage.GoodsInfo)

      else
        add(saveInfo as Admin.GoodsManage.GoodsInfo)
    }

    else { console.log(errors) }
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 600px" :title="`${info?.id ? '编辑' : '添加'}`">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="title" label="标题">
        <NInput v-model:value="model.title" type="text" placeholder="商品标题" />
      </NFormItem>

      <NFormItem path="description" label="商品描述">
        <NInput
          v-model:value="model.description"
          type="textarea"
          placeholder=""
        />
      </NFormItem>

      <NFormItem path="discount" label="活动标签内容">
        <NInput
          v-model:value="model.discount"
          type="text"
          placeholder="某某节优惠"
        />
      </NFormItem>

      <NFormItem path="price" label="价格">
        <NInputNumber
          v-model:value="model.price"
          placeholder="现价"
          :min="0"
        />
        <div class="mx-2">
          原价
        </div>
        <NInputNumber
          v-model:value="model.originalPrice"
          placeholder="商品原价"
          :min="0"
        />
      </NFormItem>

      <NFormItem path="param" label="参数(json格式)">
        <div class="w-full">
          <NInput
            v-model:value="model.param"
            type="textarea"
            placeholder=""
          />
          <!-- <div class="mt-2">
            <NButton type="info" size="tiny">
              验证JSON格式是否正确
            </NButton>
          </div> -->
        </div>
      </NFormItem>
    </NForm>

    <template #footer>
      <NButton type="success" @click="handleValidateButtonClick">
        保存
      </NButton>
    </template>
  </NModal>
</template>
