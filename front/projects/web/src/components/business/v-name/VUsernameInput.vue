<script lang="ts" setup>
import { RuleLength } from "@/models"

const props = withDefaults(
  defineProps<{
    size?: "large" | "default" | "small"
    placeholder?: string
    disabled?: boolean
    clearable?: boolean
    showWordLimit?: boolean
  }>(),
  {
    clearable: true,
    showWordLimit: true,
    placeholder: ""
  }
)

const name = defineModel<string>()
</script>

<template>
  <v-input
    v-model="name"
    :clearable="props.clearable"
    :disabled="props.disabled"
    :maxlength="RuleLength.Username.Max"
    :placeholder="props.placeholder"
    :show-word-limit="props.showWordLimit"
    :size="props.size" />
</template>

<style lang="scss" scoped></style>
