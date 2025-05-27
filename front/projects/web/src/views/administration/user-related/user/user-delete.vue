<script lang="ts" setup>
import { cloneDeep } from "lodash-es"
import { UserInfo } from "@/types"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()

const isAction = ref(false)
const isChecked = ref(false)
const dialogInfo = ref({
  show: false,
  info: {} as UserInfo
})

function open(info: UserInfo) {
  isChecked.value = false
  dialogInfo.value.info = cloneDeep(info)
  dialogInfo.value.show = true
}

async function confirmDelete() {
  if (!isChecked.value) {
    return
  }
  isAction.value = true
  return await userService
    .delete(dialogInfo.value.info.userId)
    .then(() => {
      ShowSuccessMsg(t("message.deleteSuccess"))
      dialogInfo.value.show = false
      emits("refresh")
    })
    .finally(() => (isAction.value = false))
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" width="600px">
    <template #header>{{ t("header.deleteUser") }}</template>
    <div class="my-3 f-bold">{{ t("notify.delete") }}</div>
    <v-delete-input-continue v-model="isChecked" :keywords="dialogInfo.info.username" class="mt-5" />
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :disabled="!isChecked" :loading="isAction" type="danger" @click="confirmDelete">{{ t("btn.delete") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
