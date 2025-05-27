<script lang="ts" setup>
import { FormInstance, FormRules } from "element-plus"
import { GenerateUUID, RulePleaseEnter } from "@/utils"
import { filter, find, map } from "lodash-es"
import { RuleFormat } from "@/models"
import { NodeInfo } from "@/types"
import { NewCommand } from "@/views/administration/nodes/common.ts"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()

const isAction = ref(false)
const dialogInfo = ref({
  show: false,
  success: false,
  nodes: [] as Array<{ id: string; ip: string }>,
  result: <NodeInfo[]>[],
  activeTab: ""
})

const activeNodeInfo = computed(() => find(dialogInfo.value.result, x => x.ipAddress === dialogInfo.value.activeTab))

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
  dialogInfo.value.result = []
  dialogInfo.value.success = false
  dialogInfo.value.show = true
}

function close() {
  dialogInfo.value.show = false
  if (dialogInfo.value.success) {
    emits("refresh")
  }
}

async function save() {
  if (!(await formRef.value?.validate())) {
    return
  }
  isAction.value = true
  nodeService
    .create(map(dialogInfo.value.nodes, x => x.ip))
    .then(data => {
      dialogInfo.value.result = data
      dialogInfo.value.activeTab = dialogInfo.value.result[0].ipAddress
      dialogInfo.value.success = true
      ShowSuccessMsg(t("message.addSuccess"))
    })
    .finally(() => {
      isAction.value = false
    })
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" :close-on-press-escape="false" width="660px" @closed="close()">
    <template #header>{{ t("header.addNodes") }}</template>
    <div v-if="!dialogInfo.success">
      <v-tips>{{ t("tips.addNodesTips") }}</v-tips>
      <div class="mt-5 mb-3">
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
          <el-button v-if="dialogInfo.nodes.length < 10" size="small" type="info" @click="addNode()">
            <el-icon :size="16">
              <IconMdiAdd />
            </el-icon>
            {{ t("btn.addNode") }}
          </el-button>
        </el-form>
      </div>
    </div>
    <div v-else>
      <v-alert type="success">
        <span v-html="t('tips.addNodesSuccessTips')" />
      </v-alert>
      <div class="node-view-command">
        <div>
          <div v-for="item in dialogInfo.result" :key="item.ipAddress" class="mb-1">
            <el-button :type="dialogInfo.activeTab === item.ipAddress ? 'primary' : undefined" @click="dialogInfo.activeTab = item.ipAddress">
              {{ item.ipAddress }}
            </el-button>
          </div>
        </div>
        <div class="flex-1">
          <v-node-command :command="NewCommand(activeNodeInfo?.ipAddress, activeNodeInfo?.registerInfo.token, false)" style="height: 360px" />
        </div>
      </div>
    </div>
    <template #footer>
      <div v-if="!dialogInfo.success">
        <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
        <el-button :disabled="dialogInfo.nodes.length === 0" :loading="isAction" type="primary" @click="save">{{ t("btn.save") }}</el-button>
      </div>
      <el-button v-else @click="close()">{{ t("btn.close") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped>
.edit {
  width: 100%;
  height: 500px;
}

.node-view-command {
  display: flex;
  align-items: start;
  gap: 12px;
  margin-top: 20px;
}
</style>
