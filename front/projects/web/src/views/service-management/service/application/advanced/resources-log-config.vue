<script lang="ts" setup>
import { ServiceResourcesInfo } from "@/types"
import { FormInstance, FormRules } from "element-plus"
import { RulePleaseEnter } from "@/utils"
import { filter, toLower } from "lodash-es"
import { RuleLength } from "@/models"

const props = defineProps<{ hasValid?: boolean }>()
const emits = defineEmits<{
  (e: "check"): void
}>()

const { t } = useI18n()

const resources = defineModel<ServiceResourcesInfo>("resources")
const logConfig = defineModel<{
  type: string
  options: Array<{ id: string; name: string; value: string }>
}>("logConfig")

const driverOptions = ref<string[]>(["awslogs", "fluentd", "gcplogs", "gelf", "journald", "json-file", "local", "splunk", "syslog"])

const formRef = useTemplateRef<FormInstance>("formRef")
const rules = ref<FormRules>({
  name: [
    { required: true, validator: RulePleaseEnter("label.name"), trigger: "blur" },
    { required: true, validator: checkLogConfigName, trigger: "blur" }
  ],
  value: [{ required: true, validator: RulePleaseEnter("label.value"), trigger: "blur" }]
})

function checkLogConfigName(rule: any, value: any, callback: any) {
  const name = value as string
  if (filter(logConfig.value?.options, x => toLower(x.name) === name).length > 1) {
    return callback(new Error(`${t("rules.duplicate")} ${t("label.name")}`))
  }
  return callback()
}

function formatMemoryLimit(v: number) {
  if (!v) {
    return t("label.unLimited")
  }
  return `${v} ${t("label.mb")}`
}

function formatMemoryReservation(v: number) {
  if (!v || (resources.value?.memory && v > resources.value?.memory)) {
    return t("label.unLimited")
  }
  return `${v} ${t("label.mb")}`
}

function formatMaxCpuUsage(v: number) {
  if (!v || v === 100) {
    return t("label.unLimited")
  }
  return `${v}%`
}

function addOption() {
  logConfig.value?.options?.push({ id: GenerateUUID(), name: "", value: "" })
}

function removeOption(index: number) {
  logConfig.value?.options?.splice(index, 1)
  emits("check")
}

function changeLogConfigType() {
  if (!logConfig.value?.type) {
    logConfig.value!.options = []
    emits("check")
  }
}

async function validate() {
  return await formRef.value?.validate()
}

onMounted(() => {
  if (props.hasValid) {
    validate()
  }
})

defineExpose({ validate })
</script>

