<script lang="ts" setup>
import type { ISelectProps } from "element-plus"
import { cloneDeep, omit } from "lodash-es"

type Props = Partial<
  ISelectProps & {
    outLabelWidth?: string
    outLabel?: string
  }
>

const props = withDefaults(defineProps<Props>(), {
  reserveKeyword: true,
  teleported: true,
  persistent: true,
  validateEvent: true,
  showArrow: true
})
const emits = defineEmits<{
  (e: "update:modelValue", v: string): void
  (e: "change", v: any): void
  (e: "visible-change", v: boolean): void
  (e: "remove-tag", v: any): void
  (e: "clear"): void
  (e: "blur", v: FocusEvent): void
  (e: "focus", v: FocusEvent): void
}>()

const slots = useSlots()

const attrs = computed(() => {
  const attrs: any = cloneDeep(omit(props, ["outLabelWidth", "outLabel"]))
  return Object.keys(attrs).reduce((acc, key) => {
    if (typeof attrs[key] !== "undefined") {
      acc[key] = props[key]
    }
    return acc
  }, {})
})
const labelClass = computed(() => {
  switch (props.size) {
    case "large":
      return ["select-label", "select-label--large"]
    case "small":
      return ["select-label", "select-label--small"]
    default:
      return ["select-label"]
  }
})
</script>

<template>
  <div class="select-box">
    <div v-if="props.outLabel" :class="labelClass" :style="{ width: props.outLabelWidth }">{{ props.outLabel }}</div>
    <el-select
      class="select-content"
      v-bind="{ ...attrs }"
      @blur="emits('blur', $event)"
      @change="emits('change', $event)"
      @clear="emits('clear')"
      @focus="emits('focus', $event)"
      @visible-change="emits('visible-change', $event)"
      @remove-tag="emits('remove-tag', $event)"
      @update:modelValue="emits('update:modelValue', $event)">
      <template v-if="!!slots.header" #header>
        <slot name="header" />
      </template>
      <template v-if="!!slots.default" #default>
        <slot name="default" />
      </template>
      <template v-if="!!slots.footer" #footer>
        <slot name="footer" />
      </template>
      <template v-if="!!slots.prefix" #prefix>
        <slot name="prefix" />
      </template>
      <template v-if="!!slots.empty" #empty>
        <slot name="empty" />
      </template>
      <template v-if="!!slots.tag" #tag>
        <slot name="tag" />
      </template>
      <template v-if="!!slots.loading" #loading>
        <slot name="loading" />
      </template>
      <template v-if="!!slots.label" #label>
        <slot name="label" />
      </template>
    </el-select>
  </div>
</template>

<style scoped>
.select-box {
  display: flex;
  align-items: center;
  width: 100%;

  .select-label {
    background-color: var(--el-fill-color-light);
    line-height: 24px;
    min-height: 32px;
    color: var(--el-text-color-regular);
    border: 1px solid var(--el-border-color);
    box-sizing: border-box;
    border-top-left-radius: 4px;
    border-bottom-left-radius: 4px;
    border-right: none;
    padding: 0 12px;
    font-size: 14px;
    display: flex;
    align-items: center;
  }

  .select-label--small {
    min-height: 24px;
    line-height: 20px;
    font-size: 12px;
  }

  .select-label--large {
    min-height: 40px;
    line-height: 24px;
    font-size: 14px;
  }

  .select-label + .select-content {
    flex: 1;
    --el-border-radius-base: 0 4px 4px 0;
  }
}
</style>
