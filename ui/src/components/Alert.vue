<script setup lang="ts">
defineProps<{
  msg: string,
  type: string
}>()
</script>

<template>
  <div class="alert alert-dismissible fade show" :class="`alert-${type}`" role="alert">
    {{ msg }}
    <!-- <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button> -->
  </div>
</template>