<script lang="ts" setup>
import { TimestampToTime } from "@/utils"

const props = withDefaults(defineProps<{ timestamp?: number; format?: number }>(), { format: -1 })
</script>

<template>
  <span>{{ TimestampToTime(props.timestamp, props.format) }}</span>
</template>

<style lang="scss" scoped></style>
