<script lang="ts" setup>
import { NodeStatus } from "@/models"
import { toLower } from "lodash-es"

const props = withDefaults(
  defineProps<{
    status: string
    size?: "large" | "default" | "small"
    effect?: "dark" | "light" | "plain"
    round?: boolean
  }>(),
  { size: "default", effect: "dark" }
)
const { t } = useI18n()
</script>

<template>
  <el-tag :effect="props.effect" :round="props.round" :size="props.size" :type="toLower(props.status) === toLower(NodeStatus.Online) ? 'success' : 'danger'">
    {{ toLower(props.status) === toLower(NodeStatus.Online) ? t("label.healthy") : t("label.deadly") }}
  </el-tag>
</template>

<style lang="scss" scoped></style>
