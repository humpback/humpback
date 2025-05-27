<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    content?: string
    effect?: "light" | "dark"
    iconSize?: number
    maxWidth?: string
    placement?: "auto" | "auto-start" | "auto-end" | "top" | "bottom" | "right" | "left" | "top-start" | "top-end" | "bottom-start" | "bottom-end"
  }>(),
  { iconSize: 16, effect: "dark", placement: "top-start", maxWidth: "300px" }
)

const slots = useSlots()
const attrs = useAttrs()
</script>

<template>
  <v-tooltip :content="props.content" :effect="props.effect" :max-width="props.maxWidth" :placement="props.placement">
    <template #content>
      <slot v-if="!!slots.content" name="content" />
    </template>
    <el-icon :size="props.iconSize" style="color: var(--el-text-color-primary)" v-bind="{ ...attrs }">
      <IconMdiHelpCircleOutline />
    </el-icon>
  </v-tooltip>
</template>

<style lang="scss" scoped></style>
