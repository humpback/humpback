<script lang="ts" setup>
import type { AlertProps } from "element-plus"
import { omit } from "lodash-es"

type Props = Partial<
  AlertProps & {
    size?: "small" | "default" | "large"
  }
>

const props = withDefaults(defineProps<Props>(), { size: "default", closable: false, showIcon: true, type: "info" })
const slots = useSlots()
const attr = computed(() => omit(props, ["size"]))
const className = computed(() => {
  switch (props.size) {
    case "small":
      return ["hp-alert", "hp-alert--small"]
    case "large":
      return ["hp-alert", "hp-alert--large"]
    default:
      return ["hp-alert", "hp-alert--default"]
  }
})
</script>

<template>
  <el-alert :class="className" v-bind="attr">
    <template v-if="!!slots.default" #default>
      <slot name="default" />
    </template>
    <template v-if="!!slots.title" #title>
      <slot name="title" />
    </template>
  </el-alert>
</template>

<style lang="scss" scoped>
.hp-alert.el-alert {
  padding: 6px 10px;
  display: flex;
  align-items: start;
  justify-content: left;
  gap: 8px;

  :deep(.el-alert__icon.is-big) {
    margin-right: 0;
  }
}

.hp-alert--small {
  --el-alert-icon-large-size: 16px;
  --el-alert-title-with-description-font-size: 14px;
  --el-alert-description-font-size: 12px;
  line-height: var(--el-alert-icon-large-size);

  :deep(.el-alert__title) {
    line-height: 16px;
  }
}

.hp-alert--default {
  --el-alert-icon-large-size: 18px;
  --el-alert-title-with-description-font-size: 16px;
  --el-alert-description-font-size: 14px;
  line-height: var(--el-alert-icon-large-size);

  :deep(.el-alert__title) {
    line-height: 20px;
  }
}

.hp-alert--large {
  --el-alert-icon-large-size: 21px;
  --el-alert-title-with-description-font-size: 18px;
  --el-alert-description-font-size: 16px;
  line-height: var(--el-alert-icon-large-size);

  :deep(.el-alert__title) {
    line-height: 23px;
  }
}
</style>
