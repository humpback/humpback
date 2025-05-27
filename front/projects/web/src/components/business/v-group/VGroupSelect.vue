<script lang="ts" setup>
import { GroupInfo } from "@/types"
import { map } from "lodash-es"

const props = withDefaults(
  defineProps<{
    options?: GroupInfo[]
    filterable?: boolean
    placeholder?: string
    showOutLabel?: boolean
    size?: "small" | "default" | "large"
    clearable?: boolean
    modelValue?: string
    outLabelWidth?: string
  }>(),
  {
    filterable: true,
    placeholder: "",
    clearable: true,
    outLabelWidth: "80px"
  }
)

const emits = defineEmits<{
  (e: "update:modelValue", v: string): void
  (e: "change"): void
}>()

const { t } = useI18n()

const isLoading = ref(false)
const group = computed({
  get() {
    return props.modelValue || ""
  },
  set(v: any) {
    emits("update:modelValue", v || "")
  }
})
const groupList = ref<GroupInfo[]>([])
const selectOptions = computed(() => map(props.options || groupList.value, x => ({ label: x.groupName, value: x.groupId })))

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

async function getGroups() {
  isLoading.value = true
  return await groupService
    .list()
    .then(list => (groupList.value = list))
    .finally(() => (isLoading.value = false))
}

onMounted(async () => {
  if (!props.options) {
    await getGroups()
  }
})
</script>

<template>
  <div class="select-box">
    <div v-if="props.showOutLabel" :class="labelClass" :style="{ width: props.outLabelWidth }">{{ t("label.group") }}</div>
    <el-select-v2
      v-model="group"
      :clearable="props.clearable"
      :filterable="props.filterable"
      :loading-text="t('message.loading')"
      :options="selectOptions"
      :placeholder="props.placeholder"
      :size="props.size"
      @change="emits('change')">
      <template v-if="isLoading" #prefix>
        <el-button :loading="isLoading" link />
      </template>
    </el-select-v2>
  </div>
</template>

<style lang="scss" scoped>
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

  .select-label + .el-select {
    flex: 1;
    --el-border-radius-base: 0 4px 4px 0;
  }
}
</style>
