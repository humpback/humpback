<script lang="ts" setup>
import { ConfigInfo, NewConfigEmptyInfo } from "@/types"
import { cloneDeep } from "lodash-es"
import { FormInstance, FormRules } from "element-plus"
import { RulePleaseEnter } from "@/utils"
import { RuleLength } from "@/models"

const emits = defineEmits<{
  (e: "refresh"): void
}>()

const { t } = useI18n()

const isAction = ref(false)
const dialogInfo = ref({
  show: false,
  info: {} as ConfigInfo
})

const formRef = useTemplateRef<FormInstance>("formRef")
const rules = ref<FormRules>({
  configName: [
    { required: true, validator: RulePleaseEnter("label.name"), trigger: "blur" },
    { required: true, validator: RuleLimitRange(RuleLength.ConfigName.Min, RuleLength.ConfigName.Max), trigger: "blur" }
  ],
  configValue: [
    { required: true, validator: RulePleaseEnter("label.value"), trigger: "blur" },
    { required: true, validator: checkConfigValue, trigger: "blur" }
  ],
  description: [{ validator: RuleLimitMax(RuleLength.Description.Max), trigger: "blur" }]
})

function checkConfigValue(rule: any, value: any, callback: any) {
  const data = value as string
  const limit = dialogInfo.value.info.configType === ConfigType.Static ? RuleLength.ConfigValue.Max / 2 : RuleLength.ConfigValue.Max
  if (data.length > limit) {
    return callback(new Error(t("rules.limitLengthMax", { max: limit })))
  }
  callback()
}

function changType() {
  formRef.value?.clearValidate(["configValue"])
}

function open(info?: ConfigInfo) {
  dialogInfo.value.info = info ? cloneDeep(info) : NewConfigEmptyInfo()
  dialogInfo.value.show = true
}

async function save() {
  if (!(await formRef.value?.validate())) {
    return
  }
  const body: any = {
    configName: dialogInfo.value.info.configName,
    configValue: dialogInfo.value.info.configValue,
    configType: dialogInfo.value.info.configType,
    description: dialogInfo.value.info.description
  }
  isAction.value = true
  if (dialogInfo.value.info.configId) {
    body.configId = dialogInfo.value.info.configId
    configService
      .update(body)
      .then(() => {
        ShowSuccessMsg(t("message.saveSuccess"))
        dialogInfo.value.show = false
        emits("refresh")
      })
      .finally(() => (isAction.value = false))
  } else {
    configService
      .create(body)
      .then(() => {
        ShowSuccessMsg(t("message.addSuccess"))
        dialogInfo.value.show = false
        emits("refresh")
      })
      .finally(() => (isAction.value = false))
  }
}

defineExpose({ open })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" :close-on-press-escape="false" width="1200px">
    <template #header>{{ dialogInfo.info.configId ? t("header.editConfig") : t("header.addConfig") }}</template>
    <div class="my-3">
      <el-form ref="formRef" :model="dialogInfo.info" :rules="rules" label-position="top" label-width="auto">
        <el-form-item :label="t('label.name')" prop="configName">
          <v-input v-model="dialogInfo.info.configName" :maxlength="RuleLength.ConfigName.Max" clearable show-word-limit />
        </el-form-item>

        <el-form-item :label="t('label.description')" prop="description">
          <v-description-input v-model="dialogInfo.info.description" />
        </el-form-item>

        <el-form-item :label="t('label.value')">
          <el-radio-group v-model="dialogInfo.info.configType" @change="changType">
            <el-radio :value="ConfigType.Static">{{ t("label.static") }}</el-radio>
            <el-radio :value="ConfigType.Volume">{{ t("label.volume") }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item prop="configValue">
          <v-input
            v-if="dialogInfo.info.configType === ConfigType.Static"
            v-model="dialogInfo.info.configValue"
            :autosize="{ minRows: 2, maxRows: 8 }"
            :maxlength="RuleLength.ConfigValue.Max / 2"
            resize="vertical"
            show-word-limit
            type="textarea" />
          <div v-if="dialogInfo.info.configType === ConfigType.Volume" class="edit">
            <v-monaco-edit v-model="dialogInfo.info.configValue" />
          </div>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button :loading="isAction" type="primary" @click="save">{{ t("btn.save") }}</el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped>
.edit {
  width: 100%;
  height: 500px;
}
</style>
