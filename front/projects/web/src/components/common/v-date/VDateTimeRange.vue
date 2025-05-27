<script lang="ts" setup>
import type { DatePickerProps } from "element-plus"
import { cloneDeep, omit } from "lodash-es"

type Props = Partial<
  DatePickerProps & {
    outLabelWidth?: string
    outLabel?: string
  }
>

const props = withDefaults(defineProps<Props>(), {
  type: "datetimerange",
  valueFormat: "x",
  editable: true,
  clearable: true,
  teleported: true,
  showNow: true
})

const emits = defineEmits<{
  (e: "change"): void
  (e: "blur", v: FocusEvent): void
  (e: "focus", v: FocusEvent): void
}>()

const slots = useSlots()

const startTime = defineModel<number>("startTime")
const endTime = defineModel<number>("endTime")

const attrs = computed(() => {
  const attrs: any = cloneDeep(omit(props, ["modelValue", "outLabelWidth", "outLabel"]))
  return Object.keys(attrs).reduce((acc, key) => {
    if (typeof attrs[key] !== "undefined") {
      acc[key] = props[key]
    }
    return acc
  }, {})
})

const dateTimeRange = computed<any>({
  get() {
    if (startTime.value && endTime.value) {
      return [startTime.value, endTime.value]
    }
    return []
  },
  set(value) {
    if (Array.isArray(value) && value.length == 2) {
      startTime.value = value[0]
      endTime.value = value[1]
    } else {
      startTime.value = 0
      endTime.value = 0
    }
    emits("change")
  }
})
</script>

<template>
  <div class="hp-date-picker">
    <div v-if="props.outLabel" :style="{ width: props.outLabelWidth }" class="out-label">{{ props.outLabel }}</div>
    <el-date-picker v-model="dateTimeRange" class="time-range" v-bind="{ ...attrs }" @blur="emits('blur', $event)" @focus="emits('focus', $event)">
      <template v-if="!!slots.default" #default>
        <slot name="default" />
      </template>
      <template v-if="!!slots['range-separator']" #range-separator>
        <slot name="range-separator" />
      </template>
      <template v-if="!!slots['prev-month']" #prev-month>
        <slot name="prev-month" />
      </template>
      <template v-if="!!slots['next-month']" #next-month>
        <slot name="next-month" />
      </template>
      <template v-if="!!slots['prev-year']" #prev-year>
        <slot name="prev-year" />
      </template>
      <template v-if="!!slots['next-year']" #next-year>
        <slot name="next-year" />
      </template>
    </el-date-picker>
  </div>
</template>

<style lang="scss" scoped>
.hp-date-picker {
  display: flex;
  align-items: center;

  .out-label {
    background-color: var(--el-fill-color-light);
    line-height: 24px;
    min-height: 32px;
    font-size: 14px;
    color: var(--el-text-color-regular);
    border: 1px solid var(--el-border-color);
    box-sizing: border-box;
    border-top-left-radius: 4px;
    border-bottom-left-radius: 4px;
    border-right: none;
    padding: 0 12px;
    display: flex;
    align-items: center;

    & + :deep(.time-range.el-date-editor) {
      border-bottom-left-radius: unset;
      border-top-left-radius: unset;
    }
  }
}
</style>
