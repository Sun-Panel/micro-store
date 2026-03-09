<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { NButton, NModal } from 'naive-ui'
import { SModal } from '@/components/common'
// import startDrag from '@/utils/Drag'

const showModal = ref(false)
// let startX: number, startY: number, startTop: number, startLeft: number

// const onMouseDown = (e: MouseEvent) => {
//   startX = e.clientX
//   startY = e.clientY
//   startTop = modalTop.value
//   startLeft = modalLeft.value

//   document.addEventListener('mousemove', onMouseMove)
//   document.addEventListener('mouseup', onMouseUp)
// }

// const onMouseMove = (e: MouseEvent) => {
//   const dx = e.clientX - startX
//   const dy = e.clientY - startY
//   modalTop.value = startTop + dy
//   modalLeft.value = startLeft + dx
// }

// const onMouseUp = () => {
//   document.removeEventListener('mousemove', onMouseMove)
//   document.removeEventListener('mouseup', onMouseUp)
// }

// onBeforeUnmount(() => {
//   document.removeEventListener('mousemove', onMouseMove)
//   document.removeEventListener('mouseup', onMouseUp)
// })

function openW() {
  showModal.value = true
  // nextTick(() => {
  //   const oBox = document.getElementById('move-66666')
  //   const oBar = document.getElementById('basic-modal-bar')
  //   // startDrag(oBar as HTMLElement, oBox as HTMLElement)
  // })
}

onMounted(() => {

})
</script>

<template>
  <div>
    <NButton @click="openW">
      Open Modal
    </NButton>
    <SModal v-model:show="showModal" style="width: 500px;border-radius: 1rem;" move>
      <template #header>
        卡萨丁发斯蒂芬
      </template>
      66666666666
    </SModal>
    <NModal
      id="move-66666"
      preset="card"
      title="Draggable Modal"
      :closable="false"
      style="width: 600px"
    >
      <template #header>
        <div id="basic-modal-bar" class="w-full cursor-move">
          我是标题
        </div>
      </template>
      <div
        ref="modal"
      >
        <NButton
          class="close-button"
          size="small"
          @click="showModal = false"
        >
          Close
        </NButton>
        <p>This is a draggable modal.</p>
      </div>
    </NModal>
  </div>
</template>

<style>
.modal-content {
  cursor: move;
  background: white;
  padding: 16px;
  border-radius: 8px;
  user-select: none; /* 防止文字被选中 */
}
</style>
