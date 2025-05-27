<script lang="ts" setup>
import { FormInstance, FormRules } from "element-plus"
import { GenerateUUID, RulePleaseEnter } from "@/utils"
import { filter, map } from "lodash-es"
import { RuleFormat } from "@/models"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()

const isAction = ref(false)
const dialogInfo = ref({
  show: false,
  nodes: [] as Array<{ id: string; ip: string }>
})

const formRef = useTemplateRef<FormInstance>("formRef")
const rules = ref<FormRules>({
  ip: [
    { required: true, validator: RulePleaseEnter("label.ipAddress"), trigger: "blur" },
    { required: true, validator: checkIpAddress, trigger: "blur" }
  ]
})

function checkIpAddress(rule: any, value: any, callback: any) {
  const ipAddress = value as string
  if (!new RegExp(RuleFormat.IPAddress).test(ipAddress)) {
    return callback(new Error(t("rules.invalidIpAddress")))
  }
  if (filter(dialogInfo.value.nodes, x => x.ip === ipAddress).length > 1) {
    return callback(new Error(t("rules.duplicateIPAddress")))
  }
  callback()
}

function addNode() {
  dialogInfo.value.nodes.push({ id: GenerateUUID(), ip: "" })
}

function removeNode(index: number) {
  dialogInfo.value.nodes.splice(index, 1)
}

function open() {
  dialogInfo.value.nodes = [{ id: GenerateUUID(), ip: "" }]
  dialogInfo.value.show = true
}

async function save() {
  if (!(await formRef.value?.validate())) {
    return
  }
  isAction.value = true
  nodeService
    .create(map(dialogInfo.value.nodes, x => x.ip))
    .then(() => {
      ShowSuccessMsg(t("message.addSuccess"))
      dialogInfo.value.show = false
      emits("refresh")
    })
    .finally(() => {
      isAction.value = false
    })
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" :close-on-press-escape="false" width="600px">
    <template #header>{{ t("header.addNodes") }}</template>
    <div class="my-3">
      <el-form ref="formRef" :model="dialogInfo.nodes" :rules="rules" label-position="top" label-width="auto">
        <el-form-item v-for="(node, index) in dialogInfo.nodes" :key="node.id" :prop="`${index}.ip`" :rules="rules.ip">
          <div class="d-flex gap-2 w-100">
            <v-input v-model="dialogInfo.nodes[index].ip" :placeholder="t('placeholder.enterIpAddress')" clearable />
            <el-button plain type="danger" @click="removeNode(index)">
              <el-icon :size="16">
                <IconMdiRemove />
              </el-icon>
            </el-button>
          </div>
        </el-form-item>
        <el-button size="small" type="info" @click="addNode()">
          <el-icon :size="16">
            <IconMdiAdd />
          </el-icon>
          {{ t("btn.addNode") }}
        </el-button>
      </el-form>
    </div>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :disabled="dialogInfo.nodes.length === 0" :loading="isAction" type="primary" @click="save">{{ t("btn.save") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped>
.edit {
  width: 100%;
  height: 500px;
}
</style>
