<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    modelValue: string
    size?: "large" | "default" | "small"
    clearable?: boolean
    placeholder?: string
  }>(),
  { size: "default", clearable: true, placeholder: "" }
)
const emits = defineEmits<{
  (e: "update:modelValue", data: string): void
  (e: "change", data: string): void
}>()

const { t } = useI18n()

const schedule = computed({
  get: () => props.modelValue || "",
  set: (v: string) => {
    emits("update:modelValue", v || "")
  }
})

const options = computed(() => [
  { label: "label.yes", value: "Yes" },
  { label: "label.no", value: "No" }
])

function change() {
  emits("change", schedule.value)
}
</script>

<template>
  <v-select
    v-model="schedule"
    :clearable="props.clearable"
    :out-label="t('label.schedule')"
    :placeholder="props.placeholder"
    out-label-width="100px"
    @change="change()">
    <el-option :label="t('label.all')" value="" />
    <el-option v-for="item in options" :key="item.value" :label="t(item.label)" :value="item.value" />
  </v-select>
</template>

<style lang="scss" scoped></style>
