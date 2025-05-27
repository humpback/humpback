<script lang="ts" setup>
import { ConfigInfo } from "@/types"
import { cloneDeep } from "lodash-es"

const { t } = useI18n()

const dialogInfo = ref({
  show: false,
  info: {} as ConfigInfo
})

function open(info: ConfigInfo) {
  dialogInfo.value.info = cloneDeep(info)
  dialogInfo.value.show = true
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show">
    <template #header> {{ t("header.config") }}</template>
    <div class="view-content">
      <v-monaco-edit v-model="dialogInfo.info.configValue" read-only />
    </div>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.close") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped>
.view-content {
  height: 500px;
  width: 100%;
}
</style>
