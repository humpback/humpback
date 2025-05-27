<script lang="ts" setup>
import { FormInstance, FormRules } from "element-plus"
import { GenerateUUID, RulePleaseEnter } from "@/utils"
import { NodeInfo } from "@/types"
import { filter, map, toLower } from "lodash-es"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()

const isAction = ref(false)
const dialogInfo = ref({
  show: false,
  info: {} as NodeInfo,
  labels: [] as Array<{ id: string; key: string; value: string }>
})

const formRef = useTemplateRef<FormInstance>("formRef")
const rules = ref<FormRules>({
  key: [
    { required: true, validator: RulePleaseEnter("label.labelKey"), trigger: "blur" },
    { required: true, validator: checkLabelKey, trigger: "blur" }
  ],
  value: [{ required: true, validator: RulePleaseEnter("label.labelValue"), trigger: "blur" }]
})

function checkLabelKey(rule: any, value: any, callback: any) {
  const key = value as string
  if (filter(dialogInfo.value.labels, x => toLower(x.key) === toLower(key)).length > 1) {
    return callback(new Error(t("rules.duplicateLabelKey")))
  }
  callback()
}

function addLabel() {
  dialogInfo.value.labels.push({ id: GenerateUUID(), key: "", value: "" })
}

function removeLabel(index: number) {
  dialogInfo.value.labels.splice(index, 1)
}

function open(node: NodeInfo) {
  dialogInfo.value.info = node
  dialogInfo.value.labels = map(Object.keys(dialogInfo.value.info.labels).sort(), x => ({ id: GenerateUUID(), key: x, value: dialogInfo.value.info.labels[x] }))
  if (dialogInfo.value.labels.length === 0) {
    addLabel()
  }
  dialogInfo.value.show = true
}

async function save() {
  if (!(await formRef.value?.validate())) {
    return
  }
  const body: any = {
    nodeId: dialogInfo.value.info.nodeId,
    labels: {} as any
  }
  map(dialogInfo.value.labels, x => {
    body.labels[x.key] = x.value
  })
  isAction.value = true
  nodeService
    .updateLabel(body)
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
  <v-dialog v-model="dialogInfo.show" :close-on-press-escape="false" width="600px">
    <template #header>{{ t("header.editLabel") }}</template>
    <el-text class="f-bold" style="color: var(--el-text-color-primary)">{{ dialogInfo.info.ipAddress }}</el-text>
    <div class="my-3">
      <el-scrollbar always max-height="600px" style="box-sizing: border-box">
        <el-form ref="formRef" :model="dialogInfo.labels" :rules="rules" label-position="top" label-width="auto">
          <el-form-item v-for="(node, index) in dialogInfo.labels" :key="node.id">
            <div class="d-flex gap-3 w-100 pr-3">
              <el-form-item :prop="`${index}.key`" :rules="rules.key" style="flex: 4">
                <div class="d-flex gap-2 w-100">
                  <v-input v-model="dialogInfo.labels[index].key" :placeholder="t('placeholder.labelKey')" clearable />
                </div>
              </el-form-item>
              <el-form-item :prop="`${index}.value`" :rules="rules.value" style="flex: 5">
                <div class="d-flex gap-2 w-100">
                  <v-input v-model="dialogInfo.labels[index].value" :placeholder="t('placeholder.labelValue')" clearable />
                  <el-button plain type="danger" @click="removeLabel(index)">
                    <el-icon :size="16">
                      <IconMdiRemove />
                    </el-icon>
                  </el-button>
                </div>
              </el-form-item>
            </div>
          </el-form-item>
        </el-form>
      </el-scrollbar>
      <el-button size="small" type="info" @click="addLabel()">
        <el-icon :size="16">
          <IconMdiAdd />
        </el-icon>
        {{ t("btn.addLabel") }}
      </el-button>
    </div>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :loading="isAction" type="primary" @click="save">{{ t("btn.save") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped></style>
