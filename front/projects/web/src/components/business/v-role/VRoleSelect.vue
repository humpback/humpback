<script lang="ts" setup>
import { find } from "lodash-es"

const props = withDefaults(
  defineProps<{
    modelValue: number
    size?: "large" | "default" | "small"
    clearable?: boolean
    placeholder?: string
    onlyShowRoles?: number[]
  }>(),
  { size: "default", clearable: false, placeholder: "" }
)
const emits = defineEmits<{
  (e: "update:modelValue", data: number): void
  (e: "change", data: number): void
}>()

const { t } = useI18n()

const role = computed({
  get: () => props.modelValue,
  set: (v: number) => {
    emits("update:modelValue", v)
  }
})

const options = computed(() => {
  const data = [
    { label: "role.user", value: UserRole.User },
    { label: "role.admin", value: UserRole.Admin },
    { label: "role.superAdmin", value: UserRole.SuperAdmin }
  ]
  if (Array.isArray(props.onlyShowRoles) && props.onlyShowRoles.length > 0) {
    return data.filter(x => find(props.onlyShowRoles, r => r === x.value))
  }
  return data
})

function change() {
  emits("change", role.value)
}
</script>

<template>
  <el-select v-model="role" :clearable="props.clearable" :placeholder="props.placeholder" @change="change()">
    <el-option v-for="item in options" :key="item.value" :label="t(item.label)" :value="item.value" />
  </el-select>
</template>

<style lang="scss" scoped></style>
