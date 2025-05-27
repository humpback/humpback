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

async function save() {
  isAction.value = true
  return await nodeService
    .updateSwitch({
      nodeId: dialogInfo.value.info.nodeId,
      enable: !dialogInfo.value.info.isEnable
    })
    .then(() => {
      ShowSuccessMsg(t("message.saveSuccess"))
      dialogInfo.value.show = false
      emits("refresh")
    })
    .finally(() => (isAction.value = false))
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" width="600px">
    <template #header>{{ dialogInfo.info.isEnable ? t("header.disableNode") : t("header.enableNode") }}</template>
    <div class="my-3">
      <strong v-if="dialogInfo.info.isEnable" v-html="t('notify.disableNode', { ip: dialogInfo.info.ipAddress })" />
      <strong v-else v-html="t('notify.enableNode', { ip: dialogInfo.info.ipAddress })" />
    </div>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button v-if="dialogInfo.info.isEnable" :loading="isAction" type="info" @click="save">{{ t("btn.disable") }}</el-button>
      <el-button v-else :loading="isAction" type="success" @click="save">{{ t("btn.enable") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
