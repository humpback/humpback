<script lang="ts" setup>
import { find } from "lodash-es"
import { ConfigType } from "@/models"

const props = defineProps<{ configType: number }>()
const { t } = useI18n()

const str = computed<{ label: string; value: number; color?: "primary" | "success" | "warning" | "danger" | "info" } | undefined>(() =>
  find(
    [
      { label: "label.static", value: ConfigType.Static },
      { label: "label.volume", value: ConfigType.Volume, color: "warning" }
    ],
    x => x.value === props.configType
  )
)
</script>

<template>
  <el-text v-if="str" :type="str.color">{{ t(str.label) }}</el-text>
  <span v-else>--</span>
</template>

<style lang="scss" scoped></style>
