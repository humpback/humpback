<script lang="ts" setup>
import { RuleLength } from "@/models"

const props = withDefaults(
  defineProps<{
    modelValue: string
    size?: "large" | "default" | "small"
    placeholder?: string
    showWordLimit?: boolean
    disabled?: boolean
    clearable?: boolean
  }>(),
  { showWordLimit: true }
)

const emits = defineEmits<{ (e: "update:model-value", data: string): void }>()

const phone = computed({
  get() {
    return props.modelValue
  },
  set(v: string) {
    emits("update:model-value", v.replace(/\D/g, ""))
  }
})
</script>

<template>
  <v-input
    v-model="phone"
    :clearable="props.clearable"
    :disabled="props.disabled"
    :maxlength="RuleLength.Phone.Max"
    :placeholder="props.placeholder"
    :show-word-limit="props.showWordLimit"
    :size="props.size" />
</template>

<style scoped>
.code_box {
  height: 30px;

  span {
    font-size: 12px;
    padding-bottom: 1px;
  }
}
</style>
