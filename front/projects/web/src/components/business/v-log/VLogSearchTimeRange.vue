<script lang="ts" setup>
const emits = defineEmits<{
  (e: "change"): void
}>()

const { t } = useI18n()

const startTime = defineModel<number>("startTime")
const endTime = defineModel<number>("endTime")

const shortcuts = ref([
  {
    text: t("label.last10Minutes"),
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setMinutes(start.getMinutes() - 10)
      return [start, end]
    }
  },
  {
    text: t("label.lastHour"),
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setHours(start.getHours() - 1)
      return [start, end]
    }
  },
  {
    text: t("label.lastDay"),
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setDate(start.getDay() - 1)
      return [start, end]
    }
  }
])
</script>

<template>
  <v-date-time-range
    v-model:end-time="endTime"
    v-model:start-time="startTime"
    :end-placeholder="t('label.endTime')"
    :out-label="t('label.timeRange')"
    :range-separator="t('label.to')"
    :shortcuts="shortcuts"
    :start-placeholder="t('label.startTime')"
    out-label-width="120px"
    @change="emits('change')" />
</template>

<style lang="scss" scoped></style>
