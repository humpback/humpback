<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    modelValue: number
    size?: "large" | "default" | "small"
    clearable?: boolean
    placeholder?: string
  }>(),
  { size: "default", clearable: true, placeholder: "" }
)
const emits = defineEmits<{
  (e: "update:modelValue", data: number): void
  (e: "change", data: number): void
}>()

const { t } = useI18n()

const role = computed({
  get: () => props.modelValue || 0,
  set: (v: number) => {
    emits("update:modelValue", v || 0)
  }
})

const options = computed(() => [
  { label: "role.user", value: UserRole.User },
  { label: "role.admin", value: UserRole.Admin },
  { label: "role.superAdmin", value: UserRole.SuperAdmin }
])

function change() {
  emits("change", role.value)
}
</script>

<template>
  <v-select v-model="role" :clearable="props.clearable" :out-label="t('label.role')" :placeholder="props.placeholder" out-label-width="80px" @change="change()">
    <el-option :label="t('label.all')" :value="0" />
    <el-option v-for="item in options" :key="item.value" :label="t(item.label)" :value="item.value" />
  </v-select>
</template>

<style lang="scss" scoped></style>
