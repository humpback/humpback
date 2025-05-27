<script lang="ts" setup>
import { NodeSwitch, ServiceStatus } from "@/models"

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

const status = computed({
  get: () => props.modelValue || "",
  set: (v: string) => {
    emits("update:modelValue", v || "")
  }
})

const options = computed<
  Array<{
    label: string
    value: string
    type: "info" | "success" | "warning" | "primary" | "danger"
  }>
>(() => [
  { label: "serviceStatus.disabled", value: NodeSwitch.Disabled, type: "info" },
  { label: "serviceStatus.assigning", value: ServiceStatus.ServiceStatusNotReady, type: "primary" },
  { label: "serviceStatus.warning", value: ServiceStatus.ServiceStatusFailed, type: "warning" },
  { label: "serviceStatus.running", value: ServiceStatus.ServiceStatusRunning, type: "success" }
])

function change() {
  emits("change", status.value)
}
</script>

<template>
  <v-select
    v-model="status"
    :clearable="props.clearable"
    :out-label="t('label.status')"
    :placeholder="props.placeholder"
    out-label-width="80px"
    @change="change()">
    <el-option :label="t('label.all')" value="" />
    <el-option v-for="item in options" :key="item.value" :label="t(item.label)" :value="item.value">
      <el-text :type="item.type">{{ t(item.label) }}</el-text>
    </el-option>
  </v-select>
</template>

<style lang="scss" scoped></style>
