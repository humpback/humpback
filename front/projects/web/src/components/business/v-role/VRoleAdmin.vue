<script lang="ts" setup>
import { IsAdmin, IsSuperAdmin } from "@/utils"

const props = withDefaults(defineProps<{ role: number; size?: "large" | "default" | "small" }>(), { size: "small" })
const { t } = useI18n()
</script>

<template>
  <el-tag v-if="IsSuperAdmin(props.role)" :size="props.size" class="tag-super-admin" round>{{ t("role.superAdmin") }}</el-tag>
  <el-tag v-if="IsAdmin(props.role)" :size="props.size" class="tag-admin" round>{{ t("role.admin") }}</el-tag>
</template>

<style lang="scss" scoped>
.tag-admin {
  background-image: linear-gradient(to top right, var(--el-color-warning-light-7), var(--el-color-warning-dark-2));
  color: white;
}

.tag-super-admin {
  background-image: linear-gradient(to top right, var(--el-color-success-light-7), var(--el-color-success-dark-2));
  color: white;
}
</style>
