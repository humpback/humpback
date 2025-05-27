<script lang="ts" setup>
import { RuleLength } from "@/models"

const props = defineProps<{
  modelValue: string
  clearable?: boolean
  disabled?: boolean
  size?: "large" | "default" | "small"
  placeholder?: string
  maxlength?: number | string
  minlength?: number | string
  label?: string
}>()

const emits = defineEmits<{ (e: "update:model-value", v: string): void }>()

const slots = useSlots()

const showPWD = ref(false)

function updateModelValue(v: string) {
  emits("update:model-value", v)
}

function blur(e: FocusEvent) {
  showPWD.value = false
  emits("update:model-value", props.modelValue)
}

function focus(evt: FocusEvent) {
  showPWD.value = true
}
</script>

<template>
  <el-input
    :aria-label="props.label"
    :clearable="props.clearable"
    :disabled="props.disabled"
    :maxlength="props.maxlength || RuleLength.Password.Max"
    :minlength="props.minlength || RuleLength.Password.Min"
    :model-value="props.modelValue || ''"
    :placeholder="props.placeholder"
    :show-password="showPWD"
    :size="props.size"
    type="password"
    @blur="blur"
    @focus="focus"
    @update:modelValue="updateModelValue">
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
