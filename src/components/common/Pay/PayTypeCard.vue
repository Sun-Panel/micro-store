<script setup lang="ts">
import { NCard, NCheckbox } from 'naive-ui'
</script>

<template>
  <div class="w-[250px]">
    <NCard>
      <NCheckbox size="large">
        <slot />
      </NCheckbox>
    </NCard>
  </div>
</template>
