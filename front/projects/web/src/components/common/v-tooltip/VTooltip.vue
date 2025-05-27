<script lang="ts" setup>
import { ElTooltipProps } from "element-plus"
import { omit } from "lodash-es"

type Props = Partial<
  ElTooltipProps & {
    width?: string
    maxWidth?: string
  }
>
const props = withDefaults(defineProps<Props>(), {
  effect: "light",
  placement: "top",
  teleported: true,
  showArrow: true,
  enterable: true,
  persistent: false
})
const slots = useSlots()

const attrs = computed(() => {
  return omit(props, ["popperStyle", "popper-style", "width"])
})

const style = computed(() => {
  return [
    {
      width: AddUnitPX(props.width),
      maxWidth: AddUnitPX(props.maxWidth)
    },
    props.popperStyle!
  ]
})

function AddUnitPX(value: string | number | null | undefined, defaultUnit = "px") {
  if (!value) {
    return ""
  }
  if (typeof value === "number") {
    return `${value}${defaultUnit}`
  }
  if ((value as string).toLowerCase().indexOf("px") !== -1) {
    return value
  }
  return `${value}${defaultUnit}`
}
</script>

<template>
  <el-tooltip :popperStyle="style" v-bind="{ ...attrs }">
    <template v-if="!!slots.content" #content>
      <slot name="content" />
    </template>
    <template v-if="!!slots.default" #default>
      <slot />
    </template>
  </el-tooltip>
</template>

<style lang="scss" scoped></style>
