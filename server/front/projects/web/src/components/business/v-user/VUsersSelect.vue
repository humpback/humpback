<script lang="ts" setup>
import { UserInfo } from "@/types"
import { map } from "lodash-es"

const props = withDefaults(
  defineProps<{
    options?: UserInfo[]
    filterable?: boolean
    multiple?: boolean
    placeholder?: string
    showOutLabel?: boolean
    outLabelWidth?: string
    outLabel?: string
    size?: "small" | "default" | "large"
    showFooter?: boolean
    clearable?: boolean
  }>(),
  {
    filterable: true,
    multiple: true,
    placeholder: "",
    outLabelWidth: "80px"
  }
)

const { t } = useI18n()
const userStore = useUserStore()

const isLoading = ref(false)
const users = defineModel<string[] | string>()
const userList = ref<UserInfo[]>([])
const selectOptions = computed(() => map(props.options || userList.value, x => ({ label: x.username, value: x.userId })))

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

async function getUsers() {
  isLoading.value = true
  return await userService
    .list()
    .then(res => {
      userList.value = res
    })
    .finally(() => (isLoading.value = false))
}

onMounted(async () => {
  if (!props.options) {
    await getUsers()
  }
})
</script>

<template>
  <div class="select-box">
    <div v-if="props.showOutLabel" :class="labelClass" :style="{ width: props.outLabelWidth }">{{ props.outLabel || t("label.user") }}</div>
    <el-select-v2
      v-model="users"
      :clearable="props.clearable"
      :filterable="props.filterable"
      :loading-text="t('message.loading')"
      :multiple="props.multiple"
      :options="selectOptions"
      :placeholder="props.placeholder"
      :size="props.size">
      <template v-if="isLoading" #prefix>
        <el-button :loading="isLoading" link />
      </template>
      <template v-if="userStore.isAdmin && props.showFooter" #footer>
        <div class="text-align-right">
          <el-link href="/ws/user-related/users" target="_blank" type="primary">
            <strong>{{ t("btn.goToAddUser") }}</strong>
          </el-link>
        </div>
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
