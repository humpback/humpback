<script lang="ts" setup>
import { NodeInfo } from "@/types"
import { cloneDeep } from "lodash-es"
import { NewCommand } from "./common.ts"

const { t } = useI18n()

const dialogInfo = ref({
  show: false,
  isUninstall: false,
  info: {} as NodeInfo
})

const command = computed(() => NewCommand(dialogInfo.value.info?.ipAddress, dialogInfo.value.isUninstall))

function open(info: NodeInfo) {
  dialogInfo.value.info = cloneDeep(info)
  dialogInfo.value.isUninstall = false
  dialogInfo.value.show = true
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" width="600px">
    <template #header>{{ t("header.command") }}</template>
    <div class="mt-2">
      <div class="px-1 d-flex gap-3">
        <div class="flex-1 f-bold" style="color: var(--el-text-color-primary)">{{ dialogInfo.info.ipAddress }}</div>
        <el-radio-group v-model="dialogInfo.isUninstall">
          <el-radio :value="false">{{ t("label.install") }}</el-radio>
          <el-radio :value="true">{{ t("label.uninstall") }}</el-radio>
        </el-radio-group>
      </div>
      <div class="mt-2">
        <v-node-command :command="command" />
      </div>
    </div>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.close") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
