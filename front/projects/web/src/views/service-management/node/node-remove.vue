<script lang="ts" setup>
import { cloneDeep } from "lodash-es"
import { NodeInfo } from "@/types"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()
const stateStore = useStateStore()

const isAction = ref(false)
const isChecked = ref(false)
const dialogInfo = ref({
  show: false,
  info: {} as NodeInfo
})

function open(info: NodeInfo) {
  isChecked.value = false
  dialogInfo.value.info = cloneDeep(info)
  dialogInfo.value.show = true
}

async function remove() {
  if (!isChecked.value) {
    return
  }
  isAction.value = true
  return await groupService
    .updateNodes(stateStore.getGroup()!.groupId, {
      isDelete: true,
      nodes: [dialogInfo.value.info.nodeId]
    })
    .then(() => {
      ShowSuccessMsg(t("message.removeSuccess"))
      dialogInfo.value.show = false
      emits("refresh")
    })
    .finally(() => (isAction.value = false))
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" width="600px">
    <template #header>{{ t("header.removeNode") }}</template>
    <div class="my-3 f-bold">{{ t("notify.removeNodeFromGroup") }}</div>
    <v-delete-input-continue v-model="isChecked" :keywords="dialogInfo.info.ipAddress" class="mt-5" />
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :disabled="!isChecked" :loading="isAction" type="danger" @click="remove">{{ t("btn.remove") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
