<script setup lang="ts">
import { NDivider, NH4, NTag } from 'naive-ui'
import { t } from '@/locales'

interface Props {
  goodsInfo: Goods.GoodsInfo | null
}

defineProps<Props>()
</script>

<template>
  <div>
    <div>
      <NH4 prefix="bar" type="error">
        <span class="text-slate-500">
          {{ t('order.goodsTitle') }}
        </span>
      </NH4>
      <div class=" text-xl">
        {{ goodsInfo?.title }}
        <NTag v-if="goodsInfo?.discount" size="small" :bordered="false" type="error">
          {{ goodsInfo.discount }}
        </NTag>
      </div>
    </div>
    <div>
      <NH4 prefix="bar" style="margin-top:30px" type="error">
        <span class="text-slate-500">
          {{ t('order.goodsDescription') }}
        </span>
      </NH4>
      <div v-html="goodsInfo?.description" />
    </div>
    <NDivider dashed />
    <div>
      <span>
        {{ t('common.amount') }} : (CNY)
      </span>
      <span class="text-2xl">
        {{ goodsInfo?.price }}
      </span>
      <!-- <span>
        元
      </span> -->
    </div>
  </div>
</template>
