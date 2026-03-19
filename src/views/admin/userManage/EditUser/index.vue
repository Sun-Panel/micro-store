<script setup lang="ts">
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { NButton, NCheckbox, NCheckboxGroup, NDivider, NForm, NFormItem, NInput, NModal, NSelect, NSpace, useMessage } from 'naive-ui'
import { AdminUserManageEdit, AdminUserManageUpdatePassword } from '@/api/admin'
import { apiRespErrMsgAndCustomCodeNeg1Msg } from '@/utils/cmn/apiMessage'
import { t } from '@/locales'
import { ALL_ROLES, getSelectedRoles, calculateRoleFromSelected, ROLE_USER } from '@/utils/role'

interface Props {
  visible: boolean
  userId?: number
  userInfo?: User.Info
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const message = useMessage()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', id: number): void// 创建完成
}

const formInitValue = {
  name: '',
  username: '',
  role: ROLE_USER,
  status: 3,
}

const model = ref<User.Info>({ ...formInitValue })
const formRef = ref<FormInst | null>(null)

// 选中的角色列表（用于复选框组）
const selectedRoles = ref<number[]>([ROLE_USER])

// 密码修改
const newPassword = ref<string>('')

interface Options {
  value?: number
  label?: string
  description?: string
}

const statusOptions = ref<Options[]>([
  {
    label: '启用',
    value: 1,
  },
  {
    label: '停用',
    value: 2,
  },
  {
    label: '未激活',
    value: 3,
  },
])

const rules: FormRules = {
  username: [
    {
      required: true,
      trigger: 'blur',
      message: '请输入账号且大于5个字符',
      min: 5,
    },
    {
      trigger: 'blur',
      message: '请输入邮箱作为账号',
      type: 'email',
    },
  ],
  status: {
    required: true,
    trigger: 'blur',
    type: 'number',
    message: '请选择账号状态',
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
  if (props.userInfo?.id) {
    model.value = { ...props.userInfo }
    // 从角色值解析出选中的角色列表
    selectedRoles.value = getSelectedRoles(props.userInfo.role || ROLE_USER)
  } else {
    model.value = { ...formInitValue }
    selectedRoles.value = [ROLE_USER]
  }
  // 重置密码
  newPassword.value = ''
})

const add = async () => {
  // 计算最终角色值
  model.value.role = calculateRoleFromSelected(selectedRoles.value)
  
  await AdminUserManageEdit<User.Info>(model.value).then(async (res) => {
    // 如果设置了新密码，则修改密码
    if (newPassword.value && model.value.id) {
      try {
        await AdminUserManageUpdatePassword({
          userId: model.value.id!,
          newPassword: newPassword.value,
        })
        message.success('用户信息和密码已更新')
      } catch (error) {
        message.warning('用户信息已更新，但密码修改失败')
      }
    } else {
      message.success('操作成功')
    }
    emit('done', res.data.id as number)
  }).catch((error) => {
    apiRespErrMsgAndCustomCodeNeg1Msg(error, t('common.failed'))
  })
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      add()
    else
      console.log(errors)
  })
}
</script>

<template>
  <NModal v-model:show="show" preset="card" style="width: 650px" :title="`${userInfo?.id ? '编辑' : '添加'}用户`">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="username" label="账号" style="margin-top: 20px;">
        <NInput v-model:value="model.username" type="text" placeholder="邮箱地址作为账号" />
      </NFormItem>

      <NFormItem path="name" label="昵称" style="margin-top: 20px;">
        <NInput v-model:value="model.name" type="text" placeholder="请输入昵称" />
      </NFormItem>

      <NFormItem path="role" label="角色权限">
        <NCheckboxGroup v-model:value="selectedRoles">
          <NSpace>
            <NCheckbox v-for="role in ALL_ROLES" :key="role.value" :value="role.value">
              <span>{{ role.name }}</span>
              <span class="text-gray-400 text-xs ml-1">({{ role.description }})</span>
            </NCheckbox>
          </NSpace>
        </NCheckboxGroup>
      </NFormItem>

      <NFormItem path="status" label="状态">
        <NSelect
          v-model:value="model.status"
          :options="statusOptions"
        />
      </NFormItem>

      <!-- 密码修改（仅编辑模式） -->
      <template v-if="userInfo?.id">
        <NDivider style="margin: 12px 0;">
          密码修改
        </NDivider>
        <NFormItem path="newPassword" label="新密码">
          <NInput
            v-model:value="newPassword"
            type="password"
            placeholder="留空则不修改密码"
            show-password-on="click"
          />
        </NFormItem>
        <p class="text-gray-400 text-xs -mt-2 mb-4">
          设置新密码后将自动注销用户当前登录状态
        </p>
      </template>
    </NForm>

    <template #footer>
      <NButton type="success" @click="handleValidateButtonClick">
        保存
      </NButton>
    </template>
  </NModal>
</template>
