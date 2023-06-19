<script setup lang="ts">
import type { Layer } from '@/types';
import { formatSize } from '@/utils';

defineProps<{
  layers?: Layer[],
}>()
</script>

<template>
  <div class="layers font-monospace">
    <div class="layer row" v-for="(layer,index) in layers" :key="index">
        <div class="col-1">{{ index }}</div>
        <div class="col text-sm-truncate">{{ layer.script }}</div>
        <div class="col-2 text-sm-end">{{ formatSize(layer.size) }}</div>
    </div>
  </div>
</template>