<script lang="ts" setup>
const props = defineProps<{ text?: string }>()
</script>

<template>
  {{ props.text || "--" }}
</template>

<style lang="scss" scoped></style>
