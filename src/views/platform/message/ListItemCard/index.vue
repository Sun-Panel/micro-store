<script setup lang="ts">
import { NBadge, NButton, NCard, NEllipsis, NGrid, NGridItem, NModal, NTag, useDialog } from 'naive-ui'
import { computed, ref } from 'vue'
import { buildTimeString } from '@/utils/cmn'
import { deleteByMessageId, updateReadStatus } from '@/api/system/message'
import { MessageInfo } from '@/views/components'

const props = defineProps<{
  isSend: boolean
  messageInfo: Message.MessageInfo
}>()

// 定义事件1
const emit = defineEmits<{
  (e: 'deleteDone'): void
}>()

const msgInfo = computed(() => props.messageInfo)
const messageInfoModalShow = ref(false)
const dialog = useDialog()

async function handleUpdateReadStatus(status: boolean) {
  try {
    await updateReadStatus(props.messageInfo.messageId, status)
    msgInfo.value.haveRead = status ? 1 : 0
  }
  catch (error) {

  }
}

function handleShowInfo() {
  messageInfoModalShow.value = true
  if (!props.isSend && !msgInfo.value.haveRead)
    handleUpdateReadStatus(true)
}

async function deletePost() {
  try {
    await deleteByMessageId(props.messageInfo.messageId, props.isSend)
    emit('deleteDone')
  }
  catch (error) {

  }
}

function handleDelete() {
  dialog.warning({
    title: '警告',
    content: '你确定删除？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      deletePost()
    },
  })
}
</script>

<template>
  <div>
    <NCard size="small">
      <NGrid cols="2" y-gap="10" item-responsive>
        <NGridItem span="2 800:1">
          <div class="cursor-pointer" @click="handleShowInfo">
            <NEllipsis>
              <template v-if="msgInfo.haveRead === 0 && !isSend">
                <NBadge dot />
              </template>
              {{ msgInfo.title ? msgInfo.title : msgInfo.content }}
            </NEllipsis>
          </div>
        </NGridItem>
        <NGridItem span="2 800:1">
          <div class="flex items-center ml-2">
            <span>
              <NTag size="small" round :title="isSend ? msgInfo.toUser?.username : msgInfo.fromUser?.username">
                @
                <template v-if="isSend">
                  {{ msgInfo.toUser?.name }}
                </template>
                <template v-else>
                  {{ msgInfo.fromUser?.name }}
                </template>
              </NTag>
            </span>

            <div class="ml-auto flex">
              <div class="mx-2">
                {{ buildTimeString(msgInfo.createTime, "YYYY-MM-DD HH:mm:ss") }}
              </div>

              <div class="flex items-center">
                <div>
                  <template v-if="!isSend">
                    <template v-if="msgInfo.haveRead === 0">
                      <NButton size="tiny" type="info" strong secondary @click="handleUpdateReadStatus(true)">
                        标为已读
                      </NButton>
                    </template>
                    <template v-if="msgInfo.haveRead === 1">
                      <NButton size="tiny" @click="handleUpdateReadStatus(false)">
                        标为未读
                      </NButton>
                    </template>
                  </template>
                </div>
                <div class="ml-2">
                  <NButton size="tiny" type="error" strong secondary @click="handleDelete">
                    删除
                  </NButton>
                </div>
              </div>
            </div>
          </div>
        </NGridItem>
      </NGrid>
    </NCard>

    <NModal v-model:show="messageInfoModalShow" preset="card" style="width: 800px" title="详情">
      <MessageInfo :message-info="messageInfo" />
    </NModal>
  </div>
</template>
