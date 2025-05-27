<script lang="ts" setup>
import cronstrue from "cronstrue"
import "cronstrue/locales/en"
import "cronstrue/locales/zh_CN"
import { toLower } from "lodash-es"
import { FormInstance, FormRules } from "element-plus"
import { RulePleaseEnter } from "@/utils"
import { isValidCron } from "cron-validator"

const emits = defineEmits<{
  (e: "change", corn: string, index: number): void
}>()

const { t, locale } = useI18n()

const dialogInfo = ref({
  show: false,
  index: -1,
  cron: "* * * * *"
})

const formRef = useTemplateRef<FormInstance>("formRef")
const rules = ref<FormRules>({
  cron: [
    { required: true, validator: RulePleaseEnter("label.cron"), trigger: "blur" },
    { required: true, validator: checkCron, trigger: "blur" }
  ]
})

const cronTextInfo = computed(() => parseCronToText(dialogInfo.value.cron))

function checkCron(rule: any, value: any, callback: any) {
  const cron = value as string
  if (!isValidCron(cron, { seconds: false })) {
    return callback(new Error(t("rules.formatErrCron")))
  }
  return callback()
}

function parseCronToText(corn: string) {
  try {
    const text = cronstrue.toString(corn, {
      use24HourTimeFormat: true,
      throwExceptionOnParseError: true,
      verbose: false,
      locale: toLower(locale.value) === "zh-cn" ? "zh_CN" : "en"
    })
    return {
      valid: true,
      text: text
    }
  } catch (error) {
    console.error(error)
    return {
      valid: false,
      text: t("tips.invalidCorn")
    }
  }
}

function openDialog(cron: string, index: number) {
  dialogInfo.value.cron = index === -1 ? "* * * * *" : cron
  dialogInfo.value.index = index
  dialogInfo.value.show = true
}

async function save() {
  if (!(await formRef.value?.validate())) {
    return
  }
  dialogInfo.value.show = false
  emits("change", dialogInfo.value.cron, dialogInfo.value.index)
}

defineExpose({ openDialog })
</script>

<template>
  <v-dialog v-model="dialogInfo.show" :close-on-press-escape="false" width="1000px">
    <template #header>{{ dialogInfo.index != -1 ? t("header.editSchedule") : t("header.addSchedule") }}</template>

    <el-form ref="formRef" :model="dialogInfo" :rules="rules" class="form-box mt-5" label-position="left" label-suffix=":" label-width="100px">
      <el-form-item :label="t('label.form')" required>
        <cron-element-plus
          v-model="dialogInfo.cron"
          :buttonProps="{ type: 'info', plain: true }"
          :locale="toLower(locale) === 'zh-cn' ? 'zh' : 'en'"
          format="crontab"
          show-seconds />
      </el-form-item>
      <el-form-item :label="t('label.text')" :rules="rules.cron" prop="cron">
        <v-input v-model="dialogInfo.cron" />
      </el-form-item>
      <el-form-item :label="t('label.description')">
        <el-text :type="cronTextInfo.valid ? 'success' : 'danger'">{{ cronTextInfo.text }}</el-text>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="dialogInfo.show = false">{{ t("btn.cancel") }}</el-button>
      <el-button type="primary" @click="save">
        {{ t("btn.save") }}
      </el-button>
    </template>
  </v-dialog>
</template>

<style lang="scss" scoped>
.form-box {
  :deep(.el-form-item__label) {
    font-weight: 600;
  }
}
</style>
