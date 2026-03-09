<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NButton, NCard, NForm, NFormItem, NInput, useMessage } from 'naive-ui'
import { getAboutDescription as getAboutDescriptionApi } from '@/api/openness'
import { save as saveApi } from '@/api/admin/aboutSetting'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'
import { t } from '@/locales'

const aboutContent = ref<string>('')
const ms = useMessage()

async function getAboutDescription() {
  const { data } = await getAboutDescriptionApi<string>()
  aboutContent.value = data
}

async function handleSave() {
  await saveApi(aboutContent.value)
    .then(() => {
      ms.success('更新成功')
    })
    .catch((error) => {
      apiRespErrMsgAndCustomCodeNeg1Msg(error, '更新失败')
    })
}

onMounted(() => {
  getAboutDescription()
})
</script>

<template>
  <div>
    <NCard class="max-w-[500px]">
      <NForm>
        <NFormItem label="关于页面的内容" style="margin-top: 20px;">
          <NInput v-model:value="aboutContent" type="textarea" :rows="10" placeholder="请输入关于的内容（支持html）" />
        </NFormItem>
      </NForm>
      <NButton type="success" @click="handleSave">
        {{ t('common.save') }}
      </NButton>
    </NCard>
  </div>
</template>
