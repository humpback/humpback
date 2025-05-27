<script lang="ts" setup>
import { InputNumberProps } from "element-plus"
import { omit } from "lodash-es"

type Props = Partial<InputNumberProps & { align: "left" | "center" | "right" }>

const props = withDefaults(defineProps<Props>(), {
  min: -Infinity,
  max: Infinity,
  step: 1,
  size: "default",
  controls: true,
  validateEvent: true,
  controlsPosition: "right",
  align: "left"
})

const emits = defineEmits<{
  (e: "update:modelValue", data: number | undefined): void
  (e: "change", currentValue: number | undefined, oldValue: number | undefined): void
}>()

const slots = useSlots()

function updateModelValue(v: number | undefined) {
  emits("update:modelValue", v)
}

function change(currentValue: number | undefined, oldValue: number | undefined) {
  emits("change", currentValue, oldValue)
}

function inputClass() {
  switch (props.align) {
    case "left":
      return "number-left"
    case "center":
      return "number-center"
    case "right":
      return "number-right"
  }
}
</script>
<template>
  <el-input-number :class="inputClass()" v-bind="{ ...omit(props, ['align']) }" @change="change" @update:modelValue="updateModelValue">
    <template v-if="!!slots['decrease-icon']" #decrease-icon>
      <slot name="decrease-icon"></slot>
    </template>
    <template v-if="!!slots['increase-icon']" #increase-icon>
      <slot name="increase-icon"></slot>
    </template>
    <template v-if="!!slots.prefix" #prefix>
      <slot name="prefix"></slot>
    </template>
    <template v-if="!!slots.suffix" #suffix>
      <slot name="suffix"></slot>
    </template>
  </el-input-number>
</template>

<style lang="scss" scoped>
.number-left {
  :deep(.el-input__inner) {
    text-align: left;
  }
}

.number-center {
  :deep(.el-input__inner) {
    text-align: center;
  }
}

.number-right {
  :deep(.el-input__inner) {
    text-align: right;
  }
}
</style>
