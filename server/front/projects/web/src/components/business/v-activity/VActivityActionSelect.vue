<script lang="ts" setup>
import { filter, find } from "lodash-es"

const props = defineProps<{ mode: "group" | "service" | "config" | "registry" | "node" | "user" | "team"; showOutLabel?: boolean }>()

const { t } = useI18n()

const action = defineModel<string>()

const options = computed<Array<{ i18nLabel: string; value: string; modeList: string[] }>>(() => {
  const actions = [
    { i18nLabel: "label.add", value: "Add", modeList: ["group", "service", "registry", "config", "node", "user", "team"] },
    { i18nLabel: "label.update", value: "Update", modeList: ["group", "registry", "config", "user", "team"] },
    { i18nLabel: "label.delete", value: "Delete", modeList: ["group", "service", "registry", "config", "node", "user", "team"] },
    { i18nLabel: "label.updateLabel", value: "UpdateLabel", modeList: ["node"] },
    { i18nLabel: "label.addNode", value: "AddNode", modeList: ["group"] },
    { i18nLabel: "label.removeNode", value: "RemoveNode", modeList: ["group"] },
    { i18nLabel: "label.updateBasic", value: "UpdateBasic", modeList: ["service"] },
    { i18nLabel: "label.updateApplication", value: "UpdateApplication", modeList: ["service"] },
    { i18nLabel: "label.updateDeployment", value: "UpdateDeployment", modeList: ["service"] },
    { i18nLabel: "label.enable", value: "Enable", modeList: ["service", "node"] },
    { i18nLabel: "label.disable", value: "Disable", modeList: ["service", "node"] },
    { i18nLabel: "label.start", value: "Start", modeList: ["service"] },
    { i18nLabel: "label.restart", value: "Restart", modeList: ["service"] },
    { i18nLabel: "label.stop", value: "Stop", modeList: ["service"] }
  ]
  return filter(actions, x => !!find(x.modeList, m => m === props.mode))
})
</script>

<template>
  <v-select
    v-model="action"
    :out-label="props.showOutLabel ? t('label.action') : ''"
    :placeholder="t('placeholder.all')"
    clearable
    filterable
    out-label-width="80px">
    <el-option v-for="item in options" :key="item.value" :label="t(item.i18nLabel)" :value="item.value" />
  </v-select>
</template>

<style lang="scss" scoped></style>
