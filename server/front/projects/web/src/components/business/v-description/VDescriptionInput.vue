<script lang="ts" setup>
import { RuleLength } from "@/models"

const props = withDefaults(
  defineProps<{
    size?: "large" | "default" | "small"
    placeholder?: string
    row?: number
    autosize?: { minRows?: number; maxRows?: number } | boolean | undefined
    resize?: "none" | "both" | "horizontal" | "vertical" | undefined
    showWordLimit?: boolean
  }>(),
  {
    row: 2,
    resize: "vertical",
    showWordLimit: true,
    autosize: () => {
      return { minRows: 2, maxRows: 4 }
    }
  }
)

const notes = defineModel<string>()
</script>

<template>
  <v-input
    v-model="notes"
    :autosize="props.autosize"
    :maxlength="RuleLength.Description.Max"
    :placeholder="props.placeholder"
    :resize="props.resize"
    :row="2"
    :show-word-limit="props.showWordLimit"
    :size="props.size"
    clearable
    type="textarea" />
</template>

<style lang="scss" scoped></style>
