<script lang="ts" setup>
import type { InputProps } from "element-plus"

type Props = Partial<
  InputProps & {
    type: "text" | "textarea" | "password" | "button" | "checkbox" | "file" | "radio"
    modelValue: string
  }
>

const props = withDefaults(defineProps<Props>(), { validateEvent: true })
const emits = defineEmits<{
  (e: "update:modelValue", v: string): void
  (e: "change", v: string): void
  (e: "input", v: string): void
  (e: "blur"): void
}>()

const slots = useSlots()

function updateModelValue(v: string) {
  emits("update:modelValue", v)
}

function blur(e: FocusEvent) {
  emits("update:modelValue", props.modelValue ? props.modelValue.trim() : "")
  emits("blur")
}

function change(v: string) {
  emits("change", v)
}

function input(v: string) {
  emits("input", v)
}
</script>

<template>
  <el-input v-bind="{ ...props }" @blur="blur" @change="change" @input="input" @update:modelValue="updateModelValue">
    <template v-if="!!slots.prefix" #prefix>
      <slot name="prefix"></slot>
    </template>
    <template v-if="!!slots.append" #append>
      <slot name="append"></slot>
    </template>
    <template v-if="!!slots.prepend" #prepend>
      <slot name="prepend"></slot>
    </template>
    <template v-if="!!slots.suffix" #suffix>
      <slot name="suffix"></slot>
    </template>
  </el-input>
</template>

<style scoped></style>
