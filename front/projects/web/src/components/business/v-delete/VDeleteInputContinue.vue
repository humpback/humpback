<script lang="ts" setup>
const props = defineProps<{ modelValue: boolean; keywords: string }>()
const emits = defineEmits<{
  (e: "update:modelValue", value: boolean): void
}>()

const { t } = useI18n()

const input = ref("")

function changeInput(v: string) {
  if (v === props.keywords) {
    emits("update:modelValue", true)
  } else {
    emits("update:modelValue", false)
  }
  input.value = v
}
</script>

<template>
  <div>
    <div><strong v-html="t('notify.enterKeywordsContinue', { name: props.keywords })"></strong></div>
    <el-input :modelValue="input" class="mt-3" @update:modelValue="changeInput" />
  </div>
</template>

<style lang="scss" scoped></style>
