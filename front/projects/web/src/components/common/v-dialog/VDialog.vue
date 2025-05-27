<script lang="ts" setup>
import type { DialogProps } from "element-plus"
import { omit } from "lodash-es"

type Props = Partial<
  DialogProps & {
    minWidth: string | number
    maxWidth: string | number
  }
>

const props = withDefaults(defineProps<Props>(), {
  destroyOnClose: true,
  draggable: true,
  modal: true,
  showClose: true,
  closeOnClickModal: false,
  appendToBody: true,
  closeOnPressEscape: true
})

const emits = defineEmits<{
  (e: "update:model-value", data: boolean): void
  (e: "closed"): void
}>()

const slots = useSlots()
const pageStore = usePageStore()

const attrs = computed(() => {
  return omit(props, ["width"])
})

const style = computed(() => {
  return { minWidth: props.minWidth, maxWidth: props.maxWidth }
})

const dialogWidth = computed(() => {
  if (props.width) {
    return props.width
  }
  if (pageStore.isSmallScreen) {
    return "100%"
  }
  if (pageStore.isBigScreen) {
    return "50%"
  }
  return "60%"
})

function changeModelValue(v: boolean) {
  emits("update:model-value", v)
}

function closedDialog() {
  emits("closed")
}
</script>

<template>
  <el-dialog
    v-if="props.modelValue"
    :style="style"
    :width="dialogWidth"
    class="custom-dialog"
    v-bind="{ ...attrs }"
    @closed="closedDialog()"
    @update:modelValue="changeModelValue">
    <template v-if="!!slots.header" #header>
      <strong>
        <slot name="header" />
      </strong>
    </template>
    <template v-if="!!slots.default" #default>
      <slot />
    </template>
    <template v-if="!!slots.footer" #footer>
      <slot name="footer" />
    </template>
  </el-dialog>
</template>

<style scoped></style>
