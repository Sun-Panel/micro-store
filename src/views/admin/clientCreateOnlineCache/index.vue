<script setup lang="ts">
import { NAlert, NButton, NCard, NForm, NFormItem, NSwitch } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { getAll, setAll } from '@/api/admin/clientCreateOnlineCache'
import { apiRespErrMsg, message } from '@/utils/cmn/apiMessage'

interface FuncNameSetting {
  ping: boolean
  autoLogin: boolean
  renewTempAuth: boolean
}

const createConfig = ref<FuncNameSetting>({
  ping: false,
  autoLogin: false,
  renewTempAuth: false,
})

function save() {
  setAll<FuncNameSetting>(createConfig.value).then((res) => {
    message.success('保存成功')
  }).catch((err) => {
    message.error('数据获取失败')

    apiRespErrMsg(err)
  })
}

onMounted(() => {
  getAll<FuncNameSetting>().then((res) => {
    createConfig.value = res.data
  }).catch(() => {
    message.error('数据获取失败')
  })
})
</script>

<template>
  <div>
    <NCard>
      <NAlert type="info">
        当访问接口时，如果没有客户端缓存（client_create_online_client_cache）是否允许创建客户端的缓存，一般在迁移服务器的时候开启半个月左右。
      </NAlert>
      <br>
      <NForm :modal="createConfig">
        <!-- 创建表单项三个NSwitch -->
        <NFormItem label="PING 接口 （Ping 此接口建议长期开启）">
          <NSwitch v-model:value="createConfig.ping" />
        </NFormItem>
        <NFormItem label="自动登录接口 （AutoLogin）">
          <NSwitch v-model:value="createConfig.autoLogin" />
        </NFormItem>
        <NFormItem label="续约短期授权接口 （RenewTempAuth）">
          <NSwitch v-model:value="createConfig.renewTempAuth" />
        </NFormItem>
      </NForm>
      <NButton type="success" @click="save">
        保存
      </NButton>
    </NCard>
  </div>
</template>
