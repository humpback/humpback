<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    notes?: string
    size?: "large" | "default" | "small"
    type?: "primary" | "success" | "warning" | "danger" | "info" | ""
  }>(),
  { size: "default", type: "" }
)
</script>

<template>
  <el-text :size="props.size" :type="props.type">
    {{ props.notes || "--" }}
  </el-text>
</template>

<style lang="scss" scoped></style>
