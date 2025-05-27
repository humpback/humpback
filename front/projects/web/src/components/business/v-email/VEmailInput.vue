<script lang="ts" setup>
import { RuleLength } from "@/models"

const props = withDefaults(
  defineProps<{
    size?: "large" | "default" | "small"
    placeholder?: string
    showWordLimit?: boolean
    clearable?: boolean
  }>(),
  { showWordLimit: true }
)
const email = defineModel<string>()
</script>

<template>
  <v-input
    v-model="email"
    :clearable="props.clearable"
    :maxlength="RuleLength.Email.Max"
    :placeholder="props.placeholder"
    :show-word-limit="props.showWordLimit"
    :size="size" />
</template>

<style scoped></style>
