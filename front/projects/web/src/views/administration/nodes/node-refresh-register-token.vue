<script lang="ts" setup>
import { cloneDeep } from "lodash-es"
import { NodeInfo } from "@/types"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()

const isAction = ref(false)
const dialogInfo = ref({
  show: false,
  info: {} as NodeInfo
})

function open(info: NodeInfo) {
  dialogInfo.value.info = cloneDeep(info)
  dialogInfo.value.show = true
}

async function confirm() {
  isAction.value = true
  return await nodeService
    .refreshRegisterToken(dialogInfo.value.info.nodeId)
    .then(() => {
      ShowSuccessMsg(t("message.refreshSuccess"))
      dialogInfo.value.show = false
      emits("refresh")
    })
    .finally(() => (isAction.value = false))
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" width="600px">
    <template #header>{{ t("header.refreshRegisterToken") }}</template>
    <div class="my-3 f-bold" v-html="t('notify.refreshRegisterToken', { ip: dialogInfo.info.ipAddress })" />
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :loading="isAction" type="primary" @click="confirm">{{ t("btn.confirm") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
