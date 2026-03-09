<script setup lang="ts">
import { ref } from 'vue'
import { NDivider, NH4, NH5, NImage, NImageGroup, NSpace, NTag } from 'naive-ui'
const props = defineProps<{
  messageInfo: Message.MessageInfo
}>()
const msgInfo = ref(props.messageInfo)
</script>

<template>
  <div>
    <!-- 标题 -->
    <div v-if="msgInfo.title">
      <NH5> 标题:</NH5>
      {{ msgInfo.title }}
      <NDivider dashed />
    </div>

    <!-- 发送人 -->
    <div>
      <div>
        发送人:
        <NTag size="small" round :title="msgInfo.fromUser?.username">
          @{{ msgInfo.fromUser?.name }} [{{ msgInfo.fromUser?.username }}]
        </NTag>
      </div>

      <div class="mt-2">
        接收人:
        <NTag size="small" round :title="msgInfo.toUser?.username">
          @{{ msgInfo.toUser?.name }} [{{ msgInfo.toUser?.username }}]
        </NTag>
      </div>
    </div>
    <NDivider dashed />

    <!-- 内容 -->
    <div class="mt-5 mb-10 w-full">
      <NH4>内容:</NH4>
      <div class="whitespace-pre-wrap">
        {{ msgInfo.content }}
      </div>
    </div>

    <!-- 图片 -->
    <div v-if="msgInfo.appendix.length > 0">
      <NDivider dashed />
      <NH4>附件:</NH4>
      <NImageGroup show-toolbar-tooltip>
        <NSpace>
          <div
            v-for="imgSrc, index in msgInfo.appendix"
            :key="index"
            class="border flex items-center w-[100px] h-[100px] justify-center"
          >
            <NImage
              :src="imgSrc"
            />
          </div>
        </NSpace>
      </NImageGroup>
    </div>
  </div>
</template>