<template>
  <div class="mb-5">
    <el-text class="f-bold">{{ t("label.resources") }}</el-text>
    <div class="mt-3">
      <v-tips>{{ t("tips.resourcesTips") }}</v-tips>
    </div>
  </div>
  <el-form label-position="left" label-width="auto">
    <el-form-item>
      <template #label>
        <div class="d-flex gap-1">
          <span>{{ t("label.memoryLimit") }}</span>
          <v-help-tooltip :content="t('tips.memoryLimitTips')" max-width="400px" />
        </div>
      </template>
      <div class="d-flex gap-5 w-100 px-5 flex-wrap">
        <div class="flex-1 slider-line">
          <el-slider v-model="resources!.memory" :format-tooltip="formatMemoryLimit" :max="20480" :min="0" :step="20" />
        </div>
        <div>
          <v-input-number
            v-model="resources!.memory"
            :controls="false"
            :max="RuleLength.MemoryLimit.Max"
            :min="RuleLength.MemoryLimit.Min"
            :precision="0"
            :step="20">
            <template #suffix>{{ t("label.mb") }}</template>
          </v-input-number>
        </div>
      </div>
    </el-form-item>

    <el-form-item>
      <template #label>
        <div class="d-flex gap-1">
          <span>{{ t("label.memoryReservation") }}</span>
          <v-help-tooltip :content="t('tips.memoryReservationTips')" max-width="400px" />
        </div>
      </template>
      <div class="d-flex gap-5 flex-wrap w-100 px-5">
        <div class="flex-1 slider-line">
          <el-slider v-model="resources!.memoryReservation" :format-tooltip="formatMemoryReservation" :max="20480" :min="0" :step="20" />
        </div>
        <div>
          <v-input-number
            v-model="resources!.memoryReservation"
            :controls="false"
            :max="RuleLength.MemoryReservation.Max"
            :min="RuleLength.MemoryReservation.Min"
            :precision="0"
            :step="20">
            <template #suffix>{{ t("label.mb") }}</template>
          </v-input-number>
        </div>
      </div>
    </el-form-item>

    <el-form-item>
      <template #label>
        <div class="d-flex gap-1">
          <span>{{ t("label.maximumCpuUsage") }}</span>
          <v-help-tooltip :content="t('tips.maxCpuUsageTips')" max-width="400px" />
        </div>
      </template>
      <div class="d-flex gap-5 flex-wrap w-100 px-5">
        <div class="flex-1 slider-line">
          <el-slider v-model="resources!.maxCpuUsage" :format-tooltip="formatMaxCpuUsage" :max="100" :min="0" />
        </div>
        <div>
          <v-input-number v-model="resources!.maxCpuUsage" :controls="false" :max="RuleLength.MaxCpuUsage.Max" :min="RuleLength.MaxCpuUsage.Min" :precision="0">
            <template #suffix>%</template>
          </v-input-number>
        </div>
      </div>
    </el-form-item>
  </el-form>

  <el-divider />
  <div class="my-5">
    <el-text class="f-bold">{{ t("label.logConfig") }}</el-text>
    <div class="mt-3">
      <v-tips>
        <i18n-t keypath="tips.logDriverTips" scope="global">
          <template #document>
            <el-link href="https://docs.docker.com/engine/logging/configure/#supported-logging-drivers" target="_blank" type="primary">
              {{ t("label.dockerDocumentation") }}
            </el-link>
          </template>
        </i18n-t>
      </v-tips>
    </div>
  </div>
  <el-form ref="formRef" :model="logConfig" :rules="rules" label-position="left" label-width="100px">
    <el-form-item :label="t('label.driver')">
      <el-select
        v-model="logConfig!.type"
        :placeholder="t('placeholder.defaultLoggingDriver')"
        clearable
        style="max-width: 300px"
        @change="changeLogConfigType()">
        <el-option v-for="item in driverOptions" :key="item" :label="item" :value="item" />
      </el-select>
    </el-form-item>

    <el-form-item :label="t('label.options')">
      <div>
        <div v-for="(log, index) in logConfig?.options" :key="log.id" class="d-flex gap-3 mb-4">
          <el-form-item :prop="`options.${index}.name`" :rules="rules.name">
            <v-input v-model="logConfig!.options[index].name" @blur="emits('check')">
              <template #prepend>{{ t("label.name") }}</template>
            </v-input>
          </el-form-item>
          <el-form-item :prop="`options.${index}.value`" :rules="rules.value">
            <v-input v-model="logConfig!.options[index].value" @blur="emits('check')">
              <template #prepend>{{ t("label.value") }}</template>
            </v-input>
          </el-form-item>
          <el-form-item>
            <el-button plain style="padding: 4px 12px" text type="danger" @click="removeOption(index)">
              <el-icon :size="26">
                <IconMdiTrash />
              </el-icon>
            </el-button>
          </el-form-item>
        </div>
        <div>
          <el-button :disabled="!logConfig!.type" size="small" type="info" @click="addOption()">
            <el-icon :size="16">
              <IconMdiAdd />
            </el-icon>
            {{ t("btn.addOption") }}
          </el-button>
        </div>
      </div>
    </el-form-item>
  </el-form>
</template>

<style lang="scss" scoped>
.slider-line {
  max-width: 600px;
  min-width: 160px;
}
</style>
