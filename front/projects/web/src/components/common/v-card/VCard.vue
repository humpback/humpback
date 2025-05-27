<script lang="ts" setup>
import type { CardProps } from "element-plus"
import { omit } from "lodash-es"

type Props = Partial<
  CardProps & {
    minHeight?: string
    maxHeight?: string
    maxWidth?: string
    minWidth?: string
    width?: string
    showHeaderColor?: boolean
    round?: boolean
    bodyPaddingNone?: boolean
  }
>

const props = withDefaults(defineProps<Props>(), { shadow: "never", round: true })

const slots = useSlots()

const style = computed(() => {
  return {
    minHeight: props.minHeight,
    maxHeight: props.maxHeight,
    maxWidth: props.maxWidth,
    minWidth: props.minWidth,
    width: props.width
  }
})
</script>

<template>
  <el-card
    :class="{ 'header-color': props.showHeaderColor, 'round': props.round, 'body-padding-none': props.bodyPaddingNone }"
    :style="style"
    class="custom-card"
    v-bind="{ ...omit(props, ['bodyPaddingNone']) }">
    <template v-if="!!slots.header" #header>
      <slot name="header" />
    </template>
    <template v-if="!!slots.default" #default>
      <slot v-if="!!slots.bodyTitle" name="bodyTitle" />
      <slot />
    </template>
    <template v-if="!!slots.footer" #footer>
      <slot name="footer" />
    </template>
  </el-card>
</template>

<style lang="scss" scoped>
.custom-card {
  --el-box-shadow-light: var(--el-box-shadow-lighter);
  border: none;

  :deep(.el-card__header) {
    padding: 12px 20px;
  }
}

.round {
  border-radius: 8px;
}

.header-color {
  :deep(.el-card__header) {
    background-color: var(--hp-card-header-bg-color);
    color: #2b2b2b;
  }
}

.body-padding-none {
  :deep(.el-card__body) {
    padding: 0;
  }
}

.card-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 20px;
}
</style>
