<script lang="ts" setup>
import { find } from "lodash-es"

const props = defineProps<{ role: number }>()
const { t } = useI18n()

const str = computed<
  | {
      label: string
      value: number
      color?: "primary" | "success" | "warning" | "danger" | "info"
    }
  | undefined
>(() =>
  find(
    [
      { label: "role.user", value: UserRole.User },
      { label: "role.admin", value: UserRole.Admin, color: "warning" },
      { label: "role.superAdmin", value: UserRole.SuperAdmin, color: "success" }
    ],
    x => x.value === props.role
  )
)
</script>

<template>
  <el-text v-if="str" :type="str.color">{{ t(str.label) }}</el-text>
  <span v-else>--</span>
</template>

<style lang="scss" scoped></style>
